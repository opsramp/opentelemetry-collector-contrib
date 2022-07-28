// Code generated from Sql.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // Sql

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

// Suppress unused import errors
var _ = fmt.Printf
var _ = strconv.Itoa
var _ = sync.Once{}

type SqlParser struct {
	*antlr.BaseParser
}

var sqlParserStaticData struct {
	once                   sync.Once
	serializedATN          []int32
	literalNames           []string
	symbolicNames          []string
	ruleNames              []string
	predictionContextCache *antlr.PredictionContextCache
	atn                    *antlr.ATN
	decisionToDFA          []*antlr.DFA
}

func sqlParserInit() {
	staticData := &sqlParserStaticData
	staticData.literalNames = []string{
		"", "", "", "','", "'('", "')'", "';'", "", "", "", "", "", "", "",
		"", "'='", "'>'", "'<'", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "'*'",
	}
	staticData.symbolicNames = []string{
		"", "SPACE", "WS", "COMMA", "L_BRACKET", "R_BRACKET", "EOQ", "K_SELECT",
		"K_WHERE", "K_WINDOW_TUMBLING", "K_GROUP_BY", "K_AND", "K_OR", "K_IS",
		"K_LIKE", "K_EQUAL", "K_GREATER", "K_LESS", "K_LESS_EQUAL", "K_GREATER_EQUAL",
		"K_NOT_EQUAL", "K_NULL", "K_IS_NULL", "K_IS_NOT_NULL", "K_NOT", "K_NOT_IN",
		"K_IN", "K_TRUE", "K_FALSE", "K_COUNT", "K_MIN", "K_MAX", "K_AVG", "IDENTIFIER",
		"BOOLEAN_LITERAL", "NUMERIC_LITERAL", "STRING_LITERAL", "STAR",
	}
	staticData.ruleNames = []string{
		"sqlQuery", "selectQuery", "resultColumns", "column", "whereStatement",
		"tumblingWindow", "expr", "compoundExpr", "comparisonOperator", "literalValue",
		"groupBy", "avg",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 37, 111, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 3, 1, 31, 8, 1, 1,
		1, 1, 1, 1, 2, 1, 2, 1, 2, 5, 2, 38, 8, 2, 10, 2, 12, 2, 41, 9, 2, 1, 2,
		1, 2, 1, 2, 4, 2, 46, 8, 2, 11, 2, 12, 2, 47, 1, 2, 3, 2, 51, 8, 2, 1,
		3, 1, 3, 1, 4, 1, 4, 3, 4, 57, 8, 4, 1, 4, 1, 4, 1, 4, 3, 4, 62, 8, 4,
		1, 4, 1, 4, 3, 4, 66, 8, 4, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 72, 8, 5, 1,
		6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 3, 6, 85,
		8, 6, 1, 6, 1, 6, 1, 6, 5, 6, 90, 8, 6, 10, 6, 12, 6, 93, 9, 6, 1, 7, 1,
		7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 1, 11, 1, 11,
		1, 11, 1, 11, 1, 11, 1, 11, 0, 1, 12, 12, 0, 2, 4, 6, 8, 10, 12, 14, 16,
		18, 20, 22, 0, 5, 1, 0, 11, 12, 1, 0, 22, 23, 2, 0, 13, 20, 25, 26, 1,
		0, 34, 36, 1, 0, 29, 32, 110, 0, 24, 1, 0, 0, 0, 2, 27, 1, 0, 0, 0, 4,
		50, 1, 0, 0, 0, 6, 52, 1, 0, 0, 0, 8, 65, 1, 0, 0, 0, 10, 67, 1, 0, 0,
		0, 12, 84, 1, 0, 0, 0, 14, 94, 1, 0, 0, 0, 16, 98, 1, 0, 0, 0, 18, 100,
		1, 0, 0, 0, 20, 102, 1, 0, 0, 0, 22, 105, 1, 0, 0, 0, 24, 25, 3, 2, 1,
		0, 25, 26, 5, 0, 0, 1, 26, 1, 1, 0, 0, 0, 27, 28, 5, 7, 0, 0, 28, 30, 3,
		4, 2, 0, 29, 31, 3, 8, 4, 0, 30, 29, 1, 0, 0, 0, 30, 31, 1, 0, 0, 0, 31,
		32, 1, 0, 0, 0, 32, 33, 5, 6, 0, 0, 33, 3, 1, 0, 0, 0, 34, 39, 3, 6, 3,
		0, 35, 36, 5, 3, 0, 0, 36, 38, 3, 6, 3, 0, 37, 35, 1, 0, 0, 0, 38, 41,
		1, 0, 0, 0, 39, 37, 1, 0, 0, 0, 39, 40, 1, 0, 0, 0, 40, 51, 1, 0, 0, 0,
		41, 39, 1, 0, 0, 0, 42, 45, 3, 22, 11, 0, 43, 44, 5, 3, 0, 0, 44, 46, 3,
		22, 11, 0, 45, 43, 1, 0, 0, 0, 46, 47, 1, 0, 0, 0, 47, 45, 1, 0, 0, 0,
		47, 48, 1, 0, 0, 0, 48, 51, 1, 0, 0, 0, 49, 51, 5, 37, 0, 0, 50, 34, 1,
		0, 0, 0, 50, 42, 1, 0, 0, 0, 50, 49, 1, 0, 0, 0, 51, 5, 1, 0, 0, 0, 52,
		53, 5, 33, 0, 0, 53, 7, 1, 0, 0, 0, 54, 55, 5, 8, 0, 0, 55, 57, 3, 12,
		6, 0, 56, 54, 1, 0, 0, 0, 56, 57, 1, 0, 0, 0, 57, 66, 1, 0, 0, 0, 58, 61,
		3, 10, 5, 0, 59, 60, 5, 8, 0, 0, 60, 62, 3, 12, 6, 0, 61, 59, 1, 0, 0,
		0, 61, 62, 1, 0, 0, 0, 62, 63, 1, 0, 0, 0, 63, 64, 3, 20, 10, 0, 64, 66,
		1, 0, 0, 0, 65, 56, 1, 0, 0, 0, 65, 58, 1, 0, 0, 0, 66, 9, 1, 0, 0, 0,
		67, 68, 5, 9, 0, 0, 68, 71, 5, 35, 0, 0, 69, 70, 5, 8, 0, 0, 70, 72, 3,
		12, 6, 0, 71, 69, 1, 0, 0, 0, 71, 72, 1, 0, 0, 0, 72, 11, 1, 0, 0, 0, 73,
		74, 6, 6, -1, 0, 74, 75, 5, 33, 0, 0, 75, 76, 3, 16, 8, 0, 76, 77, 3, 18,
		9, 0, 77, 85, 1, 0, 0, 0, 78, 79, 3, 14, 7, 0, 79, 80, 7, 0, 0, 0, 80,
		81, 3, 14, 7, 0, 81, 85, 1, 0, 0, 0, 82, 83, 5, 33, 0, 0, 83, 85, 7, 1,
		0, 0, 84, 73, 1, 0, 0, 0, 84, 78, 1, 0, 0, 0, 84, 82, 1, 0, 0, 0, 85, 91,
		1, 0, 0, 0, 86, 87, 10, 1, 0, 0, 87, 88, 7, 0, 0, 0, 88, 90, 3, 12, 6,
		2, 89, 86, 1, 0, 0, 0, 90, 93, 1, 0, 0, 0, 91, 89, 1, 0, 0, 0, 91, 92,
		1, 0, 0, 0, 92, 13, 1, 0, 0, 0, 93, 91, 1, 0, 0, 0, 94, 95, 5, 4, 0, 0,
		95, 96, 3, 12, 6, 0, 96, 97, 5, 5, 0, 0, 97, 15, 1, 0, 0, 0, 98, 99, 7,
		2, 0, 0, 99, 17, 1, 0, 0, 0, 100, 101, 7, 3, 0, 0, 101, 19, 1, 0, 0, 0,
		102, 103, 5, 10, 0, 0, 103, 104, 3, 6, 3, 0, 104, 21, 1, 0, 0, 0, 105,
		106, 7, 4, 0, 0, 106, 107, 5, 4, 0, 0, 107, 108, 3, 6, 3, 0, 108, 109,
		5, 5, 0, 0, 109, 23, 1, 0, 0, 0, 10, 30, 39, 47, 50, 56, 61, 65, 71, 84,
		91,
	}
	deserializer := antlr.NewATNDeserializer(nil)
	staticData.atn = deserializer.Deserialize(staticData.serializedATN)
	atn := staticData.atn
	staticData.decisionToDFA = make([]*antlr.DFA, len(atn.DecisionToState))
	decisionToDFA := staticData.decisionToDFA
	for index, state := range atn.DecisionToState {
		decisionToDFA[index] = antlr.NewDFA(state, index)
	}
}

