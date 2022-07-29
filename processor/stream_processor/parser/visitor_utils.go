package parser

import (
	"fmt"
	"strings"
)

func KeyExists(key string, resultColumns []IColumnContext) bool {
	for _, col := range resultColumns {
		if key == col.GetText() {
			return true
		}
	}
	return false
}

func compareString(ctx *SimpleConditionContext, fieldVal, comparisonVal string) bool {
	switch ctx.ComparisonOperator().GetStart().GetTokenType() {
	case SqlParserK_EQUAL:
		return fieldVal == comparisonVal
	case SqlParserK_NOT_EQUAL:
		return fieldVal != comparisonVal
	case SqlParserK_GREATER:
		fmt.Println(fieldVal, " ", comparisonVal, " ", fieldVal > comparisonVal)
		return fieldVal > comparisonVal
	case SqlParserK_GREATER_EQUAL:
		return fieldVal >= comparisonVal
	case SqlParserK_LESS:
		return fieldVal < comparisonVal
	case SqlParserK_LESS_EQUAL:
		return fieldVal <= comparisonVal
	case SqlParserK_LIKE:
		return strings.Contains(fmt.Sprint(fieldVal), fmt.Sprint(comparisonVal))
	case SqlParserK_IS_NULL:
		return len(fmt.Sprint(fieldVal)) == 0
	case SqlParserK_IS_NOT_NULL:
		return len(fmt.Sprint(fieldVal)) != 0
	default:
		return false
	}
}

func compareNumeric(ctx *SimpleConditionContext, fieldVal, comparisonVal int) bool {
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

func compareBool(ctx *SimpleConditionContext, fieldVal, comparisonVal bool) bool {
	switch ctx.ComparisonOperator().GetStart().GetTokenType() {
	case SqlParserK_EQUAL:
		return fieldVal == comparisonVal
	case SqlParserK_NOT_EQUAL:
		return fieldVal != comparisonVal
	}
	return false
}
