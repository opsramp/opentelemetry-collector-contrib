package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

// check that type implements SqlVisitor
var _ SqlVisitor = (*SQLStreamVisitor)(nil)

type SQLStreamVisitor struct {
	antlr.ParseTreeVisitor
	logger *zap.Logger

	query string

	tumblingPeriod time.Duration

	// this is what we've got
	originalRecords       plog.LogRecordSlice
	currentOriginalRecord plog.LogRecord

	// this is logs to emit
	resultRecords         plog.LogRecordSlice
	currentResultRecord   plog.LogRecord
	groupByTempLogRecords map[string]plog.LogRecordSlice

	in        <-chan plog.LogRecordSlice
	out       chan<- plog.LogRecordSlice
	outErr    chan<- error
	shutdownC chan struct{}
}

func NewSQLStreamVisitor(query string, in <-chan plog.LogRecordSlice, out chan<- plog.LogRecordSlice, outErr chan<- error, logger *zap.Logger) *SQLStreamVisitor {

	visitor := &SQLStreamVisitor{
		ParseTreeVisitor: &SQLStreamVisitor{},
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

func (v *SQLStreamVisitor) Stop() {
	close(v.shutdownC)
}

func (v *SQLStreamVisitor) startSimpleWorkerLoop() {

	for {
		select {
		case v.originalRecords = <-v.in:
			if err := v.runQuery(); err != nil {
				v.outErr <- err
				break
			}
			v.out <- v.resultRecords
			v.originalRecords = plog.NewLogRecordSlice()

		case <-v.shutdownC:
			if v.resultRecords.Len() > 0 {
				v.out <- v.resultRecords
			}
			v.logger.Debug("sql engine shutting down")
			return

		}

	}
}

func (v *SQLStreamVisitor) startWindowTumblingLoop() {
	v.originalRecords = plog.NewLogRecordSlice()
	ticker := time.NewTicker(v.tumblingPeriod)

	go func() {
		for {
			select {
			case newRec := <-v.in:
				newRec.CopyTo(v.originalRecords)
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
				if v.originalRecords.Len() > 0 {
					if err := v.runQuery(); err != nil {
						v.outErr <- err
						break
					}
					v.out <- v.resultRecords
					v.originalRecords = plog.NewLogRecordSlice()
				}
			case <-v.shutdownC:
				v.logger.Debug("sql engine shutting down")
				return

			}
		}
	}()
}

func (v *SQLStreamVisitor) runQuery() error {
	v.resultRecords = plog.NewLogRecordSlice()
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

func (v *SQLStreamVisitor) GetResult() (plog.LogRecordSlice, error) {
	return v.originalRecords, nil
}

func (v *SQLStreamVisitor) Visit(tree antlr.ParseTree) interface{} {

	switch t := tree.(type) {
	case *antlr.ErrorNodeImpl:
		return nil
	case *SqlQueryContext:
		return t.Accept(v)
	}

	return nil
}

func (v *SQLStreamVisitor) VisitSqlQuery(ctx *SqlQueryContext) interface{} {
	return ctx.SelectQuery().Accept(v)
}

func (v *SQLStreamVisitor) VisitSelectSimple(ctx *SelectSimpleContext) interface{} {
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

func (v *SQLStreamVisitor) VisitSelectTumbling(ctx *SelectTumblingContext) interface{} {
	if ctx.WhereStatement() != nil {
		switch res := ctx.WhereStatement().Accept(v).(type) {
		case error:
			return res
		case nil:

		}
	}

	return ctx.AggregationColumns().Accept(v)
}

func (v *SQLStreamVisitor) VisitSelectTumblingGroupBy(ctx *SelectTumblingGroupByContext) interface{} {
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

	return ctx.GroupByAggregationColumns().Accept(v)
}

func (v *SQLStreamVisitor) VisitGroupBy(ctx *GroupByContext) interface{} {

	v.groupByTempLogRecords = make(map[string]plog.LogRecordSlice)

	for i := 0; i < v.originalRecords.Len(); i++ {
		record := v.originalRecords.At(i)
		_, value, err := getAttributeValue(ctx.Column(), record.Attributes())
		if err != nil {
			return err
		}

		_, ok := v.groupByTempLogRecords[value.AsString()]
		if !ok {
			v.groupByTempLogRecords[value.AsString()] = plog.NewLogRecordSlice()
		}
		appendedRec := v.groupByTempLogRecords[value.AsString()].AppendEmpty()
		record.CopyTo(appendedRec)
	}
	return nil
}

func (v *SQLStreamVisitor) VisitSelectGroupByAggregations(ctx *SelectGroupByAggregationsContext) interface{} {
	resLS := plog.NewLogRecordSlice()

	for k, ls := range v.groupByTempLogRecords {
		newRec := resLS.AppendEmpty()
		if err := calculateGroupByAggregations(ctx, k, newRec, ls); err != nil {
			return err
		}

	}
	v.resultRecords = resLS
	return nil
}

func (v *SQLStreamVisitor) VisitSelectAggregations(ctx *SelectAggregationsContext) interface{} {

	resLS := plog.NewLogRecordSlice()
	newRec := resLS.AppendEmpty()
	if err := calculateAggregations(ctx, newRec, v.originalRecords); err != nil {
		return err
	}
	v.resultRecords = resLS
	return nil
}

func (v *SQLStreamVisitor) VisitAggregationColumn(ctx *AggregationColumnContext) interface{} {
	// first we check that fields listed in select (including nested) exist in log record
	fieldName := ctx.IDENTIFIER(0).GetText()
	value, ok := v.currentOriginalRecord.Attributes().Get(fieldName)
	if !ok {
		return fmt.Errorf("field %q missed", fieldName)
	}

	// this is simple field
	if len(ctx.AllIDENTIFIER()) == 1 {
		if ctx.Alias() != nil {
			fieldName = ctx.Alias().GetStop().GetText()
		}
		v.currentResultRecord.Attributes().Insert(fieldName, value)

		return nil
	}

	// here we process nested field
	nestedFieldName := ctx.IDENTIFIER(1).GetText()
	nestedVal, err := nestedFieldExistsInAttr(fieldName, nestedFieldName, v.currentOriginalRecord.Attributes())
	if err != nil {
		return err
	}
	if ctx.Alias() != nil {
		fieldName = ctx.Alias().GetStop().GetText()
		v.currentResultRecord.Attributes().Insert(fieldName, nestedVal)
		return nil
	}

	var newMapVal pcommon.Value
	newMapVal, ok = v.currentResultRecord.Attributes().Get(fieldName)
	if !ok {
		newMapVal = pcommon.NewValueMap()
	}

	newMapVal.MapVal().Insert(nestedFieldName, nestedVal)
	v.currentResultRecord.Attributes().Insert(fieldName, newMapVal)

	return nil
}

// VisitSelectColumns is called in case of missed where statement and removes missed attributes
func (v *SQLStreamVisitor) VisitSelectColumns(ctx *SelectColumnsContext) interface{} {
	var err error
	for i := 0; i < v.originalRecords.Len(); i++ {
		v.currentOriginalRecord = v.originalRecords.At(i)
		v.currentResultRecord = plog.NewLogRecord()

		for _, col := range ctx.AllColumn() {
			switch res := col.Accept(v).(type) {
			case error:
				err = res
			}
		}

		v.currentResultRecord.MoveTo(v.resultRecords.AppendEmpty())

	}

	return err

}

func (v *SQLStreamVisitor) VisitWhereStmt(ctx *WhereStmtContext) interface{} {
	var err error
	// if WHERE statement missed
	if ctx.Expr() == nil {
		return nil
	}
	selectColCtx := getSelectColumnsFromWhereCtx(ctx)
	selectAggrColCtx := getSelectAggregationsFromWhereCtx(ctx)
	selectStarCtx := getSelectStarCtx(ctx)

	v.originalRecords.RemoveIf(func(record plog.LogRecord) bool {
		v.currentOriginalRecord = record
		v.currentResultRecord = plog.NewLogRecord()

		switch res := ctx.Expr().Accept(v).(type) {
		case error:
			err = res
			return false

		case bool:
			if !res {
				return true
			}

			// select * case
			if selectStarCtx != nil {
				v.currentOriginalRecord.MoveTo(v.resultRecords.AppendEmpty())
				return false
			}

			// align result records according listed in select fields
			if selectColCtx != nil {
				err = v.applySelectColumns(selectColCtx)
			}

			// the same for  aggregation case ##selectTumbling #selectTumblingGroupBy
			if selectAggrColCtx != nil {
				err = v.applySelectAggregationColumns(selectAggrColCtx)
			}

			v.currentResultRecord.MoveTo(v.resultRecords.AppendEmpty())
			return false
		}

		return false
	})

	return err
}

func (v *SQLStreamVisitor) applySelectColumns(ctx *SelectColumnsContext) error {
	for _, col := range ctx.AllColumn() {
		switch res := col.Accept(v).(type) {
		case error:
			return res
		}
	}

	return nil
}

func (v *SQLStreamVisitor) applySelectAggregationColumns(ctx *SelectAggregationsContext) error {
	for _, col := range ctx.AllAggregationColumn() {
		switch res := col.Accept(v).(type) {
		case error:
			return res
		}
	}

	return nil
}

func (v *SQLStreamVisitor) VisitSelectStar(ctx *SelectStarContext) interface{} {
	v.originalRecords.CopyTo(v.resultRecords)
	return nil
}

func (v *SQLStreamVisitor) VisitIdentifierColumn(ctx *IdentifierColumnContext) interface{} {
	// first we check that fields listed in select (including nested) exist in log record
	fieldName := ctx.IDENTIFIER(0).GetText()
	value, ok := v.currentOriginalRecord.Attributes().Get(fieldName)
	if !ok {
		return fmt.Errorf("field %q missed", fieldName)
	}

	// this is simple field
	if len(ctx.AllIDENTIFIER()) == 1 {
		if ctx.Alias() != nil {
			fieldName = ctx.Alias().GetStop().GetText()
		}
		v.currentResultRecord.Attributes().Insert(fieldName, value)

		return nil
	}

	// here we process nested field
	nestedFieldName := ctx.IDENTIFIER(1).GetText()
	nestedVal, err := nestedFieldExistsInAttr(fieldName, nestedFieldName, v.currentOriginalRecord.Attributes())
	if err != nil {
		return err
	}
	if ctx.Alias() != nil {
		fieldName = ctx.Alias().GetStop().GetText()
		v.currentResultRecord.Attributes().Insert(fieldName, nestedVal)
		return nil
	}

	var newMapVal pcommon.Value
	newMapVal, ok = v.currentResultRecord.Attributes().Get(fieldName)
	if !ok {
		newMapVal = pcommon.NewValueMap()
	}

	newMapVal.MapVal().Insert(nestedFieldName, nestedVal)
	v.currentResultRecord.Attributes().Insert(fieldName, newMapVal)

	return nil
}

func (v *SQLStreamVisitor) VisitFunctionColumn(ctx *FunctionColumnContext) interface{} {

	var simpleFuncCtx *SimpleFunctionContext
	res := ctx.Function().Accept(v)

	switch res.(type) {
	case error:
		return res
	case pcommon.Value:
	default:
		return fmt.Errorf("undefined result type")
	}

	simpleFuncCtx, ok := ctx.GetChild(0).(*SimpleFunctionContext)
	if !ok {
		simpleFuncCtx = getSimpleFunctionContext(ctx)
	}

	fieldName := simpleFuncCtx.IDENTIFIER(0).GetText()
	if ctx.Alias() != nil {
		fieldName = ctx.Alias().GetStop().GetText()
	}

	// if simple field
	if len(simpleFuncCtx.AllIDENTIFIER()) == 1 {
		v.currentResultRecord.Attributes().Insert(fieldName, res.(pcommon.Value))
		return nil
	}

	//this is nested field
	nestedFieldName := simpleFuncCtx.IDENTIFIER(1).GetText()
	newMapVal := pcommon.NewValueMap()
	newMapVal.MapVal().Insert(nestedFieldName, res.(pcommon.Value))
	v.currentResultRecord.Attributes().Insert(nestedFieldName, newMapVal)

	return res
}

func (v *SQLStreamVisitor) VisitSimpleFunction(ctx *SimpleFunctionContext) interface{} {

	fieldName := ctx.IDENTIFIER(0).GetText()
	functionMame := ctx.FunctionName().GetText()
	args := ctx.AllLiteralValue()

	value, ok := v.currentOriginalRecord.Attributes().Get(fieldName)
	if !ok {
		return fmt.Errorf("field %q missed", fieldName)
	}

	// this is simple field
	if len(ctx.AllIDENTIFIER()) == 1 {

		newValue, err := v.applyFunction(value, functionMame, args...)
		if err != nil {
			return err
		}

		// v.currentResultRecord.Attributes().Insert(fieldName, newValue)
		return newValue
	}

	// here we process nested field
	nestedFieldName := ctx.IDENTIFIER(1).GetText()
	nestedVal, err := nestedFieldExistsInAttr(fieldName, nestedFieldName, v.currentOriginalRecord.Attributes())
	if err != nil {
		return err
	}
	newValue, err := v.applyFunction(value, functionMame, args...)

	newMapVal := pcommon.NewValueMap()
	newMapVal.MapVal().Insert(nestedFieldName, nestedVal)
	//v.currentResultRecord.Attributes().Insert(fieldName, newValue)
	return newValue
}

func (v *SQLStreamVisitor) VisitRecursiveFunction(ctx *RecursiveFunctionContext) interface{} {

	res := ctx.Function().Accept(v)

	switch res.(type) {
	case error:
		return res
	case pcommon.Value:
	default:
		return fmt.Errorf("undefined return type")

	}
	tmpValue, err := v.applyFunction(res.(pcommon.Value), ctx.FunctionName().GetText(), ctx.AllLiteralValue()...)
	if err != nil {
		return nil
	}

	return tmpValue
}

func (v *SQLStreamVisitor) applyFunction(value pcommon.Value, functionName string, args ...ILiteralValueContext) (pcommon.Value, error) {
	switch functionName {
	case upper:
		return pcommon.NewValueString(strings.ToUpper(value.AsString())), nil
	case lower:
		return pcommon.NewValueString(strings.ToLower(value.AsString())), nil
	case substr:
		if len(args) != 2 {
			return pcommon.NewValueEmpty(), fmt.Errorf("substr accept two arguments")
		}
		res, err := substring(value, args[0].GetText(), args[1].GetText())
		return res, err

	default:
		return pcommon.NewValueEmpty(), fmt.Errorf("function %q isn't available", functionName)

	}
}

func (v *SQLStreamVisitor) VisitSimpleCondition(ctx *SimpleConditionContext) interface{} {
	return ctx.SimpleExpr().Accept(v)
}

func (v *SQLStreamVisitor) VisitSimpleExpression(ctx *SimpleExpressionContext) interface{} {
	fieldValue, ok := v.currentOriginalRecord.Attributes().Get(ctx.IDENTIFIER().GetText())
	if !ok {
		return fmt.Errorf("field %q missed in log record", ctx.IDENTIFIER().GetText())
	}

	res, err := compareExpression(ctx.ComparisonOperator(), ctx.LiteralValue(), fieldValue)
	if err != nil {
		return err
	}
	return res

}

func (v *SQLStreamVisitor) VisitNestedExpression(ctx *NestedExpressionContext) interface{} {
	fieldName := ctx.IDENTIFIER(0).GetText()
	fieldValue, ok := v.currentOriginalRecord.Attributes().Get(fieldName)
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

func (v *SQLStreamVisitor) VisitSimpleRecursiveCondition(ctx *SimpleRecursiveConditionContext) interface{} {

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

func (v *SQLStreamVisitor) VisitWindowTumbling(ctx *WindowTumblingContext) interface{} {
	return nil
}
func (v *SQLStreamVisitor) VisitCompoundRecursiveCondition(ctx *CompoundRecursiveConditionContext) interface{} {
	var left, right interface{}

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

func (v *SQLStreamVisitor) VisitSimpleCompoundCondition(ctx *SimpleCompoundConditionContext) interface{} {
	return ctx.CompoundExpr().Accept(v)
}

func (v *SQLStreamVisitor) VisitCompoundExpression(ctx *CompoundExpressionContext) interface{} {
	return ctx.Expr().Accept(v)
}

func (v *SQLStreamVisitor) VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *SQLStreamVisitor) VisitLiteralValue(ctx *LiteralValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *SQLStreamVisitor) VisitAlias(ctx *AliasContext) interface{} {
	return nil
}

func (v *SQLStreamVisitor) VisitFunction(ctx *FunctionContext) interface{} {
	return nil
}

func (v *SQLStreamVisitor) VisitFunctionName(ctx *FunctionNameContext) interface{} {
	return nil
}
