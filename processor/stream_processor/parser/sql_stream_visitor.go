package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

// check that type implements SqlVisitor
var _ SqlVisitor = (*SqlStreamVisitor)(nil)

type SqlStreamVisitor struct {
	antlr.ParseTreeVisitor
	logger *zap.Logger

	query string

	tumblingPeriod time.Duration

	logRecords    plog.LogRecordSlice
	currentRecord plog.LogRecord

	in        <-chan plog.LogRecordSlice
	out       chan<- plog.LogRecordSlice
	outErr    chan<- error
	shutdownC chan struct{}
}

func NewSqlStreamVisitor(query string, in <-chan plog.LogRecordSlice, out chan<- plog.LogRecordSlice, outErr chan<- error, logger *zap.Logger) *SqlStreamVisitor {

	visitor := &SqlStreamVisitor{
		ParseTreeVisitor: &SqlStreamVisitor{},
		query:            query,
		logger:           logger,
		in:               in,
		out:              out,
		outErr:           outErr,
		shutdownC:        make(chan struct{}, 1),
	}

	isTumbling, val := IsTumblingQuery(query)
	if isTumbling {
		visitor.tumblingPeriod = time.Millisecond * time.Duration(val)
		visitor.startWindowTumblingLoop()
		return visitor
	}

	go visitor.startSimpleWorkerLoop()
	return visitor
}

func (v *SqlStreamVisitor) Stop() {
	close(v.shutdownC)
}

func (v *SqlStreamVisitor) startSimpleWorkerLoop() {

	for {
		select {
		case v.logRecords = <-v.in:
			if err := v.runQuery(); err != nil {
				v.outErr <- err
				break
			}
			v.out <- v.logRecords
			v.logRecords = plog.NewLogRecordSlice()

		case <-v.shutdownC:
			v.logger.Debug("sql engine shutting down")
			return

		}

	}
}

func (v *SqlStreamVisitor) startWindowTumblingLoop() {
	v.logRecords = plog.NewLogRecordSlice()
	ticker := time.NewTicker(v.tumblingPeriod)

	go func() {
		for {
			select {
			case ls := <-v.in:
				ls.MoveAndAppendTo(v.logRecords)
			case <-v.shutdownC:
				v.logger.Debug("sql engine shutting down")
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ticker.C:
				if v.logRecords.Len() > 0 {
					if err := v.runQuery(); err != nil {
						v.outErr <- err
						break
					}
					v.out <- v.logRecords
					v.logRecords = plog.NewLogRecordSlice()
				}
			case <-v.shutdownC:
				v.logger.Debug("sql engine shutting down")
				return

			}
		}
	}()
}

func (v *SqlStreamVisitor) runQuery() error {
	is := antlr.NewInputStream(v.query)
	lexer := NewSqlLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := NewSqlParser(stream)
	res := v.Visit(parser.SqlQuery())
	switch res := res.(type) {
	case error:
		return res
	}
	return nil
}

func (v *SqlStreamVisitor) GetResult() (plog.LogRecordSlice, error) {
	return v.logRecords, nil
}

func (v *SqlStreamVisitor) Visit(tree antlr.ParseTree) interface{} {

	switch t := tree.(type) {
	case *antlr.ErrorNodeImpl:
		return nil
	case *SqlQueryContext:
		return t.Accept(v)
	}

	return nil
}

func (v *SqlStreamVisitor) VisitSqlQuery(ctx *SqlQueryContext) interface{} {
	return ctx.SelectQuery().Accept(v)
}

func (v *SqlStreamVisitor) VisitSelectSimple(ctx *SelectSimpleContext) interface{} {
	if ctx.WhereStatement() == nil {
		return ctx.ResultColumns().Accept(v)
	}
	switch res := ctx.WhereStatement().Accept(v).(type) {
	case error:
		return res
	case nil:

	}
	return ctx.ResultColumns().Accept(v)

}

func (v *SqlStreamVisitor) VisitSelectTumbling(ctx *SelectTumblingContext) interface{} {
	if ctx.WhereStatement() == nil {
		return ctx.ResultColumns().Accept(v)
	}
	switch res := ctx.WhereStatement().Accept(v).(type) {
	case error:
		return res
	case nil:

	}
	return ctx.ResultColumns().Accept(v)
}

func (v *SqlStreamVisitor) VisitSelectColumns(ctx *SelectColumnsContext) interface{} {
	for _, column := range ctx.AllColumn() {
		v.logRecords.RemoveIf(func(record plog.LogRecord) bool {
			_, ok := record.Attributes().Get(column.GetText())

			if ok {
				// Remove attributes which are not listed in value columns statement
				// We can't use RemoveIf as it takes only values
				removed := make([]string, 0, record.Attributes().Len())
				record.Attributes().Range(func(k string, value pcommon.Value) bool {
					if !KeyExists(k, ctx.AllColumn()) {
						removed = append(removed, k)
					}
					return true
				})

				for _, removedKey := range removed {
					record.Attributes().Remove(removedKey)
				}

				return false
			}
			v.logger.Error("field is missing", zap.String("field", column.GetText()))
			return true
		})

	}
	return nil

}

