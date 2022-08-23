package parser

import (
	"fmt"
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

	logRecords        plog.LogRecordSlice
	currentRecord     plog.LogRecord
	groupByLogRecords map[string]plog.LogRecordSlice

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

	if ctx.WhereStatement() != nil {
		switch res := ctx.WhereStatement().Accept(v).(type) {
		case error:
			return res
		case nil:

		}
	}

	return ctx.AggregationColumns().Accept(v)
}

func (v *SqlStreamVisitor) VisitSelectTumblingGroupBy(ctx *SelectTumblingGroupByContext) interface{} {
	if ctx.WhereStatement() != nil {
		switch res := ctx.WhereStatement().Accept(v).(type) {
		case error:
			return res
		case nil:
		}
	}

	switch res := ctx.GroupBy().Accept(v).(type) {
	case error:
		return res
	case nil:
	}

	return ctx.AggregationColumns().Accept(v)
}

func (v *SqlStreamVisitor) VisitGroupBy(ctx *GroupByContext) interface{} {

	v.groupByLogRecords = make(map[string]plog.LogRecordSlice)
	for i := 0; i < v.logRecords.Len(); i++ {
		record := v.logRecords.At(i)
		groupByField, ok := record.Attributes().Get(ctx.Column().GetText())
		if !ok {
			return fmt.Errorf("group by field %q missed", ctx.Column().GetText())
		}
		_, ok = v.groupByLogRecords[groupByField.AsString()]
		if !ok {
			v.groupByLogRecords[groupByField.AsString()] = plog.NewLogRecordSlice()
		}
		appendedRec := v.groupByLogRecords[groupByField.AsString()].AppendEmpty()
		record.CopyTo(appendedRec)
	}

	return nil
}

func (v *SqlStreamVisitor) VisitSelectColumns(ctx *SelectColumnsContext) interface{} {

	for _, column := range ctx.AllColumn() {
		v.logRecords.RemoveIf(func(record plog.LogRecord) bool {
			// apply scalar function
			v.currentRecord = record
			identifier := getColumnIdentifier(column)
			_, ok := record.Attributes().Get(identifier)

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
				column.Accept(v)

				return false
			}

			v.logger.Error("field is missing", zap.String("field", column.GetText()))
			return true
		})

	}
	return nil

}

func (v *SqlStreamVisitor) VisitIdentifierCol(ctx *IdentifierColContext) interface{} {
	return nil
}

