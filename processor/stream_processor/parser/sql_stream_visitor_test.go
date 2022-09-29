package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/plog"
)

func TestResultColumnsSelectColumns(t *testing.T) {

	tests := []struct {
		name          string
		query         string
		expectedCount int
		err           error
	}{
		{
			name:          "correct value columns",
			query:         `SELECT name, price, is_alive;`,
			expectedCount: 100,
			err:           nil,
		},
		{
			name:          "incorrect value columns",
			query:         `SELECT field, is_alive;`,
			err:           errors.New(`field "field" missed`),
			expectedCount: 0,
		},
		{
			name:          "one value columns",
			query:         `SELECT name ;`,
			err:           nil,
			expectedCount: 100,
		},
		{
			name:          "one incorrect value columns",
			query:         `SELECT field ;`,
			err:           errors.New(`field "field" missed`),
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			visitor := NewSQLStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- GenerateTestLogs()
			var ls plog.LogRecordSlice
			var err error
			select {
			case ls = <-out:
			case err = <-outErr:
			}

			assert.Equal(t, tt.err, err)
			if tt.err == nil {
				assert.Equal(t, tt.expectedCount, ls.Len())
			}
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
			query:         `SELECT price, is_alive ;`,
			expectedAttr:  []string{"price", "is_alive"},
			expectedCount: 100,
		},
		{
			name:          "as clause",
			query:         `SELECT price as NewPrice, is_alive as DEAD ;`,
			expectedAttr:  []string{"NewPrice", "DEAD"},
			expectedCount: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			visitor := NewSQLStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- GenerateTestLogs()

			ls := <-out
			assert.Equal(t, tt.expectedCount, ls.Len())

			for i := 0; i < ls.Len(); i++ {
				rec := ls.At(i)
				assert.Equal(t, len(tt.expectedAttr), rec.Attributes().Len())
				for _, expectedAttr := range tt.expectedAttr {
					_, ok := rec.Attributes().Get(expectedAttr)
					assert.True(t, ok)
				}
			}
		})
	}
}

func TestWhereCondition(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected error
	}{

		{
			name:     "where fields exists ",
			query:    `SELECT name, is_alive WHERE name = 'test' and is_alive = 'true';`,
			expected: nil,
		},

		{
			name:     "where fields missed ",
			query:    `SELECT name, is_alive WHERE name = 'test' and non_exists = 'true';`,
			expected: errors.New(`field "non_exists" missed in log record`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			visitor := NewSQLStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- GenerateTestLogs()
			var err error

			select {
			case <-out:
				return
			case err = <-outErr:
			}

			assert.Equal(t, tt.expected, err)
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
			query:         `SELECT name, is_alive WHERE is_alive = 'true';`,
			expectedAttr:  []string{"name, price"},
			expectedCount: 50,
		},
		{
			name:          "not equal true ",
			query:         `SELECT name, is_alive WHERE is_alive != 'true';`,
			expectedAttr:  []string{"name, is_alive"},
			expectedCount: 50,
		},
		{
			name:          "equal false ",
			query:         `SELECT name, is_alive WHERE is_alive = 'false';`,
			expectedAttr:  []string{"name, is_alive"},
			expectedCount: 50,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			visitor := NewSQLStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- GenerateTestLogs()
			ls := <-out
			assert.Equal(t, tt.expectedCount, ls.Len())

		})
	}

}

func TestRecursiveCondition(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expectedCount int
	}{
		{
			name:          "where and condition ",
			query:         `SELECT name, is_alive WHERE name != 'test' and is_alive = 'true';`,
			expectedCount: 50,
		},
		{
			name:          "where or 1 ",
			query:         `SELECT name, is_alive WHERE name != 'test' or is_alive = 'true';`,
			expectedCount: 100,
		},
		{
			name:          "test or 2",
			query:         `SELECT name WHERE name = 'Test name 10' or is_alive = 'false';`,
			expectedCount: 51,
		},
		{
			name:          "test and 1",
			query:         `SELECT name WHERE name = 'Test name 10' and is_alive = 'false';`,
			expectedCount: 0,
		},
		{
			name:          "test and 2",
			query:         `SELECT name WHERE name = 'Test name 10' and is_alive = 'true';`,
			expectedCount: 1,
		},
		{
			name:          "test and 3",
			query:         `SELECT name WHERE name = 'Test name 10' and price = 10;`,
			expectedCount: 1,
		},
		{
			name:          "test and 4",
			query:         `SELECT name WHERE name = 'Test name 10' and price = 9;`,
			expectedCount: 0,
		},
		{
			name:          "like 1",
			query:         `SELECT name WHERE name like 'Test name 1' and is_alive = 'false';`,
			expectedCount: 6,
		},
		{
			name:          "like 2",
			query:         `SELECT name WHERE name like 'Test name' and price > 70';`,
			expectedCount: 29,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			visitor := NewSQLStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- GenerateTestLogs()
			ls := <-out
			assert.Equal(t, tt.expectedCount, ls.Len())
		})
	}

}

