// Code generated from sql.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // sql

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BasesqlListener is a complete listener for a parse tree produced by sqlParser.
type BasesqlListener struct{}

var _ sqlListener = &BasesqlListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BasesqlListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BasesqlListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BasesqlListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BasesqlListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSqlQuery is called when production sqlQuery is entered.
func (s *BasesqlListener) EnterSqlQuery(ctx *SqlQueryContext) {}

// ExitSqlQuery is called when production sqlQuery is exited.
func (s *BasesqlListener) ExitSqlQuery(ctx *SqlQueryContext) {}

// EnterSelectQuery is called when production selectQuery is entered.
func (s *BasesqlListener) EnterSelectQuery(ctx *SelectQueryContext) {}

// ExitSelectQuery is called when production selectQuery is exited.
func (s *BasesqlListener) ExitSelectQuery(ctx *SelectQueryContext) {}

// EnterResultColumns is called when production resultColumns is entered.
func (s *BasesqlListener) EnterResultColumns(ctx *ResultColumnsContext) {}

// ExitResultColumns is called when production resultColumns is exited.
func (s *BasesqlListener) ExitResultColumns(ctx *ResultColumnsContext) {}

// EnterColumn is called when production column is entered.
func (s *BasesqlListener) EnterColumn(ctx *ColumnContext) {}

// ExitColumn is called when production column is exited.
func (s *BasesqlListener) ExitColumn(ctx *ColumnContext) {}

// EnterWhereStatement is called when production whereStatement is entered.
func (s *BasesqlListener) EnterWhereStatement(ctx *WhereStatementContext) {}

// ExitWhereStatement is called when production whereStatement is exited.
func (s *BasesqlListener) ExitWhereStatement(ctx *WhereStatementContext) {}

// EnterGroupBy is called when production groupBy is entered.
func (s *BasesqlListener) EnterGroupBy(ctx *GroupByContext) {}

// ExitGroupBy is called when production groupBy is exited.
func (s *BasesqlListener) ExitGroupBy(ctx *GroupByContext) {}

// EnterAvg is called when production avg is entered.
func (s *BasesqlListener) EnterAvg(ctx *AvgContext) {}

// ExitAvg is called when production avg is exited.
func (s *BasesqlListener) ExitAvg(ctx *AvgContext) {}

// EnterExpr is called when production expr is entered.
func (s *BasesqlListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BasesqlListener) ExitExpr(ctx *ExprContext) {}

// EnterComparisonOperator is called when production comparisonOperator is entered.
func (s *BasesqlListener) EnterComparisonOperator(ctx *ComparisonOperatorContext) {}

// ExitComparisonOperator is called when production comparisonOperator is exited.
func (s *BasesqlListener) ExitComparisonOperator(ctx *ComparisonOperatorContext) {}

// EnterLiteralValue is called when production literalValue is entered.
func (s *BasesqlListener) EnterLiteralValue(ctx *LiteralValueContext) {}

// ExitLiteralValue is called when production literalValue is exited.
func (s *BasesqlListener) ExitLiteralValue(ctx *LiteralValueContext) {}
