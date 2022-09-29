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
		"", "", "", "','", "'('", "')'", "'.'", "';'", "", "", "", "", "", "",
		"", "", "", "", "'='", "'>'", "'<'", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "'*'",
	}
	staticData.symbolicNames = []string{
		"", "SPACE", "WS", "COMMA", "L_BRACKET", "R_BRACKET", "DOT", "EOQ",
		"K_SELECT", "K_WHERE", "K_WINDOW_TUMBLING", "K_GROUP_BY", "K_AND", "K_AS",
		"K_OR", "K_IS", "K_LIKE", "K_NOT_LIKE", "K_EQUAL", "K_GREATER", "K_LESS",
		"K_LESS_EQUAL", "K_GREATER_EQUAL", "K_NOT_EQUAL", "K_NULL", "K_IS_NULL",
		"K_IS_NOT_NULL", "K_NOT", "K_NOT_IN", "K_IN", "K_COUNT", "K_SUM", "K_MIN",
		"K_MAX", "K_AVG", "K_TRUE", "K_FALSE", "IDENTIFIER", "NUMERIC_LITERAL",
		"STRING_LITERAL", "BOOLEAN_LITERAL", "STAR",
	}
	staticData.ruleNames = []string{
		"sqlQuery", "selectQuery", "windowTumbling", "resultColumns", "aggregationColumns",
		"groupByAggregationColumns", "aggregationColumn", "column", "alias",
		"function", "functionName", "whereStatement", "expr", "simpleExpr",
		"compoundExpr", "comparisonOperator", "literalValue", "groupBy",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 41, 230, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 2, 17, 7, 17, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 3, 1, 43,
		8, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 51, 8, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 3, 1, 59, 8, 1, 1, 1, 1, 1, 1, 1, 3, 1, 64, 8, 1,
		1, 2, 1, 2, 1, 2, 1, 3, 1, 3, 1, 3, 5, 3, 72, 8, 3, 10, 3, 12, 3, 75, 9,
		3, 1, 3, 3, 3, 78, 8, 3, 1, 4, 1, 4, 1, 4, 5, 4, 83, 8, 4, 10, 4, 12, 4,
		86, 9, 4, 1, 5, 1, 5, 3, 5, 90, 8, 5, 1, 5, 1, 5, 1, 5, 3, 5, 95, 8, 5,
		5, 5, 97, 8, 5, 10, 5, 12, 5, 100, 9, 5, 1, 6, 1, 6, 1, 6, 1, 6, 1, 6,
		1, 6, 3, 6, 108, 8, 6, 1, 6, 1, 6, 3, 6, 112, 8, 6, 1, 7, 1, 7, 1, 7, 1,
		7, 3, 7, 118, 8, 7, 1, 7, 3, 7, 121, 8, 7, 1, 7, 1, 7, 3, 7, 125, 8, 7,
		3, 7, 127, 8, 7, 1, 8, 1, 8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9,
		3, 9, 138, 8, 9, 1, 9, 1, 9, 5, 9, 142, 8, 9, 10, 9, 12, 9, 145, 9, 9,
		1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9, 155, 8, 9, 1, 9,
		1, 9, 5, 9, 159, 8, 9, 10, 9, 12, 9, 162, 9, 9, 1, 9, 1, 9, 1, 9, 1, 9,
		1, 9, 1, 9, 1, 9, 5, 9, 171, 8, 9, 10, 9, 12, 9, 174, 9, 9, 1, 9, 1, 9,
		3, 9, 178, 8, 9, 1, 10, 1, 10, 1, 11, 1, 11, 1, 11, 1, 12, 1, 12, 1, 12,
		1, 12, 3, 12, 189, 8, 12, 1, 12, 1, 12, 1, 12, 3, 12, 194, 8, 12, 1, 12,
		3, 12, 197, 8, 12, 1, 12, 1, 12, 1, 12, 5, 12, 202, 8, 12, 10, 12, 12,
		12, 205, 9, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13, 1, 13,
		1, 13, 1, 13, 3, 13, 217, 8, 13, 1, 14, 1, 14, 1, 14, 1, 14, 1, 15, 1,
		15, 1, 16, 1, 16, 1, 17, 1, 17, 1, 17, 1, 17, 0, 1, 24, 18, 0, 2, 4, 6,
		8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 0, 4, 1, 0, 30,
		34, 2, 0, 12, 12, 14, 14, 2, 0, 15, 23, 28, 29, 1, 0, 38, 40, 241, 0, 36,
		1, 0, 0, 0, 2, 63, 1, 0, 0, 0, 4, 65, 1, 0, 0, 0, 6, 77, 1, 0, 0, 0, 8,
		79, 1, 0, 0, 0, 10, 89, 1, 0, 0, 0, 12, 101, 1, 0, 0, 0, 14, 126, 1, 0,
		0, 0, 16, 128, 1, 0, 0, 0, 18, 177, 1, 0, 0, 0, 20, 179, 1, 0, 0, 0, 22,
		181, 1, 0, 0, 0, 24, 196, 1, 0, 0, 0, 26, 216, 1, 0, 0, 0, 28, 218, 1,
		0, 0, 0, 30, 222, 1, 0, 0, 0, 32, 224, 1, 0, 0, 0, 34, 226, 1, 0, 0, 0,
		36, 37, 3, 2, 1, 0, 37, 38, 5, 0, 0, 1, 38, 1, 1, 0, 0, 0, 39, 40, 5, 8,
		0, 0, 40, 42, 3, 6, 3, 0, 41, 43, 3, 22, 11, 0, 42, 41, 1, 0, 0, 0, 42,
		43, 1, 0, 0, 0, 43, 44, 1, 0, 0, 0, 44, 45, 5, 7, 0, 0, 45, 64, 1, 0, 0,
		0, 46, 47, 5, 8, 0, 0, 47, 48, 3, 8, 4, 0, 48, 50, 3, 4, 2, 0, 49, 51,
		3, 22, 11, 0, 50, 49, 1, 0, 0, 0, 50, 51, 1, 0, 0, 0, 51, 52, 1, 0, 0,
		0, 52, 53, 5, 7, 0, 0, 53, 64, 1, 0, 0, 0, 54, 55, 5, 8, 0, 0, 55, 56,
		3, 10, 5, 0, 56, 58, 3, 4, 2, 0, 57, 59, 3, 22, 11, 0, 58, 57, 1, 0, 0,
		0, 58, 59, 1, 0, 0, 0, 59, 60, 1, 0, 0, 0, 60, 61, 3, 34, 17, 0, 61, 62,
		5, 7, 0, 0, 62, 64, 1, 0, 0, 0, 63, 39, 1, 0, 0, 0, 63, 46, 1, 0, 0, 0,
		63, 54, 1, 0, 0, 0, 64, 3, 1, 0, 0, 0, 65, 66, 5, 10, 0, 0, 66, 67, 5,
		38, 0, 0, 67, 5, 1, 0, 0, 0, 68, 73, 3, 14, 7, 0, 69, 70, 5, 3, 0, 0, 70,
		72, 3, 14, 7, 0, 71, 69, 1, 0, 0, 0, 72, 75, 1, 0, 0, 0, 73, 71, 1, 0,
		0, 0, 73, 74, 1, 0, 0, 0, 74, 78, 1, 0, 0, 0, 75, 73, 1, 0, 0, 0, 76, 78,
		5, 41, 0, 0, 77, 68, 1, 0, 0, 0, 77, 76, 1, 0, 0, 0, 78, 7, 1, 0, 0, 0,
		79, 84, 3, 12, 6, 0, 80, 81, 5, 3, 0, 0, 81, 83, 3, 12, 6, 0, 82, 80, 1,
		0, 0, 0, 83, 86, 1, 0, 0, 0, 84, 82, 1, 0, 0, 0, 84, 85, 1, 0, 0, 0, 85,
		9, 1, 0, 0, 0, 86, 84, 1, 0, 0, 0, 87, 90, 3, 14, 7, 0, 88, 90, 3, 12,
		6, 0, 89, 87, 1, 0, 0, 0, 89, 88, 1, 0, 0, 0, 90, 98, 1, 0, 0, 0, 91, 94,
		5, 3, 0, 0, 92, 95, 3, 12, 6, 0, 93, 95, 3, 14, 7, 0, 94, 92, 1, 0, 0,
		0, 94, 93, 1, 0, 0, 0, 95, 97, 1, 0, 0, 0, 96, 91, 1, 0, 0, 0, 97, 100,
		1, 0, 0, 0, 98, 96, 1, 0, 0, 0, 98, 99, 1, 0, 0, 0, 99, 11, 1, 0, 0, 0,
		100, 98, 1, 0, 0, 0, 101, 102, 7, 0, 0, 0, 102, 107, 5, 4, 0, 0, 103, 108,
		5, 37, 0, 0, 104, 105, 5, 37, 0, 0, 105, 106, 5, 6, 0, 0, 106, 108, 5,
		37, 0, 0, 107, 103, 1, 0, 0, 0, 107, 104, 1, 0, 0, 0, 108, 109, 1, 0, 0,
		0, 109, 111, 5, 5, 0, 0, 110, 112, 3, 16, 8, 0, 111, 110, 1, 0, 0, 0, 111,
		112, 1, 0, 0, 0, 112, 13, 1, 0, 0, 0, 113, 118, 5, 37, 0, 0, 114, 115,
		5, 37, 0, 0, 115, 116, 5, 6, 0, 0, 116, 118, 5, 37, 0, 0, 117, 113, 1,
		0, 0, 0, 117, 114, 1, 0, 0, 0, 118, 120, 1, 0, 0, 0, 119, 121, 3, 16, 8,
		0, 120, 119, 1, 0, 0, 0, 120, 121, 1, 0, 0, 0, 121, 127, 1, 0, 0, 0, 122,
		124, 3, 18, 9, 0, 123, 125, 3, 16, 8, 0, 124, 123, 1, 0, 0, 0, 124, 125,
		1, 0, 0, 0, 125, 127, 1, 0, 0, 0, 126, 117, 1, 0, 0, 0, 126, 122, 1, 0,
		0, 0, 127, 15, 1, 0, 0, 0, 128, 129, 5, 13, 0, 0, 129, 130, 5, 37, 0, 0,
		130, 17, 1, 0, 0, 0, 131, 132, 3, 20, 10, 0, 132, 137, 5, 4, 0, 0, 133,
		138, 5, 37, 0, 0, 134, 135, 5, 37, 0, 0, 135, 136, 5, 6, 0, 0, 136, 138,
		5, 37, 0, 0, 137, 133, 1, 0, 0, 0, 137, 134, 1, 0, 0, 0, 138, 143, 1, 0,
		0, 0, 139, 140, 5, 3, 0, 0, 140, 142, 3, 32, 16, 0, 141, 139, 1, 0, 0,
		0, 142, 145, 1, 0, 0, 0, 143, 141, 1, 0, 0, 0, 143, 144, 1, 0, 0, 0, 144,
		146, 1, 0, 0, 0, 145, 143, 1, 0, 0, 0, 146, 147, 5, 5, 0, 0, 147, 178,
		1, 0, 0, 0, 148, 149, 3, 20, 10, 0, 149, 154, 5, 4, 0, 0, 150, 155, 5,
		37, 0, 0, 151, 152, 5, 37, 0, 0, 152, 153, 5, 6, 0, 0, 153, 155, 5, 37,
		0, 0, 154, 150, 1, 0, 0, 0, 154, 151, 1, 0, 0, 0, 155, 160, 1, 0, 0, 0,
		156, 157, 5, 3, 0, 0, 157, 159, 3, 32, 16, 0, 158, 156, 1, 0, 0, 0, 159,
		162, 1, 0, 0, 0, 160, 158, 1, 0, 0, 0, 160, 161, 1, 0, 0, 0, 161, 163,
		1, 0, 0, 0, 162, 160, 1, 0, 0, 0, 163, 164, 5, 5, 0, 0, 164, 178, 1, 0,
		0, 0, 165, 166, 3, 20, 10, 0, 166, 167, 5, 4, 0, 0, 167, 172, 3, 18, 9,
		0, 168, 169, 5, 3, 0, 0, 169, 171, 3, 32, 16, 0, 170, 168, 1, 0, 0, 0,
		171, 174, 1, 0, 0, 0, 172, 170, 1, 0, 0, 0, 172, 173, 1, 0, 0, 0, 173,
		175, 1, 0, 0, 0, 174, 172, 1, 0, 0, 0, 175, 176, 5, 5, 0, 0, 176, 178,
		1, 0, 0, 0, 177, 131, 1, 0, 0, 0, 177, 148, 1, 0, 0, 0, 177, 165, 1, 0,
		0, 0, 178, 19, 1, 0, 0, 0, 179, 180, 5, 37, 0, 0, 180, 21, 1, 0, 0, 0,
		181, 182, 5, 9, 0, 0, 182, 183, 3, 24, 12, 0, 183, 23, 1, 0, 0, 0, 184,
		185, 6, 12, -1, 0, 185, 197, 3, 26, 13, 0, 186, 189, 3, 28, 14, 0, 187,
		189, 3, 26, 13, 0, 188, 186, 1, 0, 0, 0, 188, 187, 1, 0, 0, 0, 189, 190,
		1, 0, 0, 0, 190, 193, 7, 1, 0, 0, 191, 194, 3, 28, 14, 0, 192, 194, 3,
		26, 13, 0, 193, 191, 1, 0, 0, 0, 193, 192, 1, 0, 0, 0, 194, 197, 1, 0,
		0, 0, 195, 197, 3, 28, 14, 0, 196, 184, 1, 0, 0, 0, 196, 188, 1, 0, 0,
		0, 196, 195, 1, 0, 0, 0, 197, 203, 1, 0, 0, 0, 198, 199, 10, 2, 0, 0, 199,
		200, 7, 1, 0, 0, 200, 202, 3, 24, 12, 3, 201, 198, 1, 0, 0, 0, 202, 205,
		1, 0, 0, 0, 203, 201, 1, 0, 0, 0, 203, 204, 1, 0, 0, 0, 204, 25, 1, 0,
		0, 0, 205, 203, 1, 0, 0, 0, 206, 207, 5, 37, 0, 0, 207, 208, 3, 30, 15,
		0, 208, 209, 3, 32, 16, 0, 209, 217, 1, 0, 0, 0, 210, 211, 5, 37, 0, 0,
		211, 212, 5, 6, 0, 0, 212, 213, 5, 37, 0, 0, 213, 214, 3, 30, 15, 0, 214,
		215, 3, 32, 16, 0, 215, 217, 1, 0, 0, 0, 216, 206, 1, 0, 0, 0, 216, 210,
		1, 0, 0, 0, 217, 27, 1, 0, 0, 0, 218, 219, 5, 4, 0, 0, 219, 220, 3, 24,
		12, 0, 220, 221, 5, 5, 0, 0, 221, 29, 1, 0, 0, 0, 222, 223, 7, 2, 0, 0,
		223, 31, 1, 0, 0, 0, 224, 225, 7, 3, 0, 0, 225, 33, 1, 0, 0, 0, 226, 227,
		5, 11, 0, 0, 227, 228, 3, 14, 7, 0, 228, 35, 1, 0, 0, 0, 27, 42, 50, 58,
		63, 73, 77, 84, 89, 94, 98, 107, 111, 117, 120, 124, 126, 137, 143, 154,
		160, 172, 177, 188, 193, 196, 203, 216,
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
	SqlParserDOT               = 6
	SqlParserEOQ               = 7
	SqlParserK_SELECT          = 8
	SqlParserK_WHERE           = 9
	SqlParserK_WINDOW_TUMBLING = 10
	SqlParserK_GROUP_BY        = 11
	SqlParserK_AND             = 12
	SqlParserK_AS              = 13
	SqlParserK_OR              = 14
	SqlParserK_IS              = 15
	SqlParserK_LIKE            = 16
	SqlParserK_NOT_LIKE        = 17
	SqlParserK_EQUAL           = 18
	SqlParserK_GREATER         = 19
	SqlParserK_LESS            = 20
	SqlParserK_LESS_EQUAL      = 21
	SqlParserK_GREATER_EQUAL   = 22
	SqlParserK_NOT_EQUAL       = 23
	SqlParserK_NULL            = 24
	SqlParserK_IS_NULL         = 25
	SqlParserK_IS_NOT_NULL     = 26
	SqlParserK_NOT             = 27
	SqlParserK_NOT_IN          = 28
	SqlParserK_IN              = 29
	SqlParserK_COUNT           = 30
	SqlParserK_SUM             = 31
	SqlParserK_MIN             = 32
	SqlParserK_MAX             = 33
	SqlParserK_AVG             = 34
	SqlParserK_TRUE            = 35
	SqlParserK_FALSE           = 36
	SqlParserIDENTIFIER        = 37
	SqlParserNUMERIC_LITERAL   = 38
	SqlParserSTRING_LITERAL    = 39
	SqlParserBOOLEAN_LITERAL   = 40
	SqlParserSTAR              = 41
)

