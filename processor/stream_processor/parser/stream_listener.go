package parser

import (
	"context"
	"fmt"
	"go.opentelemetry.io/collector/pdata/plog"
)

type StreamListener struct {
	*BaseSqlListener
	in, out chan plog.LogRecord
	ctx     context.Context
	log     plog.LogRecord
}

func NewStreamListener(in chan plog.LogRecord, out chan plog.LogRecord, ctx context.Context) StreamListener {
	return StreamListener{
		in:  in,
		out: out,
		ctx: ctx,
	}
}

func (s *StreamListener) EnterResultColumns(ctx *ResultColumnsContext) {
	fmt.Println(ctx.Column(0).GetText())
	//for col := ctx.get
}

func (s *StreamListener) EnterSqlQuery(ctx *SqlQueryContext) {
	//fmt.Println(ctx.GetPayload())
}

func (s *StreamListener) EnterSelectQuery(ctx *SelectQueryContext) {
	//fmt.Println(ctx.GetParser().GetCurrentToken().GetText())
	fmt.Println(ctx.GetToken(SqlParserK_SELECT, 0).GetText())

}

func (s *StreamListener) EnterWhereStatement(ctx *WhereStatementContext) {
	fmt.Println(ctx.GetParser().GetCurrentToken().GetText())
}

func (s *StreamListener) EnterExpr(ctx *ExprContext) {
	fmt.Println(ctx.GetRuleContext().GetText())
	fmt.Println(ctx.GetRuleIndex())
}
