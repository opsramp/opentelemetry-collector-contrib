// Code generated from Sql.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Sql

import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseSqlVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseSqlVisitor) VisitSqlQuery(ctx *SqlQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSelectSimple(ctx *SelectSimpleContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSelectTumbling(ctx *SelectTumblingContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSelectTumblingGroupBy(ctx *SelectTumblingGroupByContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitWindowTumbling(ctx *WindowTumblingContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSelectColumns(ctx *SelectColumnsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSelectStar(ctx *SelectStarContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSelectAggregations(ctx *SelectAggregationsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSelectGroupByAggregations(ctx *SelectGroupByAggregationsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitAggregationColumn(ctx *AggregationColumnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitIdentifierColumn(ctx *IdentifierColumnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitFunctionColumn(ctx *FunctionColumnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitAlias(ctx *AliasContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSimpleFunction(ctx *SimpleFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitRecursiveFunction(ctx *RecursiveFunctionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitFunctionName(ctx *FunctionNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitWhereStmt(ctx *WhereStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSimpleCondition(ctx *SimpleConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitCompoundRecursiveCondition(ctx *CompoundRecursiveConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSimpleCompoundCondition(ctx *SimpleCompoundConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSimpleRecursiveCondition(ctx *SimpleRecursiveConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSimpleExpression(ctx *SimpleExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitNestedExpression(ctx *NestedExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitCompoundExpression(ctx *CompoundExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitLiteralValue(ctx *LiteralValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitGroupBy(ctx *GroupByContext) interface{} {
	return v.VisitChildren(ctx)
}
