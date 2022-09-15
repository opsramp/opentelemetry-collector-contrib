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
		"column", "alias", "function", "functionName", "aggregationColumn",
		"whereStatement", "expr", "simpleExpr", "compoundExpr", "comparisonOperator",
		"literalValue", "groupBy",
	}
	staticData.predictionContextCache = antlr.NewPredictionContextCache()
	staticData.serializedATN = []int32{
		4, 1, 41, 233, 2, 0, 7, 0, 2, 1, 7, 1, 2, 2, 7, 2, 2, 3, 7, 3, 2, 4, 7,
		4, 2, 5, 7, 5, 2, 6, 7, 6, 2, 7, 7, 7, 2, 8, 7, 8, 2, 9, 7, 9, 2, 10, 7,
		10, 2, 11, 7, 11, 2, 12, 7, 12, 2, 13, 7, 13, 2, 14, 7, 14, 2, 15, 7, 15,
		2, 16, 7, 16, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 1, 3, 1, 41, 8, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 49, 8, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 3, 1, 57, 8, 1, 1, 1, 1, 1, 1, 1, 3, 1, 62, 8, 1, 1, 2, 1, 2,
		1, 2, 1, 3, 1, 3, 1, 3, 5, 3, 70, 8, 3, 10, 3, 12, 3, 73, 9, 3, 1, 3, 3,
		3, 76, 8, 3, 1, 4, 1, 4, 1, 4, 3, 4, 81, 8, 4, 1, 4, 1, 4, 1, 4, 5, 4,
		86, 8, 4, 10, 4, 12, 4, 89, 9, 4, 1, 5, 1, 5, 1, 5, 1, 5, 3, 5, 95, 8,
		5, 1, 5, 3, 5, 98, 8, 5, 1, 5, 1, 5, 3, 5, 102, 8, 5, 3, 5, 104, 8, 5,
		1, 6, 1, 6, 1, 6, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 3, 7, 115, 8, 7,
		1, 7, 1, 7, 5, 7, 119, 8, 7, 10, 7, 12, 7, 122, 9, 7, 1, 7, 1, 7, 1, 7,
		1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 3, 7, 132, 8, 7, 1, 7, 1, 7, 5, 7, 136, 8,
		7, 10, 7, 12, 7, 139, 9, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 1, 7, 5,
		7, 148, 8, 7, 10, 7, 12, 7, 151, 9, 7, 1, 7, 1, 7, 3, 7, 155, 8, 7, 1,
		8, 1, 8, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9, 165, 8, 9, 1, 9, 1,
		9, 3, 9, 169, 8, 9, 1, 9, 1, 9, 1, 9, 1, 9, 3, 9, 175, 8, 9, 1, 9, 1, 9,
		1, 9, 1, 9, 3, 9, 181, 8, 9, 1, 10, 1, 10, 1, 10, 1, 11, 1, 11, 1, 11,
		1, 11, 3, 11, 190, 8, 11, 1, 11, 1, 11, 1, 11, 3, 11, 195, 8, 11, 1, 11,
		3, 11, 198, 8, 11, 1, 11, 1, 11, 1, 11, 5, 11, 203, 8, 11, 10, 11, 12,
		11, 206, 9, 11, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12, 1, 12,
		1, 12, 1, 12, 3, 12, 218, 8, 12, 1, 13, 1, 13, 1, 13, 1, 13, 1, 14, 1,
		14, 1, 15, 1, 15, 1, 16, 1, 16, 1, 16, 3, 16, 231, 8, 16, 1, 16, 0, 1,
		22, 17, 0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32,
		0, 4, 1, 0, 30, 34, 2, 0, 12, 12, 14, 14, 2, 0, 15, 23, 28, 29, 1, 0, 38,
		40, 247, 0, 34, 1, 0, 0, 0, 2, 61, 1, 0, 0, 0, 4, 63, 1, 0, 0, 0, 6, 75,
		1, 0, 0, 0, 8, 80, 1, 0, 0, 0, 10, 103, 1, 0, 0, 0, 12, 105, 1, 0, 0, 0,
		14, 154, 1, 0, 0, 0, 16, 156, 1, 0, 0, 0, 18, 180, 1, 0, 0, 0, 20, 182,
		1, 0, 0, 0, 22, 197, 1, 0, 0, 0, 24, 217, 1, 0, 0, 0, 26, 219, 1, 0, 0,
		0, 28, 223, 1, 0, 0, 0, 30, 225, 1, 0, 0, 0, 32, 230, 1, 0, 0, 0, 34, 35,
		3, 2, 1, 0, 35, 36, 5, 0, 0, 1, 36, 1, 1, 0, 0, 0, 37, 38, 5, 8, 0, 0,
		38, 40, 3, 6, 3, 0, 39, 41, 3, 20, 10, 0, 40, 39, 1, 0, 0, 0, 40, 41, 1,
		0, 0, 0, 41, 42, 1, 0, 0, 0, 42, 43, 5, 7, 0, 0, 43, 62, 1, 0, 0, 0, 44,
		45, 5, 8, 0, 0, 45, 46, 3, 8, 4, 0, 46, 48, 3, 4, 2, 0, 47, 49, 3, 20,
		10, 0, 48, 47, 1, 0, 0, 0, 48, 49, 1, 0, 0, 0, 49, 50, 1, 0, 0, 0, 50,
		51, 5, 7, 0, 0, 51, 62, 1, 0, 0, 0, 52, 53, 5, 8, 0, 0, 53, 54, 3, 8, 4,
		0, 54, 56, 3, 4, 2, 0, 55, 57, 3, 20, 10, 0, 56, 55, 1, 0, 0, 0, 56, 57,
		1, 0, 0, 0, 57, 58, 1, 0, 0, 0, 58, 59, 3, 32, 16, 0, 59, 60, 5, 7, 0,
		0, 60, 62, 1, 0, 0, 0, 61, 37, 1, 0, 0, 0, 61, 44, 1, 0, 0, 0, 61, 52,
		1, 0, 0, 0, 62, 3, 1, 0, 0, 0, 63, 64, 5, 10, 0, 0, 64, 65, 5, 38, 0, 0,
		65, 5, 1, 0, 0, 0, 66, 71, 3, 10, 5, 0, 67, 68, 5, 3, 0, 0, 68, 70, 3,
		10, 5, 0, 69, 67, 1, 0, 0, 0, 70, 73, 1, 0, 0, 0, 71, 69, 1, 0, 0, 0, 71,
		72, 1, 0, 0, 0, 72, 76, 1, 0, 0, 0, 73, 71, 1, 0, 0, 0, 74, 76, 5, 41,
		0, 0, 75, 66, 1, 0, 0, 0, 75, 74, 1, 0, 0, 0, 76, 7, 1, 0, 0, 0, 77, 78,
		3, 10, 5, 0, 78, 79, 5, 3, 0, 0, 79, 81, 1, 0, 0, 0, 80, 77, 1, 0, 0, 0,
		80, 81, 1, 0, 0, 0, 81, 82, 1, 0, 0, 0, 82, 87, 3, 18, 9, 0, 83, 84, 5,
		3, 0, 0, 84, 86, 3, 18, 9, 0, 85, 83, 1, 0, 0, 0, 86, 89, 1, 0, 0, 0, 87,
		85, 1, 0, 0, 0, 87, 88, 1, 0, 0, 0, 88, 9, 1, 0, 0, 0, 89, 87, 1, 0, 0,
		0, 90, 95, 5, 37, 0, 0, 91, 92, 5, 37, 0, 0, 92, 93, 5, 6, 0, 0, 93, 95,
		5, 37, 0, 0, 94, 90, 1, 0, 0, 0, 94, 91, 1, 0, 0, 0, 95, 97, 1, 0, 0, 0,
		96, 98, 3, 12, 6, 0, 97, 96, 1, 0, 0, 0, 97, 98, 1, 0, 0, 0, 98, 104, 1,
		0, 0, 0, 99, 101, 3, 14, 7, 0, 100, 102, 3, 12, 6, 0, 101, 100, 1, 0, 0,
		0, 101, 102, 1, 0, 0, 0, 102, 104, 1, 0, 0, 0, 103, 94, 1, 0, 0, 0, 103,
		99, 1, 0, 0, 0, 104, 11, 1, 0, 0, 0, 105, 106, 5, 13, 0, 0, 106, 107, 5,
		37, 0, 0, 107, 13, 1, 0, 0, 0, 108, 109, 3, 16, 8, 0, 109, 114, 5, 4, 0,
		0, 110, 115, 5, 37, 0, 0, 111, 112, 5, 37, 0, 0, 112, 113, 5, 6, 0, 0,
		113, 115, 5, 37, 0, 0, 114, 110, 1, 0, 0, 0, 114, 111, 1, 0, 0, 0, 115,
		120, 1, 0, 0, 0, 116, 117, 5, 3, 0, 0, 117, 119, 3, 30, 15, 0, 118, 116,
		1, 0, 0, 0, 119, 122, 1, 0, 0, 0, 120, 118, 1, 0, 0, 0, 120, 121, 1, 0,
		0, 0, 121, 123, 1, 0, 0, 0, 122, 120, 1, 0, 0, 0, 123, 124, 5, 5, 0, 0,
		124, 155, 1, 0, 0, 0, 125, 126, 3, 16, 8, 0, 126, 131, 5, 4, 0, 0, 127,
		132, 5, 37, 0, 0, 128, 129, 5, 37, 0, 0, 129, 130, 5, 6, 0, 0, 130, 132,
		5, 37, 0, 0, 131, 127, 1, 0, 0, 0, 131, 128, 1, 0, 0, 0, 132, 137, 1, 0,
		0, 0, 133, 134, 5, 3, 0, 0, 134, 136, 3, 30, 15, 0, 135, 133, 1, 0, 0,
		0, 136, 139, 1, 0, 0, 0, 137, 135, 1, 0, 0, 0, 137, 138, 1, 0, 0, 0, 138,
		140, 1, 0, 0, 0, 139, 137, 1, 0, 0, 0, 140, 141, 5, 5, 0, 0, 141, 155,
		1, 0, 0, 0, 142, 143, 3, 16, 8, 0, 143, 144, 5, 4, 0, 0, 144, 149, 3, 14,
		7, 0, 145, 146, 5, 3, 0, 0, 146, 148, 3, 30, 15, 0, 147, 145, 1, 0, 0,
		0, 148, 151, 1, 0, 0, 0, 149, 147, 1, 0, 0, 0, 149, 150, 1, 0, 0, 0, 150,
		152, 1, 0, 0, 0, 151, 149, 1, 0, 0, 0, 152, 153, 5, 5, 0, 0, 153, 155,
		1, 0, 0, 0, 154, 108, 1, 0, 0, 0, 154, 125, 1, 0, 0, 0, 154, 142, 1, 0,
		0, 0, 155, 15, 1, 0, 0, 0, 156, 157, 5, 37, 0, 0, 157, 17, 1, 0, 0, 0,
		158, 159, 7, 0, 0, 0, 159, 164, 5, 4, 0, 0, 160, 165, 5, 37, 0, 0, 161,
		162, 5, 37, 0, 0, 162, 163, 5, 6, 0, 0, 163, 165, 5, 37, 0, 0, 164, 160,
		1, 0, 0, 0, 164, 161, 1, 0, 0, 0, 165, 166, 1, 0, 0, 0, 166, 168, 5, 5,
		0, 0, 167, 169, 3, 12, 6, 0, 168, 167, 1, 0, 0, 0, 168, 169, 1, 0, 0, 0,
		169, 181, 1, 0, 0, 0, 170, 175, 5, 37, 0, 0, 171, 172, 5, 37, 0, 0, 172,
		173, 5, 6, 0, 0, 173, 175, 5, 37, 0, 0, 174, 170, 1, 0, 0, 0, 174, 171,
		1, 0, 0, 0, 174, 175, 1, 0, 0, 0, 175, 176, 1, 0, 0, 0, 176, 177, 5, 30,
		0, 0, 177, 178, 5, 4, 0, 0, 178, 179, 5, 41, 0, 0, 179, 181, 5, 5, 0, 0,
		180, 158, 1, 0, 0, 0, 180, 174, 1, 0, 0, 0, 181, 19, 1, 0, 0, 0, 182, 183,
		5, 9, 0, 0, 183, 184, 3, 22, 11, 0, 184, 21, 1, 0, 0, 0, 185, 186, 6, 11,
		-1, 0, 186, 198, 3, 24, 12, 0, 187, 190, 3, 26, 13, 0, 188, 190, 3, 24,
		12, 0, 189, 187, 1, 0, 0, 0, 189, 188, 1, 0, 0, 0, 190, 191, 1, 0, 0, 0,
		191, 194, 7, 1, 0, 0, 192, 195, 3, 26, 13, 0, 193, 195, 3, 24, 12, 0, 194,
		192, 1, 0, 0, 0, 194, 193, 1, 0, 0, 0, 195, 198, 1, 0, 0, 0, 196, 198,
		3, 26, 13, 0, 197, 185, 1, 0, 0, 0, 197, 189, 1, 0, 0, 0, 197, 196, 1,
		0, 0, 0, 198, 204, 1, 0, 0, 0, 199, 200, 10, 2, 0, 0, 200, 201, 7, 1, 0,
		0, 201, 203, 3, 22, 11, 3, 202, 199, 1, 0, 0, 0, 203, 206, 1, 0, 0, 0,
		204, 202, 1, 0, 0, 0, 204, 205, 1, 0, 0, 0, 205, 23, 1, 0, 0, 0, 206, 204,
		1, 0, 0, 0, 207, 208, 5, 37, 0, 0, 208, 209, 3, 28, 14, 0, 209, 210, 3,
		30, 15, 0, 210, 218, 1, 0, 0, 0, 211, 212, 5, 37, 0, 0, 212, 213, 5, 6,
		0, 0, 213, 214, 5, 37, 0, 0, 214, 215, 3, 28, 14, 0, 215, 216, 3, 30, 15,
		0, 216, 218, 1, 0, 0, 0, 217, 207, 1, 0, 0, 0, 217, 211, 1, 0, 0, 0, 218,
		25, 1, 0, 0, 0, 219, 220, 5, 4, 0, 0, 220, 221, 3, 22, 11, 0, 221, 222,
		5, 5, 0, 0, 222, 27, 1, 0, 0, 0, 223, 224, 7, 2, 0, 0, 224, 29, 1, 0, 0,
		0, 225, 226, 7, 3, 0, 0, 226, 31, 1, 0, 0, 0, 227, 228, 5, 11, 0, 0, 228,
		231, 3, 10, 5, 0, 229, 231, 1, 0, 0, 0, 230, 227, 1, 0, 0, 0, 230, 229,
		1, 0, 0, 0, 231, 33, 1, 0, 0, 0, 28, 40, 48, 56, 61, 71, 75, 80, 87, 94,
		97, 101, 103, 114, 120, 131, 137, 149, 154, 164, 168, 174, 180, 189, 194,
		197, 204, 217, 230,
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
	SqlParserRULE_sqlQuery           = 0
	SqlParserRULE_selectQuery        = 1
	SqlParserRULE_windowTumbling     = 2
	SqlParserRULE_resultColumns      = 3
	SqlParserRULE_aggregationColumns = 4
	SqlParserRULE_column             = 5
	SqlParserRULE_alias              = 6
	SqlParserRULE_function           = 7
	SqlParserRULE_functionName       = 8
	SqlParserRULE_aggregationColumn  = 9
	SqlParserRULE_whereStatement     = 10
	SqlParserRULE_expr               = 11
	SqlParserRULE_simpleExpr         = 12
	SqlParserRULE_compoundExpr       = 13
	SqlParserRULE_comparisonOperator = 14
	SqlParserRULE_literalValue       = 15
	SqlParserRULE_groupBy            = 16
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
		p.SetState(34)
		p.SelectQuery()
	}
	{
		p.SetState(35)
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

func (s *SelectTumblingGroupByContext) AggregationColumns() IAggregationColumnsContext {
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

	p.SetState(61)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 3, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSelectSimpleContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(37)
			p.Match(SqlParserK_SELECT)
		}
		{
			p.SetState(38)
			p.ResultColumns()
		}
		p.SetState(40)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_WHERE {
			{
				p.SetState(39)
				p.WhereStatement()
			}

		}
		{
			p.SetState(42)
			p.Match(SqlParserEOQ)
		}

	case 2:
		localctx = NewSelectTumblingContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(44)
			p.Match(SqlParserK_SELECT)
		}
		{
			p.SetState(45)
			p.AggregationColumns()
		}
		{
			p.SetState(46)
			p.WindowTumbling()
		}
		p.SetState(48)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_WHERE {
			{
				p.SetState(47)
				p.WhereStatement()
			}

		}
		{
			p.SetState(50)
			p.Match(SqlParserEOQ)
		}

	case 3:
		localctx = NewSelectTumblingGroupByContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(52)
			p.Match(SqlParserK_SELECT)
		}
		{
			p.SetState(53)
			p.AggregationColumns()
		}
		{
			p.SetState(54)
			p.WindowTumbling()
		}
		p.SetState(56)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_WHERE {
			{
				p.SetState(55)
				p.WhereStatement()
			}

		}
		{
			p.SetState(58)
			p.GroupBy()
		}
		{
			p.SetState(59)
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
		p.SetState(63)
		p.Match(SqlParserK_WINDOW_TUMBLING)
	}
	{
		p.SetState(64)
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

	p.SetState(75)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SqlParserIDENTIFIER:
		localctx = NewSelectColumnsContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(66)
			p.Column()
		}
		p.SetState(71)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == SqlParserCOMMA {
			{
				p.SetState(67)
				p.Match(SqlParserCOMMA)
			}
			{
				p.SetState(68)
				p.Column()
			}

			p.SetState(73)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}

	case SqlParserSTAR:
		localctx = NewSelectStarContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(74)
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

