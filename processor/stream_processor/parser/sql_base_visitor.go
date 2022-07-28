// Code generated from Sql.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Sql

import "github.com/antlr/antlr4/runtime/Go/antlr"

type BaseSqlVisitor struct {
	*antlr.BaseParseTreeVisitor
}

func (v *BaseSqlVisitor) VisitSqlQuery(ctx *SqlQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSelectQuery(ctx *SelectQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSelectColumns(ctx *SelectColumnsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSelectAVG(ctx *SelectAVGContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSelectStar(ctx *SelectStarContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitColumn(ctx *ColumnContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitWhereStmt(ctx *WhereStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitTumblingStmt(ctx *TumblingStmtContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitTumblingWindow(ctx *TumblingWindowContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSimpleCondition(ctx *SimpleConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitNullCondition(ctx *NullConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitCompoundRecursiveCondition(ctx *CompoundRecursiveConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitSimpleRecursiveCondition(ctx *SimpleRecursiveConditionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *BaseSqlVisitor) VisitCompoundExpr(ctx *CompoundExprContext) interface{} {
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

func (v *BaseSqlVisitor) VisitAvg(ctx *AvgContext) interface{} {
	return v.VisitChildren(ctx)
}
