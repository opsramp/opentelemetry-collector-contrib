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
		"", "", "','", "'('", "')'", "';'", "", "", "", "", "", "", "", "",
		"'='", "'>'", "'<'", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "'*'",
	}
	staticData.symbolicNames = []string{
		"", "SPACE", "COMMA", "L_BRACKET", "R_BRACKET", "EOQ", "K_SELECT", "K_WHERE",
		"K_WINDOW_TUMBLING", "K_GROUP_BY", "K_AND", "K_OR", "K_IS", "K_LIKE",
		"K_EQUAL", "K_GREATER", "K_LESS", "K_LESS_EQUAL", "K_GREATER_EQUAL",
		"K_NOT_EQUAL", "K_NULL", "K_IS_NULL", "K_IS_NOT_NULL", "K_NOT", "K_NOT_IN",
		"K_IN", "K_COUNT", "K_MIN", "K_MAX", "K_AVG", "IDENTIFIER", "NUMERIC_LITERAL",
		"STRING_LITERAL", "STAR",
	}
	staticData.ruleNames = []string{
		"sqlQuery", "selectQuery", "resultColumns", "column", "whereStatement",
		"tumblingWindow", "groupBy", "avg", "expr", "comparisonOperator", "literalValue",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 33, 100, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 3, 1, 29, 8, 1, 1, 1, 1, 1, 1,
		2, 1, 2, 1, 2, 4, 2, 36, 8, 2, 11, 2, 12, 2, 37, 1, 2, 1, 2, 1, 2, 4, 2,
		43, 8, 2, 11, 2, 12, 2, 44, 1, 2, 3, 2, 48, 8, 2, 1, 3, 1, 3, 1, 4, 1,
		4, 3, 4, 54, 8, 4, 1, 4, 1, 4, 1, 4, 3, 4, 59, 8, 4, 1, 4, 1, 4, 3, 4,
		63, 8, 4, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 69, 8, 5, 1, 6, 1, 6, 1, 6, 1,
		7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 1, 8, 3,
		8, 86, 8, 8, 1, 8, 1, 8, 1, 8, 5, 8, 91, 8, 8, 10, 8, 12, 8, 94, 9, 8,
		1, 9, 1, 9, 1, 10, 1, 10, 1, 10, 0, 1, 16, 11, 0, 2, 4, 6, 8, 10, 12, 14,
		16, 18, 20, 0, 5, 1, 0, 26, 29, 1, 0, 21, 22, 1, 0, 10, 11, 2, 0, 12, 19,
		24, 25, 1, 0, 31, 32, 99, 0, 22, 1, 0, 0, 0, 2, 25, 1, 0, 0, 0, 4, 47,
		1, 0, 0, 0, 6, 49, 1, 0, 0, 0, 8, 62, 1, 0, 0, 0, 10, 64, 1, 0, 0, 0, 12,
		70, 1, 0, 0, 0, 14, 73, 1, 0, 0, 0, 16, 85, 1, 0, 0, 0, 18, 95, 1, 0, 0,
		0, 20, 97, 1, 0, 0, 0, 22, 23, 3, 2, 1, 0, 23, 24, 5, 0, 0, 1, 24, 1, 1,
		0, 0, 0, 25, 26, 5, 6, 0, 0, 26, 28, 3, 4, 2, 0, 27, 29, 3, 8, 4, 0, 28,
		27, 1, 0, 0, 0, 28, 29, 1, 0, 0, 0, 29, 30, 1, 0, 0, 0, 30, 31, 5, 5, 0,
		0, 31, 3, 1, 0, 0, 0, 32, 35, 3, 6, 3, 0, 33, 34, 5, 2, 0, 0, 34, 36, 3,
		6, 3, 0, 35, 33, 1, 0, 0, 0, 36, 37, 1, 0, 0, 0, 37, 35, 1, 0, 0, 0, 37,
		38, 1, 0, 0, 0, 38, 48, 1, 0, 0, 0, 39, 42, 3, 14, 7, 0, 40, 41, 5, 2,
		0, 0, 41, 43, 3, 14, 7, 0, 42, 40, 1, 0, 0, 0, 43, 44, 1, 0, 0, 0, 44,
		42, 1, 0, 0, 0, 44, 45, 1, 0, 0, 0, 45, 48, 1, 0, 0, 0, 46, 48, 5, 33,
		0, 0, 47, 32, 1, 0, 0, 0, 47, 39, 1, 0, 0, 0, 47, 46, 1, 0, 0, 0, 48, 5,
		1, 0, 0, 0, 49, 50, 5, 30, 0, 0, 50, 7, 1, 0, 0, 0, 51, 52, 5, 7, 0, 0,
		52, 54, 3, 16, 8, 0, 53, 51, 1, 0, 0, 0, 53, 54, 1, 0, 0, 0, 54, 63, 1,
		0, 0, 0, 55, 58, 3, 10, 5, 0, 56, 57, 5, 7, 0, 0, 57, 59, 3, 16, 8, 0,
		58, 56, 1, 0, 0, 0, 58, 59, 1, 0, 0, 0, 59, 60, 1, 0, 0, 0, 60, 61, 3,
		12, 6, 0, 61, 63, 1, 0, 0, 0, 62, 53, 1, 0, 0, 0, 62, 55, 1, 0, 0, 0, 63,
		9, 1, 0, 0, 0, 64, 65, 5, 8, 0, 0, 65, 68, 5, 31, 0, 0, 66, 67, 5, 7, 0,
		0, 67, 69, 3, 16, 8, 0, 68, 66, 1, 0, 0, 0, 68, 69, 1, 0, 0, 0, 69, 11,
		1, 0, 0, 0, 70, 71, 5, 9, 0, 0, 71, 72, 3, 6, 3, 0, 72, 13, 1, 0, 0, 0,
		73, 74, 7, 0, 0, 0, 74, 75, 5, 3, 0, 0, 75, 76, 3, 6, 3, 0, 76, 77, 5,
		4, 0, 0, 77, 15, 1, 0, 0, 0, 78, 79, 6, 8, -1, 0, 79, 80, 5, 30, 0, 0,
		80, 81, 3, 18, 9, 0, 81, 82, 3, 20, 10, 0, 82, 86, 1, 0, 0, 0, 83, 84,
		5, 30, 0, 0, 84, 86, 7, 1, 0, 0, 85, 78, 1, 0, 0, 0, 85, 83, 1, 0, 0, 0,
		86, 92, 1, 0, 0, 0, 87, 88, 10, 2, 0, 0, 88, 89, 7, 2, 0, 0, 89, 91, 3,
		16, 8, 3, 90, 87, 1, 0, 0, 0, 91, 94, 1, 0, 0, 0, 92, 90, 1, 0, 0, 0, 92,
		93, 1, 0, 0, 0, 93, 17, 1, 0, 0, 0, 94, 92, 1, 0, 0, 0, 95, 96, 7, 3, 0,
		0, 96, 19, 1, 0, 0, 0, 97, 98, 7, 4, 0, 0, 98, 21, 1, 0, 0, 0, 10, 28,
		37, 44, 47, 53, 58, 62, 68, 85, 92,
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
	SqlParserCOMMA             = 2
	SqlParserL_BRACKET         = 3
	SqlParserR_BRACKET         = 4
	SqlParserEOQ               = 5
	SqlParserK_SELECT          = 6
	SqlParserK_WHERE           = 7
	SqlParserK_WINDOW_TUMBLING = 8
	SqlParserK_GROUP_BY        = 9
	SqlParserK_AND             = 10
	SqlParserK_OR              = 11
	SqlParserK_IS              = 12
	SqlParserK_LIKE            = 13
	SqlParserK_EQUAL           = 14
	SqlParserK_GREATER         = 15
	SqlParserK_LESS            = 16
	SqlParserK_LESS_EQUAL      = 17
	SqlParserK_GREATER_EQUAL   = 18
	SqlParserK_NOT_EQUAL       = 19
	SqlParserK_NULL            = 20
	SqlParserK_IS_NULL         = 21
	SqlParserK_IS_NOT_NULL     = 22
	SqlParserK_NOT             = 23
	SqlParserK_NOT_IN          = 24
	SqlParserK_IN              = 25
	SqlParserK_COUNT           = 26
	SqlParserK_MIN             = 27
	SqlParserK_MAX             = 28
	SqlParserK_AVG             = 29
	SqlParserIDENTIFIER        = 30
	SqlParserNUMERIC_LITERAL   = 31
	SqlParserSTRING_LITERAL    = 32
	SqlParserSTAR              = 33
)

