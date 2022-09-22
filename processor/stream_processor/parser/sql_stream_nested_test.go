package parser

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

func TestWhereNestedCorrectFIelds(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected error
	}{

		{
			name:     "simple nested",
			query:    `SELECT * WHERE provider.type = 'middle';`,
			expected: nil,
		},

		{
			name:     "missed nested error",
			query:    `SELECT* WHERE name.type = 'middle';`,
			expected: errors.New(`field "name" is not nested map`),
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

func TestWhereNestedSimple(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		expected      int
		expectedValue pcommon.Value
	}{
		{
			name:          "simple nested",
			query:         `SELECT * WHERE provider.type = 'middle';`,
			expected:      34,
			expectedValue: pcommon.NewValueString("middle"),
		},
		{
			name:          "simple nested 2",
			query:         `SELECT * WHERE provider.type = 'small';`,
			expected:      33,
			expectedValue: pcommon.NewValueString("small"),
		},
		{
			name:          "simple nested 3",
			query:         `SELECT * WHERE provider.type = 'big';`,
			expected:      33,
			expectedValue: pcommon.NewValueString("big"),
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

			assert.Equal(t, tt.expected, ls.Len())
			for i := 0; i < ls.Len(); i++ {
				val, ok := ls.At(i).Attributes().Get("provider")
				assert.True(t, ok)
				nestedVal, ok := val.MapVal().Get("type")
				assert.True(t, ok)
				assert.Equal(t, tt.expectedValue, nestedVal)
			}
		})
	}

}

func TestWhereNestedCompound(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "nested 1",
			query:    `SELECT * WHERE provider.type like 'l';`,
			expected: 67,
		},
		{
			name:     "nested 2",
			query:    `SELECT * WHERE provider.type like 'bi';`,
			expected: 33,
		},
		{
			name:     "nested 3",
			query:    `SELECT * WHERE provider.type like 'kwa';`,
			expected: 0,
		},
		{
			name:     "nested 4",
			query:    `SELECT * WHERE provider.type = 'middle' or provider.type = 'small';`,
			expected: 67,
		},
		{
			name:     "nested 5",
			query:    `SELECT * WHERE provider.type like 'midd' and provider.type like 'le';`,
			expected: 34,
		},
		{
			name:     "nested 6",
			query:    `SELECT * WHERE provider.type like 'midd' and price < 10;`,
			expected: 4,
		},
		{
			name:     "nested 7",
			query:    `SELECT * WHERE provider.type = 'middle' and price < 30;`,
			expected: 10,
		},
		{
			name:     "nested 8",
			query:    `SELECT * WHERE provider.type = 'middle' and price < 30 and is_alive='true';`,
			expected: 5,
		},
		{
			name:     "nested 9",
			query:    `SELECT * WHERE provider.type = 'middle' and price < 30 and is_alive='true';`,
			expected: 5,
		},
		{
			name:     "nested 10",
			query:    `SELECT * WHERE (provider.type = 'middle' and price < 30) or  (provider.type = 'middle' and price > 90);`,
			expected: 13,
		},
		{
			name:     "nested 11",
			query:    `SELECT * WHERE (provider.type = 'middle' and price < 30) and  (provider.type = 'middle' and price > 90);`,
			expected: 0,
		},
		{
			name:     "nested 12",
			query:    `SELECT * WHERE (provider.type = 'middle' and price < 10) or (provider.type = 'middle' and price > 90) or (provider.type = 'big' and price > 90);`,
			expected: 10,
		},
		{
			name:     "nested 13",
			query:    `SELECT * WHERE (provider.type = 'middle' and price < 10 and is_alive='true');`,
			expected: 2,
		},
		{
			name:     "nested 14",
			query:    `SELECT * WHERE (provider.type = 'middle' and price < 10 and is_alive='true') or (provider.type = 'small' and price < 10 and is_alive='false');`,
			expected: 4,
		},
		{
			name:     "nested 15",
			query:    `SELECT * WHERE provider.type = 'middle' and provider.number > 50 and price > 50;`,
			expected: 17,
		},
		{
			name:     "nested 16",
			query:    `SELECT * WHERE provider.type = 'big' and provider.number > 50 and price > 50;`,
			expected: 16,
		},
		{
			name:     "nested 17",
			query:    `SELECT * WHERE provider.type = 'big' and provider.number > 50 or price > 80 ;`,
			expected: 29,
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

			assert.Equal(t, tt.expected, ls.Len())

		})
	}

}