// SqlParser rules.
const (
	SqlParserRULE_sqlQuery                  = 0
	SqlParserRULE_selectQuery               = 1
	SqlParserRULE_windowTumbling            = 2
	SqlParserRULE_resultColumns             = 3
	SqlParserRULE_aggregationColumns        = 4
	SqlParserRULE_groupByAggregationColumns = 5
	SqlParserRULE_aggregationColumn         = 6
	SqlParserRULE_column                    = 7
	SqlParserRULE_alias                     = 8
	SqlParserRULE_function                  = 9
	SqlParserRULE_functionName              = 10
	SqlParserRULE_whereStatement            = 11
	SqlParserRULE_expr                      = 12
	SqlParserRULE_simpleExpr                = 13
	SqlParserRULE_compoundExpr              = 14
	SqlParserRULE_comparisonOperator        = 15
	SqlParserRULE_literalValue              = 16
	SqlParserRULE_groupBy                   = 17
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
		p.SetState(36)
		p.SelectQuery()
	}
	{
		p.SetState(37)
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

func (s *SelectQueryContext) CopyFrom(ctx *SelectQueryContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *SelectQueryContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectQueryContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SelectTumblingContext struct {
	*SelectQueryContext
}

func NewSelectTumblingContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectTumblingContext {
	var p = new(SelectTumblingContext)

	p.SelectQueryContext = NewEmptySelectQueryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SelectQueryContext))

	return p
}