// SqlParser rules.
const (
	SqlParserRULE_sqlQuery           = 0
	SqlParserRULE_selectQuery        = 1
	SqlParserRULE_resultColumns      = 2
	SqlParserRULE_column             = 3
	SqlParserRULE_whereStatement     = 4
	SqlParserRULE_tumblingWindow     = 5
	SqlParserRULE_groupBy            = 6
	SqlParserRULE_avg                = 7
	SqlParserRULE_expr               = 8
	SqlParserRULE_comparisonOperator = 9
	SqlParserRULE_literalValue       = 10
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

func (s *SqlQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.EnterSqlQuery(s)
	}
}

func (s *SqlQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.ExitSqlQuery(s)
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
		p.SetState(22)
		p.SelectQuery()
	}
	{
		p.SetState(23)
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

func (s *SelectQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.EnterSelectQuery(s)
	}
}

func (s *SelectQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.ExitSelectQuery(s)
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
		p.SetState(25)
		p.Match(SqlParserK_SELECT)
	}

	{
		p.SetState(26)
		p.ResultColumns()
	}

	p.SetState(28)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(27)
			p.WhereStatement()
		}

	}
	{
		p.SetState(30)
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

func (s *ResultColumnsContext) AllColumn() []IColumnContext {
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

func (s *ResultColumnsContext) Column(i int) IColumnContext {
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

func (s *ResultColumnsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SqlParserCOMMA)
}

func (s *ResultColumnsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserCOMMA, i)
}

func (s *ResultColumnsContext) AllAvg() []IAvgContext {
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

func (s *ResultColumnsContext) Avg(i int) IAvgContext {
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

func (s *ResultColumnsContext) STAR() antlr.TerminalNode {
	return s.GetToken(SqlParserSTAR, 0)
}

func (s *ResultColumnsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ResultColumnsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ResultColumnsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.EnterResultColumns(s)
	}
}

func (s *ResultColumnsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.ExitResultColumns(s)
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

	p.SetState(47)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SqlParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(32)
			p.Column()
		}
		p.SetState(35)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == SqlParserCOMMA {
			{
				p.SetState(33)
				p.Match(SqlParserCOMMA)
			}
			{
				p.SetState(34)
				p.Column()
			}

			p.SetState(37)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case SqlParserK_COUNT, SqlParserK_MIN, SqlParserK_MAX, SqlParserK_AVG:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(39)
			p.Avg()
		}
		p.SetState(42)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == SqlParserCOMMA {
			{
				p.SetState(40)
				p.Match(SqlParserCOMMA)
			}
			{
				p.SetState(41)
				p.Avg()
			}

			p.SetState(44)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case SqlParserSTAR:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(46)
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

func (s *ColumnContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.EnterColumn(s)
	}
}

func (s *ColumnContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.ExitColumn(s)
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
		p.SetState(49)
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

func (s *WhereStatementContext) K_WHERE() antlr.TerminalNode {
	return s.GetToken(SqlParserK_WHERE, 0)
}

func (s *WhereStatementContext) Expr() IExprContext {
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

func (s *WhereStatementContext) TumblingWindow() ITumblingWindowContext {
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

func (s *WhereStatementContext) GroupBy() IGroupByContext {
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

func (s *WhereStatementContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WhereStatementContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WhereStatementContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.EnterWhereStatement(s)
	}
}

func (s *WhereStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.ExitWhereStatement(s)
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

	p.SetState(62)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SqlParserEOQ, SqlParserK_WHERE:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(53)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_WHERE {
			{
				p.SetState(51)
				p.Match(SqlParserK_WHERE)
			}
			{
				p.SetState(52)
				p.expr(0)
			}

		}

	case SqlParserK_WINDOW_TUMBLING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(55)
			p.TumblingWindow()
		}

		p.SetState(58)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_WHERE {
			{
				p.SetState(56)
				p.Match(SqlParserK_WHERE)
			}
			{
				p.SetState(57)
				p.expr(0)
			}

		}

		{
			p.SetState(60)
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

func (s *TumblingWindowContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.EnterTumblingWindow(s)
	}
}

func (s *TumblingWindowContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.ExitTumblingWindow(s)
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
		p.SetState(64)
		p.Match(SqlParserK_WINDOW_TUMBLING)
	}
	{
		p.SetState(65)
		p.Match(SqlParserNUMERIC_LITERAL)
	}
	p.SetState(68)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(66)
			p.Match(SqlParserK_WHERE)
		}
		{
			p.SetState(67)
			p.expr(0)
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

func (s *GroupByContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.EnterGroupBy(s)
	}
}

func (s *GroupByContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.ExitGroupBy(s)
	}
}

func (p *SqlParser) GroupBy() (localctx IGroupByContext) {
	this := p
	_ = this

	localctx = NewGroupByContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SqlParserRULE_groupBy)

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
		p.SetState(70)
		p.Match(SqlParserK_GROUP_BY)
	}
	{
		p.SetState(71)
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

func (s *AvgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.EnterAvg(s)
	}
}

func (s *AvgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.ExitAvg(s)
	}
}

func (p *SqlParser) Avg() (localctx IAvgContext) {
	this := p
	_ = this

	localctx = NewAvgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SqlParserRULE_avg)
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
		p.SetState(73)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SqlParserK_COUNT)|(1<<SqlParserK_MIN)|(1<<SqlParserK_MAX)|(1<<SqlParserK_AVG))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(74)
		p.Match(SqlParserL_BRACKET)
	}
	{
		p.SetState(75)
		p.Column()
	}
	{
		p.SetState(76)
		p.Match(SqlParserR_BRACKET)
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

func (s *ExprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, 0)
}

