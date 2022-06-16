// Code generated from sql.g4 by ANTLR 4.10.1. DO NOT EDIT.

package parser // sql

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

type sqlParser struct {
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
		"groupBy", "avg", "expr", "comparisonOperator", "literalValue",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 33, 92, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 1, 0, 1,
		0, 1, 0, 1, 1, 1, 1, 1, 1, 3, 1, 27, 8, 1, 1, 1, 1, 1, 1, 2, 1, 2, 1, 2,
		4, 2, 34, 8, 2, 11, 2, 12, 2, 35, 1, 2, 1, 2, 1, 2, 4, 2, 41, 8, 2, 11,
		2, 12, 2, 42, 1, 2, 3, 2, 46, 8, 2, 1, 3, 1, 3, 1, 4, 1, 4, 3, 4, 52, 8,
		4, 1, 4, 1, 4, 1, 4, 1, 4, 3, 4, 58, 8, 4, 1, 4, 3, 4, 61, 8, 4, 1, 5,
		1, 5, 1, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7,
		1, 7, 1, 7, 3, 7, 78, 8, 7, 1, 7, 1, 7, 1, 7, 5, 7, 83, 8, 7, 10, 7, 12,
		7, 86, 9, 7, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 0, 1, 14, 10, 0, 2, 4, 6, 8,
		10, 12, 14, 16, 18, 0, 5, 1, 0, 26, 29, 1, 0, 21, 22, 1, 0, 10, 11, 2,
		0, 12, 19, 24, 25, 1, 0, 31, 32, 91, 0, 20, 1, 0, 0, 0, 2, 23, 1, 0, 0,
		0, 4, 45, 1, 0, 0, 0, 6, 47, 1, 0, 0, 0, 8, 60, 1, 0, 0, 0, 10, 62, 1,
		0, 0, 0, 12, 65, 1, 0, 0, 0, 14, 77, 1, 0, 0, 0, 16, 87, 1, 0, 0, 0, 18,
		89, 1, 0, 0, 0, 20, 21, 3, 2, 1, 0, 21, 22, 5, 0, 0, 1, 22, 1, 1, 0, 0,
		0, 23, 24, 5, 6, 0, 0, 24, 26, 3, 4, 2, 0, 25, 27, 3, 8, 4, 0, 26, 25,
		1, 0, 0, 0, 26, 27, 1, 0, 0, 0, 27, 28, 1, 0, 0, 0, 28, 29, 5, 5, 0, 0,
		29, 3, 1, 0, 0, 0, 30, 33, 3, 6, 3, 0, 31, 32, 5, 2, 0, 0, 32, 34, 3, 6,
		3, 0, 33, 31, 1, 0, 0, 0, 34, 35, 1, 0, 0, 0, 35, 33, 1, 0, 0, 0, 35, 36,
		1, 0, 0, 0, 36, 46, 1, 0, 0, 0, 37, 40, 3, 12, 6, 0, 38, 39, 5, 2, 0, 0,
		39, 41, 3, 12, 6, 0, 40, 38, 1, 0, 0, 0, 41, 42, 1, 0, 0, 0, 42, 40, 1,
		0, 0, 0, 42, 43, 1, 0, 0, 0, 43, 46, 1, 0, 0, 0, 44, 46, 5, 33, 0, 0, 45,
		30, 1, 0, 0, 0, 45, 37, 1, 0, 0, 0, 45, 44, 1, 0, 0, 0, 46, 5, 1, 0, 0,
		0, 47, 48, 5, 30, 0, 0, 48, 7, 1, 0, 0, 0, 49, 50, 5, 7, 0, 0, 50, 52,
		3, 14, 7, 0, 51, 49, 1, 0, 0, 0, 51, 52, 1, 0, 0, 0, 52, 61, 1, 0, 0, 0,
		53, 54, 5, 8, 0, 0, 54, 57, 5, 31, 0, 0, 55, 56, 5, 7, 0, 0, 56, 58, 3,
		14, 7, 0, 57, 55, 1, 0, 0, 0, 57, 58, 1, 0, 0, 0, 58, 59, 1, 0, 0, 0, 59,
		61, 3, 10, 5, 0, 60, 51, 1, 0, 0, 0, 60, 53, 1, 0, 0, 0, 61, 9, 1, 0, 0,
		0, 62, 63, 5, 9, 0, 0, 63, 64, 3, 6, 3, 0, 64, 11, 1, 0, 0, 0, 65, 66,
		7, 0, 0, 0, 66, 67, 5, 3, 0, 0, 67, 68, 3, 6, 3, 0, 68, 69, 5, 4, 0, 0,
		69, 13, 1, 0, 0, 0, 70, 71, 6, 7, -1, 0, 71, 72, 5, 30, 0, 0, 72, 73, 3,
		16, 8, 0, 73, 74, 3, 18, 9, 0, 74, 78, 1, 0, 0, 0, 75, 76, 5, 30, 0, 0,
		76, 78, 7, 1, 0, 0, 77, 70, 1, 0, 0, 0, 77, 75, 1, 0, 0, 0, 78, 84, 1,
		0, 0, 0, 79, 80, 10, 2, 0, 0, 80, 81, 7, 2, 0, 0, 81, 83, 3, 14, 7, 3,
		82, 79, 1, 0, 0, 0, 83, 86, 1, 0, 0, 0, 84, 82, 1, 0, 0, 0, 84, 85, 1,
		0, 0, 0, 85, 15, 1, 0, 0, 0, 86, 84, 1, 0, 0, 0, 87, 88, 7, 3, 0, 0, 88,
		17, 1, 0, 0, 0, 89, 90, 7, 4, 0, 0, 90, 19, 1, 0, 0, 0, 9, 26, 35, 42,
		45, 51, 57, 60, 77, 84,
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

// sqlParserInit initializes any static state used to implement sqlParser. By default the
// static state used to implement the parser is lazily initialized during the first call to
// NewsqlParser(). You can call this function if you wish to initialize the static state ahead
// of time.
func SqlParserInit() {
	staticData := &sqlParserStaticData
	staticData.once.Do(sqlParserInit)
}

// NewsqlParser produces a new parser instance for the optional input antlr.TokenStream.
func NewsqlParser(input antlr.TokenStream) *sqlParser {
	SqlParserInit()
	this := new(sqlParser)
	this.BaseParser = antlr.NewBaseParser(input)
	staticData := &sqlParserStaticData
	this.Interpreter = antlr.NewParserATNSimulator(this, staticData.atn, staticData.decisionToDFA, staticData.predictionContextCache)
	this.RuleNames = staticData.ruleNames
	this.LiteralNames = staticData.literalNames
	this.SymbolicNames = staticData.symbolicNames
	this.GrammarFileName = "sql.g4"

	return this
}

// sqlParser tokens.
const (
	sqlParserEOF               = antlr.TokenEOF
	sqlParserSPACE             = 1
	sqlParserCOMMA             = 2
	sqlParserL_BRACKET         = 3
	sqlParserR_BRACKET         = 4
	sqlParserEOQ               = 5
	sqlParserK_SELECT          = 6
	sqlParserK_WHERE           = 7
	sqlParserK_WINDOW_TUMBLING = 8
	sqlParserK_GROUP_BY        = 9
	sqlParserK_AND             = 10
	sqlParserK_OR              = 11
	sqlParserK_IS              = 12
	sqlParserK_LIKE            = 13
	sqlParserK_EQUAL           = 14
	sqlParserK_GREATER         = 15
	sqlParserK_LESS            = 16
	sqlParserK_LESS_EQUAL      = 17
	sqlParserK_GREATER_EQUAL   = 18
	sqlParserK_NOT_EQUAL       = 19
	sqlParserK_NULL            = 20
	sqlParserK_IS_NULL         = 21
	sqlParserK_IS_NOT_NULL     = 22
	sqlParserK_NOT             = 23
	sqlParserK_NOT_IN          = 24
	sqlParserK_IN              = 25
	sqlParserK_COUNT           = 26
	sqlParserK_MIN             = 27
	sqlParserK_MAX             = 28
	sqlParserK_AVG             = 29
	sqlParserIDENTIFIER        = 30
	sqlParserNUMERIC_LITERAL   = 31
	sqlParserSTRING_LITERAL    = 32
	sqlParserSTAR              = 33
)

// sqlParser rules.
const (
	sqlParserRULE_sqlQuery           = 0
	sqlParserRULE_selectQuery        = 1
	sqlParserRULE_resultColumns      = 2
	sqlParserRULE_column             = 3
	sqlParserRULE_whereStatement     = 4
	sqlParserRULE_groupBy            = 5
	sqlParserRULE_avg                = 6
	sqlParserRULE_expr               = 7
	sqlParserRULE_comparisonOperator = 8
	sqlParserRULE_literalValue       = 9
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
	p.RuleIndex = sqlParserRULE_sqlQuery
	return p
}

func (*SqlQueryContext) IsSqlQueryContext() {}

func NewSqlQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SqlQueryContext {
	var p = new(SqlQueryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_sqlQuery

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
	return s.GetToken(sqlParserEOF, 0)
}

func (s *SqlQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SqlQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *SqlQueryContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterSqlQuery(s)
	}
}

func (s *SqlQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitSqlQuery(s)
	}
}