func TestCompoundCondition(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expectedCount int
	}{
		{
			name:          "and 1 ",
			query:         `SELECT name, is_alive WHERE (name != 'test' and is_alive = 'true') and (is_alive = 'true' and price > 50);`,
			expectedCount: 24,
		},
		{
			name:          "and 2 ",
			query:         `SELECT name, is_alive WHERE (name = 'Test name 10' and price > 5) and (is_alive = 'true' and price > 5);`,
			expectedCount: 1,
		},
		{
			name:          "and 3 ",
			query:         `SELECT name, is_alive WHERE (name = 'Test name 10' and price > 5) and (is_alive = 'false' and price > 5);`,
			expectedCount: 0,
		},
		{
			name:          "or 1 ",
			query:         `SELECT name, is_alive WHERE (name != 'test' and is_alive = 'true') or (is_alive = 'true' and price > 50);`,
			expectedCount: 50,
		},
		{
			name:          "or 2 ",
			query:         `SELECT name, is_alive WHERE (name = 'Test name 10' and price > 5 ) or (is_alive = 'true');`,
			expectedCount: 50,
		},
		{
			name:          "or 2 ",
			query:         `SELECT name, is_alive WHERE (name = 'Test name 10' and price > 5 ) or (is_alive = 'false');`,
			expectedCount: 51,
		},
		{
			name:          "like 1  ",
			query:         `SELECT name, is_alive WHERE (name like '2' and price > 5 ) or (is_alive = 'false');`,
			expectedCount: 63,
		},
		{
			name:          "like 2  ",
			query:         `SELECT name, is_alive WHERE (name like '2' and price > 5 ) or (price < 30 or is_alive = 'false');`,
			expectedCount: 72,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			visitor := NewSQLStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- GenerateTestLogs()
			ls := <-out
			assert.Equal(t, tt.expectedCount, ls.Len())
		})
	}

}

func TestComplexCompoundCondition(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expectedCount int
	}{
		{
			name:          "3 compounds ",
			query:         `SELECT * WHERE (name like '2' and price < 3) or (name like '3' and price = 3) or (name = 'Test name 4' );`,
			expectedCount: 3,
		},
		{
			name:          "3 compounds without brackets",
			query:         `SELECT * WHERE (name like '2' and price < 3) or (name like '3' and price = 3) or name = 'Test name 4' ;`,
			expectedCount: 3,
		},
		{
			name:          "3 and compounds",
			query:         `SELECT * WHERE (name like '2' and price < 3) or (name like '3' and price = 3) or (name = 'Test name 4' and price = 4) ;`,
			expectedCount: 3,
		},
		{
			name:          "4 compounds",
			query:         `SELECT * WHERE (name like '2' and price < 3) or (name like '3' and price = 3) or (name = 'Test name 4' and price = 4) or (name = 'Test name 5' and price = 5) ;`,
			expectedCount: 4,
		},
		{
			name:          "5 compounds",
			query:         `SELECT * WHERE (name like '2' and price < 3) or (name like '3' and price = 3) or (name = 'Test name 4' and price = 4) or (name = 'Test name 5' and price = 5) or (name = 'Test name 6' and price = 6) ;`,
			expectedCount: 5,
		},
		{
			name:          "5 compounds",
			query:         `SELECT * WHERE (name like '2' and price < 3) or (name like '3' and price = 3) or (name = 'Test name 4' and price = 4) or (name = 'Test name 5' and price = 5) and (name = 'Test name 6' and price = 6) ;`,
			expectedCount: 0,
		},
		{
			name:          "3 and",
			query:         `SELECT * WHERE name like '2' and price < 3 and is_alive = 'true' ;`,
			expectedCount: 1,
		},
		{
			name:          "3 and",
			query:         `SELECT * WHERE name like '2' and price < 3 and is_alive != 'true' ;`,
			expectedCount: 0,
		},
		{
			name:          "3 and",
			query:         `SELECT * WHERE (name like '2' and price < 3 and is_alive = 'true') ;`,
			expectedCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			visitor := NewSQLStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- GenerateTestLogs()
			ls := <-out
			assert.Equal(t, tt.expectedCount, ls.Len())
		})
	}

}

func TestWriteTestLogsToCSV(t *testing.T) {
	t.Skip()
	f, _ := os.Create("test.csv")
	w := csv.NewWriter(f)
	defer w.Flush()
	for i := 0; i < 100; i++ {
		isAlive := strconv.FormatBool(i%2 == 0)
		err := w.Write([]string{"Test name " + strconv.Itoa(i), isAlive, strconv.Itoa(i)})
		assert.Nil(t, err)

	}

}

func TestWriteTestLogsGroupByToCSV(t *testing.T) {
	t.Skip()
	f, _ := os.Create("testGroupBy.csv")
	w := csv.NewWriter(f)
	defer w.Flush()
	for i := 0; i < 100; i++ {
		isAlive := strconv.FormatBool(i%2 == 0)
		name := strconv.Itoa(i)
		err := w.Write([]string{fmt.Sprint("Test name ", string(name[0])), isAlive, strconv.Itoa(i)})
		assert.Nil(t, err)

	}

}
