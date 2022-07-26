package parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"testing"
)

func TestSelectExpr(t *testing.T) {
	is := antlr.NewInputStream(`SELECT field1, min(field2) WHERE field3 <= 3 and field4 is not null `)
	lexer := NewSqlLexer(is)

	for {
		t := lexer.NextToken()
		if t.GetTokenType() == antlr.TokenEOF {
			break
		}
		fmt.Printf("%s (%q)\n",
			lexer.SymbolicNames[t.GetTokenType()], t.GetText())
	}

}

/*
func TestParser(t *testing.T) {
	is := antlr.NewInputStream("SELEct field1, field2 WHERE field3 > 3 and field4 IS 'test';")

	lexer := NewSqlLexer(is)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	// Create the Parser
	p := NewSqlParser(stream)
	in := make(chan plog.LogRecord)
	out := make(chan plog.LogRecord)
	ctx := context.TODO()

	listener := NewStreamListener(in, out, ctx)

	antlr.ParseTreeWalkerDefault.Walk(&listener, p.SelectQuery())

}
*/

func createSampleLog() plog.LogRecord {
	record := plog.NewLogRecord()
	record.Attributes().Insert("name", pcommon.NewValueString("kwak"))
	record.Attributes().Insert("source", pcommon.NewValueString("some"))
	record.Attributes().InsertBool("alive", true)
	return record

}