func (p *sqlParser) SqlQuery() (localctx ISqlQueryContext) {
	this := p
	_ = this

	localctx = NewSqlQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 0, sqlParserRULE_sqlQuery)

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
		p.SetState(20)
		p.SelectQuery()
	}
	{
		p.SetState(21)
		p.Match(sqlParserEOF)
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
	p.RuleIndex = sqlParserRULE_selectQuery
	return p
}

func (*SelectQueryContext) IsSelectQueryContext() {}

func NewSelectQueryContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SelectQueryContext {
	var p = new(SelectQueryContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_selectQuery

	return p
}

func (s *SelectQueryContext) GetParser() antlr.Parser { return s.parser }

func (s *SelectQueryContext) K_SELECT() antlr.TerminalNode {
	return s.GetToken(sqlParserK_SELECT, 0)
}

func (s *SelectQueryContext) EOQ() antlr.TerminalNode {
	return s.GetToken(sqlParserEOQ, 0)
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
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterSelectQuery(s)
	}
}

func (s *SelectQueryContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitSelectQuery(s)
	}
}

func (p *sqlParser) SelectQuery() (localctx ISelectQueryContext) {
	this := p
	_ = this

	localctx = NewSelectQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, sqlParserRULE_selectQuery)

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
		p.SetState(23)
		p.Match(sqlParserK_SELECT)
	}

	{
		p.SetState(24)
		p.ResultColumns()
	}

	p.SetState(26)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 0, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(25)
			p.WhereStatement()
		}

	}
	{
		p.SetState(28)
		p.Match(sqlParserEOQ)
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
	p.RuleIndex = sqlParserRULE_resultColumns
	return p
}