// SqlParserInit initializes any static state used to implement SqlParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewSqlParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func SqlParserInit() {
	staticData := &sqlParserStaticData
	staticData.once.Do(sqlParserInit)
}

// NewSqlParser produces a new parser instance for the optional input antlr.TokenStream.
func NewSqlParser(input antlr.TokenStream) *SqlParser {
	SqlParserInit()
	this := new(SqlParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &sqlParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "Sql.g4"

	return this
}

// SqlParser tokens.
const (
	SqlParserEOF               = antlr.TokenEOF
	SqlParserSPACE             = 1
	SqlParserWS                = 2
	SqlParserCOMMA             = 3
	SqlParserL_BRACKET         = 4
	SqlParserR_BRACKET         = 5
	SqlParserEOQ               = 6
	SqlParserK_SELECT          = 7
	SqlParserK_WHERE           = 8
	SqlParserK_WINDOW_TUMBLING = 9
	SqlParserK_GROUP_BY        = 10
	SqlParserK_AND             = 11
	SqlParserK_OR              = 12
	SqlParserK_IS              = 13
	SqlParserK_LIKE            = 14
	SqlParserK_EQUAL           = 15
	SqlParserK_GREATER         = 16
	SqlParserK_LESS            = 17
	SqlParserK_LESS_EQUAL      = 18
	SqlParserK_GREATER_EQUAL   = 19
	SqlParserK_NOT_EQUAL       = 20
	SqlParserK_NULL            = 21
	SqlParserK_IS_NULL         = 22
	SqlParserK_IS_NOT_NULL     = 23
	SqlParserK_NOT             = 24
	SqlParserK_NOT_IN          = 25
	SqlParserK_IN              = 26
	SqlParserK_TRUE            = 27
	SqlParserK_FALSE           = 28
	SqlParserK_COUNT           = 29
	SqlParserK_MIN             = 30
	SqlParserK_MAX             = 31
	SqlParserK_AVG             = 32
	SqlParserIDENTIFIER        = 33
	SqlParserBOOLEAN_LITERAL   = 34
	SqlParserNUMERIC_LITERAL   = 35
	SqlParserSTRING_LITERAL    = 36
	SqlParserSTAR              = 37
)

// SqlParser rules.
const (
	SqlParserRULE_sqlQuery           = 0
	SqlParserRULE_selectQuery        = 1
	SqlParserRULE_resultColumns      = 2
	SqlParserRULE_column             = 3
	SqlParserRULE_whereStatement     = 4
	SqlParserRULE_tumblingWindow     = 5
	SqlParserRULE_expr               = 6
	SqlParserRULE_compoundExpr       = 7
	SqlParserRULE_comparisonOperator = 8
	SqlParserRULE_literalValue       = 9
	SqlParserRULE_groupBy            = 10
	SqlParserRULE_avg                = 11
)

// ISqlQueryContext is an interface to support dynamic dispatch.
type ISqlQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSqlQueryContext differentiates from other interfaces.
	IsSqlQueryContext()
}

type SqlQueryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySqlQueryContext() *SqlQueryContext {
	var p = new(SqlQueryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_sqlQuery
	return p
}

func (*SqlQueryContext) IsSqlQueryContext() {}

func NewSqlQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SqlQueryContext {
	var p = new(SqlQueryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_sqlQuery

	return p
}

func (s *SqlQueryContext) GetParser() antlr.Parser { return s.parser }

func (s *SqlQueryContext) SelectQuery() ISelectQueryContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISelectQueryContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISelectQueryContext)
}

func (s *SqlQueryContext) EOF() antlr.TerminalNode {
	return s.GetToken(SqlParserEOF, 0)
}

func (s *SqlQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SqlQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SqlQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSqlQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) SqlQuery() (localctx ISqlQueryContext) {
	this := p
	_ = this

	localctx = NewSqlQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, SqlParserRULE_sqlQuery)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(24)
		p.SelectQuery()
	}
	{
		p.SetState(25)
		p.Match(SqlParserEOF)
	}

	return localctx
}

// ISelectQueryContext is an interface to support dynamic dispatch.
type ISelectQueryContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSelectQueryContext differentiates from other interfaces.
	IsSelectQueryContext()
}

type SelectQueryContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySelectQueryContext() *SelectQueryContext {
	var p = new(SelectQueryContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_selectQuery
	return p
}

func (*SelectQueryContext) IsSelectQueryContext() {}

func NewSelectQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectQueryContext {
	var p = new(SelectQueryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_selectQuery

	return p
}

func (s *SelectQueryContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectQueryContext) K_SELECT() antlr.TerminalNode {
	return s.GetToken(SqlParserK_SELECT, 0)
}

func (s *SelectQueryContext) EOQ() antlr.TerminalNode {
	return s.GetToken(SqlParserEOQ, 0)
}

func (s *SelectQueryContext) ResultColumns() IResultColumnsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IResultColumnsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IResultColumnsContext)
}

func (s *SelectQueryContext) WhereStatement() IWhereStatementContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWhereStatementContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWhereStatementContext)
}

func (s *SelectQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SelectQueryContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSelectQuery(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) SelectQuery() (localctx ISelectQueryContext) {
	this := p
	_ = this

	localctx = NewSelectQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, SqlParserRULE_selectQuery)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(27)
		p.Match(SqlParserK_SELECT)
	}

	{
		p.SetState(28)
		p.ResultColumns()
	}

	p.SetState(30)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(29)
			p.WhereStatement()
		}

	}
	{
		p.SetState(32)
		p.Match(SqlParserEOQ)
	}

	return localctx
}

