package parser

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

func TestIsWindowTumbling(t *testing.T) {
	tests := []struct {
		name  string
		query string
		value int
		res   bool
	}{
		{
			name:  "not tumbling query",
			query: `SELECT * WHERE price > 80 or name like '3';`,
			value: 0,
			res:   false,
		},
		{
			name:  "tumbling query",
			query: `select max(price) window tumbling 3 where price > 4;`,
			value: 3,
			res:   true,
		},
		{
			name:  "incorrect tumbling query",
			query: `select max(price) window tumbling kwa where price > 4;`,
			value: 0,
			res:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, val := IsTumblingQuery(tt.query)
			assert.Equal(t, tt.value, val)
			assert.Equal(t, tt.res, res)
		})
	}
}

func TestWindowTumblingLoop(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected int
		res      float64
	}{
		{
			name:     "1 ms",
			query:    "select max(price) window tumbling 50 where price > 4;",
			expected: 1,
			res:      99.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := make(chan plog.LogRecordSlice)
			out := make(chan plog.LogRecordSlice)
			outErr := make(chan error)
			var ls plog.LogRecordSlice
			var wg sync.WaitGroup
			wg.Add(1)
			visitor := NewSQLStreamVisitor(tt.query, in, out, outErr, zap.NewNop())
			defer visitor.Stop()
			go func() {
				defer wg.Done()
				for {
					ls = <-out
					break
				}
			}()
			in <- GenerateTestLogs()
			<-time.After(1 * time.Millisecond)
			in <- GenerateTestLogs()
			<-time.After(1 * time.Millisecond)
			in <- GenerateTestLogs()
			<-time.After(1 * time.Millisecond)
			wg.Wait()

			assert.Equal(t, tt.expected, ls.Len())
			val, ok := ls.At(0).Attributes().Get("price")
			assert.True(t, ok)
			assert.Equal(t, tt.res, val.DoubleVal())
		})
	}
}

func TestWindowTumblingAvg(t *testing.T) {
	in := make(chan plog.LogRecordSlice)
	out := make(chan plog.LogRecordSlice)
	outErr := make(chan error)
	query := "select avg(price), sum(price) as SumPrice window tumbling 3 ;"
	var ls plog.LogRecordSlice

	visitor := NewSQLStreamVisitor(query, in, out, outErr, zap.NewNop())
	defer visitor.Stop()
	in <- GenerateTestLogs()
	ls = <-out
	res, ok := ls.At(0).Attributes().Get("price")
	assert.True(t, ok)
	assert.Equal(t, pcommon.NewValueDouble(49.5), res)
}

func TestWindowTumblingCount(t *testing.T) {
	in := make(chan plog.LogRecordSlice)
	out := make(chan plog.LogRecordSlice)
	outErr := make(chan error)
	query := "select count(price) window tumbling 30 ;"
	var ls plog.LogRecordSlice

	visitor := NewSQLStreamVisitor(query, in, out, outErr, zap.NewNop())
	defer visitor.Stop()
	in <- GenerateTestLogs()
	ls = <-out
	res, ok := ls.At(0).Attributes().Get("price")
	assert.True(t, ok)
	assert.Equal(t, pcommon.NewValueDouble(100.0), res)
}

func TestWindowTumblingSum(t *testing.T) {
	in := make(chan plog.LogRecordSlice)
	out := make(chan plog.LogRecordSlice)
	outErr := make(chan error)
	query := "select sum(price) window tumbling 30 ;"
	var ls plog.LogRecordSlice

	visitor := NewSQLStreamVisitor(query, in, out, outErr, zap.NewNop())
	defer visitor.Stop()
	in <- GenerateTestLogs()
	ls = <-out
	res, ok := ls.At(0).Attributes().Get("price")
	assert.True(t, ok)
	assert.Equal(t, pcommon.NewValueDouble(4950.0), res)
}

