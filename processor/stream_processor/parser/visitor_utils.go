package parser

import (
	"fmt"
	"strconv"
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
)

const (
	lower  = "lower"
	upper  = "upper"
	substr = "substr"
)

func compareExpression(operator IComparisonOperatorContext, literalValue ILiteralValueContext, fieldValue pcommon.Value) (bool, error) {
	switch literalValue.GetStart().GetTokenType() {
	case SqlParserNUMERIC_LITERAL:
		parsedValue, err := strconv.ParseFloat(fieldValue.AsString(), 64)
		if err != nil {
			return false, fmt.Errorf("can't convert record field value %q to numeric; %w", fieldValue.AsString(), err)
		}
		comparisonValue, err := strconv.ParseFloat(literalValue.GetText(), 64)
		if err != nil {
			return false, fmt.Errorf("can't convert comparison value %q to numeric; %w", literalValue.GetText(), err)
		}
		return compareNumeric(operator.GetStart().GetTokenType(), parsedValue, comparisonValue), nil
	case SqlParserSTRING_LITERAL:
		// we need to remove quotes
		value := strings.TrimSuffix(strings.TrimPrefix(literalValue.GetText(), `'`), `'`)
		return compareString(operator.GetStart().GetTokenType(), fieldValue.AsString(), value), nil
	case SqlParserBOOLEAN_LITERAL:
		parsedValue, err := strconv.ParseBool(fieldValue.AsString())
		if err != nil {
			return false, fmt.Errorf("can't convert field value %q to boolean; %w", fieldValue.AsString(), err)
		}
		comparisonValue, err := strconv.ParseBool(literalValue.GetText())
		if err != nil {
			return false, fmt.Errorf("can't convert comparison value %q to boolean; %w", literalValue.GetText(), err)
		}

		return compareBool(operator.GetStart().GetTokenType(), parsedValue, comparisonValue), nil

	default:
		return false, fmt.Errorf("missed literal value type %q", literalValue.GetText())
	}
}

func compareString(comparisonToken int, fieldVal, comparisonVal string) bool {
	switch comparisonToken {
	case SqlParserK_EQUAL:
		return fieldVal == comparisonVal
	case SqlParserK_NOT_EQUAL:
		return fieldVal != comparisonVal
	case SqlParserK_GREATER:
		return fieldVal > comparisonVal
	case SqlParserK_GREATER_EQUAL:
		return fieldVal >= comparisonVal
	case SqlParserK_LESS:
		return fieldVal < comparisonVal
	case SqlParserK_LESS_EQUAL:
		return fieldVal <= comparisonVal
	case SqlParserK_LIKE:
		return strings.Contains(fmt.Sprint(fieldVal), fmt.Sprint(comparisonVal))
	case SqlParserK_NOT_LIKE:
		return !strings.Contains(fmt.Sprint(fieldVal), fmt.Sprint(comparisonVal))
	case SqlParserK_IS_NULL:
		return len(fmt.Sprint(fieldVal)) == 0
	case SqlParserK_IS_NOT_NULL:
		return len(fmt.Sprint(fieldVal)) != 0
	default:
		return false
	}
}

func compareNumeric(comparisonToken int, fieldVal, comparisonVal float64) bool {
	switch comparisonToken {
	case SqlParserK_EQUAL:
		return fieldVal == comparisonVal
	case SqlParserK_NOT_EQUAL:
		return fieldVal != comparisonVal
	case SqlParserK_GREATER:
		return fieldVal > comparisonVal
	case SqlParserK_GREATER_EQUAL:
		return fieldVal >= comparisonVal
	case SqlParserK_LESS:
		return fieldVal < comparisonVal
	case SqlParserK_LESS_EQUAL:
		return fieldVal <= comparisonVal
	case SqlParserK_IS_NULL:
		return fieldVal == 0
	case SqlParserK_IS_NOT_NULL:
		return fieldVal != 0
	default:
		return false
	}
}