type SelectAggregationContext struct {
	*AggregationColumnsContext
}

func NewSelectAggregationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *SelectAggregationContext {
	var p = new(SelectAggregationContext)

	p.AggregationColumnsContext = NewEmptyAggregationColumnsContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AggregationColumnsContext))

	return p
}

func (s *SelectAggregationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *SelectAggregationContext) AllAggregationColumn() []IAggregationColumnContext {
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

func (s *SelectAggregationContext) AggregationColumn(i int) IAggregationColumnContext {
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

func (s *SelectAggregationContext) Column() IColumnContext {
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

func (s *SelectAggregationContext) AllCOMMA() []antlr.TerminalNode {
	return s.GetTokens(SqlParserCOMMA)
}

func (s *SelectAggregationContext) COMMA(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserCOMMA, i)
}

func (s *SelectAggregationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitSelectAggregation(s)

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

	localctx = NewSelectAggregationContext(p, localctx)
	p.EnterOuterAlt(localctx, 1)
	p.SetState(80)
	p.GetErrorHandler().Sync(p)

	if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 6, p.GetParserRuleContext()) == 1 {
		{
			p.SetState(77)
			p.Column()
		}
		{
			p.SetState(78)
			p.Match(SqlParserCOMMA)
		}

	}
	{
		p.SetState(82)
		p.AggregationColumn()
	}
	p.SetState(87)
	p.GetErrorHandler().Sync(p)
	_la = p.GetTokenStream().LA(1)

	for _la == SqlParserCOMMA {
		{
			p.SetState(83)
			p.Match(SqlParserCOMMA)
		}
		{
			p.SetState(84)
			p.AggregationColumn()
		}

		p.SetState(89)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)
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
	p.EnterRule(localctx, 10, SqlParserRULE_column)
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

	p.SetState(103)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 11, p.GetParserRuleContext()) {
	case 1:
		localctx = NewIdentifierColumnContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		p.SetState(94)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 8, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(90)
				p.Match(SqlParserIDENTIFIER)
			}

		case 2:
			{
				p.SetState(91)
				p.Match(SqlParserIDENTIFIER)
			}
			{
				p.SetState(92)
				p.Match(SqlParserDOT)
			}
			{
				p.SetState(93)
				p.Match(SqlParserIDENTIFIER)
			}

		}
		p.SetState(97)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_AS {
			{
				p.SetState(96)
				p.Alias()
			}

		}

	case 2:
		localctx = NewFunctionColumnContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(99)
			p.Function()
		}
		p.SetState(101)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_AS {
			{
				p.SetState(100)
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
	p.EnterRule(localctx, 12, SqlParserRULE_alias)

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
		p.Match(SqlParserK_AS)
	}
	{
		p.SetState(106)
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
	p.EnterRule(localctx, 14, SqlParserRULE_function)
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

	p.SetState(154)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 17, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSimpleFunctionContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(108)
			p.FunctionName()
		}
		{
			p.SetState(109)
			p.Match(SqlParserL_BRACKET)
		}
		p.SetState(114)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 12, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(110)
				p.Match(SqlParserIDENTIFIER)
			}

		case 2:
			{
				p.SetState(111)
				p.Match(SqlParserIDENTIFIER)
			}
			{
				p.SetState(112)
				p.Match(SqlParserDOT)
			}
			{
				p.SetState(113)
				p.Match(SqlParserIDENTIFIER)
			}

		}
		p.SetState(120)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == SqlParserCOMMA {
			{
				p.SetState(116)
				p.Match(SqlParserCOMMA)
			}
			{
				p.SetState(117)
				p.LiteralValue()
			}

			p.SetState(122)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(123)
			p.Match(SqlParserR_BRACKET)
		}

	case 2:
		localctx = NewSimpleFunctionContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(125)
			p.FunctionName()
		}
		{
			p.SetState(126)
			p.Match(SqlParserL_BRACKET)
		}
		p.SetState(131)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 14, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(127)
				p.Match(SqlParserIDENTIFIER)
			}

		case 2:
			{
				p.SetState(128)
				p.Match(SqlParserIDENTIFIER)
			}
			{
				p.SetState(129)
				p.Match(SqlParserDOT)
			}
			{
				p.SetState(130)
				p.Match(SqlParserIDENTIFIER)
			}

		}
		p.SetState(137)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == SqlParserCOMMA {
			{
				p.SetState(133)
				p.Match(SqlParserCOMMA)
			}
			{
				p.SetState(134)
				p.LiteralValue()
			}

			p.SetState(139)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(140)
			p.Match(SqlParserR_BRACKET)
		}

	case 3:
		localctx = NewRecursiveFunctionContext(p, localctx)
		p.EnterOuterAlt(localctx, 3)
		{
			p.SetState(142)
			p.FunctionName()
		}
		{
			p.SetState(143)
			p.Match(SqlParserL_BRACKET)
		}
		{
			p.SetState(144)
			p.Function()
		}
		p.SetState(149)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		for _la == SqlParserCOMMA {
			{
				p.SetState(145)
				p.Match(SqlParserCOMMA)
			}
			{
				p.SetState(146)
				p.LiteralValue()
			}

			p.SetState(151)
			p.GetErrorHandler().Sync(p)
			_la = p.GetTokenStream().LA(1)
		}
		{
			p.SetState(152)
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
	p.EnterRule(localctx, 16, SqlParserRULE_functionName)

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
		p.SetState(156)
		p.Match(SqlParserIDENTIFIER)
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

func (s *AggregationColumnContext) CopyFrom(ctx *AggregationColumnContext) {
	s.BaseParserRuleContext.CopyFrom(ctx.BaseParserRuleContext)
}

func (s *AggregationColumnContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *AggregationColumnContext) ToStringTree(ruleNames []string, recog antlr.Recognizer) string {
	return antlr.TreesStringTree(s, ruleNames, recog)
}

type ColumnAggregationContext struct {
	*AggregationColumnContext
}

func NewColumnAggregationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ColumnAggregationContext {
	var p = new(ColumnAggregationContext)

	p.AggregationColumnContext = NewEmptyAggregationColumnContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AggregationColumnContext))

	return p
}