// IResultColumnsContext is an interface to support dynamic dispatch.
type IResultColumnsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsResultColumnsContext differentiates from other interfaces.
	IsResultColumnsContext()
}

type ResultColumnsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyResultColumnsContext() *ResultColumnsContext {
	var p = new(ResultColumnsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_resultColumns
	return p
}

func (*ResultColumnsContext) IsResultColumnsContext() {}

func NewResultColumnsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ResultColumnsContext {
	var p = new(ResultColumnsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_resultColumns

	return p
}

func (s *ResultColumnsContext) GetParser() antlr.Parser { return s.parser }

func (s *ResultColumnsContext) CopyFrom(ctx *ResultColumnsContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ResultColumnsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ResultColumnsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SelectAVGContext struct {
	*ResultColumnsContext
}

func NewSelectAVGContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectAVGContext {
	var p = new(SelectAVGContext)

	p.ResultColumnsContext = NewEmptyResultColumnsContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ResultColumnsContext))

	return p
}

func (s *SelectAVGContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectAVGContext) AllAvg() []IAvgContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAvgContext); ok {
			len++
		}
	}

	tst := make([]IAvgContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAvgContext); ok {
			tst[i] = t.(IAvgContext)
			i++
		}
	}

	return tst
}

func (s *SelectAVGContext) Avg(i int) IAvgContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAvgContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAvgContext)
}

func (s *SelectAVGContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SqlParserCOMMA)
}

func (s *SelectAVGContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserCOMMA, i)
}

func (s *SelectAVGContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSelectAVG(s)

	default:
		return t.VisitChildren(s)
	}
}

type SelectStarContext struct {
	*ResultColumnsContext
}

func NewSelectStarContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectStarContext {
	var p = new(SelectStarContext)

	p.ResultColumnsContext = NewEmptyResultColumnsContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ResultColumnsContext))

	return p
}

func (s *SelectStarContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectStarContext) STAR() antlr.TerminalNode {
	return s.GetToken(SqlParserSTAR, 0)
}

func (s *SelectStarContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSelectStar(s)

	default:
		return t.VisitChildren(s)
	}
}

type SelectColumnsContext struct {
	*ResultColumnsContext
}

func NewSelectColumnsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectColumnsContext {
	var p = new(SelectColumnsContext)

	p.ResultColumnsContext = NewEmptyResultColumnsContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ResultColumnsContext))

	return p
}

func (s *SelectColumnsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectColumnsContext) AllColumn() []IColumnContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IColumnContext); ok {
			len++
		}
	}

	tst := make([]IColumnContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IColumnContext); ok {
			tst[i] = t.(IColumnContext)
			i++
		}
	}

	return tst
}

func (s *SelectColumnsContext) Column(i int) IColumnContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnContext)
}

func (s *SelectColumnsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SqlParserCOMMA)
}

func (s *SelectColumnsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserCOMMA, i)
}

func (s *SelectColumnsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSelectColumns(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) ResultColumns() (localctx IResultColumnsContext) {
	this := p
	_ = this

	localctx = NewResultColumnsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SqlParserRULE_resultColumns)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(50)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SqlParserIDENTIFIER:
		localctx = NewSelectColumnsContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(34)
			p.Column()
		}
		p.SetState(39)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == SqlParserCOMMA {
			{
				p.SetState(35)
				p.Match(SqlParserCOMMA)
			}
			{
				p.SetState(36)
				p.Column()
			}

			p.SetState(41)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case SqlParserK_COUNT, SqlParserK_MIN, SqlParserK_MAX, SqlParserK_AVG:
		localctx = NewSelectAVGContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(42)
			p.Avg()
		}
		p.SetState(45)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == SqlParserCOMMA {
			{
				p.SetState(43)
				p.Match(SqlParserCOMMA)
			}
			{
				p.SetState(44)
				p.Avg()
			}

			p.SetState(47)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case SqlParserSTAR:
		localctx = NewSelectStarContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(49)
			p.Match(SqlParserSTAR)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IColumnContext is an interface to support dynamic dispatch.
type IColumnContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsColumnContext differentiates from other interfaces.
	IsColumnContext()
}

type ColumnContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyColumnContext() *ColumnContext {
	var p = new(ColumnContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_column
	return p
}

func (*ColumnContext) IsColumnContext() {}

func NewColumnContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnContext {
	var p = new(ColumnContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_column

	return p
}

func (s *ColumnContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, 0)
}

func (s *ColumnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitColumn(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) Column() (localctx IColumnContext) {
	this := p
	_ = this

	localctx = NewColumnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, SqlParserRULE_column)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(52)
		p.Match(SqlParserIDENTIFIER)
	}

	return localctx
}

// IWhereStatementContext is an interface to support dynamic dispatch.
type IWhereStatementContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWhereStatementContext differentiates from other interfaces.
	IsWhereStatementContext()
}

type WhereStatementContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWhereStatementContext() *WhereStatementContext {
	var p = new(WhereStatementContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_whereStatement
	return p
}

func (*WhereStatementContext) IsWhereStatementContext() {}

func NewWhereStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereStatementContext {
	var p = new(WhereStatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_whereStatement

	return p
}

func (s *WhereStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *WhereStatementContext) CopyFrom(ctx *WhereStatementContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *WhereStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type WhereStmtContext struct {
	*WhereStatementContext
}

func NewWhereStmtContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *WhereStmtContext {
	var p = new(WhereStmtContext)

	p.WhereStatementContext = NewEmptyWhereStatementContext()
	p.parser = parser
	p.CopyFrom(ctx.(*WhereStatementContext))

	return p
}

func (s *WhereStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereStmtContext) K_WHERE() antlr.TerminalNode {
	return s.GetToken(SqlParserK_WHERE, 0)
}

func (s *WhereStmtContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *WhereStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitWhereStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

type TumblingStmtContext struct {
	*WhereStatementContext
}

func NewTumblingStmtContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *TumblingStmtContext {
	var p = new(TumblingStmtContext)

	p.WhereStatementContext = NewEmptyWhereStatementContext()
	p.parser = parser
	p.CopyFrom(ctx.(*WhereStatementContext))

	return p
}

func (s *TumblingStmtContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TumblingStmtContext) TumblingWindow() ITumblingWindowContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ITumblingWindowContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ITumblingWindowContext)
}

func (s *TumblingStmtContext) GroupBy() IGroupByContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupByContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupByContext)
}

func (s *TumblingStmtContext) K_WHERE() antlr.TerminalNode {
	return s.GetToken(SqlParserK_WHERE, 0)
}

func (s *TumblingStmtContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *TumblingStmtContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitTumblingStmt(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) WhereStatement() (localctx IWhereStatementContext) {
	this := p
	_ = this

	localctx = NewWhereStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, SqlParserRULE_whereStatement)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.SetState(65)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SqlParserEOQ, SqlParserK_WHERE:
		localctx = NewWhereStmtContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(56)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_WHERE {
			{
				p.SetState(54)
				p.Match(SqlParserK_WHERE)
			}
			{
				p.SetState(55)
				p.expr(0)
			}

		}

	case SqlParserK_WINDOW_TUMBLING:
		localctx = NewTumblingStmtContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(58)
			p.TumblingWindow()
		}

		p.SetState(61)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_WHERE {
			{
				p.SetState(59)
				p.Match(SqlParserK_WHERE)
			}
			{
				p.SetState(60)
				p.expr(0)
			}

		}

		{
			p.SetState(63)
			p.GroupBy()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// ITumblingWindowContext is an interface to support dynamic dispatch.
type ITumblingWindowContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsTumblingWindowContext differentiates from other interfaces.
	IsTumblingWindowContext()
}

type TumblingWindowContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyTumblingWindowContext() *TumblingWindowContext {
	var p = new(TumblingWindowContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_tumblingWindow
	return p
}

func (*TumblingWindowContext) IsTumblingWindowContext() {}

func NewTumblingWindowContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *TumblingWindowContext {
	var p = new(TumblingWindowContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_tumblingWindow

	return p
}

func (s *TumblingWindowContext) GetParser() antlr.Parser { return s.parser }

func (s *TumblingWindowContext) K_WINDOW_TUMBLING() antlr.TerminalNode {
	return s.GetToken(SqlParserK_WINDOW_TUMBLING, 0)
}

func (s *TumblingWindowContext) NUMERIC_LITERAL() antlr.TerminalNode {
	return s.GetToken(SqlParserNUMERIC_LITERAL, 0)
}

func (s *TumblingWindowContext) K_WHERE() antlr.TerminalNode {
	return s.GetToken(SqlParserK_WHERE, 0)
}

func (s *TumblingWindowContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *TumblingWindowContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *TumblingWindowContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *TumblingWindowContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitTumblingWindow(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) TumblingWindow() (localctx ITumblingWindowContext) {
	this := p
	_ = this

	localctx = NewTumblingWindowContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SqlParserRULE_tumblingWindow)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(67)
		p.Match(SqlParserK_WINDOW_TUMBLING)
	}
	{
		p.SetState(68)
		p.Match(SqlParserNUMERIC_LITERAL)
	}
	p.SetState(71)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(69)
			p.Match(SqlParserK_WHERE)
		}
		{
			p.SetState(70)
			p.expr(0)
		}

	}

	return localctx
}

// IExprContext is an interface to support dynamic dispatch.
type IExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsExprContext differentiates from other interfaces.
	IsExprContext()
}

type ExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyExprContext() *ExprContext {
	var p = new(ExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_expr
	return p
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) CopyFrom(ctx *ExprContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SimpleConditionContext struct {
	*ExprContext
}

func NewSimpleConditionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SimpleConditionContext {
	var p = new(SimpleConditionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *SimpleConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpleConditionContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, 0)
}

func (s *SimpleConditionContext) ComparisonOperator() IComparisonOperatorContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IComparisonOperatorContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IComparisonOperatorContext)
}

func (s *SimpleConditionContext) LiteralValue() ILiteralValueContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralValueContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ILiteralValueContext)
}

func (s *SimpleConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSimpleCondition(s)

	default:
		return t.VisitChildren(s)
	}
}

type NullConditionContext struct {
	*ExprContext
}

func NewNullConditionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NullConditionContext {
	var p = new(NullConditionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *NullConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NullConditionContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, 0)
}

func (s *NullConditionContext) K_IS_NULL() antlr.TerminalNode {
	return s.GetToken(SqlParserK_IS_NULL, 0)
}

func (s *NullConditionContext) K_IS_NOT_NULL() antlr.TerminalNode {
	return s.GetToken(SqlParserK_IS_NOT_NULL, 0)
}

func (s *NullConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitNullCondition(s)

	default:
		return t.VisitChildren(s)
	}
}

type CompoundRecursiveConditionContext struct {
	*ExprContext
}

func NewCompoundRecursiveConditionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CompoundRecursiveConditionContext {
	var p = new(CompoundRecursiveConditionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *CompoundRecursiveConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompoundRecursiveConditionContext) AllCompoundExpr() []ICompoundExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ICompoundExprContext); ok {
			len++
		}
	}

	tst := make([]ICompoundExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ICompoundExprContext); ok {
			tst[i] = t.(ICompoundExprContext)
			i++
		}
	}

	return tst
}

func (s *CompoundRecursiveConditionContext) CompoundExpr(i int) ICompoundExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompoundExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompoundExprContext)
}

func (s *CompoundRecursiveConditionContext) K_AND() antlr.TerminalNode {
	return s.GetToken(SqlParserK_AND, 0)
}

func (s *CompoundRecursiveConditionContext) K_OR() antlr.TerminalNode {
	return s.GetToken(SqlParserK_OR, 0)
}

func (s *CompoundRecursiveConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitCompoundRecursiveCondition(s)

	default:
		return t.VisitChildren(s)
	}
}

type SimpleRecursiveConditionContext struct {
	*ExprContext
}

func NewSimpleRecursiveConditionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SimpleRecursiveConditionContext {
	var p = new(SimpleRecursiveConditionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *SimpleRecursiveConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpleRecursiveConditionContext) AllExpr() []IExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IExprContext); ok {
			len++
		}
	}

	tst := make([]IExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IExprContext); ok {
			tst[i] = t.(IExprContext)
			i++
		}
	}

	return tst
}

func (s *SimpleRecursiveConditionContext) Expr(i int) IExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			if j == i {
				t = ctx.(antlr.RuleContext)
				break
			}
			j++
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *SimpleRecursiveConditionContext) K_AND() antlr.TerminalNode {
	return s.GetToken(SqlParserK_AND, 0)
}

func (s *SimpleRecursiveConditionContext) K_OR() antlr.TerminalNode {
	return s.GetToken(SqlParserK_OR, 0)
}

func (s *SimpleRecursiveConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSimpleRecursiveCondition(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) Expr() (localctx IExprContext) {
	return p.expr(0)
}

func (p *SqlParser) expr(_p int) (localctx IExprContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 12
	p.EnterRecursionRule(localctx, 12, SqlParserRULE_expr, _p)
	var _la int

	defer func() {
		p.UnrollRecursionContexts(_parentctx)
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	var _alt int

	p.EnterOuterAlt(localctx, 1)
	p.SetState(84)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSimpleConditionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(74)
			p.Match(SqlParserIDENTIFIER)
		}
		{
			p.SetState(75)
			p.ComparisonOperator()
		}
		{
			p.SetState(76)
			p.LiteralValue()
		}

	case 2:
		localctx = NewCompoundRecursiveConditionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(78)
			p.CompoundExpr()
		}
		{
			p.SetState(79)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SqlParserK_AND || _la == SqlParserK_OR) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(80)
			p.CompoundExpr()
		}

	case 3:
		localctx = NewNullConditionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(82)
			p.Match(SqlParserIDENTIFIER)
		}
		{
			p.SetState(83)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SqlParserK_IS_NULL || _la == SqlParserK_IS_NOT_NULL) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(91)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewSimpleRecursiveConditionContext(p, NewExprContext(p, _parentctx, _parentState))
			p.PushNewRecursionContext(localctx, _startState, SqlParserRULE_expr)
			p.SetState(86)

			if !(p.Precpred(p.GetParserRuleContext(), 1)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 1)", ""))
			}
			{
				p.SetState(87)
				_la = p.GetTokenStream().LA(1)

				if !(_la == SqlParserK_AND || _la == SqlParserK_OR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(88)
				p.expr(2)
			}

		}
		p.SetState(93)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
	}

	return localctx
}

