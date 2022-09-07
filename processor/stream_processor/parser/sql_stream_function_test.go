package parser

import (
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"strings"
	"testing"
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
			appliedFunc:   strings.ToLower,
			expectedCount: 1,
		},
		{
			name:          "recursive 2",
			query:         `SELECT upper(lower(upper(substr(name,2,3)))) WHERE name = 'Test name 50';`,
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
			visitor := NewSqlStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
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