func (s *ColumnAggregationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnAggregationContext) L_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserL_BRACKET, 0)
}

func (s *ColumnAggregationContext) R_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserR_BRACKET, 0)
}

func (s *ColumnAggregationContext) K_MIN() antlr.TerminalNode {
	return s.GetToken(SqlParserK_MIN, 0)
}

func (s *ColumnAggregationContext) K_MAX() antlr.TerminalNode {
	return s.GetToken(SqlParserK_MAX, 0)
}

func (s *ColumnAggregationContext) K_COUNT() antlr.TerminalNode {
	return s.GetToken(SqlParserK_COUNT, 0)
}

func (s *ColumnAggregationContext) K_AVG() antlr.TerminalNode {
	return s.GetToken(SqlParserK_AVG, 0)
}

func (s *ColumnAggregationContext) K_SUM() antlr.TerminalNode {
	return s.GetToken(SqlParserK_SUM, 0)
}

func (s *ColumnAggregationContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(SqlParserIDENTIFIER)
}

func (s *ColumnAggregationContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, i)
}

func (s *ColumnAggregationContext) DOT() antlr.TerminalNode {
	return s.GetToken(SqlParserDOT, 0)
}

func (s *ColumnAggregationContext) Alias() IAliasContext {
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

func (s *ColumnAggregationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitColumnAggregation(s)

	default:
		return t.VisitChildren(s)
	}
}

