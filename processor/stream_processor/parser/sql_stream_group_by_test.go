package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"sync"
	"testing"
	"time"
)

func TestWindowTumblingGroupBy(t *testing.T) {
	in := make(chan plog.LogRecordSlice)
	out := make(chan plog.LogRecordSlice)
	outErr := make(chan error)
	query := "select name, sum(price) window tumbling 30 group by name;"

	visitor := NewSqlStreamVisitor(query, in, out, outErr, zap.NewNop())
	defer visitor.Stop()
	in <- generateGroupByTestLogs()
	ls := <-out
	assert.Equal(t, ls.Len(), 10)
}

func TestWindowTumblingGroupByFlow(t *testing.T) {
	t.Skip()
	in := make(chan plog.LogRecordSlice)
	out := make(chan plog.LogRecordSlice)
	outErr := make(chan error)
	query := "select name, avg(price) window tumbling 3000 group by name;"

	visitor := NewSqlStreamVisitor(query, in, out, outErr, zap.NewNop())
	defer visitor.Stop()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for {
			in <- generateGroupByTestLogs()
			<-time.After(500 * time.Millisecond)
		}
	}()

	go func() {
		for {
			ls := <-out
			for i := 0; i < ls.Len(); i++ {
				ls.At(i).Attributes().Range(func(k string, v pcommon.Value) bool {
					fmt.Printf("%q: %q ", k, v.AsString())
					return true
				})
				fmt.Print("\n")
			}
		}
	}()

	wg.Wait()
}