func (s *SelectTumblingContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectTumblingContext) K_SELECT() antlr.TerminalNode {
	return s.GetToken(SqlParserK_SELECT, 0)
}

func (s *SelectTumblingContext) AggregationColumns() IAggregationColumnsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationColumnsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAggregationColumnsContext)
}

func (s *SelectTumblingContext) WindowTumbling() IWindowTumblingContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWindowTumblingContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWindowTumblingContext)
}

func (s *SelectTumblingContext) EOQ() antlr.TerminalNode {
	return s.GetToken(SqlParserEOQ, 0)
}

func (s *SelectTumblingContext) WhereStatement() IWhereStatementContext {
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

func (s *SelectTumblingContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSelectTumbling(s)

	default:
		return t.VisitChildren(s)
	}
}

type SelectSimpleContext struct {
	*SelectQueryContext
}

func NewSelectSimpleContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectSimpleContext {
	var p = new(SelectSimpleContext)

	p.SelectQueryContext = NewEmptySelectQueryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SelectQueryContext))

	return p
}

func (s *SelectSimpleContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectSimpleContext) K_SELECT() antlr.TerminalNode {
	return s.GetToken(SqlParserK_SELECT, 0)
}

func (s *SelectSimpleContext) ResultColumns() IResultColumnsContext {
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

func (s *SelectSimpleContext) EOQ() antlr.TerminalNode {
	return s.GetToken(SqlParserEOQ, 0)
}

func (s *SelectSimpleContext) WhereStatement() IWhereStatementContext {
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

func (s *SelectSimpleContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSelectSimple(s)

	default:
		return t.VisitChildren(s)
	}
}

type SelectTumblingGroupByContext struct {
	*SelectQueryContext
}

func NewSelectTumblingGroupByContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectTumblingGroupByContext {
	var p = new(SelectTumblingGroupByContext)

	p.SelectQueryContext = NewEmptySelectQueryContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SelectQueryContext))

	return p
}

func (s *SelectTumblingGroupByContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectTumblingGroupByContext) K_SELECT() antlr.TerminalNode {
	return s.GetToken(SqlParserK_SELECT, 0)
}

func (s *SelectTumblingGroupByContext) GroupByAggregationColumns() IGroupByAggregationColumnsContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IGroupByAggregationColumnsContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IGroupByAggregationColumnsContext)
}