func compareBool(comparisonToken int, fieldVal, comparisonVal bool) bool {
	switch comparisonToken {
	case SqlParserK_EQUAL:
		return fieldVal == comparisonVal
	case SqlParserK_NOT_EQUAL:
		return fieldVal != comparisonVal
	}
	return false
}

func getFieldAggregatedValue(ctx *AggregationColumnContext, ls plog.LogRecordSlice) (float64, error) {

	var res float64
	var err error

	if ctx.K_AVG() != nil {
		res, err = avg(ls, ctx)
	}
	if ctx.K_SUM() != nil {
		res, err = sum(ls, ctx)
	}
	if ctx.K_COUNT() != nil {
		res = float64(count(ls))
	}
	if ctx.K_MAX() != nil {
		res, err = max(ls, ctx)
	}
	if ctx.K_MIN() != nil {
		res, err = min(ls, ctx)
	}

	return res, err
}

func insertIdentifierColumnToRecord(ctx *IdentifierColumnContext, attr pcommon.Map, value string) error {
	if ctx.Alias() != nil {
		attr.Insert(ctx.Alias().GetStop().GetText(), pcommon.NewValueString(value))
		return nil
	}

	fieldName := ctx.IDENTIFIER(0).GetText()
	if len(ctx.AllIDENTIFIER()) > 1 {
		nested, ok := attr.Get(fieldName)

		if ok && nested.Type() != pcommon.ValueTypeMap {
			return fmt.Errorf("field %q already exists and isn't nested", fieldName)
		}

		if !ok {
			nested = pcommon.NewValueMap()
			nested.MapVal().Insert(ctx.IDENTIFIER(1).GetText(), pcommon.NewValueString(value))
		}

		attr.Insert(fieldName, nested)
		return nil
	}

	_, ok := attr.Get(fieldName)
	if ok {
		return fmt.Errorf("field %q is duplicated. Use an alias", fieldName)
	}
	attr.Insert(fieldName, pcommon.NewValueString(value))

	return nil
}
func insertAggregationColumnToRecord(ctx *AggregationColumnContext, attr pcommon.Map, value float64) error {
	if ctx.Alias() != nil {
		attr.Insert(ctx.Alias().GetStop().GetText(), pcommon.NewValueDouble(value))
		return nil
	}

	fieldName := ctx.IDENTIFIER(0).GetText()
	if len(ctx.AllIDENTIFIER()) > 1 {
		nested, ok := attr.Get(fieldName)

		if ok && nested.Type() != pcommon.ValueTypeMap {
			return fmt.Errorf("field %q already exists and isn't nested", fieldName)
		}

		if !ok {
			nested = pcommon.NewValueMap()
			nested.MapVal().Insert(ctx.IDENTIFIER(1).GetText(), pcommon.NewValueDouble(value))
		}

		attr.Insert(fieldName, nested)
		return nil
	}

	_, ok := attr.Get(fieldName)
	if ok {
		return fmt.Errorf("field %q is duplicated. Use an alias", fieldName)
	}
	attr.Insert(fieldName, pcommon.NewValueDouble(value))

	return nil
}

func sum(ls plog.LogRecordSlice, ctx *AggregationColumnContext) (float64, error) {
	var sum float64
	for i := 0; i < ls.Len(); i++ {
		curRec := ls.At(i)
		_, val, err := getAttributeValueForAggregation(ctx, curRec.Attributes())
		if err != nil {
			return 0, err
		}
		convertedVal, err := strconv.ParseFloat(val.AsString(), 64)
		if err != nil {
			return 0.0, err
		}
		sum += convertedVal
	}

	return sum, nil
}

func min(ls plog.LogRecordSlice, ctx *AggregationColumnContext) (float64, error) {
	var conErr error
	ls.Sort(func(a, b plog.LogRecord) bool {
		_, leftVal, err := getAttributeValueForAggregation(ctx, a.Attributes())
		if err != nil {
			conErr = err
		}

		_, rightVal, err := getAttributeValueForAggregation(ctx, b.Attributes())

		if err != nil {
			conErr = err
		}

		switch leftVal.Type() {
		case pcommon.ValueTypeInt:
			return leftVal.IntVal() < rightVal.IntVal()
		case pcommon.ValueTypeDouble:
			return leftVal.DoubleVal() < rightVal.DoubleVal()

		}

		return false
	})

	if ls.Len() > 0 {
		_, res, err := getAttributeValueForAggregation(ctx, ls.At(0).Attributes())
		convertedRes, _ := strconv.ParseFloat(res.AsString(), 64)
		return convertedRes, err
	}

	return 0, conErr
}

