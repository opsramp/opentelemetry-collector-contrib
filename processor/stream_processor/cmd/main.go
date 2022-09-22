package main

import (
	"bufio"
	"fmt"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/stream_processor/parser"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
	"os"
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

		in := make(chan plog.LogRecordSlice)
		out := make(chan plog.LogRecordSlice)
		outErr := make(chan error)
		visitor := parser.NewSQLStreamVisitor(query, in, out, outErr, zap.NewNop())

		in <- parser.GenerateTestLogs()

		select {
		case ls := <-out:
			for i := 0; i < ls.Len(); i++ {
				ls.At(i).Attributes().Range(func(k string, v pcommon.Value) bool {
					fmt.Printf("%q: %q ", k, v.AsString())
					return true
				})

				fmt.Print("\n")
			}

		case err := <-outErr:
			fmt.Println(err)
			continue
		}
		visitor.Stop()

	}

}
