package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

func TestScalarFunctions(t *testing.T) {

	tests := []struct {
		name         string
		query        string
		appliedFunc  func(string) string
		expectedAttr []string

		expectedCount int
	}{
		{
			name:          "upper",
			query:         `SELECT lower(name) WHERE name = 'Test name 50';`,
			expectedAttr:  []string{"name"},
			appliedFunc:   strings.ToLower,
			expectedCount: 1,
		},
		{
			name:          "substr",
			query:         `SELECT substr(name,2,3) WHERE name = 'Test name 50';`,
			expectedAttr:  []string{"name"},
			appliedFunc:   strings.ToLower,
			expectedCount: 1,
		},
		{
			name:          "recursive 1",
			query:         `SELECT upper(lower(name)) WHERE name = 'Test name 50';`,
			expectedAttr:  []string{"name"},
			appliedFunc:   strings.ToUpper,
			expectedCount: 1,
		},
		{
			name:          "recursive 2",
			query:         `SELECT lower(upper(lower(upper(substr(name,2,3))))) WHERE name = 'Test name 50';`,
			expectedAttr:  []string{"name"},
			appliedFunc:   strings.ToLower,
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
			for i := 0; i < ls.Len(); i++ {
				rec := ls.At(i)
				val, ok := rec.Attributes().Get("name")
				assert.True(t, ok)
				assert.Equal(t, tt.appliedFunc(val.AsString()), val.AsString())
			}

		})
	}
}

func TestRecursiveScalarFunctions(t *testing.T) {

	tests := []struct {
		name, query, attr string
		expectedAttr      map[string]string
		expectedResult    int
	}{
		{
			name:         "recursive 1",
			query:        `SELECT substr(name,2,3) WHERE name = 'Test name 50';`,
			expectedAttr: map[string]string{"name": "st "},
		},
		{
			name:         "recursive 2",
			query:        `SELECT upper(lower(upper(substr(name,2,3)))) WHERE name = 'Test name 50';`,
			expectedAttr: map[string]string{"name": "ST "},
		},
		{
			name:         "recursive with alias",
			query:        `SELECT upper(substr(name,2,3)), upper(name) as UpperName WHERE name = 'Test name 50';`,
			expectedAttr: map[string]string{"name": "ST ", "UpperName": "TEST NAME 50"},
		},
		{
			name:         "recursive with alias",
			query:        `SELECT upper(substr(name,2,3)), upper(name) as UpperName, name as JustName  WHERE name = 'Test name 50';`,
			expectedAttr: map[string]string{"name": "ST ", "UpperName": "TEST NAME 50", "JustName": "Test name 50"},
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
			for k, v := range tt.expectedAttr {
				val, ok := res.Get(k)
				assert.True(t, ok)
				assert.Equal(t, v, val.AsString())
			}

		})
	}
}
