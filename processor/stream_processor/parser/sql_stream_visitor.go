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

	// this is what we've got
	origLogRecords    plog.LogRecordSlice
	currentOrigRecord plog.LogRecord

	// this is logs to emin
	resultLogRecords    plog.LogRecordSlice
	currentResultRecord plog.LogRecord

	groupByLogRecords map[string]plog.LogRecordSlice

	in        <-chan plog.LogRecordSlice
	out       chan<- plog.LogRecordSlice
	outErr    chan<- error
	shutdownC chan struct{}
}

func NewSqlStreamVisitor(query string, in <-chan plog.LogRecordSlice, out chan<- plog.LogRecordSlice, outErr chan<- error, logger *zap.Logger) *SqlStreamVisitor {

	visitor := &SqlStreamVisitor{
		ParseTreeVisitor:  &SqlStreamVisitor{},
		query:             query,
		logger:            logger,
		groupByLogRecords: map[string]plog.LogRecordSlice{},
		in:                in,
		out:               out,
		outErr:            outErr,
		shutdownC:         make(chan struct{}, 1),
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
		case v.origLogRecords = <-v.in:
			v.resultLogRecords = plog.NewLogRecordSlice()
			if err := v.runQuery(); err != nil {
				v.outErr <- err
				break
			}
			v.out <- v.resultLogRecords
			v.origLogRecords = plog.NewLogRecordSlice()
			v.resultLogRecords = plog.NewLogRecordSlice()

		case <-v.shutdownC:
			if v.resultLogRecords.Len() > 0 {
				v.out <- v.resultLogRecords
			}
			v.logger.Debug("sql engine shutting down")
			return

		}

	}
}