func max(ls plog.LogRecordSlice, ctx *AggregationColumnContext) (float64, error) {
	var conErr error
	sorted := ls.Sort(func(a, b plog.LogRecord) bool {
		_, leftVal, err := getAttributeValueForAggregation(ctx, a.Attributes())
		if err != nil {
			conErr = err
		}

		_, rightVal, err := getAttributeValueForAggregation(ctx, b.Attributes())
		if err != nil {
			conErr = err
		}

		switch leftVal.Type() {
		case pcommon.ValueTypeInt:
			return leftVal.IntVal() > rightVal.IntVal()
		case pcommon.ValueTypeDouble:
			return leftVal.DoubleVal() > rightVal.DoubleVal()

		}

		return false
	})

	if sorted.Len() > 0 {
		_, res, err := getAttributeValueForAggregation(ctx, sorted.At(0).Attributes())
		convertedRes, _ := strconv.ParseFloat(res.AsString(), 64)
		return convertedRes, err
	}

	return 0.0, conErr
}

func avg(ls plog.LogRecordSlice, ctx *AggregationColumnContext) (float64, error) {
	var sum float64
	for i := 0; i < ls.Len(); i++ {
		curRec := ls.At(i)
		_, val, err := getAttributeValueForAggregation(ctx, curRec.Attributes())
		if err != nil {
			return 0, err
		}
		convertedVal, err := strconv.ParseFloat(val.AsString(), 64)
		if err != nil {
			return 0.0, err
		}
		sum += convertedVal
	}
	return sum / float64(ls.Len()), nil
}

func count(ls plog.LogRecordSlice) int {
	return ls.Len()
}

// check if we need remove attribute missed in select list
func fieldExists(key string, value pcommon.Value, allColumns []interface{}) bool {
	//simple attribute
	if value.Type() != pcommon.ValueTypeMap {
		for _, col := range allColumns {
			var id string
			switch typedCol := col.(type) {
			case *IdentifierColumnContext:
				id = typedCol.IDENTIFIER(0).GetText()
			case *SimpleFunctionContext:
				id = typedCol.IDENTIFIER(0).GetText()
			}
			if key == id {
				return false
			}
		}
		return true
	}

	//find if nested attr existed in select list
	for _, col := range allColumns {
		switch typedCol := col.(type) {
		case *IdentifierColumnContext:
			if len(typedCol.AllIDENTIFIER()) == 1 {
				if key == typedCol.IDENTIFIER(0).GetText() {
					return false
				}
			}
		case *SimpleFunctionContext:
			if len(typedCol.AllIDENTIFIER()) == 1 {
				if key == typedCol.IDENTIFIER(0).GetText() {
					return false
				}
			}
		}
	}

	nestedFound := false
	//now check is we have only nested field in select list e.g. field.nestedField
	value.MapVal().RemoveIf(func(nestedKey string, value pcommon.Value) bool {
		for _, col := range allColumns {
			switch typedCol := col.(type) {
			case *IdentifierColumnContext:
				if len(typedCol.AllIDENTIFIER()) > 1 {
					if key == typedCol.IDENTIFIER(0).GetText() && nestedKey == typedCol.IDENTIFIER(1).GetText() {
						nestedFound = true
						return false
					}
				}
			case *SimpleFunctionContext:
				if len(typedCol.AllIDENTIFIER()) > 1 {
					if key == typedCol.IDENTIFIER(0).GetText() && nestedKey == typedCol.IDENTIFIER(1).GetText() {
						nestedFound = true
						return false
					}
				}
			}
		}
		return true
	})

	return !nestedFound
}

