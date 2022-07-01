// Code generated from Sql.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Sql

import "github.com/antlr/antlr4/runtime/Go/antlr"

// SqlListener is a complete listener for a parse tree produced by SqlParser.
type SqlListener interface {
	antlr.ParseTreeListener

	// EnterSqlQuery is called when entering the sqlQuery production.
	EnterSqlQuery(c *SqlQueryContext)

	// EnterSelectQuery is called when entering the selectQuery production.
	EnterSelectQuery(c *SelectQueryContext)

	// EnterResultColumns is called when entering the resultColumns production.
	EnterResultColumns(c *ResultColumnsContext)

	// EnterColumn is called when entering the column production.
	EnterColumn(c *ColumnContext)

	// EnterWhereStatement is called when entering the whereStatement production.
	EnterWhereStatement(c *WhereStatementContext)

	// EnterTumblingWindow is called when entering the tumblingWindow production.
	EnterTumblingWindow(c *TumblingWindowContext)

	// EnterGroupBy is called when entering the groupBy production.
	EnterGroupBy(c *GroupByContext)

	// EnterAvg is called when entering the avg production.
	EnterAvg(c *AvgContext)

	// EnterExpr is called when entering the expr production.
	EnterExpr(c *ExprContext)

	// EnterComparisonOperator is called when entering the comparisonOperator production.
	EnterComparisonOperator(c *ComparisonOperatorContext)

	// EnterLiteralValue is called when entering the literalValue production.
	EnterLiteralValue(c *LiteralValueContext)

	// ExitSqlQuery is called when exiting the sqlQuery production.
	ExitSqlQuery(c *SqlQueryContext)

	// ExitSelectQuery is called when exiting the selectQuery production.
	ExitSelectQuery(c *SelectQueryContext)

	// ExitResultColumns is called when exiting the resultColumns production.
	ExitResultColumns(c *ResultColumnsContext)

	// ExitColumn is called when exiting the column production.
	ExitColumn(c *ColumnContext)

	// ExitWhereStatement is called when exiting the whereStatement production.
	ExitWhereStatement(c *WhereStatementContext)

	// ExitTumblingWindow is called when exiting the tumblingWindow production.
	ExitTumblingWindow(c *TumblingWindowContext)

	// ExitGroupBy is called when exiting the groupBy production.
	ExitGroupBy(c *GroupByContext)

	// ExitAvg is called when exiting the avg production.
	ExitAvg(c *AvgContext)

	// ExitExpr is called when exiting the expr production.
	ExitExpr(c *ExprContext)

	// ExitComparisonOperator is called when exiting the comparisonOperator production.
	ExitComparisonOperator(c *ComparisonOperatorContext)

	// ExitLiteralValue is called when exiting the literalValue production.
	ExitLiteralValue(c *LiteralValueContext)
}