type ColumnCountAggregationContext struct {
	*AggregationColumnContext
}

func NewColumnCountAggregationContext(parser antlr.Parser, ctx antlr.ParserRuleContext) *ColumnCountAggregationContext {
	var p = new(ColumnCountAggregationContext)

	p.AggregationColumnContext = NewEmptyAggregationColumnContext()
	p.parser = parser
	p.CopyFrom(ctx.(*AggregationColumnContext))

	return p
}

func (s *ColumnCountAggregationContext) GetRuleContext() antlr.RuleContext {
	return s
}

func (s *ColumnCountAggregationContext) K_COUNT() antlr.TerminalNode {
	return s.GetToken(SqlParserK_COUNT, 0)
}

func (s *ColumnCountAggregationContext) L_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserL_BRACKET, 0)
}

func (s *ColumnCountAggregationContext) STAR() antlr.TerminalNode {
	return s.GetToken(SqlParserSTAR, 0)
}

func (s *ColumnCountAggregationContext) R_BRACKET() antlr.TerminalNode {
	return s.GetToken(SqlParserR_BRACKET, 0)
}

func (s *ColumnCountAggregationContext) AllIDENTIFIER() []antlr.TerminalNode {
	return s.GetTokens(SqlParserIDENTIFIER)
}

func (s *ColumnCountAggregationContext) IDENTIFIER(i int) antlr.TerminalNode {
	return s.GetToken(SqlParserIDENTIFIER, i)
}

