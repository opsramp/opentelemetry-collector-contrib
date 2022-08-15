package parser

import (
	"fmt"
	"go.opentelemetry.io/collector/pdata/plog"
	"strconv"
	"strings"
)

func compareString(ctx *SimpleExpressionContext, fieldVal, comparisonVal string) bool {
	switch ctx.ComparisonOperator().GetStart().GetTokenType() {
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

func compareNumeric(ctx *SimpleExpressionContext, fieldVal, comparisonVal float64) bool {
	switch ctx.ComparisonOperator().GetStart().GetTokenType() {
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

func compareBool(ctx *SimpleExpressionContext, fieldVal, comparisonVal bool) bool {
	switch ctx.ComparisonOperator().GetStart().GetTokenType() {
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

func getColumnIdentifier(column IColumnContext) string {
	switch col := column.GetRuleContext().(type) {
	case *FunctionColContext:
		return col.IDENTIFIER().GetText()
	case *IdentifierColContext:
		return col.IDENTIFIER().GetText()
	default:
		return ""
	}
}

func KeyExists(key string, resultColumns []IColumnContext) bool {
	for _, col := range resultColumns {
		var id string
		switch col := col.GetRuleContext().(type) {
		case *FunctionColContext:
			id = col.IDENTIFIER().GetText()
		case *IdentifierColContext:
			id = col.IDENTIFIER().GetText()
		default:
			return false
		}
		if key == id {
			return true
		}
	}
	return false
}