func (v *SqlStreamVisitor) VisitFunctionCol(ctx *FunctionColContext) interface{} {
	if ctx.K_LOWER() != nil {
		return lower(v.currentRecord, ctx.IDENTIFIER().GetText())
	}
	if ctx.K_UPPER() != nil {
		return upper(v.currentRecord, ctx.IDENTIFIER().GetText())
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
	dest := plog.NewLogRecordSlice()
	if v.groupByLogRecords != nil {
		for key, recSlice := range v.groupByLogRecords {
			res, err := v.getSimpleAggregatedValue(ctx, ctx.Column(1).GetText(), recSlice)
			if err != nil {
				return err
			}

			res.Attributes().Insert(ctx.Column(0).GetText(), pcommon.NewValueString(key))
			res.CopyTo(dest.AppendEmpty())
		}
		v.logRecords = dest
		return nil
	}

	res, err := v.getSimpleAggregatedValue(ctx, ctx.Column(0).GetText(), v.logRecords)
	if err != nil {
		return err
	}
	res.CopyTo(dest.AppendEmpty())
	v.logRecords = dest

	return nil
}

func (v *SqlStreamVisitor) getSimpleAggregatedValue(ctx *SelectAVGContext, fieldName string, ls plog.LogRecordSlice) (plog.LogRecord, error) {

	resLs := plog.NewLogRecordSlice()
	aggrRecord := resLs.AppendEmpty()
	var res float64
	var err error

	if ctx.K_AVG() != nil {
		res, err = avg(ls, fieldName)
	}
	if ctx.K_SUM() != nil {
		res, err = sum(ls, fieldName)
	}
	if ctx.K_COUNT() != nil {
		res = float64(count(ls))
	}
	if ctx.K_MAX() != nil {
		res, err = max(ls, fieldName)
	}
	if ctx.K_MIN() != nil {
		res, err = min(ls, fieldName)
	}
	if err != nil {
		return plog.LogRecord{}, err
	}

	aggrRecord.Attributes().Insert(fieldName, pcommon.NewValueDouble(res))
	return aggrRecord, nil
}

func (v *SqlStreamVisitor) VisitSelectStar(ctx *SelectStarContext) interface{} {
	return nil
}

func (v *SqlStreamVisitor) VisitColumn(ctx *ColumnContext) interface{} {
	return nil
}

func (v *SqlStreamVisitor) VisitSimpleCondition(ctx *SimpleConditionContext) interface{} {
	return ctx.SimpleExpr().Accept(v)
}

func (v *SqlStreamVisitor) VisitSimpleExpression(ctx *SimpleExpressionContext) interface{} {
	fieldValue, ok := v.currentRecord.Attributes().Get(ctx.IDENTIFIER().GetText())
	if !ok {
		return fmt.Errorf("field %q missed in log record", ctx.IDENTIFIER().GetText())
	}

	res, err := compareExpression(ctx.ComparisonOperator(), ctx.LiteralValue(), fieldValue)
	if err != nil {
		return err
	}
	return res

}

func (v *SqlStreamVisitor) VisitNestedExpression(ctx *NestedExpressionContext) interface{} {
	fieldName := ctx.IDENTIFIER(0).GetText()
	fieldValue, ok := v.currentRecord.Attributes().Get(fieldName)
	if !ok {
		return fmt.Errorf("field %q missed in log record", fieldValue.AsString())
	}
	if fieldValue.Type() != pcommon.ValueTypeMap {
		return fmt.Errorf("field %q is not nested map", fieldName)
	}

	nestedFieldName := ctx.IDENTIFIER(1).GetText()
	nestedFieldValue, ok := fieldValue.MapVal().Get(nestedFieldName)
	if !ok {
		v.logger.Error("nested key missed", zap.String("name: ", nestedFieldName))
		return false
	}

	res, err := compareExpression(ctx.ComparisonOperator(), ctx.LiteralValue(), nestedFieldValue)
	if err != nil {
		return err
	}
	return res

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
	var left, right interface{}

	/*if ctx.CompoundExpr(0) != nil {
		fmt.Println("left: ", ctx.CompoundExpr(0).GetText())
	}
	if ctx.CompoundExpr(1) != nil {
		fmt.Println("right: ", ctx.CompoundExpr(1).GetText())
	}
	if ctx.SimpleExpr(0) != nil {
		fmt.Println("simple: ", ctx.SimpleExpr(0).GetText())
	}

	*/

	if ctx.CompoundExpr(0) != nil {
		left = ctx.CompoundExpr(0).Accept(v)
		if ctx.CompoundExpr(1) != nil {
			right = ctx.CompoundExpr(1).Accept(v)
		} else {
			right = ctx.SimpleExpr(0).Accept(v)
		}
	} else {
		left = ctx.SimpleExpr(0).Accept(v)
		if ctx.SimpleExpr(1) != nil {
			right = ctx.SimpleExpr(1).Accept(v)
		} else {
			right = ctx.CompoundExpr(1).Accept(v)
		}

	}

	switch left.(type) {
	case error:
		return left
	case bool:
	}

	switch right.(type) {
	case error:
		return left
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

func (v *SqlStreamVisitor) VisitSimpleCompoundCondition(ctx *SimpleCompoundConditionContext) interface{} {
	return ctx.CompoundExpr().Accept(v)
}

func (v *SqlStreamVisitor) VisitCompoundExpression(ctx *CompoundExpressionContext) interface{} {
	return ctx.Expr().Accept(v)
}

func (v *SqlStreamVisitor) VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *SqlStreamVisitor) VisitLiteralValue(ctx *LiteralValueContext) interface{} {
	return v.VisitChildren(ctx)
}
