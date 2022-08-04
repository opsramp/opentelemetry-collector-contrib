package main

import (
	"bufio"
	"fmt"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/stream_processor/parser"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("-> ")
		query, _ := reader.ReadString('\n')
		if strings.TrimRight(strings.ToLower(query), "\r\n\\") == "quit" {
			break
		}

		logs := generateTestLogs()
		in := make(chan plog.LogRecordSlice)
		out := make(chan plog.LogRecordSlice)
		outErr := make(chan error)
		visitor := parser.NewSqlStreamVisitor(query, in, out, outErr, zap.NewNop())
		in <- logs

		select {
		case ld := <-out:
			for i := 0; i < ld.Len(); i++ {
				name, ok := ld.At(i).Attributes().Get("name")
				if ok {
					fmt.Printf("name: %q ", name.AsString())
				}

				isAlive, ok := ld.At(i).Attributes().Get("IsAlive")
				if ok {
					fmt.Printf("IsAlive: %q ", isAlive.AsString())
				}
				price, ok := ld.At(i).Attributes().Get("price")
				if ok {
					fmt.Printf("price:%q ", price.AsString())
				}
				fmt.Print("\n")
			}

		case err := <-outErr:
			fmt.Println(err)
			continue
		}
		visitor.Stop()

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
