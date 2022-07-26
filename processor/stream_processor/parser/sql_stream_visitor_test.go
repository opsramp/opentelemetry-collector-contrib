package parser

import (
	"fmt"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/plog"
	"testing"
)

func TestNewResultColumnsContext(t *testing.T) {

	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "correct result columns",
			query:    `SELECT name, price, IsAlive WHERE name = 'test' and isAlive = 'true';`,
			expected: 100,
		},
		{
			name:     "incorrect result columns",
			query:    `SELECT field, isAlive WHERE name = 'test' and isAlive = 'true';`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := antlr.NewInputStream(tt.query)

			lexer := NewSqlLexer(is)
			stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

			parser := NewSqlParser(stream)

			logs := generateTestLogs()
			visitor := NewSqlStreamVisitor(logs)
			visitor.Visit(parser.SqlQuery())

			logs, _ = visitor.GetResult()

			assert.Equal(t, tt.expected, logs.Len())
		})
	}

}

func generateTestLogs() plog.LogRecordSlice {

	ld := plog.NewLogs()
	sc := ld.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty()

	for i := 0; i < 100; i++ {
		record := sc.LogRecords().AppendEmpty()
		record.Attributes().InsertString("name", fmt.Sprintf("Test log # %q", i))
		record.Attributes().InsertBool("IsAlive", i%2 == 0)
		record.Attributes().InsertInt("price", int64(i))
	}

	return ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
}