func (s *SelectTumblingGroupByContext) WindowTumbling() IWindowTumblingContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IWindowTumblingContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IWindowTumblingContext)
}

func (s *SelectTumblingGroupByContext) GroupBy() IGroupByContext {
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

func (s *SelectTumblingGroupByContext) EOQ() antlr.TerminalNode {
	return s.GetToken(SqlParserEOQ, 0)
}

func (s *SelectTumblingGroupByContext) WhereStatement() IWhereStatementContext {
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

func (s *SelectTumblingGroupByContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSelectTumblingGroupBy(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) SelectQuery() (localctx ISelectQueryContext) {
	this := p
	_ = this

	localctx = NewSelectQueryContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 2, SqlParserRULE_selectQuery)
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

	p.SetState(63)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSelectSimpleContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(39)
			p.Match(SqlParserK_SELECT)
		}
		{
			p.SetState(40)
			p.ResultColumns()
		}
		p.SetState(42)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_WHERE {
			{
				p.SetState(41)
				p.WhereStatement()
			}

		}
		{
			p.SetState(44)
			p.Match(SqlParserEOQ)
		}

	case 2:
		localctx = NewSelectTumblingContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(46)
			p.Match(SqlParserK_SELECT)
		}
		{
			p.SetState(47)
			p.AggregationColumns()
		}
		{
			p.SetState(48)
			p.WindowTumbling()
		}
		p.SetState(50)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_WHERE {
			{
				p.SetState(49)
				p.WhereStatement()
			}

		}
		{
			p.SetState(52)
			p.Match(SqlParserEOQ)
		}

	case 3:
		localctx = NewSelectTumblingGroupByContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(54)
			p.Match(SqlParserK_SELECT)
		}
		{
			p.SetState(55)
			p.GroupByAggregationColumns()
		}
		{
			p.SetState(56)
			p.WindowTumbling()
		}
		p.SetState(58)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_WHERE {
			{
				p.SetState(57)
				p.WhereStatement()
			}

		}
		{
			p.SetState(60)
			p.GroupBy()
		}
		{
			p.SetState(61)
			p.Match(SqlParserEOQ)
		}

	}

	return localctx
}

// IWindowTumblingContext is an interface to support dynamic dispatch.
type IWindowTumblingContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsWindowTumblingContext differentiates from other interfaces.
	IsWindowTumblingContext()
}

type WindowTumblingContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyWindowTumblingContext() *WindowTumblingContext {
	var p = new(WindowTumblingContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_windowTumbling
	return p
}

func (*WindowTumblingContext) IsWindowTumblingContext() {}

func NewWindowTumblingContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *WindowTumblingContext {
	var p = new(WindowTumblingContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_windowTumbling

	return p
}

func (s *WindowTumblingContext) GetParser() antlr.Parser { return s.parser }

func (s *WindowTumblingContext) K_WINDOW_TUMBLING() antlr.TerminalNode {
	return s.GetToken(SqlParserK_WINDOW_TUMBLING, 0)
}

func (s *WindowTumblingContext) NUMERIC_LITERAL() antlr.TerminalNode {
	return s.GetToken(SqlParserNUMERIC_LITERAL, 0)
}

func (s *WindowTumblingContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *WindowTumblingContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *WindowTumblingContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitWindowTumbling(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) WindowTumbling() (localctx IWindowTumblingContext) {
	this := p
	_ = this

	localctx = NewWindowTumblingContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 4, SqlParserRULE_windowTumbling)

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
		p.Match(SqlParserK_WINDOW_TUMBLING)
	}
	{
		p.SetState(66)
		p.Match(SqlParserNUMERIC_LITERAL)
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
	p.EnterRule(localctx, 6, SqlParserRULE_resultColumns)
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

	p.SetState(77)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SqlParserIDENTIFIER:
		localctx = NewSelectColumnsContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(68)
			p.Column()
		}
		p.SetState(73)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == SqlParserCOMMA {
			{
				p.SetState(69)
				p.Match(SqlParserCOMMA)
			}
			{
				p.SetState(70)
				p.Column()
			}

			p.SetState(75)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case SqlParserSTAR:
		localctx = NewSelectStarContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(76)
			p.Match(SqlParserSTAR)
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

// IAggregationColumnsContext is an interface to support dynamic dispatch.
type IAggregationColumnsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAggregationColumnsContext differentiates from other interfaces.
	IsAggregationColumnsContext()
}

type AggregationColumnsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAggregationColumnsContext() *AggregationColumnsContext {
	var p = new(AggregationColumnsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_aggregationColumns
	return p
}

func (*AggregationColumnsContext) IsAggregationColumnsContext() {}

func NewAggregationColumnsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AggregationColumnsContext {
	var p = new(AggregationColumnsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_aggregationColumns

	return p
}

func (s *AggregationColumnsContext) GetParser() antlr.Parser { return s.parser }

func (s *AggregationColumnsContext) CopyFrom(ctx *AggregationColumnsContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *AggregationColumnsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregationColumnsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SelectAggregationsContext struct {
	*AggregationColumnsContext
}

func NewSelectAggregationsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectAggregationsContext {
	var p = new(SelectAggregationsContext)

	p.AggregationColumnsContext = NewEmptyAggregationColumnsContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AggregationColumnsContext))

	return p
}

func (s *SelectAggregationsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectAggregationsContext) AllAggregationColumn() []IAggregationColumnContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAggregationColumnContext); ok {
			len++
		}
	}

	tst := make([]IAggregationColumnContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAggregationColumnContext); ok {
			tst[i] = t.(IAggregationColumnContext)
			i++
		}
	}

	return tst
}

func (s *SelectAggregationsContext) AggregationColumn(i int) IAggregationColumnContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationColumnContext); ok {
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

	return t.(IAggregationColumnContext)
}

func (s *SelectAggregationsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SqlParserCOMMA)
}

func (s *SelectAggregationsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserCOMMA, i)
}

func (s *SelectAggregationsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSelectAggregations(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) AggregationColumns() (localctx IAggregationColumnsContext) {
	this := p
	_ = this

	localctx = NewAggregationColumnsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 8, SqlParserRULE_aggregationColumns)
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

	localctx = NewSelectAggregationsContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(79)
		p.AggregationColumn()
	}
	p.SetState(84)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SqlParserCOMMA {
		{
			p.SetState(80)
			p.Match(SqlParserCOMMA)
		}
		{
			p.SetState(81)
			p.AggregationColumn()
		}

		p.SetState(86)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IGroupByAggregationColumnsContext is an interface to support dynamic dispatch.
type IGroupByAggregationColumnsContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsGroupByAggregationColumnsContext differentiates from other interfaces.
	IsGroupByAggregationColumnsContext()
}

type GroupByAggregationColumnsContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyGroupByAggregationColumnsContext() *GroupByAggregationColumnsContext {
	var p = new(GroupByAggregationColumnsContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_groupByAggregationColumns
	return p
}

func (*GroupByAggregationColumnsContext) IsGroupByAggregationColumnsContext() {}

func NewGroupByAggregationColumnsContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *GroupByAggregationColumnsContext {
	var p = new(GroupByAggregationColumnsContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_groupByAggregationColumns

	return p
}

func (s *GroupByAggregationColumnsContext) GetParser() antlr.Parser { return s.parser }

func (s *GroupByAggregationColumnsContext) CopyFrom(ctx *GroupByAggregationColumnsContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *GroupByAggregationColumnsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *GroupByAggregationColumnsContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type SelectGroupByAggregationsContext struct {
	*GroupByAggregationColumnsContext
}

func NewSelectGroupByAggregationsContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectGroupByAggregationsContext {
	var p = new(SelectGroupByAggregationsContext)

	p.GroupByAggregationColumnsContext = NewEmptyGroupByAggregationColumnsContext()
	p.parser = parser
	p.CopyFrom(ctx.(*GroupByAggregationColumnsContext))

	return p
}

func (s *SelectGroupByAggregationsContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectGroupByAggregationsContext) AllColumn() []IColumnContext {
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

func (s *SelectGroupByAggregationsContext) Column(i int) IColumnContext {
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

func (s *SelectGroupByAggregationsContext) AllAggregationColumn() []IAggregationColumnContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(IAggregationColumnContext); ok {
			len++
		}
	}

	tst := make([]IAggregationColumnContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(IAggregationColumnContext); ok {
			tst[i] = t.(IAggregationColumnContext)
			i++
		}
	}

	return tst
}

func (s *SelectGroupByAggregationsContext) AggregationColumn(i int) IAggregationColumnContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAggregationColumnContext); ok {
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

	return t.(IAggregationColumnContext)
}

func (s *SelectGroupByAggregationsContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SqlParserCOMMA)
}

func (s *SelectGroupByAggregationsContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserCOMMA, i)
}

func (s *SelectGroupByAggregationsContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSelectGroupByAggregations(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) GroupByAggregationColumns() (localctx IGroupByAggregationColumnsContext) {
	this := p
	_ = this

	localctx = NewGroupByAggregationColumnsContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 10, SqlParserRULE_groupByAggregationColumns)
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

	localctx = NewSelectGroupByAggregationsContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	p.SetState(89)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SqlParserIDENTIFIER:
		{
			p.SetState(87)
			p.Column()
		}

	case SqlParserK_COUNT, SqlParserK_SUM, SqlParserK_MIN, SqlParserK_MAX, SqlParserK_AVG:
		{
			p.SetState(88)
			p.AggregationColumn()
		}

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}
	p.SetState(98)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SqlParserCOMMA {
		{
			p.SetState(91)
			p.Match(SqlParserCOMMA)
		}
		p.SetState(94)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SqlParserK_COUNT, SqlParserK_SUM, SqlParserK_MIN, SqlParserK_MAX, SqlParserK_AVG:
			{
				p.SetState(92)
				p.AggregationColumn()
			}

		case SqlParserIDENTIFIER:
			{
				p.SetState(93)
				p.Column()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

		p.SetState(100)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
	}

	return localctx
}

// IAggregationColumnContext is an interface to support dynamic dispatch.
type IAggregationColumnContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAggregationColumnContext differentiates from other interfaces.
	IsAggregationColumnContext()
}

type AggregationColumnContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAggregationColumnContext() *AggregationColumnContext {
	var p = new(AggregationColumnContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_aggregationColumn
	return p
}

func (*AggregationColumnContext) IsAggregationColumnContext() {}

func NewAggregationColumnContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AggregationColumnContext {
	var p = new(AggregationColumnContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_aggregationColumn

	return p
}

func (s *AggregationColumnContext) GetParser() antlr.Parser { return s.parser }

func (s *AggregationColumnContext) L_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserL_BRACKET, 0)
}

func (s *AggregationColumnContext) R_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserR_BRACKET, 0)
}

func (s *AggregationColumnContext) K_MIN() antlr.TerminalNode {
	return s.GetToken(SqlParserK_MIN, 0)
}

func (s *AggregationColumnContext) K_MAX() antlr.TerminalNode {
	return s.GetToken(SqlParserK_MAX, 0)
}

func (s *AggregationColumnContext) K_COUNT() antlr.TerminalNode {
	return s.GetToken(SqlParserK_COUNT, 0)
}

func (s *AggregationColumnContext) K_AVG() antlr.TerminalNode {
	return s.GetToken(SqlParserK_AVG, 0)
}

func (s *AggregationColumnContext) K_SUM() antlr.TerminalNode {
	return s.GetToken(SqlParserK_SUM, 0)
}

func (s *AggregationColumnContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(SqlParserIDENTIFIER)
}

func (s *AggregationColumnContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, i)
}

func (s *AggregationColumnContext) DOT() antlr.TerminalNode {
	return s.GetToken(SqlParserDOT, 0)
}

func (s *AggregationColumnContext) Alias() IAliasContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAliasContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAliasContext)
}

func (s *AggregationColumnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregationColumnContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AggregationColumnContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitAggregationColumn(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) AggregationColumn() (localctx IAggregationColumnContext) {
	this := p
	_ = this

	localctx = NewAggregationColumnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 12, SqlParserRULE_aggregationColumn)
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
		p.SetState(101)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-30)&-(0x1f+1)) == 0 && ((1<<uint((_la-30)))&((1<<(SqlParserK_COUNT-30))|(1<<(SqlParserK_SUM-30))|(1<<(SqlParserK_MIN-30))|(1<<(SqlParserK_MAX-30))|(1<<(SqlParserK_AVG-30)))) != 0) {
			p.GetErrorHandler().RecoverInline(p)
		} else {
			p.GetErrorHandler().ReportMatch(p)
			p.Consume()
		}
	}
	{
		p.SetState(102)
		p.Match(SqlParserL_BRACKET)
	}
	p.SetState(107)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 10, p.GetParserRuleContext()) {
	case 1:
		{
			p.SetState(103)
			p.Match(SqlParserIDENTIFIER)
		}

	case 2:
		{
			p.SetState(104)
			p.Match(SqlParserIDENTIFIER)
		}
		{
			p.SetState(105)
			p.Match(SqlParserDOT)
		}
		{
			p.SetState(106)
			p.Match(SqlParserIDENTIFIER)
		}

	}
	{
		p.SetState(109)
		p.Match(SqlParserR_BRACKET)
	}
	p.SetState(111)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	if _la == SqlParserK_AS {
		{
			p.SetState(110)
			p.Alias()
		}

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

func (s *ColumnContext) CopyFrom(ctx *ColumnContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *ColumnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type IdentifierColumnContext struct {
	*ColumnContext
}

func NewIdentifierColumnContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *IdentifierColumnContext {
	var p = new(IdentifierColumnContext)

	p.ColumnContext = NewEmptyColumnContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ColumnContext))

	return p
}