func (v *SqlStreamVisitor) startWindowTumblingLoop() {
	v.origLogRecords = plog.NewLogRecordSlice()
	ticker := time.NewTicker(v.tumblingPeriod)

	go func() {
		for {
			select {
			case ls := <-v.in:
				ls.MoveAndAppendTo(v.origLogRecords)
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
				if v.origLogRecords.Len() > 0 {
					if err := v.runQuery(); err != nil {
						v.outErr <- err
						break
					}
					v.out <- v.origLogRecords
					v.origLogRecords = plog.NewLogRecordSlice()
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
	return v.origLogRecords, nil
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

	return nil
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

	//v.groupByLogRecords = make(map[string]plog.LogRecordSlice)
	for i := 0; i < v.origLogRecords.Len(); i++ {
		record := v.origLogRecords.At(i)
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

// VisitSelectColumns is called in case of missed where statement and removes missed attributes
func (v *SqlStreamVisitor) VisitSelectColumns(ctx *SelectColumnsContext) interface{} {
	var err error
	for i := 0; i < v.origLogRecords.Len(); i++ {
		v.currentOrigRecord = v.origLogRecords.At(i)
		v.currentResultRecord = plog.NewLogRecord()

		for _, col := range ctx.AllColumn() {
			switch res := col.Accept(v).(type) {
			case error:
				err = res
			}
		}

		v.currentResultRecord.MoveTo(v.resultLogRecords.AppendEmpty())

	}

	return err

}

func (v *SqlStreamVisitor) VisitWhereStmt(ctx *WhereStmtContext) interface{} {
	var err error
	//if WHERE statement missed
	if ctx.Expr() == nil {
		return nil
	}
	resColumnCtx := getSelectColumnsFromWhereCtx(ctx)

	v.origLogRecords.RemoveIf(func(record plog.LogRecord) bool {
		v.currentOrigRecord = record
		v.currentResultRecord = plog.NewLogRecord()

		switch res := ctx.Expr().Accept(v).(type) {
		case error:
			err = res
			return false

		case bool:
			if !res {
				return true
			}

			if errRes := v.applySelectColumns(resColumnCtx); errRes != nil {
				err = errRes
			}

			return false
		}
		return false
	})

	return err
}

func (v *SqlStreamVisitor) applySelectColumns(ctx *SelectColumnsContext) error {
	//if * in select list
	if ctx == nil {
		v.currentOrigRecord.MoveTo(v.resultLogRecords.AppendEmpty())
		return nil
	}

	for _, col := range ctx.AllColumn() {
		switch res := col.Accept(v).(type) {
		case error:
			return res
		}
	}
	v.currentResultRecord.MoveTo(v.resultLogRecords.AppendEmpty())
	return nil
}

//adjustResultColumns this is called if where statement exists
func (v *SqlStreamVisitor) adjustResultColumns(ctx *SelectColumnsContext) error {

	// Remove attributes which are not listed in value columns statement
	// We can't use RemoveIf as it takes only values
	removed := make([]string, 0, v.currentOrigRecord.Attributes().Len())

	v.currentOrigRecord.Attributes().RemoveIf(func(k string, value pcommon.Value) bool {
		return fieldExists(k, value, ctx.AllColumn())
	})

	for _, removedKey := range removed {
		v.currentOrigRecord.Attributes().Remove(removedKey)
	}

	return nil
}

func (v *SqlStreamVisitor) VisitSelectAggregation(ctx *SelectAggregationContext) interface{} {
	return ctx.AggregationColumn().Accept(v)
}

func (v *SqlStreamVisitor) VisitColumnAggregation(ctx *ColumnAggregationContext) interface{} {
	/*dest := plog.NewLogRecordSlice()

	if v.groupByLogRecords != nil {
		for key, recSlice := range v.groupByLogRecords {
			res, err := v.getSimpleAggregatedValue(ctx, ctx.().GetText(), recSlice)
			if err != nil {
				return err
			}

			res.Attributes().Insert(ctx.Column().GetText(), pcommon.NewValueString(key))
			res.CopyTo(dest.AppendEmpty())
		}
		v.origLogRecords = dest
		return nil
	}

	res, err := v.getSimpleAggregatedValue(ctx, ctx.Column().GetText(), v.origLogRecords)
	if err != nil {
		return err
	}
	res.CopyTo(dest.AppendEmpty())
	v.origLogRecords = dest

	*/

	return nil
}

func (v *SqlStreamVisitor) VisitColumnCountAggregation(ctx *ColumnCountAggregationContext) interface{} {
	//TODO implement me
	panic("implement me")
}

func (v *SqlStreamVisitor) getSimpleAggregatedValue(ctx *ColumnAggregationContext, fieldName string, ls plog.LogRecordSlice) (plog.LogRecord, error) {

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

func (v *SqlStreamVisitor) getCountAggregatedValue(ctx *ColumnAggregationContext, fieldName string, ls plog.LogRecordSlice) (plog.LogRecord, error) {

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
	v.origLogRecords.CopyTo(v.resultLogRecords)
	return nil
}

func (v *SqlStreamVisitor) VisitIdentifierColumn(ctx *IdentifierColumnContext) interface{} {
	//first we check that fields listed in select (including nested) exist in log record
	fieldName := ctx.IDENTIFIER(0).GetText()
	value, ok := v.currentOrigRecord.Attributes().Get(fieldName)
	if !ok {
		return fmt.Errorf("field %q missed", fieldName)
	}

	//this is simple field
	if len(ctx.AllIDENTIFIER()) == 1 {
		if ctx.Alias() != nil {
			fieldName = ctx.Alias().GetStop().GetText()
		}
		v.currentResultRecord.Attributes().Insert(fieldName, value)

		return nil
	}

	// here we process nested field
	nestedFieldName := ctx.IDENTIFIER(1).GetText()
	nestedVal, err := nestedFieldExistsInAttr(fieldName, nestedFieldName, v.currentOrigRecord.Attributes())
	if err != nil {
		return err
	}
	if ctx.Alias() != nil {
		fieldName = ctx.Alias().GetStop().GetText()
		v.currentResultRecord.Attributes().Insert(fieldName, nestedVal)
		return nil
	}
	newMapVal := pcommon.NewValueMap()
	newMapVal.MapVal().Insert(nestedFieldName, nestedVal)
	v.currentResultRecord.Attributes().Insert(fieldName, newMapVal)

	return nil
}

func (v *SqlStreamVisitor) VisitFunctionColumn(ctx *FunctionColumnContext) interface{} {
	//first we check that fields listed in select (including nested) exist in log record
	fieldName := ctx.IDENTIFIER(0).GetText()
	functionMame := ctx.Function().GetText()
	args := ctx.AllLiteralValue()

	value, ok := v.currentOrigRecord.Attributes().Get(fieldName)
	if !ok {
		return fmt.Errorf("field %q missed", fieldName)
	}

	//this is simple field
	if len(ctx.AllIDENTIFIER()) == 1 {

		if ctx.Alias() != nil {
			fieldName = ctx.Alias().GetStop().GetText()
		}
		newValue, err := v.applyFunction(value, functionMame, args...)
		if err != nil {
			return err
		}

		v.currentResultRecord.Attributes().Insert(fieldName, newValue)

		return nil
	}

	// here we process nested field
	nestedFieldName := ctx.IDENTIFIER(1).GetText()
	nestedVal, err := nestedFieldExistsInAttr(fieldName, nestedFieldName, v.currentOrigRecord.Attributes())
	if err != nil {
		return err
	}
	if ctx.Alias() != nil {
		fieldName = ctx.Alias().GetStop().GetText()
		v.currentResultRecord.Attributes().Insert(fieldName, nestedVal)
		return nil
	}
	newMapVal := pcommon.NewValueMap()
	newMapVal.MapVal().Insert(nestedFieldName, nestedVal)
	v.currentResultRecord.Attributes().Insert(fieldName, newMapVal)

	return nil
}
func (v *SqlStreamVisitor) applyFunction(value pcommon.Value, functionName string, args ...ILiteralValueContext) (pcommon.Value, error) {
	switch functionName {
	case "upper":
	case "lower":
	case "substr":
	default:
		return pcommon.NewValueEmpty(), fmt.Errorf("function %q isn't available", functionName)

	}
	return pcommon.NewValueEmpty(), nil
}

func (v *SqlStreamVisitor) VisitSimpleCondition(ctx *SimpleConditionContext) interface{} {
	return ctx.SimpleExpr().Accept(v)
}

func (v *SqlStreamVisitor) VisitSimpleExpression(ctx *SimpleExpressionContext) interface{} {
	fieldValue, ok := v.currentOrigRecord.Attributes().Get(ctx.IDENTIFIER().GetText())
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
	fieldValue, ok := v.currentOrigRecord.Attributes().Get(fieldName)
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

func (v *SqlStreamVisitor) VisitAlias(ctx *AliasContext) interface{} {
	return nil
}

func (v *SqlStreamVisitor) VisitFunction(ctx *FunctionContext) interface{} {

	return nil
}