func (s *ExprContext) ComparisonOperator() IComparisonOperatorContext {
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

func (s *ExprContext) LiteralValue() ILiteralValueContext {
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

func (s *ExprContext) K_IS_NULL() antlr.TerminalNode {
	return s.GetToken(SqlParserK_IS_NULL, 0)
}

func (s *ExprContext) K_IS_NOT_NULL() antlr.TerminalNode {
	return s.GetToken(SqlParserK_IS_NOT_NULL, 0)
}

func (s *ExprContext) AllExpr() []IExprContext {
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

func (s *ExprContext) Expr(i int) IExprContext {
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

func (s *ExprContext) K_AND() antlr.TerminalNode {
	return s.GetToken(SqlParserK_AND, 0)
}

func (s *ExprContext) K_OR() antlr.TerminalNode {
	return s.GetToken(SqlParserK_OR, 0)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.EnterExpr(s)
	}
}

func (s *ExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.ExitExpr(s)
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
	_startState := 16
	p.EnterRecursionRule(localctx, 16, SqlParserRULE_expr, _p)
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
	p.SetState(85)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(79)
			p.Match(SqlParserIDENTIFIER)
		}
		{
			p.SetState(80)
			p.ComparisonOperator()
		}
		{
			p.SetState(81)
			p.LiteralValue()
		}

	case 2:
		{
			p.SetState(83)
			p.Match(SqlParserIDENTIFIER)
		}
		{
			p.SetState(84)
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
	p.SetState(92)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, SqlParserRULE_expr)
			p.SetState(87)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(88)
				_la = p.GetTokenStream().LA(1)

				if !(_la == SqlParserK_AND || _la == SqlParserK_OR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(89)
				p.expr(3)
			}

		}
		p.SetState(94)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 9, p.GetParserRuleContext())
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

func (s *ComparisonOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.EnterComparisonOperator(s)
	}
}

func (s *ComparisonOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.ExitComparisonOperator(s)
	}
}

func (p *SqlParser) ComparisonOperator() (localctx IComparisonOperatorContext) {
	this := p
	_ = this

	localctx = NewComparisonOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SqlParserRULE_comparisonOperator)
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
		p.SetState(95)
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

func (s *LiteralValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LiteralValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.EnterLiteralValue(s)
	}
}

func (s *LiteralValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(SqlListener); ok {
		listenerT.ExitLiteralValue(s)
	}
}

func (p *SqlParser) LiteralValue() (localctx ILiteralValueContext) {
	this := p
	_ = this

	localctx = NewLiteralValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, SqlParserRULE_literalValue)
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
		p.SetState(97)
		_la = p.GetTokenStream().LA(1)

		if !(_la == SqlParserNUMERIC_LITERAL || _la == SqlParserSTRING_LITERAL) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

func (p *SqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 8:
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
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