func (s *IdentifierColumnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *IdentifierColumnContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(SqlParserIDENTIFIER)
}

func (s *IdentifierColumnContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, i)
}

func (s *IdentifierColumnContext) DOT() antlr.TerminalNode {
	return s.GetToken(SqlParserDOT, 0)
}

func (s *IdentifierColumnContext) Alias() IAliasContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAliasContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAliasContext)
}

func (s *IdentifierColumnContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitIdentifierColumn(s)

	default:
		return t.VisitChildren(s)
	}
}

type FunctionColumnContext struct {
	*ColumnContext
}

func NewFunctionColumnContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *FunctionColumnContext {
	var p = new(FunctionColumnContext)

	p.ColumnContext = NewEmptyColumnContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ColumnContext))

	return p
}

func (s *FunctionColumnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionColumnContext) Function() IFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionContext)
}

func (s *FunctionColumnContext) Alias() IAliasContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IAliasContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IAliasContext)
}

func (s *FunctionColumnContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitFunctionColumn(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) Column() (localctx IColumnContext) {
	this := p
	_ = this

	localctx = NewColumnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 14, SqlParserRULE_column)
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

	p.SetState(126)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 15, p.GetParserRuleContext()) {
	case 1:
		localctx = NewIdentifierColumnContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(117)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(113)
				p.Match(SqlParserIDENTIFIER)
			}

		case 2:
			{
				p.SetState(114)
				p.Match(SqlParserIDENTIFIER)
			}
			{
				p.SetState(115)
				p.Match(SqlParserDOT)
			}
			{
				p.SetState(116)
				p.Match(SqlParserIDENTIFIER)
			}

		}
		p.SetState(120)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_AS {
			{
				p.SetState(119)
				p.Alias()
			}

		}

	case 2:
		localctx = NewFunctionColumnContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(122)
			p.Function()
		}
		p.SetState(124)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_AS {
			{
				p.SetState(123)
				p.Alias()
			}

		}

	}

	return localctx
}

// IAliasContext is an interface to support dynamic dispatch.
type IAliasContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsAliasContext differentiates from other interfaces.
	IsAliasContext()
}

type AliasContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyAliasContext() *AliasContext {
	var p = new(AliasContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_alias
	return p
}

func (*AliasContext) IsAliasContext() {}

func NewAliasContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *AliasContext {
	var p = new(AliasContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_alias

	return p
}

func (s *AliasContext) GetParser() antlr.Parser { return s.parser }

func (s *AliasContext) K_AS() antlr.TerminalNode {
	return s.GetToken(SqlParserK_AS, 0)
}

func (s *AliasContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, 0)
}

func (s *AliasContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AliasContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *AliasContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitAlias(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) Alias() (localctx IAliasContext) {
	this := p
	_ = this

	localctx = NewAliasContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 16, SqlParserRULE_alias)

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
		p.SetState(128)
		p.Match(SqlParserK_AS)
	}
	{
		p.SetState(129)
		p.Match(SqlParserIDENTIFIER)
	}

	return localctx
}

// IFunctionContext is an interface to support dynamic dispatch.
type IFunctionContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFunctionContext differentiates from other interfaces.
	IsFunctionContext()
}

type FunctionContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionContext() *FunctionContext {
	var p = new(FunctionContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_function
	return p
}

func (*FunctionContext) IsFunctionContext() {}

func NewFunctionContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionContext {
	var p = new(FunctionContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_function

	return p
}

func (s *FunctionContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionContext) CopyFrom(ctx *FunctionContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *FunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type RecursiveFunctionContext struct {
	*FunctionContext
}

func NewRecursiveFunctionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *RecursiveFunctionContext {
	var p = new(RecursiveFunctionContext)

	p.FunctionContext = NewEmptyFunctionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*FunctionContext))

	return p
}

func (s *RecursiveFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *RecursiveFunctionContext) FunctionName() IFunctionNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionNameContext)
}

func (s *RecursiveFunctionContext) L_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserL_BRACKET, 0)
}

func (s *RecursiveFunctionContext) Function() IFunctionContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionContext)
}

func (s *RecursiveFunctionContext) R_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserR_BRACKET, 0)
}

func (s *RecursiveFunctionContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SqlParserCOMMA)
}

func (s *RecursiveFunctionContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserCOMMA, i)
}

func (s *RecursiveFunctionContext) AllLiteralValue() []ILiteralValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILiteralValueContext); ok {
			len++
		}
	}

	tst := make([]ILiteralValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILiteralValueContext); ok {
			tst[i] = t.(ILiteralValueContext)
			i++
		}
	}

	return tst
}

func (s *RecursiveFunctionContext) LiteralValue(i int) ILiteralValueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralValueContext); ok {
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

	return t.(ILiteralValueContext)
}

func (s *RecursiveFunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitRecursiveFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

type SimpleFunctionContext struct {
	*FunctionContext
}

func NewSimpleFunctionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SimpleFunctionContext {
	var p = new(SimpleFunctionContext)

	p.FunctionContext = NewEmptyFunctionContext()
	p.parser = parser
	p.CopyFrom(ctx.(*FunctionContext))

	return p
}

func (s *SimpleFunctionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpleFunctionContext) FunctionName() IFunctionNameContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(IFunctionNameContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(IFunctionNameContext)
}

func (s *SimpleFunctionContext) L_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserL_BRACKET, 0)
}

func (s *SimpleFunctionContext) R_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserR_BRACKET, 0)
}

func (s *SimpleFunctionContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(SqlParserIDENTIFIER)
}

func (s *SimpleFunctionContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, i)
}

func (s *SimpleFunctionContext) DOT() antlr.TerminalNode {
	return s.GetToken(SqlParserDOT, 0)
}

func (s *SimpleFunctionContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SqlParserCOMMA)
}

func (s *SimpleFunctionContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserCOMMA, i)
}

func (s *SimpleFunctionContext) AllLiteralValue() []ILiteralValueContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ILiteralValueContext); ok {
			len++
		}
	}

	tst := make([]ILiteralValueContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ILiteralValueContext); ok {
			tst[i] = t.(ILiteralValueContext)
			i++
		}
	}

	return tst
}

func (s *SimpleFunctionContext) LiteralValue(i int) ILiteralValueContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ILiteralValueContext); ok {
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

	return t.(ILiteralValueContext)
}