// ICompoundExprContext is an interface to support dynamic dispatch.
type ICompoundExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsCompoundExprContext differentiates from other interfaces.
	IsCompoundExprContext()
}

type CompoundExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyCompoundExprContext() *CompoundExprContext {
	var p = new(CompoundExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_compoundExpr
	return p
}

func (*CompoundExprContext) IsCompoundExprContext() {}

func NewCompoundExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *CompoundExprContext {
	var p = new(CompoundExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_compoundExpr

	return p
}

func (s *CompoundExprContext) GetParser() antlr.Parser { return s.parser }

func (s *CompoundExprContext) L_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserL_BRACKET, 0)
}

func (s *CompoundExprContext) Expr() IExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IExprContext)
}

func (s *CompoundExprContext) R_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserR_BRACKET, 0)
}

func (s *CompoundExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompoundExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *CompoundExprContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitCompoundExpr(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) CompoundExpr() (localctx ICompoundExprContext) {
	this := p
	_ = this

	localctx = NewCompoundExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SqlParserRULE_compoundExpr)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(94)
		p.Match(SqlParserL_BRACKET)
	}
	{
		p.SetState(95)
		p.expr(0)
	}
	{
		p.SetState(96)
		p.Match(SqlParserR_BRACKET)
	}

	return localctx
}

// IComparisonOperatorContext is an interface to support dynamic dispatch.
type IComparisonOperatorContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsComparisonOperatorContext differentiates from other interfaces.
	IsComparisonOperatorContext()
}

type ComparisonOperatorContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyComparisonOperatorContext() *ComparisonOperatorContext {
	var p = new(ComparisonOperatorContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_comparisonOperator
	return p
}

func (*ComparisonOperatorContext) IsComparisonOperatorContext() {}

func NewComparisonOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonOperatorContext {
	var p = new(ComparisonOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_comparisonOperator

	return p
}

func (s *ComparisonOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonOperatorContext) K_EQUAL() antlr.TerminalNode {
	return s.GetToken(SqlParserK_EQUAL, 0)
}

func (s *ComparisonOperatorContext) K_GREATER() antlr.TerminalNode {
	return s.GetToken(SqlParserK_GREATER, 0)
}

func (s *ComparisonOperatorContext) K_LESS() antlr.TerminalNode {
	return s.GetToken(SqlParserK_LESS, 0)
}

func (s *ComparisonOperatorContext) K_LESS_EQUAL() antlr.TerminalNode {
	return s.GetToken(SqlParserK_LESS_EQUAL, 0)
}

func (s *ComparisonOperatorContext) K_GREATER_EQUAL() antlr.TerminalNode {
	return s.GetToken(SqlParserK_GREATER_EQUAL, 0)
}

func (s *ComparisonOperatorContext) K_NOT_EQUAL() antlr.TerminalNode {
	return s.GetToken(SqlParserK_NOT_EQUAL, 0)
}

func (s *ComparisonOperatorContext) K_LIKE() antlr.TerminalNode {
	return s.GetToken(SqlParserK_LIKE, 0)
}

func (s *ComparisonOperatorContext) K_IN() antlr.TerminalNode {
	return s.GetToken(SqlParserK_IN, 0)
}

func (s *ComparisonOperatorContext) K_IS() antlr.TerminalNode {
	return s.GetToken(SqlParserK_IS, 0)
}

func (s *ComparisonOperatorContext) K_NOT_IN() antlr.TerminalNode {
	return s.GetToken(SqlParserK_NOT_IN, 0)
}

func (s *ComparisonOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComparisonOperatorContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitComparisonOperator(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) ComparisonOperator() (localctx IComparisonOperatorContext) {
	this := p
	_ = this

	localctx = NewComparisonOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SqlParserRULE_comparisonOperator)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(98)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SqlParserK_IS)|(1<<SqlParserK_LIKE)|(1<<SqlParserK_EQUAL)|(1<<SqlParserK_GREATER)|(1<<SqlParserK_LESS)|(1<<SqlParserK_LESS_EQUAL)|(1<<SqlParserK_GREATER_EQUAL)|(1<<SqlParserK_NOT_EQUAL)|(1<<SqlParserK_NOT_IN)|(1<<SqlParserK_IN))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// ILiteralValueContext is an interface to support dynamic dispatch.
type ILiteralValueContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsLiteralValueContext differentiates from other interfaces.
	IsLiteralValueContext()
}

type LiteralValueContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyLiteralValueContext() *LiteralValueContext {
	var p = new(LiteralValueContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_literalValue
	return p
}

func (*LiteralValueContext) IsLiteralValueContext() {}

func NewLiteralValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralValueContext {
	var p = new(LiteralValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_literalValue

	return p
}

func (s *LiteralValueContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralValueContext) NUMERIC_LITERAL() antlr.TerminalNode {
	return s.GetToken(SqlParserNUMERIC_LITERAL, 0)
}

func (s *LiteralValueContext) STRING_LITERAL() antlr.TerminalNode {
	return s.GetToken(SqlParserSTRING_LITERAL, 0)
}

func (s *LiteralValueContext) BOOLEAN_LITERAL() antlr.TerminalNode {
	return s.GetToken(SqlParserBOOLEAN_LITERAL, 0)
}

func (s *LiteralValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LiteralValueContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitLiteralValue(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) LiteralValue() (localctx ILiteralValueContext) {
	this := p
	_ = this

	localctx = NewLiteralValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SqlParserRULE_literalValue)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(100)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-34)&-(0x1f+1)) == 0 && ((1<<uint((_la-34)))&((1<<(SqlParserBOOLEAN_LITERAL-34))|(1<<(SqlParserNUMERIC_LITERAL-34))|(1<<(SqlParserSTRING_LITERAL-34)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

// IGroupByContext is an interface to support dynamic dispatch.
type IGroupByContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGroupByContext differentiates from other interfaces.
	IsGroupByContext()
}

type GroupByContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupByContext() *GroupByContext {
	var p = new(GroupByContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_groupBy
	return p
}

func (*GroupByContext) IsGroupByContext() {}

func NewGroupByContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByContext {
	var p = new(GroupByContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_groupBy

	return p
}

func (s *GroupByContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByContext) K_GROUP_BY() antlr.TerminalNode {
	return s.GetToken(SqlParserK_GROUP_BY, 0)
}

func (s *GroupByContext) Column() IColumnContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnContext)
}

func (s *GroupByContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *GroupByContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitGroupBy(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) GroupBy() (localctx IGroupByContext) {
	this := p
	_ = this

	localctx = NewGroupByContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, SqlParserRULE_groupBy)

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(102)
		p.Match(SqlParserK_GROUP_BY)
	}
	{
		p.SetState(103)
		p.Column()
	}

	return localctx
}

// IAvgContext is an interface to support dynamic dispatch.
type IAvgContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAvgContext differentiates from other interfaces.
	IsAvgContext()
}

type AvgContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAvgContext() *AvgContext {
	var p = new(AvgContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_avg
	return p
}

func (*AvgContext) IsAvgContext() {}

func NewAvgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AvgContext {
	var p = new(AvgContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_avg

	return p
}

func (s *AvgContext) GetParser() antlr.Parser { return s.parser }

func (s *AvgContext) L_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserL_BRACKET, 0)
}

func (s *AvgContext) Column() IColumnContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IColumnContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IColumnContext)
}

func (s *AvgContext) R_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserR_BRACKET, 0)
}

func (s *AvgContext) K_MIN() antlr.TerminalNode {
	return s.GetToken(SqlParserK_MIN, 0)
}

func (s *AvgContext) K_MAX() antlr.TerminalNode {
	return s.GetToken(SqlParserK_MAX, 0)
}

func (s *AvgContext) K_COUNT() antlr.TerminalNode {
	return s.GetToken(SqlParserK_COUNT, 0)
}

func (s *AvgContext) K_AVG() antlr.TerminalNode {
	return s.GetToken(SqlParserK_AVG, 0)
}

func (s *AvgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AvgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AvgContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitAvg(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) Avg() (localctx IAvgContext) {
	this := p
	_ = this

	localctx = NewAvgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, SqlParserRULE_avg)
	var _la int

	defer func() {
		p.ExitRule()
	}()

	defer func() {
		if err := recover(); err != nil {
			if v, ok := err.(antlr.RecognitionException); ok {
				localctx.SetException(v)
				p.GetErrorHandler().ReportError(p, v)
				p.GetErrorHandler().Recover(p, v)
			} else {
				panic(err)
			}
		}
	}()

	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(105)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-29)&-(0x1f+1)) == 0 && ((1<<uint((_la-29)))&((1<<(SqlParserK_COUNT-29))|(1<<(SqlParserK_MIN-29))|(1<<(SqlParserK_MAX-29))|(1<<(SqlParserK_AVG-29)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(106)
		p.Match(SqlParserL_BRACKET)
	}
	{
		p.SetState(107)
		p.Column()
	}
	{
		p.SetState(108)
		p.Match(SqlParserR_BRACKET)
	}

	return localctx
}

func (p *SqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 6:
		var t *ExprContext = nil
		if localctx != nil {
			t = localctx.(*ExprContext)
		}
		return p.Expr_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *SqlParser) Expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 1)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