func (*ResultColumnsContext) IsResultColumnsContext() {}

func NewResultColumnsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ResultColumnsContext {
	var p = new(ResultColumnsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_resultColumns

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
	return s.GetTokens(sqlParserCOMMA)
}

func (s *ResultColumnsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(sqlParserCOMMA, i)
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
	return s.GetToken(sqlParserSTAR, 0)
}

func (s *ResultColumnsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ResultColumnsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ResultColumnsContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterResultColumns(s)
	}
}

func (s *ResultColumnsContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitResultColumns(s)
	}
}

func (p *sqlParser) ResultColumns() (localctx IResultColumnsContext) {
	this := p
	_ = this

	localctx = NewResultColumnsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, sqlParserRULE_resultColumns)
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

	p.SetState(45)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case sqlParserIDENTIFIER:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(30)
			p.Column()
		}
		p.SetState(33)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == sqlParserCOMMA {
			{
				p.SetState(31)
				p.Match(sqlParserCOMMA)
			}
			{
				p.SetState(32)
				p.Column()
			}

			p.SetState(35)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case sqlParserK_COUNT, sqlParserK_MIN, sqlParserK_MAX, sqlParserK_AVG:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(37)
			p.Avg()
		}
		p.SetState(40)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for ok := true; ok; ok = _la == sqlParserCOMMA {
			{
				p.SetState(38)
				p.Match(sqlParserCOMMA)
			}
			{
				p.SetState(39)
				p.Avg()
			}

			p.SetState(42)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case sqlParserSTAR:
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(44)
			p.Match(sqlParserSTAR)
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
	p.RuleIndex = sqlParserRULE_column
	return p
}

func (*ColumnContext) IsColumnContext() {}

func NewColumnContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ColumnContext {
	var p = new(ColumnContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_column

	return p
}

func (s *ColumnContext) GetParser() antlr.Parser { return s.parser }

func (s *ColumnContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(sqlParserIDENTIFIER, 0)
}

func (s *ColumnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ColumnContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterColumn(s)
	}
}

func (s *ColumnContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitColumn(s)
	}
}

func (p *sqlParser) Column() (localctx IColumnContext) {
	this := p
	_ = this

	localctx = NewColumnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 6, sqlParserRULE_column)

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
		p.SetState(47)
		p.Match(sqlParserIDENTIFIER)
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
	p.RuleIndex = sqlParserRULE_whereStatement
	return p
}

func (*WhereStatementContext) IsWhereStatementContext() {}

func NewWhereStatementContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WhereStatementContext {
	var p = new(WhereStatementContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_whereStatement

	return p
}

func (s *WhereStatementContext) GetParser() antlr.Parser { return s.parser }

func (s *WhereStatementContext) K_WHERE() antlr.TerminalNode {
	return s.GetToken(sqlParserK_WHERE, 0)
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

func (s *WhereStatementContext) K_WINDOW_TUMBLING() antlr.TerminalNode {
	return s.GetToken(sqlParserK_WINDOW_TUMBLING, 0)
}

func (s *WhereStatementContext) NUMERIC_LITERAL() antlr.TerminalNode {
	return s.GetToken(sqlParserNUMERIC_LITERAL, 0)
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
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterWhereStatement(s)
	}
}

func (s *WhereStatementContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitWhereStatement(s)
	}
}

func (p *sqlParser) WhereStatement() (localctx IWhereStatementContext) {
	this := p
	_ = this

	localctx = NewWhereStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, sqlParserRULE_whereStatement)
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

	p.SetState(60)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case sqlParserEOQ, sqlParserK_WHERE:
		p.EnterOuterAlt(localctx, 1)
		p.SetState(51)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserK_WHERE {
			{
				p.SetState(49)
				p.Match(sqlParserK_WHERE)
			}
			{
				p.SetState(50)
				p.expr(0)
			}

		}

	case sqlParserK_WINDOW_TUMBLING:
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(53)
			p.Match(sqlParserK_WINDOW_TUMBLING)
		}
		{
			p.SetState(54)
			p.Match(sqlParserNUMERIC_LITERAL)
		}
		p.SetState(57)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == sqlParserK_WHERE {
			{
				p.SetState(55)
				p.Match(sqlParserK_WHERE)
			}
			{
				p.SetState(56)
				p.expr(0)
			}

		}

		{
			p.SetState(59)
			p.GroupBy()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
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
	p.RuleIndex = sqlParserRULE_groupBy
	return p
}

func (*GroupByContext) IsGroupByContext() {}

func NewGroupByContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByContext {
	var p = new(GroupByContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_groupBy

	return p
}

func (s *GroupByContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByContext) K_GROUP_BY() antlr.TerminalNode {
	return s.GetToken(sqlParserK_GROUP_BY, 0)
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
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterGroupBy(s)
	}
}

func (s *GroupByContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitGroupBy(s)
	}
}

func (p *sqlParser) GroupBy() (localctx IGroupByContext) {
	this := p
	_ = this

	localctx = NewGroupByContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, sqlParserRULE_groupBy)

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
		p.SetState(62)
		p.Match(sqlParserK_GROUP_BY)
	}
	{
		p.SetState(63)
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
	p.RuleIndex = sqlParserRULE_avg
	return p
}

func (*AvgContext) IsAvgContext() {}

func NewAvgContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AvgContext {
	var p = new(AvgContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_avg

	return p
}

func (s *AvgContext) GetParser() antlr.Parser { return s.parser }

func (s *AvgContext) L_BRACKET() antlr.TerminalNode {
	return s.GetToken(sqlParserL_BRACKET, 0)
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
	return s.GetToken(sqlParserR_BRACKET, 0)
}

func (s *AvgContext) K_MIN() antlr.TerminalNode {
	return s.GetToken(sqlParserK_MIN, 0)
}

func (s *AvgContext) K_MAX() antlr.TerminalNode {
	return s.GetToken(sqlParserK_MAX, 0)
}

func (s *AvgContext) K_COUNT() antlr.TerminalNode {
	return s.GetToken(sqlParserK_COUNT, 0)
}

func (s *AvgContext) K_AVG() antlr.TerminalNode {
	return s.GetToken(sqlParserK_AVG, 0)
}

func (s *AvgContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AvgContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AvgContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterAvg(s)
	}
}