func TestSelectNested(t *testing.T) {
	tests := []struct {
		name, query string

		expected       int
		expectedNested func() pcommon.Value
	}{
		{
			name:     "nested 2 fields",
			query:    `SELECT provider.source, provider.type where provider.number = 2;`,
			expected: 1,
			expectedNested: func() pcommon.Value {
				val := pcommon.NewValueMap()
				val.MapVal().Insert("source", pcommon.NewValueString("Source 2"))
				val.MapVal().Insert("type", pcommon.NewValueString("big"))
				return val
			},
		},
		{
			name:     "nested 3 fields",
			query:    `SELECT provider.source, provider.type, provider.number where provider.number = 2;`,
			expected: 1,
			expectedNested: func() pcommon.Value {
				val := pcommon.NewValueMap()
				val.MapVal().Insert("source", pcommon.NewValueString("Source 2"))
				val.MapVal().Insert("type", pcommon.NewValueString("big"))
				val.MapVal().Insert("number", pcommon.NewValueDouble(2))
				return val
			},
		},
		{
			name:     "nested all fields",
			query:    `SELECT provider where provider.number = 2;`,
			expected: 1,
			expectedNested: func() pcommon.Value {
				val := pcommon.NewValueMap()
				val.MapVal().Insert("source", pcommon.NewValueString("Source 2"))
				val.MapVal().Insert("type", pcommon.NewValueString("big"))
				val.MapVal().Insert("number", pcommon.NewValueDouble(2))
				return val
			},
		},
		{
			name:     "nested all unnecessary fields",
			query:    `SELECT provider.source, provider.type, provider.number, provider.source, name where provider.number = 2;`,
			expected: 1,
			expectedNested: func() pcommon.Value {
				val := pcommon.NewValueMap()
				val.MapVal().Insert("source", pcommon.NewValueString("Source 2"))
				val.MapVal().Insert("type", pcommon.NewValueString("big"))
				val.MapVal().Insert("number", pcommon.NewValueDouble(2))
				return val
			},
		},
		{
			name:     "nested without where",
			query:    `SELECT provider.source;`,
			expected: 100,
			expectedNested: func() pcommon.Value {
				val := pcommon.NewValueMap()
				val.MapVal().Insert("source", pcommon.NewValueString("Source 0"))
				return val
			},
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
			res, ok := ls.At(0).Attributes().Get("provider")
			assert.True(t, ok)
			assert.True(t, reflect.DeepEqual(tt.expectedNested(), res))

		})
	}

}

func TestSelectNestedWIthAliases(t *testing.T) {
	tests := []struct {
		name, query, key string
		expected         int
		expectedNested   func() pcommon.Map
	}{
		{
			name:     "nested 2 with alias",
			query:    `SELECT provider.source as NewSource, provider.type as NewType where provider.number = 2;`,
			expected: 1,
			expectedNested: func() pcommon.Map {
				val := pcommon.NewMap()
				val.Insert("NewSource", pcommon.NewValueString("Source 2"))
				val.Insert("NewType", pcommon.NewValueString("big"))
				return val
			},
		},
		{
			name:     "nested 2 with alias and nested",
			query:    `SELECT provider.source as NewSource, provider.type as NewType, provider.number where provider.number = 2;`,
			expected: 1,
			expectedNested: func() pcommon.Map {
				val := pcommon.NewMap()
				val.Insert("NewSource", pcommon.NewValueString("Source 2"))
				val.Insert("NewType", pcommon.NewValueString("big"))
				mapValue := pcommon.NewValueMap()
				mapValue.MapVal().Insert("number", pcommon.NewValueDouble(2))
				val.Insert("provider", mapValue)
				return val
			},
		},
		{
			name:     "nested 2 with alias and nested ans simple",
			query:    `SELECT provider.source as NewSource, provider.type as NewType, price, provider.number where provider.number = 2;`,
			expected: 1,
			expectedNested: func() pcommon.Map {
				val := pcommon.NewMap()
				val.Insert("NewSource", pcommon.NewValueString("Source 2"))
				val.Insert("NewType", pcommon.NewValueString("big"))
				val.Insert("price", pcommon.NewValueInt(2))
				mapValue := pcommon.NewValueMap()
				mapValue.MapVal().Insert("number", pcommon.NewValueDouble(2))
				val.Insert("provider", mapValue)
				return val
			},
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
			res := ls.At(0).Attributes()
			assert.True(t, reflect.DeepEqual(tt.expectedNested(), res))

		})
	}

}
