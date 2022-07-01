// Code generated from Sql.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Sql

import "github.com/antlr/antlr4/runtime/Go/antlr"

// BaseSqlListener is a complete listener for a parse tree produced by SqlParser.
type BaseSqlListener struct{}

var _ SqlListener = &BaseSqlListener{}

// VisitTerminal is called when a terminal node is visited.
func (s *BaseSqlListener) VisitTerminal(node antlr.TerminalNode) {}

// VisitErrorNode is called when an error node is visited.
func (s *BaseSqlListener) VisitErrorNode(node antlr.ErrorNode) {}

// EnterEveryRule is called when any rule is entered.
func (s *BaseSqlListener) EnterEveryRule(ctx antlr.ParserRuleContext) {}

// ExitEveryRule is called when any rule is exited.
func (s *BaseSqlListener) ExitEveryRule(ctx antlr.ParserRuleContext) {}

// EnterSqlQuery is called when production sqlQuery is entered.
func (s *BaseSqlListener) EnterSqlQuery(ctx *SqlQueryContext) {}

// ExitSqlQuery is called when production sqlQuery is exited.
func (s *BaseSqlListener) ExitSqlQuery(ctx *SqlQueryContext) {}

// EnterSelectQuery is called when production selectQuery is entered.
func (s *BaseSqlListener) EnterSelectQuery(ctx *SelectQueryContext) {}

// ExitSelectQuery is called when production selectQuery is exited.
func (s *BaseSqlListener) ExitSelectQuery(ctx *SelectQueryContext) {}

// EnterResultColumns is called when production resultColumns is entered.
func (s *BaseSqlListener) EnterResultColumns(ctx *ResultColumnsContext) {}

// ExitResultColumns is called when production resultColumns is exited.
func (s *BaseSqlListener) ExitResultColumns(ctx *ResultColumnsContext) {}

// EnterColumn is called when production column is entered.
func (s *BaseSqlListener) EnterColumn(ctx *ColumnContext) {}

// ExitColumn is called when production column is exited.
func (s *BaseSqlListener) ExitColumn(ctx *ColumnContext) {}

// EnterWhereStatement is called when production whereStatement is entered.
func (s *BaseSqlListener) EnterWhereStatement(ctx *WhereStatementContext) {}

// ExitWhereStatement is called when production whereStatement is exited.
func (s *BaseSqlListener) ExitWhereStatement(ctx *WhereStatementContext) {}

// EnterTumblingWindow is called when production tumblingWindow is entered.
func (s *BaseSqlListener) EnterTumblingWindow(ctx *TumblingWindowContext) {}

// ExitTumblingWindow is called when production tumblingWindow is exited.
func (s *BaseSqlListener) ExitTumblingWindow(ctx *TumblingWindowContext) {}

// EnterGroupBy is called when production groupBy is entered.
func (s *BaseSqlListener) EnterGroupBy(ctx *GroupByContext) {}

// ExitGroupBy is called when production groupBy is exited.
func (s *BaseSqlListener) ExitGroupBy(ctx *GroupByContext) {}

// EnterAvg is called when production avg is entered.
func (s *BaseSqlListener) EnterAvg(ctx *AvgContext) {}

// ExitAvg is called when production avg is exited.
func (s *BaseSqlListener) ExitAvg(ctx *AvgContext) {}

// EnterExpr is called when production expr is entered.
func (s *BaseSqlListener) EnterExpr(ctx *ExprContext) {}

// ExitExpr is called when production expr is exited.
func (s *BaseSqlListener) ExitExpr(ctx *ExprContext) {}

// EnterComparisonOperator is called when production comparisonOperator is entered.
func (s *BaseSqlListener) EnterComparisonOperator(ctx *ComparisonOperatorContext) {}

// ExitComparisonOperator is called when production comparisonOperator is exited.
func (s *BaseSqlListener) ExitComparisonOperator(ctx *ComparisonOperatorContext) {}

// EnterLiteralValue is called when production literalValue is entered.
func (s *BaseSqlListener) EnterLiteralValue(ctx *LiteralValueContext) {}

// ExitLiteralValue is called when production literalValue is exited.
func (s *BaseSqlListener) ExitLiteralValue(ctx *LiteralValueContext) {}