func (s *SimpleFunctionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSimpleFunction(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) Function() (localctx IFunctionContext) {
	this := p
	_ = this

	localctx = NewFunctionContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SqlParserRULE_function)
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

	p.SetState(177)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSimpleFunctionContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(131)
			p.FunctionName()
		}
		{
			p.SetState(132)
			p.Match(SqlParserL_BRACKET)
		}
		p.SetState(137)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 16, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(133)
				p.Match(SqlParserIDENTIFIER)
			}

		case 2:
			{
				p.SetState(134)
				p.Match(SqlParserIDENTIFIER)
			}
			{
				p.SetState(135)
				p.Match(SqlParserDOT)
			}
			{
				p.SetState(136)
				p.Match(SqlParserIDENTIFIER)
			}

		}
		p.SetState(143)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == SqlParserCOMMA {
			{
				p.SetState(139)
				p.Match(SqlParserCOMMA)
			}
			{
				p.SetState(140)
				p.LiteralValue()
			}

			p.SetState(145)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(146)
			p.Match(SqlParserR_BRACKET)
		}

	case 2:
		localctx = NewSimpleFunctionContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(148)
			p.FunctionName()
		}
		{
			p.SetState(149)
			p.Match(SqlParserL_BRACKET)
		}
		p.SetState(154)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 18, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(150)
				p.Match(SqlParserIDENTIFIER)
			}

		case 2:
			{
				p.SetState(151)
				p.Match(SqlParserIDENTIFIER)
			}
			{
				p.SetState(152)
				p.Match(SqlParserDOT)
			}
			{
				p.SetState(153)
				p.Match(SqlParserIDENTIFIER)
			}

		}
		p.SetState(160)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == SqlParserCOMMA {
			{
				p.SetState(156)
				p.Match(SqlParserCOMMA)
			}
			{
				p.SetState(157)
				p.LiteralValue()
			}

			p.SetState(162)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(163)
			p.Match(SqlParserR_BRACKET)
		}

	case 3:
		localctx = NewRecursiveFunctionContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(165)
			p.FunctionName()
		}
		{
			p.SetState(166)
			p.Match(SqlParserL_BRACKET)
		}
		{
			p.SetState(167)
			p.Function()
		}
		p.SetState(172)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == SqlParserCOMMA {
			{
				p.SetState(168)
				p.Match(SqlParserCOMMA)
			}
			{
				p.SetState(169)
				p.LiteralValue()
			}

			p.SetState(174)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(175)
			p.Match(SqlParserR_BRACKET)
		}

	}

	return localctx
}

// IFunctionNameContext is an interface to support dynamic dispatch.
type IFunctionNameContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsFunctionNameContext differentiates from other interfaces.
	IsFunctionNameContext()
}

type FunctionNameContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptyFunctionNameContext() *FunctionNameContext {
	var p = new(FunctionNameContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_functionName
	return p
}

func (*FunctionNameContext) IsFunctionNameContext() {}

func NewFunctionNameContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *FunctionNameContext {
	var p = new(FunctionNameContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_functionName

	return p
}

func (s *FunctionNameContext) GetParser() antlr.Parser { return s.parser }

func (s *FunctionNameContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, 0)
}

func (s *FunctionNameContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *FunctionNameContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

func (s *FunctionNameContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitFunctionName(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) FunctionName() (localctx IFunctionNameContext) {
	this := p
	_ = this

	localctx = NewFunctionNameContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 20, SqlParserRULE_functionName)

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
		p.SetState(179)
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

func (p *SqlParser) WhereStatement() (localctx IWhereStatementContext) {
	this := p
	_ = this

	localctx = NewWhereStatementContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 22, SqlParserRULE_whereStatement)

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

	localctx = NewWhereStmtContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(181)
		p.Match(SqlParserK_WHERE)
	}
	{
		p.SetState(182)
		p.expr(0)
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

func (s *SimpleConditionContext) SimpleExpr() ISimpleExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISimpleExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ISimpleExprContext)
}

func (s *SimpleConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSimpleCondition(s)

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

func (s *CompoundRecursiveConditionContext) K_AND() antlr.TerminalNode {
	return s.GetToken(SqlParserK_AND, 0)
}

func (s *CompoundRecursiveConditionContext) K_OR() antlr.TerminalNode {
	return s.GetToken(SqlParserK_OR, 0)
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

func (s *CompoundRecursiveConditionContext) AllSimpleExpr() []ISimpleExprContext {
	children := s.GetChildren()
	len := 0
	for _, ctx := range children {
		if _, ok := ctx.(ISimpleExprContext); ok {
			len++
		}
	}

	tst := make([]ISimpleExprContext, len)
	i := 0
	for _, ctx := range children {
		if t, ok := ctx.(ISimpleExprContext); ok {
			tst[i] = t.(ISimpleExprContext)
			i++
		}
	}

	return tst
}

func (s *CompoundRecursiveConditionContext) SimpleExpr(i int) ISimpleExprContext {
	var t antlr.RuleContext
	j := 0
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ISimpleExprContext); ok {
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

	return t.(ISimpleExprContext)
}

func (s *CompoundRecursiveConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitCompoundRecursiveCondition(s)

	default:
		return t.VisitChildren(s)
	}
}

type SimpleCompoundConditionContext struct {
	*ExprContext
}

func NewSimpleCompoundConditionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SimpleCompoundConditionContext {
	var p = new(SimpleCompoundConditionContext)

	p.ExprContext = NewEmptyExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*ExprContext))

	return p
}

func (s *SimpleCompoundConditionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpleCompoundConditionContext) CompoundExpr() ICompoundExprContext {
	var t antlr.RuleContext
	for _, ctx := range s.GetChildren() {
		if _, ok := ctx.(ICompoundExprContext); ok {
			t = ctx.(antlr.RuleContext)
			break
		}
	}

	if t == nil {
		return nil
	}

	return t.(ICompoundExprContext)
}

func (s *SimpleCompoundConditionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSimpleCompoundCondition(s)

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
	_startState := 24
	p.EnterRecursionRule(localctx, 24, SqlParserRULE_expr, _p)
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
	p.SetState(196)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSimpleConditionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(185)
			p.SimpleExpr()
		}

	case 2:
		localctx = NewCompoundRecursiveConditionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		p.SetState(188)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SqlParserL_BRACKET:
			{
				p.SetState(186)
				p.CompoundExpr()
			}

		case SqlParserIDENTIFIER:
			{
				p.SetState(187)
				p.SimpleExpr()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(190)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SqlParserK_AND || _la == SqlParserK_OR) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(193)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SqlParserL_BRACKET:
			{
				p.SetState(191)
				p.CompoundExpr()
			}

		case SqlParserIDENTIFIER:
			{
				p.SetState(192)
				p.SimpleExpr()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}

	case 3:
		localctx = NewSimpleCompoundConditionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		{
			p.SetState(195)
			p.CompoundExpr()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(203)
	p.GetErrorHandler().Sync(p)
	_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext())

	for _alt != 2 && _alt != antlr.ATNInvalidAltNumber {
		if _alt == 1 {
			if p.GetParseListeners() != nil {
				p.TriggerExitRuleEvent()
			}
			_prevctx = localctx
			localctx = NewSimpleRecursiveConditionContext(p, NewExprContext(p, _parentctx, _parentState))
			p.PushNewRecursionContext(localctx, _startState, SqlParserRULE_expr)
			p.SetState(198)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(199)
				_la = p.GetTokenStream().LA(1)

				if !(_la == SqlParserK_AND || _la == SqlParserK_OR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(200)
				p.expr(3)
			}

		}
		p.SetState(205)
		p.GetErrorHandler().Sync(p)
		_alt = p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 25, p.GetParserRuleContext())
	}

	return localctx
}

