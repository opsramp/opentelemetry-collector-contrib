// Code generated from Sql.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Sql

import "github.com/antlr/antlr4/runtime/Go/antlr"

// A complete Visitor for a parse tree produced by SqlParser.
type SqlVisitor interface {
	antlr.ParseTreeVisitor

	// Visit a parse tree produced by SqlParser#sqlQuery.
	VisitSqlQuery(ctx *SqlQueryContext) interface{}

	// Visit a parse tree produced by SqlParser#selectQuery.
	VisitSelectQuery(ctx *SelectQueryContext) interface{}

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

	// Visit a parse tree produced by SqlParser#tumblingStmt.
	VisitTumblingStmt(ctx *TumblingStmtContext) interface{}

	// Visit a parse tree produced by SqlParser#tumblingWindow.
	VisitTumblingWindow(ctx *TumblingWindowContext) interface{}

	// Visit a parse tree produced by SqlParser#groupBy.
	VisitGroupBy(ctx *GroupByContext) interface{}

	// Visit a parse tree produced by SqlParser#avg.
	VisitAvg(ctx *AvgContext) interface{}

	// Visit a parse tree produced by SqlParser#expr.
	VisitExpr(ctx *ExprContext) interface{}

	// Visit a parse tree produced by SqlParser#comparisonOperator.
	VisitComparisonOperator(ctx *ComparisonOperatorContext) interface{}

	// Visit a parse tree produced by SqlParser#literalValue.
	VisitLiteralValue(ctx *LiteralValueContext) interface{}
}
