package parser

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

func TestWindowTumblingGroupBy(t *testing.T) {
	in := make(chan plog.LogRecordSlice)
	out := make(chan plog.LogRecordSlice)
	outErr := make(chan error)
	query := "select name as TEST, sum(price), avg(price) as AVGPrice window tumbling 30 group by name ;"

	visitor := NewSQLStreamVisitor(query, in, out, outErr, zap.NewNop())
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
	query := "select name, max(price) window tumbling 3000 group by name;"

	visitor := NewSQLStreamVisitor(query, in, out, outErr, zap.NewNop())
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