func getAttributeValueForAggregation(ctx *AggregationColumnContext, attr pcommon.Map) (string, pcommon.Value, error) {
	if len(ctx.AllIDENTIFIER()) > 1 {
		value, err := nestedFieldExistsInAttr(ctx.IDENTIFIER(0).GetText(), ctx.IDENTIFIER(1).GetText(), attr)
		return ctx.IDENTIFIER(0).GetText() + "," + ctx.IDENTIFIER(1).GetText(), value, err
	}
	value, err := fieldExistsInAttr(ctx.IDENTIFIER(0).GetText(), attr)
	return ctx.IDENTIFIER(0).GetText(), value, err

}

func getAttributeValue(column IColumnContext, attr pcommon.Map) (string, pcommon.Value, error) {
	switch col := column.(type) {
	case *IdentifierColumnContext:
		if len(col.AllIDENTIFIER()) > 1 {
			value, err := nestedFieldExistsInAttr(col.IDENTIFIER(0).GetText(), col.IDENTIFIER(1).GetText(), attr)
			return col.IDENTIFIER(0).GetText() + "," + col.IDENTIFIER(1).GetText(), value, err
		}
		value, err := fieldExistsInAttr(col.IDENTIFIER(0).GetText(), attr)
		return col.IDENTIFIER(0).GetText(), value, err
	default:
		return "", pcommon.NewValueEmpty(), fmt.Errorf("unknown column type")
	}

}

// check if field exists in attr
func fieldExistsInAttr(fieldName string, attr pcommon.Map) (pcommon.Value, error) {
	value, ok := attr.Get(fieldName)
	if !ok {
		return pcommon.NewValueEmpty(), fmt.Errorf("field %q missed", fieldName)
	}
	return value, nil
}

// check is nested field exists
func nestedFieldExistsInAttr(fieldName, nestedFieldName string, attr pcommon.Map) (pcommon.Value, error) {
	value, ok := attr.Get(fieldName)
	if !ok {
		return pcommon.NewValueEmpty(), fmt.Errorf("field %q missed", fieldName)
	}

	// if this nested field
	if value.Type() != pcommon.ValueTypeMap {
		return pcommon.NewValueEmpty(), fmt.Errorf("field %q isn't nested", fieldName)
	}

	nestedValue, ok := value.MapVal().Get(nestedFieldName)
	if !ok {
		return pcommon.NewValueEmpty(), fmt.Errorf("field %q missed nested field %q", fieldName, nestedFieldName)
	}

	return nestedValue, nil
}

func getSelectColumnsFromWhereCtx(ctx *WhereStmtContext) *SelectColumnsContext {
	for _, resCtx := range ctx.GetParent().GetChildren() {
		resColumnCtx, ok := resCtx.(*SelectColumnsContext)
		if !ok {
			continue
		}
		return resColumnCtx
	}

	return nil
}

func getSelectAggregationsFromWhereCtx(ctx *WhereStmtContext) *SelectAggregationsContext {
	for _, resCtx := range ctx.GetParent().GetChildren() {
		resColumnCtx, ok := resCtx.(*SelectAggregationsContext)
		if !ok {
			continue
		}
		return resColumnCtx
	}

	return nil
}

func getSelectStarCtx(ctx *WhereStmtContext) *SelectStarContext {
	for _, resCtx := range ctx.GetParent().GetChildren() {
		resColumnCtx, ok := resCtx.(*SelectStarContext)
		if !ok {
			continue
		}
		return resColumnCtx
	}

	return nil
}

func substring(value pcommon.Value, start, length string) (pcommon.Value, error) {

	startN, err := strconv.Atoi(start)
	if err != nil {
		return pcommon.NewValueEmpty(), fmt.Errorf("first arg must be numeric")
	}

	lengthN, err := strconv.Atoi(length)
	if err != nil {
		return pcommon.NewValueEmpty(), fmt.Errorf("second arg must be numeric")
	}

	runes := []rune(value.AsString())
	if startN > len(runes) {
		return pcommon.NewValueString(""), nil
	}

	if startN+lengthN > len(runes) {
		lengthN = len(runes) - startN
	}

	return pcommon.NewValueString(string(runes[startN : startN+lengthN])), nil

}

