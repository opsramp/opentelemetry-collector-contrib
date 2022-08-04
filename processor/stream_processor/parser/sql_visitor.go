// Code generated from Sql.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Sql

import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by SqlParser.
type SqlVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by SqlParser#sqlQuery.
	VisitSqlQuery(ctx *SqlQueryContext) interface{}

	// Visit a parse tree produced by SqlParser#selectSimple.
	VisitSelectSimple(ctx *SelectSimpleContext) interface{}

	// Visit a parse tree produced by SqlParser#selectTumbling.
	VisitSelectTumbling(ctx *SelectTumblingContext) interface{}

	// Visit a parse tree produced by SqlParser#selectColumns.
	VisitSelectColumns(ctx *SelectColumnsContext) interface{}

	// Visit a parse tree produced by SqlParser#selectAVG.
	VisitSelectAVG(ctx *SelectAVGContext) interface{}

	// Visit a parse tree produced by SqlParser#selectStar.
	VisitSelectStar(ctx *SelectStarContext) interface{}

	// Visit a parse tree produced by SqlParser#column.
	VisitColumn(ctx *ColumnContext) interface{}

	// Visit a parse tree produced by SqlParser#whereStmt.
	VisitWhereStmt(ctx *WhereStmtContext) interface{}

	// Visit a parse tree produced by SqlParser#windowTumbling.
	VisitWindowTumbling(ctx *WindowTumblingContext) interface{}

	// Visit a parse tree produced by SqlParser#simpleCondition.
	VisitSimpleCondition(ctx *SimpleConditionContext) interface{}

	// Visit a parse tree produced by SqlParser#compoundRecursiveCondition.
	VisitCompoundRecursiveCondition(ctx *CompoundRecursiveConditionContext) interface{}

	// Visit a parse tree produced by SqlParser#simpleRecursiveCondition.
	VisitSimpleRecursiveCondition(ctx *SimpleRecursiveConditionContext) interface{}

	// Visit a parse tree produced by SqlParser#compoundExpr.
	VisitCompoundExpr(ctx *CompoundExprContext) interface{}

	// Visit a parse tree produced by SqlParser#comparisonOperator.
	VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{}

	// Visit a parse tree produced by SqlParser#literalValue.
	VisitLiteralValue(ctx *LiteralValueContext) interface{}

	// Visit a parse tree produced by SqlParser#groupBy.
	VisitGroupBy(ctx *GroupByContext) interface{}

	// Visit a parse tree produced by SqlParser#avg.
	VisitAvg(ctx *AvgContext) interface{}
}
