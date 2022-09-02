package parser

import (
	"fmt"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"strconv"
	"strings"
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

func sum(ls plog.LogRecordSlice, fieldName string) (float64, error) {
	var sum float64
	for i := 0; i < ls.Len(); i++ {
		curRec := ls.At(i)
		val, ok := curRec.Attributes().Get(fieldName)
		if !ok {
			return 0, fmt.Errorf("field %q missed", fieldName)
		}
		convertedVal, err := strconv.ParseFloat(val.AsString(), 64)
		if err != nil {
			return 0.0, err
		}
		sum = sum + convertedVal
	}

	return sum, nil
}

func min(ls plog.LogRecordSlice, fieldName string) (float64, error) {
	var conErr error
	ls.Sort(func(a, b plog.LogRecord) bool {
		leftVal, ok := a.Attributes().Get(fieldName)
		if !ok {
			conErr = fmt.Errorf("field %q missed", fieldName)
		}
		convertedLeft := leftVal.DoubleVal()

		rightVal, ok := b.Attributes().Get(fieldName)

		if !ok {
			conErr = fmt.Errorf("field %q missed", fieldName)
		}

		convertedRight := rightVal.DoubleVal()

		return convertedLeft < convertedRight
	})

	if ls.Len() > 0 {
		res, _ := ls.At(0).Attributes().Get(fieldName)
		convertedRes, _ := strconv.ParseFloat(res.AsString(), 64)
		return convertedRes, nil
	}

	return 0, conErr
}

func max(ls plog.LogRecordSlice, fieldName string) (float64, error) {
	var conErr error
	ls.Sort(func(a, b plog.LogRecord) bool {
		leftVal, ok := a.Attributes().Get(fieldName)
		if !ok {
			conErr = fmt.Errorf("field %q missed", fieldName)
		}
		convertedLeft, err := strconv.ParseFloat(leftVal.AsString(), 64)
		if err != nil {
			conErr = err
		}

		rightVal, ok := b.Attributes().Get(fieldName)

		if !ok {
			conErr = fmt.Errorf("field %q missed", fieldName)
		}

		convertedRight, err := strconv.ParseFloat(rightVal.AsString(), 64)
		if err != nil {
			conErr = err
		}

		return convertedLeft > convertedRight
	})

	if ls.Len() > 0 {
		res, _ := ls.At(0).Attributes().Get(fieldName)
		convertedRes, _ := strconv.ParseFloat(res.AsString(), 64)
		return convertedRes, nil
	}

	return 0.0, conErr
}

func avg(ls plog.LogRecordSlice, fieldName string) (float64, error) {
	var sum float64
	for i := 0; i < ls.Len(); i++ {
		curRec := ls.At(i)
		val, ok := curRec.Attributes().Get(fieldName)
		if !ok {
			return 0, fmt.Errorf("field %q missed", fieldName)
		}
		convertedVal, err := strconv.ParseFloat(val.AsString(), 64)
		if err != nil {
			return 0.0, err
		}
		sum = sum + convertedVal
	}
	return sum / float64(ls.Len()), nil
}

func count(ls plog.LogRecordSlice) int {
	return ls.Len()
}

func lower(record plog.LogRecord, fieldName string) error {
	val, ok := record.Attributes().Get(fieldName)
	if !ok {
		return fmt.Errorf("field %q missed", fieldName)
	}
	record.Attributes().UpdateString(fieldName, strings.ToLower(val.AsString()))
	return nil
}

func upper(record plog.LogRecord, fieldName string) error {
	val, ok := record.Attributes().Get(fieldName)
	if !ok {
		return fmt.Errorf("field %q missed", fieldName)
	}
	record.Attributes().UpdateString(fieldName, strings.ToUpper(val.AsString()))
	return nil
}

// check if we need remove attribute missed in select list
func fieldExists(key string, value pcommon.Value, allColumns []IColumnContext) bool {
	//simple attribute
	if value.Type() != pcommon.ValueTypeMap {
		for _, col := range allColumns {
			var id string
			switch typedCol := col.(type) {
			case *IdentifierColumnContext:
				id = typedCol.IDENTIFIER(0).GetText()
			case *FunctionColumnContext:
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
		case *FunctionColumnContext:
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
			case *FunctionColumnContext:
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

// check if field exists in attr
func fieldExistsInAttr(fieldName string, attr pcommon.Map) error {
	_, ok := attr.Get(fieldName)
	if !ok {
		return fmt.Errorf("field %q missed", fieldName)
	}
	return nil
}

// check is nested field exists
func nestedFieldExistsInAttr(fieldName, nestedFieldName string, attr pcommon.Map) error {
	value, ok := attr.Get(fieldName)
	if !ok {
		return fmt.Errorf("field %q missed", fieldName)
	}

	//if this nested field
	if value.Type() != pcommon.ValueTypeMap {
		return fmt.Errorf("field %q isn't nested", fieldName)
	}

	_, ok = value.MapVal().Get(nestedFieldName)
	if !ok {
		return fmt.Errorf("field %q missed nested field %q", fieldName, nestedFieldName)
	}

	return nil
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