func (v *SqlStreamVisitor) VisitWhereStmt(ctx *WhereStmtContext) interface{} {
	var err error
	//if WHERE statement missed
	if ctx.Expr() == nil {
		return nil
	}

	v.logRecords.RemoveIf(func(record plog.LogRecord) bool {
		v.currentRecord = record
		switch res := ctx.Expr().Accept(v).(type) {
		case error:
			err = res
			return false

		case bool:
			return !res

		}
		return false
	})

	return err
}

func (v *SqlStreamVisitor) VisitSelectAVG(ctx *SelectAVGContext) interface{} {

	fieldName := ctx.Column().GetText()
	ls := plog.NewLogRecordSlice()
	agrRecord := ls.AppendEmpty()
	var res float64
	var err error

	if ctx.K_AVG() != nil {
		res, err = avg(v.logRecords, fieldName)
	}
	if ctx.K_SUM() != nil {
		res, err = sum(v.logRecords, fieldName)
	}
	if ctx.K_COUNT() != nil {
		res = float64(count(v.logRecords))
	}
	if ctx.K_MAX() != nil {
		res, err = max(v.logRecords, fieldName)
	}
	if ctx.K_MIN() != nil {
		res, err = min(v.logRecords, fieldName)
	}
	if err != nil {
		return err
	}
	fmt.Println(fieldName, res)

	agrRecord.Attributes().Insert(fieldName, pcommon.NewValueDouble(res))
	v.logRecords = ls
	return nil

}

func (v *SqlStreamVisitor) VisitSelectStar(ctx *SelectStarContext) interface{} {
	return nil
}

func (v *SqlStreamVisitor) VisitColumn(ctx *ColumnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *SqlStreamVisitor) VisitSimpleCondition(ctx *SimpleConditionContext) interface{} {

	fieldNameValue, ok := v.currentRecord.Attributes().Get(ctx.IDENTIFIER().GetText())
	if !ok {
		return fmt.Errorf("field %q missed in log record", ctx.IDENTIFIER().GetText())
	}

	switch ctx.LiteralValue().GetStart().GetTokenType() {
	case SqlParserNUMERIC_LITERAL:
		fieldValue, err := strconv.ParseFloat(fieldNameValue.AsString(), 64)
		if err != nil {
			return fmt.Errorf("can't convert record field value %q to numeric; %w", fieldNameValue.AsString(), err)
		}
		comparisonValue, err := strconv.ParseFloat(ctx.LiteralValue().GetText(), 64)
		if err != nil {
			return fmt.Errorf("can't convert comparison value %q to numeric; %w", ctx.LiteralValue().GetText(), err)
		}
		return compareNumeric(ctx, fieldValue, comparisonValue)
	case SqlParserSTRING_LITERAL:
		// we need to remove quotes
		value := strings.TrimSuffix(strings.TrimPrefix(ctx.LiteralValue().GetText(), `'`), `'`)
		return compareString(ctx, fieldNameValue.AsString(), value)
	case SqlParserK_TRUE, SqlParserK_FALSE:
		fieldValue, err := strconv.ParseBool(fieldNameValue.AsString())
		if err != nil {
			return fmt.Errorf("can't convert field value %q to boolean; %w", fieldNameValue.AsString(), err)
		}
		comparisonValue, err := strconv.ParseBool(ctx.LiteralValue().GetText())
		if err != nil {
			return fmt.Errorf("can't convert comparison value %q to boolean; %w", ctx.LiteralValue().GetText(), err)
		}

		return compareBool(ctx, fieldValue, comparisonValue)

	default:
		return fmt.Errorf("missed literal value type %q", ctx.LiteralValue().GetText())
	}

}

func (v *SqlStreamVisitor) VisitSimpleRecursiveCondition(ctx *SimpleRecursiveConditionContext) interface{} {

	left := ctx.Expr(0).Accept(v)
	switch left := ctx.Expr(0).Accept(v).(type) {
	case error:
		return left
	case bool:
	}

	right := ctx.Expr(1).Accept(v)
	switch right.(type) {
	case error:
		return right
	case bool:
	}

	if ctx.K_OR() != nil {
		if left == true || right == true {
			return true
		}
		return false
	}
	if left == true && right == true {
		return true
	}
	return false
}

func (v *SqlStreamVisitor) VisitWindowTumbling(ctx *WindowTumblingContext) interface{} {
	return nil
}
func (v *SqlStreamVisitor) VisitCompoundRecursiveCondition(ctx *CompoundRecursiveConditionContext) interface{} {

	left := ctx.CompoundExpr(0).Accept(v)
	switch left := ctx.CompoundExpr(0).Accept(v).(type) {
	case error:
		return left
	case bool:
	}

	right := ctx.CompoundExpr(1).Accept(v)
	switch right.(type) {
	case error:
		return right
	case bool:
	}

	if ctx.K_OR() != nil {
		if left == true || right == true {
			return true
		}
		return false
	}
	if left == true && right == true {
		return true
	}
	return false
}

func (v *SqlStreamVisitor) VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *SqlStreamVisitor) VisitCompoundExpr(ctx *CompoundExprContext) interface{} {
	return ctx.Expr().Accept(v)
}

func (v *SqlStreamVisitor) VisitLiteralValue(ctx *LiteralValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *SqlStreamVisitor) VisitGroupBy(ctx *GroupByContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *SqlStreamVisitor) VisitAvg(ctx *AvgContext) interface{} {
	return v.VisitChildren(ctx)
}