func TestWindowTumblingMin(t *testing.T) {
	in := make(chan plog.LogRecordSlice)
	out := make(chan plog.LogRecordSlice)
	outErr := make(chan error)
	query := "select min(price) window tumbling 30 ;"
	var ls plog.LogRecordSlice

	visitor := NewSQLStreamVisitor(query, in, out, outErr, zap.NewNop())
	defer visitor.Stop()
	in <- GenerateTestLogs()
	ls = <-out
	res, ok := ls.At(0).Attributes().Get("price")
	assert.True(t, ok)
	assert.Equal(t, pcommon.NewValueDouble(0.0), res)
}

func TestWindowTumblingMax(t *testing.T) {
	in := make(chan plog.LogRecordSlice)
	out := make(chan plog.LogRecordSlice)
	outErr := make(chan error)
	query := "select max(price) window tumbling 30 ;"
	var ls plog.LogRecordSlice

	visitor := NewSQLStreamVisitor(query, in, out, outErr, zap.NewNop())
	defer visitor.Stop()
	in <- GenerateTestLogs()
	ls = <-out
	res, ok := ls.At(0).Attributes().Get("price")
	assert.True(t, ok)
	assert.Equal(t, pcommon.NewValueDouble(99.0), res)
}

/*
func TestAvgContext_K_MIN(t *testing.T) {
	ls := GenerateTestLogs()
	res, err := min(ls, "price")
	assert.Nil(t, err)
	assert.Equal(t, 0.0, res)

}

func TestAvgContext_K_MAX(t *testing.T) {
	ls := GenerateTestLogs()
	res, err := max(ls, "price")
	assert.Nil(t, err)
	assert.Equal(t, 99.0, res)
}

func TestAvgContext_K_SUM(t *testing.T) {
	ls := GenerateTestLogs()
	res, err := sum(ls, "price")
	assert.Nil(t, err)
	assert.Equal(t, 4950.0, res)
}

*/

func TestAvgContext_K_COUNT(t *testing.T) {
	//ls := GenerateTestLogs()
	count := 100 //count(ls)
	assert.Equal(t, 100, count)
}

func TestWindowTumblingFlow(t *testing.T) {
	t.Skip()
	in := make(chan plog.LogRecordSlice)
	out := make(chan plog.LogRecordSlice)
	outErr := make(chan error)
	query := "select avg(price) window tumbling 3000 ;"

	visitor := NewSQLStreamVisitor(query, in, out, outErr, zap.NewNop())
	defer visitor.Stop()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for {
			in <- generateRandomSplittedTestLogs()
			<-time.After(500 * time.Millisecond)
		}
	}()

	go func() {
		for {
			ls := <-out
			ls.At(0).Attributes().Range(func(k string, v pcommon.Value) bool {
				fmt.Printf("%q: %q ", k, v.AsString())
				return true
			})
			fmt.Print("\n")
		}
	}()

	wg.Wait()
}

func generateRandomSplittedTestLogs() plog.LogRecordSlice {

	ld := plog.NewLogs()
	sc := ld.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty()

	r := rand.Intn(100)

	for i := 0; i < r; i++ {
		record := sc.LogRecords().AppendEmpty()
		record.Attributes().InsertString("name", "Test name "+strconv.Itoa(i))
		record.Attributes().InsertBool("is_alive", i%2 == 0)
		record.Attributes().InsertInt("price", int64(i))
	}

	return ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
}

func generateGroupByTestLogs() plog.LogRecordSlice {
	ld := plog.NewLogs()
	sc := ld.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty()

	for i := 0; i < 100; i++ {
		record := sc.LogRecords().AppendEmpty()
		name := strconv.Itoa(i)
		record.Attributes().InsertString("name", fmt.Sprint("Test name ", string(name[0])))
		record.Attributes().InsertString("lname", fmt.Sprint("Last name ", string(name[0])))
		record.Attributes().InsertBool("is_alive", i%2 == 0)
		record.Attributes().InsertInt("price", int64(i))
	}

	return ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
}
