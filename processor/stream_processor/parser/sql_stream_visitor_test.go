package parser

import (
	"encoding/csv"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strconv"
	"testing"

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
			name:          "correct value columns",
			query:         `SELECT name, price, IsAlive;`,
			expectedCount: 100,
		},
		{
			name:          "incorrect value columns",
			query:         `SELECT field, isAlive;`,
			expectedCount: 0,
		},
		{
			name:          "one value columns",
			query:         `SELECT name ;`,
			expectedCount: 100,
		},
		{
			name:          "one incorrect value columns",
			query:         `SELECT field ;`,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			visitor := NewSqlStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- generateTestLogs()
			ls := <-out
			assert.Equal(t, tt.expectedCount, ls.Len())
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

			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			visitor := NewSqlStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- generateTestLogs()

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
			query:    `SELECT name, IsAlive WHERE name = 'test' and IsAlive = 'true';`,
			expected: nil,
		},

		{
			name:     "where fields missed ",
			query:    `SELECT name, IsAlive WHERE name = 'test' and non_exists = 'true';`,
			expected: errors.New(`field "non_exists" missed in log record`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			visitor := NewSqlStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- generateTestLogs()
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
			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			visitor := NewSqlStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- generateTestLogs()
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
			query:         `SELECT name, IsAlive WHERE name != 'test' and IsAlive = 'true';`,
			expectedCount: 50,
		},
		{
			name:          "where or 1 ",
			query:         `SELECT name, IsAlive WHERE name != 'test' or IsAlive = 'true';`,
			expectedCount: 100,
		},
		{
			name:          "test or 2",
			query:         `SELECT name WHERE name = 'Test name 10' or IsAlive = 'false';`,
			expectedCount: 51,
		},
		{
			name:          "test and 1",
			query:         `SELECT name WHERE name = 'Test name 10' and IsAlive = 'false';`,
			expectedCount: 0,
		},
		{
			name:          "test and 2",
			query:         `SELECT name WHERE name = 'Test name 10' and IsAlive = 'true';`,
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
			query:         `SELECT name WHERE name like 'Test name 1' and IsAlive = 'false';`,
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
			visitor := NewSqlStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- generateTestLogs()
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
			query:         `SELECT name, IsAlive WHERE (name != 'test' and IsAlive = 'true') and (IsAlive = 'true' and price > 50);`,
			expectedCount: 24,
		},
		{
			name:          "and 2 ",
			query:         `SELECT name, IsAlive WHERE (name = 'Test name 10' and price > 5) and (IsAlive = 'true' and price > 5);`,
			expectedCount: 1,
		},
		{
			name:          "and 3 ",
			query:         `SELECT name, IsAlive WHERE (name = 'Test name 10' and price > 5) and (IsAlive = 'false' and price > 5);`,
			expectedCount: 0,
		},
		{
			name:          "or 1 ",
			query:         `SELECT name, IsAlive WHERE (name != 'test' and IsAlive = 'true') or (IsAlive = 'true' and price > 50);`,
			expectedCount: 50,
		},
		{
			name:          "or 2 ",
			query:         `SELECT name, IsAlive WHERE (name = 'Test name 10' and price > 5 ) or (IsAlive = 'true');`,
			expectedCount: 50,
		},
		{
			name:          "or 2 ",
			query:         `SELECT name, IsAlive WHERE (name = 'Test name 10' and price > 5 ) or (IsAlive = 'false');`,
			expectedCount: 51,
		},
		{
			name:          "like 1  ",
			query:         `SELECT name, IsAlive WHERE (name like '2' and price > 5 ) or (IsAlive = 'false');`,
			expectedCount: 63,
		},
		{
			name:          "like 2  ",
			query:         `SELECT name, IsAlive WHERE (name like '2' and price > 5 ) or (price < 30 or IsAlive = 'false');`,
			expectedCount: 72,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			visitor := NewSqlStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			in <- generateTestLogs()
			ls := <-out
			assert.Equal(t, tt.expectedCount, ls.Len())
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

func TestWriteTestLogsToCSV(t *testing.T) {
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