func (s *AvgContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitAvg(s)
	}
}

func (p *sqlParser) Avg() (localctx IAvgContext) {
	this := p
	_ = this

	localctx = NewAvgContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, sqlParserRULE_avg)
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
		p.SetState(65)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<sqlParserK_COUNT)|(1<<sqlParserK_MIN)|(1<<sqlParserK_MAX)|(1<<sqlParserK_AVG))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(66)
		p.Match(sqlParserL_BRACKET)
	}
	{
		p.SetState(67)
		p.Column()
	}
	{
		p.SetState(68)
		p.Match(sqlParserR_BRACKET)
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
	p.RuleIndex = sqlParserRULE_expr
	return p
}

func (*ExprContext) IsExprContext() {}

func NewExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ExprContext {
	var p = new(ExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_expr

	return p
}

func (s *ExprContext) GetParser() antlr.Parser { return s.parser }

func (s *ExprContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(sqlParserIDENTIFIER, 0)
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
	return s.GetToken(sqlParserK_IS_NULL, 0)
}

func (s *ExprContext) K_IS_NOT_NULL() antlr.TerminalNode {
	return s.GetToken(sqlParserK_IS_NOT_NULL, 0)
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
	return s.GetToken(sqlParserK_AND, 0)
}

func (s *ExprContext) K_OR() antlr.TerminalNode {
	return s.GetToken(sqlParserK_OR, 0)
}

func (s *ExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ExprContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterExpr(s)
	}
}

func (s *ExprContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitExpr(s)
	}
}

func (p *sqlParser) Expr() (localctx IExprContext) {
	return p.expr(0)
}

func (p *sqlParser) expr(_p int) (localctx IExprContext) {
	this := p
	_ = this

	var _parentctx antlr.ParserRuleContext = p.GetParserRuleContext()
	_parentState := p.GetState()
	localctx = NewExprContext(p, p.GetParserRuleContext(), _parentState)
	var _prevctx IExprContext = localctx
	var _ antlr.ParserRuleContext = _prevctx // TODO: To prevent unused variable warning.
	_startState := 14
	p.EnterRecursionRule(localctx, 14, sqlParserRULE_expr, _p)
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
	p.SetState(77)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 7, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(71)
			p.Match(sqlParserIDENTIFIER)
		}
		{
			p.SetState(72)
			p.ComparisonOperator()
		}
		{
			p.SetState(73)
			p.LiteralValue()
		}

	case 2:
		{
			p.SetState(75)
			p.Match(sqlParserIDENTIFIER)
		}
		{
			p.SetState(76)
			_la = p.GetTokenStream().LA(1)

			if !(_la == sqlParserK_IS_NULL || _la == sqlParserK_IS_NOT_NULL) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(84)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewExprContext(p, _parentctx, _parentState)
			p.PushNewRecursionContext(localctx, _startState, sqlParserRULE_expr)
			p.SetState(79)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(80)
				_la = p.GetTokenStream().LA(1)

				if !(_la == sqlParserK_AND || _la == sqlParserK_OR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(81)
				p.expr(3)
			}

		}
		p.SetState(86)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext())
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
	p.RuleIndex = sqlParserRULE_comparisonOperator
	return p
}

func (*ComparisonOperatorContext) IsComparisonOperatorContext() {}

func NewComparisonOperatorContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *ComparisonOperatorContext {
	var p = new(ComparisonOperatorContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_comparisonOperator

	return p
}

func (s *ComparisonOperatorContext) GetParser() antlr.Parser { return s.parser }

func (s *ComparisonOperatorContext) K_EQUAL() antlr.TerminalNode {
	return s.GetToken(sqlParserK_EQUAL, 0)
}

func (s *ComparisonOperatorContext) K_GREATER() antlr.TerminalNode {
	return s.GetToken(sqlParserK_GREATER, 0)
}

func (s *ComparisonOperatorContext) K_LESS() antlr.TerminalNode {
	return s.GetToken(sqlParserK_LESS, 0)
}

func (s *ComparisonOperatorContext) K_LESS_EQUAL() antlr.TerminalNode {
	return s.GetToken(sqlParserK_LESS_EQUAL, 0)
}

func (s *ComparisonOperatorContext) K_GREATER_EQUAL() antlr.TerminalNode {
	return s.GetToken(sqlParserK_GREATER_EQUAL, 0)
}

func (s *ComparisonOperatorContext) K_NOT_EQUAL() antlr.TerminalNode {
	return s.GetToken(sqlParserK_NOT_EQUAL, 0)
}

func (s *ComparisonOperatorContext) K_LIKE() antlr.TerminalNode {
	return s.GetToken(sqlParserK_LIKE, 0)
}

func (s *ComparisonOperatorContext) K_IN() antlr.TerminalNode {
	return s.GetToken(sqlParserK_IN, 0)
}

func (s *ComparisonOperatorContext) K_IS() antlr.TerminalNode {
	return s.GetToken(sqlParserK_IS, 0)
}

func (s *ComparisonOperatorContext) K_NOT_IN() antlr.TerminalNode {
	return s.GetToken(sqlParserK_NOT_IN, 0)
}

func (s *ComparisonOperatorContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ComparisonOperatorContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *ComparisonOperatorContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterComparisonOperator(s)
	}
}

func (s *ComparisonOperatorContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitComparisonOperator(s)
	}
}

func (p *sqlParser) ComparisonOperator() (localctx IComparisonOperatorContext) {
	this := p
	_ = this

	localctx = NewComparisonOperatorContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, sqlParserRULE_comparisonOperator)
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
		p.SetState(87)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<sqlParserK_IS)|(1<<sqlParserK_LIKE)|(1<<sqlParserK_EQUAL)|(1<<sqlParserK_GREATER)|(1<<sqlParserK_LESS)|(1<<sqlParserK_LESS_EQUAL)|(1<<sqlParserK_GREATER_EQUAL)|(1<<sqlParserK_NOT_EQUAL)|(1<<sqlParserK_NOT_IN)|(1<<sqlParserK_IN))) != 0) {
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
	p.RuleIndex = sqlParserRULE_literalValue
	return p
}

func (*LiteralValueContext) IsLiteralValueContext() {}

func NewLiteralValueContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *LiteralValueContext {
	var p = new(LiteralValueContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = sqlParserRULE_literalValue

	return p
}

func (s *LiteralValueContext) GetParser() antlr.Parser { return s.parser }

func (s *LiteralValueContext) NUMERIC_LITERAL() antlr.TerminalNode {
	return s.GetToken(sqlParserNUMERIC_LITERAL, 0)
}

func (s *LiteralValueContext) STRING_LITERAL() antlr.TerminalNode {
	return s.GetToken(sqlParserSTRING_LITERAL, 0)
}

func (s *LiteralValueContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *LiteralValueContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *LiteralValueContext) EnterRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.EnterLiteralValue(s)
	}
}

func (s *LiteralValueContext) ExitRule(listener antlr.ParseTreeListener) {
	if listenerT, ok := listener.(sqlListener); ok {
		listenerT.ExitLiteralValue(s)
	}
}

func (p *sqlParser) LiteralValue() (localctx ILiteralValueContext) {
	this := p
	_ = this

	localctx = NewLiteralValueContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, sqlParserRULE_literalValue)
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
		p.SetState(89)
		_la = p.GetTokenStream().LA(1)

		if !(_la == sqlParserNUMERIC_LITERAL || _la == sqlParserSTRING_LITERAL) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}

	return localctx
}

func (p *sqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 7:
		var t *ExprContext = nil
		if localctx != nil {
			t = localctx.(*ExprContext)
		}
		return p.Expr_Sempred(t, predIndex)

	default:
		panic("No predicate with index: " + fmt.Sprint(ruleIndex))
	}
}

func (p *sqlParser) Expr_Sempred(localctx antlr.RuleContext, predIndex int) bool {
	this := p
	_ = this

	switch predIndex {
	case 0:
		return p.Precpred(p.GetParserRuleContext(), 2)

	default:
		panic("No predicate with index: " + fmt.Sprint(predIndex))
	}
}
