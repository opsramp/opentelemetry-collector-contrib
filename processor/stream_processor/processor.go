package stream_processor

import (
	"context"
	"fmt"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/stream_processor/parser"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type sqlStreamProcessor struct {
	logger       *zap.Logger
	config       *Config
	nextConsumer consumer.Logs
	in, out      chan plog.LogRecordSlice
	outErr       chan error
	processor    *parser.SQLStreamVisitor
}

func newSqlStreamProcessor(next consumer.Logs, logger *zap.Logger, cfg *Config) *sqlStreamProcessor {
	in := make(chan plog.LogRecordSlice)
	out := make(chan plog.LogRecordSlice)
	outErr := make(chan error)
	processor := parser.NewSQLStreamVisitor(cfg.Query, in, out, outErr, logger)

	return &sqlStreamProcessor{
		logger:       logger,
		config:       cfg,
		nextConsumer: next,
		in:           in,
		out:          out,
		outErr:       outErr,
		processor:    processor,
	}
}

func (sp *sqlStreamProcessor) processorLoop() {
	for {
		select {
		case ls := <-sp.out:
			printLogs(ls)
		case err := <-sp.outErr:
			sp.logger.Error("error", zap.Error(err))
			sp.processor.Stop()
		}
	}
}

func printLogs(ls plog.LogRecordSlice) {
	if ls.Len() == 0 {
		return
	}
	for i := 0; i < ls.Len(); i++ {
		fmt.Println("EVENT: ")
		ls.At(i).Attributes().Range(func(k string, v pcommon.Value) bool {
			if v.Type() == pcommon.ValueTypeMap {
				fmt.Printf(" Nested key %q \n", k)
				v.MapVal().Range(func(key string, val pcommon.Value) bool {
					fmt.Printf("Key: %q: Value: %q \n", key, val.AsString())
					return true
				})
				fmt.Println()
				return true
			}
			fmt.Printf("%q: %q \n", k, v.AsString())
			return true
		})
		fmt.Print("\n")
	}
	//fmt.Println("END OF THIS BATCH")
}

func (sp *sqlStreamProcessor) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	//<-time.After(1 * time.Second)
	sp.in <- ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords()
	return nil
}

func (sp *sqlStreamProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// Start is invoked during service startup.
func (sp *sqlStreamProcessor) Start(context.Context, component.Host) error {
	go sp.processorLoop()
	return nil
}

// Shutdown is invoked during service shutdown.
func (sp *sqlStreamProcessor) Shutdown(context.Context) error {
	return nil
}