// ISimpleExprContext is an interface to support dynamic dispatch.
type ISimpleExprContext interface {
	antlr.ParserRuleContext

	// GetParser returns the parser.
	GetParser() antlr.Parser

	// IsSimpleExprContext differentiates from other interfaces.
	IsSimpleExprContext()
}

type SimpleExprContext struct {
	*antlr.BaseParserRuleContext
	parser antlr.Parser
}

func NewEmptySimpleExprContext() *SimpleExprContext {
	var p = new(SimpleExprContext)
	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(nil, -1)
	p.RuleIndex = SqlParserRULE_simpleExpr
	return p
}

func (*SimpleExprContext) IsSimpleExprContext() {}

func NewSimpleExprContext(parser antlr.Parser, parent antlr.ParserRuleContext, invokingState int) *SimpleExprContext {
	var p = new(SimpleExprContext)

	p.BaseParserRuleContext = antlr.NewBaseParserRuleContext(parent, invokingState)

	p.parser = parser
	p.RuleIndex = SqlParserRULE_simpleExpr

	return p
}

func (s *SimpleExprContext) GetParser() antlr.Parser { return s.parser }

func (s *SimpleExprContext) CopyFrom(ctx *SimpleExprContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *SimpleExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpleExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type NestedExpressionContext struct {
	*SimpleExprContext
}

func NewNestedExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *NestedExpressionContext {
	var p = new(NestedExpressionContext)

	p.SimpleExprContext = NewEmptySimpleExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SimpleExprContext))

	return p
}

func (s *NestedExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *NestedExpressionContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(SqlParserIDENTIFIER)
}

func (s *NestedExpressionContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, i)
}

func (s *NestedExpressionContext) DOT() antlr.TerminalNode {
	return s.GetToken(SqlParserDOT, 0)
}

func (s *NestedExpressionContext) ComparisonOperator() IComparisonOperatorContext {
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

func (s *NestedExpressionContext) LiteralValue() ILiteralValueContext {
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

func (s *NestedExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitNestedExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

type SimpleExpressionContext struct {
	*SimpleExprContext
}

func NewSimpleExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SimpleExpressionContext {
	var p = new(SimpleExpressionContext)

	p.SimpleExprContext = NewEmptySimpleExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*SimpleExprContext))

	return p
}

func (s *SimpleExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SimpleExpressionContext) IDENTIFIER() antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, 0)
}

func (s *SimpleExpressionContext) ComparisonOperator() IComparisonOperatorContext {
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

func (s *SimpleExpressionContext) LiteralValue() ILiteralValueContext {
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

func (s *SimpleExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSimpleExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) SimpleExpr() (localctx ISimpleExprContext) {
	this := p
	_ = this

	localctx = NewSimpleExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 26, SqlParserRULE_simpleExpr)

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

	p.SetState(216)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSimpleExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(206)
			p.Match(SqlParserIDENTIFIER)
		}
		{
			p.SetState(207)
			p.ComparisonOperator()
		}
		{
			p.SetState(208)
			p.LiteralValue()
		}

	case 2:
		localctx = NewNestedExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(210)
			p.Match(SqlParserIDENTIFIER)
		}
		{
			p.SetState(211)
			p.Match(SqlParserDOT)
		}
		{
			p.SetState(212)
			p.Match(SqlParserIDENTIFIER)
		}
		{
			p.SetState(213)
			p.ComparisonOperator()
		}
		{
			p.SetState(214)
			p.LiteralValue()
		}

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

func (s *CompoundExprContext) CopyFrom(ctx *CompoundExprContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *CompoundExprContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompoundExprContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type CompoundExpressionContext struct {
	*CompoundExprContext
}

func NewCompoundExpressionContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *CompoundExpressionContext {
	var p = new(CompoundExpressionContext)

	p.CompoundExprContext = NewEmptyCompoundExprContext()
	p.parser = parser
	p.CopyFrom(ctx.(*CompoundExprContext))

	return p
}

func (s *CompoundExpressionContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *CompoundExpressionContext) L_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserL_BRACKET, 0)
}

func (s *CompoundExpressionContext) Expr() IExprContext {
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

func (s *CompoundExpressionContext) R_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserR_BRACKET, 0)
}

func (s *CompoundExpressionContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitCompoundExpression(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) CompoundExpr() (localctx ICompoundExprContext) {
	this := p
	_ = this

	localctx = NewCompoundExprContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 28, SqlParserRULE_compoundExpr)

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

	localctx = NewCompoundExpressionContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	{
		p.SetState(218)
		p.Match(SqlParserL_BRACKET)
	}
	{
		p.SetState(219)
		p.expr(0)
	}
	{
		p.SetState(220)
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

func (s *ComparisonOperatorContext) K_NOT_LIKE() antlr.TerminalNode {
	return s.GetToken(SqlParserK_NOT_LIKE, 0)
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
	p.EnterRule(localctx, 30, SqlParserRULE_comparisonOperator)
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
		p.SetState(222)
		_la = p.GetTokenStream().LA(1)

		if !(((_la)&-(0x1f+1)) == 0 && ((1<<uint(_la))&((1<<SqlParserK_IS)|(1<<SqlParserK_LIKE)|(1<<SqlParserK_NOT_LIKE)|(1<<SqlParserK_EQUAL)|(1<<SqlParserK_GREATER)|(1<<SqlParserK_LESS)|(1<<SqlParserK_LESS_EQUAL)|(1<<SqlParserK_GREATER_EQUAL)|(1<<SqlParserK_NOT_EQUAL)|(1<<SqlParserK_NOT_IN)|(1<<SqlParserK_IN))) != 0) {
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

func (s *LiteralValueContext) BOOLEAN_LITERAL() antlr.TerminalNode {
	return s.GetToken(SqlParserBOOLEAN_LITERAL, 0)
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
	p.EnterRule(localctx, 32, SqlParserRULE_literalValue)
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
		p.SetState(224)
		_la = p.GetTokenStream().LA(1)

		if !(((_la-38)&-(0x1f+1)) == 0 && ((1<<uint((_la-38)))&((1<<(SqlParserNUMERIC_LITERAL-38))|(1<<(SqlParserSTRING_LITERAL-38))|(1<<(SqlParserBOOLEAN_LITERAL-38)))) != 0) {
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
	p.EnterRule(localctx, 34, SqlParserRULE_groupBy)

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
		p.SetState(226)
		p.Match(SqlParserK_GROUP_BY)
	}
	{
		p.SetState(227)
		p.Column()
	}

	return localctx
}

func (p *SqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 12:
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
