package parser

import (
	"errors"
	"strconv"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

//This is just to determine if it is window tumbling query and get time period value
//TODO find better way
var _ SqlVisitor = (*TumblingVisitor)(nil)

type TumblingVisitor struct {
	*BaseSqlVisitor
}

func IsTumblingQuery(query string) (bool, int) {
	is := antlr.NewInputStream(query)
	lexer := NewSqlLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	stream.Size()
	parser := NewSqlParser(stream)

	t := &TumblingVisitor{
		BaseSqlVisitor: &BaseSqlVisitor{},
	}
	switch res := t.Visit(parser.SqlQuery()).(type) {
	case int:
		return true, res
	case error:
		return false, 0
	}

	return false, 0
}

func (v *TumblingVisitor) Visit(tree antlr.ParseTree) interface{} {
	return tree.(*SqlQueryContext).Accept(v)
}

func (v *TumblingVisitor) VisitSqlQuery(ctx *SqlQueryContext) interface{} {
	return ctx.SelectQuery().Accept(v)
}

func (v *TumblingVisitor) VisitSelectTumbling(ctx *SelectTumblingContext) interface{} {
	return ctx.WindowTumbling().Accept(v)
}

func (v *TumblingVisitor) VisitWindowTumbling(ctx *WindowTumblingContext) interface{} {
	if ctx.NUMERIC_LITERAL() != nil {
		res, _ := strconv.Atoi(ctx.NUMERIC_LITERAL().GetText())
		return res
	}
	return errors.New("this is not window tumbling query")
}

func (v *TumblingVisitor) VisitGroupBy(ctx *GroupByContext) interface{} {
	return nil
}

func (v *TumblingVisitor) VisitSelectTumblingGroupBy(ctx *SelectTumblingGroupByContext) interface{} {
	return ctx.WindowTumbling().Accept(v)
}