func getSimpleFunctionContext(ctx *FunctionColumnContext) *SimpleFunctionContext {
	var simpleFuncCtx *SimpleFunctionContext
	for _, childRecCtx := range ctx.Function().GetChildren() {
		switch childRecCtx.(type) {
		case *SimpleFunctionContext:
			simpleFuncCtx = childRecCtx.(*SimpleFunctionContext)
		case *RecursiveFunctionContext:
			simpleFuncCtx = findSimpleFuncCtx(childRecCtx.(*RecursiveFunctionContext))
		default:
			continue
		}
	}

	return simpleFuncCtx
}

func findSimpleFuncCtx(ctx *RecursiveFunctionContext) *SimpleFunctionContext {
	for _, childCtx := range ctx.GetChildren() {
		switch childCtx.(type) {
		case *RecursiveFunctionContext:
			return findSimpleFuncCtx(childCtx.(*RecursiveFunctionContext))
		case *SimpleFunctionContext:
			return childCtx.(*SimpleFunctionContext)
		default:
			continue
		}
	}

	return nil
}

func getGroupByFieldFromAggregationCtx(ctx *SelectGroupByAggregationsContext) (string, error) {
	for _, node := range ctx.GetParent().GetChildren() {
		groupByCtx, ok := node.(*GroupByContext)
		if !ok {
			continue
		}
		idCol, ok := groupByCtx.Column().(*IdentifierColumnContext)
		if !ok {
			return ",", fmt.Errorf("unknown group by clause field type")
		}
		return getFieldNameFromIDColumnCtx(idCol), nil
	}

	return "", fmt.Errorf("group by field missed")

}

func getFieldNameFromIDColumnCtx(col *IdentifierColumnContext) string {
	fieldName := col.IDENTIFIER(0).GetText()
	if len(col.AllIDENTIFIER()) > 1 {
		return fieldName + "." + col.IDENTIFIER(1).GetText()
	}
	return fieldName
}

func calculateGroupByAggregations(ctx *SelectGroupByAggregationsContext, key string, rec plog.LogRecord, ls plog.LogRecordSlice) error {
	for _, aggregation := range ctx.AllAggregationColumn() {

		col, ok := aggregation.(*AggregationColumnContext)
		if !ok {
			return fmt.Errorf("unknown aggragation column type")
		}
		res, err := getFieldAggregatedValue(col, ls)
		if err != nil {
			return err
		}
		if err = insertAggregationColumnToRecord(col, rec.Attributes(), res); err != nil {
			return err
		}
	}

	if len(ctx.AllColumn()) == 0 {
		return nil
	}

	for _, nonAgrCol := range ctx.AllColumn() {
		groupByFieldName, err := getGroupByFieldFromAggregationCtx(ctx)
		if err != nil {
			return err
		}

		selectField, ok := nonAgrCol.(*IdentifierColumnContext)
		if !ok {
			return fmt.Errorf("field %q in select list is inknown type", nonAgrCol.GetText())
		}

		selectFieldName := getFieldNameFromIDColumnCtx(selectField)
		if strings.Compare(groupByFieldName, selectFieldName) != 0 {
			return fmt.Errorf("you must use the same field for group by %q clause and select list", groupByFieldName)
		}

		if err := insertIdentifierColumnToRecord(selectField, rec.Attributes(), key); err != nil {
			return err
		}
	}

	return nil
}

func calculateAggregations(ctx *SelectAggregationsContext, rec plog.LogRecord, ls plog.LogRecordSlice) error {
	for _, aggregation := range ctx.AllAggregationColumn() {
		col, ok := aggregation.(*AggregationColumnContext)
		if !ok {
			return fmt.Errorf("unknown aggragation column type")
		}
		res, err := getFieldAggregatedValue(col, ls)
		if err != nil {
			return err
		}
		if err = insertAggregationColumnToRecord(col, rec.Attributes(), res); err != nil {
			return err
		}
	}
	return nil
}
