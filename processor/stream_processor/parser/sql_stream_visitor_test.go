package parser

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/plog"
)

func TestResultColumnsSelectColumns(t *testing.T) {

	tests := []struct {
		name          string
		query         string
		expectedCount int
	}{
		{
			name:          "correct result columns",
			query:         `SELECT name, price, IsAlive;`,
			expectedCount: 100,
		},
		{
			name:          "incorrect result columns",
			query:         `SELECT field, isAlive;`,
			expectedCount: 0,
		},
		{
			name:          "one result columns",
			query:         `SELECT name ;`,
			expectedCount: 100,
		},
		{
			name:          "one incorrect result columns",
			query:         `SELECT field ;`,
			expectedCount: 0,
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

			assert.Equal(t, tt.expectedCount, logs.Len())
		})
	}
}

func TestResultColumnsSelectColumnsAttributes(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expectedAttr  []string
		expectedCount int
	}{
		{
			name:          "one attributes must be left",
			query:         `SELECT name;`,
			expectedAttr:  []string{"name"},
			expectedCount: 100,
		},
		{
			name:          "two attributes must be left",
			query:         `SELECT price, IsAlive ;`,
			expectedAttr:  []string{"price", "IsAlive"},
			expectedCount: 100,
		},
		{
			name:          "incorrect attributes ",
			query:         `SELECT field, field2 ;`,
			expectedAttr:  []string{},
			expectedCount: 0,
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
			assert.Equal(t, tt.expectedCount, logs.Len())
			for i := 0; i < logs.Len(); i++ {
				rec := logs.At(i)
				assert.Equal(t, len(tt.expectedAttr), rec.Attributes().Len())
				for _, expectedAttr := range tt.expectedAttr {
					_, ok := rec.Attributes().Get(expectedAttr)
					assert.True(t, ok)
				}
			}
		})
	}
}

func TestSimpleCondition(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expectedAttr  []string
		expectedCount int
	}{
		{
			name:          "equal condition ",
			query:         `SELECT name WHERE name = 'Test name 50';`,
			expectedAttr:  []string{"name"},
			expectedCount: 1,
		},
		{
			name:          "greater  condition ",
			query:         `SELECT name WHERE name > 'Test name 50';`,
			expectedAttr:  []string{"name"},
			expectedCount: 53,
		},
		{
			name:          "not equal condition ",
			query:         `SELECT name WHERE name != 'Test name 50';`,
			expectedAttr:  []string{"name"},
			expectedCount: 99,
		},
		{
			name:          "greater numeric ",
			query:         `SELECT name, price WHERE price > 20;`,
			expectedAttr:  []string{"name"},
			expectedCount: 79,
		},
		{
			name:          "greater equal numeric ",
			query:         `SELECT name, price WHERE price >= 20;`,
			expectedAttr:  []string{"name"},
			expectedCount: 80,
		},
		{
			name:          "less  numeric ",
			query:         `SELECT name, price WHERE price < 20;`,
			expectedAttr:  []string{"name", "price"},
			expectedCount: 20,
		},
		{
			name:          "less equal numeric ",
			query:         `SELECT name, price WHERE price <= 20;`,
			expectedAttr:  []string{"name"},
			expectedCount: 21,
		},
		{
			name:          "equal true ",
			query:         `SELECT name, IsAlive WHERE IsAlive = true;`,
			expectedAttr:  []string{"name, price"},
			expectedCount: 50,
		},
		{
			name:          "not equal true ",
			query:         `SELECT name, IsAlive WHERE IsAlive != true;`,
			expectedAttr:  []string{"name, IsAlive"},
			expectedCount: 50,
		},
		{
			name:          "equal false ",
			query:         `SELECT name, IsAlive WHERE IsAlive = false;`,
			expectedAttr:  []string{"name, IsAlive"},
			expectedCount: 50,
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
			assert.Equal(t, tt.expectedCount, logs.Len())

		})
	}

}

func TestRecursiveCondition(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expectedAttr  []string
		expectedCount int
	}{

		{
			name:          "two where condition ",
			query:         `SELECT name, IsAlive WHERE name = 'test' and IsAlive = 'true';`,
			expectedAttr:  []string{"name, IsAlive"},
			expectedCount: 100,
		},
		{
			name:          "two  where condition ",
			query:         `SELECT name WHERE name IS NULL and IsAlive = 'true';`,
			expectedAttr:  []string{"name"},
			expectedCount: 100,
		},
		{
			name:          "compound condition ",
			query:         `SELECT name WHERE (price > 3 ) OR (name = 'test');`,
			expectedAttr:  []string{"name"},
			expectedCount: 100,
		},
		/*
			{
				name:          "compound condition 2 ",
				query:         `SELECT name WHERE (name IS NULL and IsAlive = 'true') OR (name = 'test' AND IsAlive = 'false');`,
				expectedAttr:  []string{"name"},
				expectedCount: 100,
			},*/
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
			assert.Equal(t, tt.expectedCount, logs.Len())

		})
	}

}

func generateTestLogs() plog.LogRecordSlice {

	ld := plog.NewLogs()
	sc := ld.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty()

	for i := 0; i < 100; i++ {
		record := sc.LogRecords().AppendEmpty()
		record.Attributes().InsertString("name", "Test name "+strconv.Itoa(i))
		record.Attributes().InsertBool("IsAlive", i%2 == 0)
		record.Attributes().InsertInt("price", int64(i))
	}

	return ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
}

func TestSmt(t *testing.T) {
	fmt.Println("Test name 50" > "Test name 6")
}