func (s *ColumnCountAggregationContext) DOT() antlr.TerminalNode {
	return s.GetToken(SqlParserDOT, 0)
}

func (s *ColumnCountAggregationContext) Accept(visitor antlr.ParseTreeVisitor) interface{} {
	switch t := visitor.(type) {
	case SqlVisitor:
		return t.VisitColumnCountAggregation(s)

	default:
		return t.VisitChildren(s)
	}
}

func (p *SqlParser) AggregationColumn() (localctx IAggregationColumnContext) {
	this := p
	_ = this

	localctx = NewAggregationColumnContext(p, p.GetParserRuleContext(), p.GetState())
	p.EnterRule(localctx, 18, SqlParserRULE_aggregationColumn)
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

	p.SetState(180)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 21, p.GetParserRuleContext()) {
	case 1:
		localctx = NewColumnAggregationContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(158)
			_la = p.GetTokenStream().LA(1)

			if !(((_la-30)&-(0x1f+1)) == 0 && ((1<<uint((_la-30)))&((1<<(SqlParserK_COUNT-30))|(1<<(SqlParserK_SUM-30))|(1<<(SqlParserK_MIN-30))|(1<<(SqlParserK_MAX-30))|(1<<(SqlParserK_AVG-30)))) != 0) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		{
			p.SetState(159)
			p.Match(SqlParserL_BRACKET)
		}
		p.SetState(164)
		p.GetErrorHandler().Sync(p)
		switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 18, p.GetParserRuleContext()) {
		case 1:
			{
				p.SetState(160)
				p.Match(SqlParserIDENTIFIER)
			}

		case 2:
			{
				p.SetState(161)
				p.Match(SqlParserIDENTIFIER)
			}
			{
				p.SetState(162)
				p.Match(SqlParserDOT)
			}
			{
				p.SetState(163)
				p.Match(SqlParserIDENTIFIER)
			}

		}
		{
			p.SetState(166)
			p.Match(SqlParserR_BRACKET)
		}
		p.SetState(168)
		p.GetErrorHandler().Sync(p)
		_la = p.GetTokenStream().LA(1)

		if _la == SqlParserK_AS {
			{
				p.SetState(167)
				p.Alias()
			}

		}

	case 2:
		localctx = NewColumnCountAggregationContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		p.SetState(174)
		p.GetErrorHandler().Sync(p)

		if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext()) == 1 {
			{
				p.SetState(170)
				p.Match(SqlParserIDENTIFIER)
			}

		} else if p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 20, p.GetParserRuleContext()) == 2 {
			{
				p.SetState(171)
				p.Match(SqlParserIDENTIFIER)
			}
			{
				p.SetState(172)
				p.Match(SqlParserDOT)
			}
			{
				p.SetState(173)
				p.Match(SqlParserIDENTIFIER)
			}

		}
		{
			p.SetState(176)
			p.Match(SqlParserK_COUNT)
		}
		{
			p.SetState(177)
			p.Match(SqlParserL_BRACKET)
		}
		{
			p.SetState(178)
			p.Match(SqlParserSTAR)
		}
		{
			p.SetState(179)
			p.Match(SqlParserR_BRACKET)
		}

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
	p.EnterRule(localctx, 20, SqlParserRULE_whereStatement)

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
		p.SetState(182)
		p.Match(SqlParserK_WHERE)
	}
	{
		p.SetState(183)
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
	_startState := 22
	p.EnterRecursionRule(localctx, 22, SqlParserRULE_expr, _p)
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
	p.SetState(197)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 24, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSimpleConditionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx

		{
			p.SetState(186)
			p.SimpleExpr()
		}

	case 2:
		localctx = NewCompoundRecursiveConditionContext(p, localctx)
		p.SetParserRuleContext(localctx)
		_prevctx = localctx
		p.SetState(189)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SqlParserL_BRACKET:
			{
				p.SetState(187)
				p.CompoundExpr()
			}

		case SqlParserIDENTIFIER:
			{
				p.SetState(188)
				p.SimpleExpr()
			}

		default:
			panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
		}
		{
			p.SetState(191)
			_la = p.GetTokenStream().LA(1)

			if !(_la == SqlParserK_AND || _la == SqlParserK_OR) {
				p.GetErrorHandler().RecoverInline(p)
			} else {
				p.GetErrorHandler().ReportMatch(p)
				p.Consume()
			}
		}
		p.SetState(194)
		p.GetErrorHandler().Sync(p)

		switch p.GetTokenStream().LA(1) {
		case SqlParserL_BRACKET:
			{
				p.SetState(192)
				p.CompoundExpr()
			}

		case SqlParserIDENTIFIER:
			{
				p.SetState(193)
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
			p.SetState(196)
			p.CompoundExpr()
		}

	}
	p.GetParserRuleContext().SetStop(p.GetTokenStream().LT(-1))
	p.SetState(204)
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
			p.SetState(199)

			if !(p.Precpred(p.GetParserRuleContext(), 2)) {
				panic(antlr.NewFailedPredicateException(p, "p.Precpred(p.GetParserRuleContext(), 2)", ""))
			}
			{
				p.SetState(200)
				_la = p.GetTokenStream().LA(1)

				if !(_la == SqlParserK_AND || _la == SqlParserK_OR) {
					p.GetErrorHandler().RecoverInline(p)
				} else {
					p.GetErrorHandler().ReportMatch(p)
					p.Consume()
				}
			}
			{
				p.SetState(201)
				p.expr(3)
			}

		}
		p.SetState(206)
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
	p.EnterRule(localctx, 24, SqlParserRULE_simpleExpr)

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

	p.SetState(217)
	p.GetErrorHandler().Sync(p)
	switch p.GetInterpreter().AdaptivePredict(p.GetTokenStream(), 26, p.GetParserRuleContext()) {
	case 1:
		localctx = NewSimpleExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(207)
			p.Match(SqlParserIDENTIFIER)
		}
		{
			p.SetState(208)
			p.ComparisonOperator()
		}
		{
			p.SetState(209)
			p.LiteralValue()
		}

	case 2:
		localctx = NewNestedExpressionContext(p, localctx)
		p.EnterOuterAlt(localctx, 2)
		{
			p.SetState(211)
			p.Match(SqlParserIDENTIFIER)
		}
		{
			p.SetState(212)
			p.Match(SqlParserDOT)
		}
		{
			p.SetState(213)
			p.Match(SqlParserIDENTIFIER)
		}
		{
			p.SetState(214)
			p.ComparisonOperator()
		}
		{
			p.SetState(215)
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
	p.EnterRule(localctx, 26, SqlParserRULE_compoundExpr)

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
		p.SetState(219)
		p.Match(SqlParserL_BRACKET)
	}
	{
		p.SetState(220)
		p.expr(0)
	}
	{
		p.SetState(221)
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
	p.EnterRule(localctx, 28, SqlParserRULE_comparisonOperator)
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
		p.SetState(223)
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
	p.EnterRule(localctx, 30, SqlParserRULE_literalValue)
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
		p.SetState(225)
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
	p.EnterRule(localctx, 32, SqlParserRULE_groupBy)

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

	p.SetState(230)
	p.GetErrorHandler().Sync(p)

	switch p.GetTokenStream().LA(1) {
	case SqlParserK_GROUP_BY:
		p.EnterOuterAlt(localctx, 1)
		{
			p.SetState(227)
			p.Match(SqlParserK_GROUP_BY)
		}
		{
			p.SetState(228)
			p.Column()
		}

	case SqlParserEOQ:
		p.EnterOuterAlt(localctx, 2)

	default:
		panic(antlr.NewNoViableAltException(p, nil, nil, nil, nil, nil))
	}

	return localctx
}

func (p *SqlParser) Sempred(localctx antlr.RuleContext, ruleIndex, predIndex int) bool {
	switch ruleIndex {
	case 11:
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
