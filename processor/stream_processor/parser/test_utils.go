package parser

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"strconv"
)

func GenerateTestLogs() plog.LogRecordSlice {

	ld := plog.NewLogs()
	sc := ld.ResourceLogs().AppendEmpty().ScopeLogs().AppendEmpty()
	types := []string{"middle", "small", "big"}

	for i := 0; i < 100; i++ {
		record := sc.LogRecords().AppendEmpty()
		record.Attributes().InsertString("name", "Test name "+strconv.Itoa(i))
		record.Attributes().InsertBool("is_alive", i%2 == 0)
		record.Attributes().InsertInt("price", int64(i))
		nested := pcommon.NewValueMap()
		typeIndex := i % 3
		nested.MapVal().InsertString("source", "Source "+strconv.Itoa(i))
		nested.MapVal().InsertString("type", types[typeIndex])
		nested.MapVal().InsertDouble("number", float64(i))
		record.Attributes().Insert("provider", nested)
	}

	return ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
}
