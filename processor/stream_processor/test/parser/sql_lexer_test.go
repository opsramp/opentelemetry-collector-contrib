package parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"testing"
)

func TestSelectExpr(t *testing.T) {
	is := antlr.NewInputStream(`SELEct field1, field2 WHERE field3 > 3 and field4 IS 'test';`)
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
