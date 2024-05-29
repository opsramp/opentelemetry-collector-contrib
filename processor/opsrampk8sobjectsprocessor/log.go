// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampk8sobjectsprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/opsrampk8sobjectsprocessor"

import (
	"context"
	"encoding/json"
	"fmt"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type opsrampK8sObjectsProcessor struct {
	logger *zap.Logger
}

// newOpsrampK8sObjectsProcessor returns a processor
func newOpsrampK8sObjectsProcessor(logger *zap.Logger) *opsrampK8sObjectsProcessor {
	return &opsrampK8sObjectsProcessor{
		logger: logger,
	}
}

func (a *opsrampK8sObjectsProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	out := make(map[string]interface{})
	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		rs := rls.At(i)
		ilss := rs.ScopeLogs()
		//resource := rs.Resource()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			logs := ils.LogRecords()
			//library := ils.Scope()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)

				var watchevent struct {
					Object struct {
						Metadata struct {
							Uid string `json:"uid"`
						} `json:"metadata"`
					} `json:"object"`
					Type string `json:"type"`
				}

				err := json.Unmarshal([]byte(lr.Body().AsString()), &watchevent)
				if err != nil {
					fmt.Println("Error unmarshalling JSON:", err)
				}

				out[watchevent.Object.Metadata.Uid] = lr

				fmt.Printf("After unmarshal opsramp resourceLog #%d, scopeLog #%d, logRecord #%d UID %s \n", i, j, k, watchevent.Object.Metadata.Uid)

			}
		}
	}

	/*
		fmt.Println("Distinct event size\n", len(out))
		for k, v := range out {
			fmt.Println("Hi")
			lr := v.(plog.LogRecord)
			fmt.Printf("UID : %s, Data : %v\n", k, lr.Body().AsString())
		}
		// Sleep required because some times the print logs are not seen. Mosly processor is died before it prints/flushes.
		time.Sleep(1 * time.Second)
	*/

	outld := plog.NewLogs()
	resourceLogs := outld.ResourceLogs()

	for _, v := range out {
		rl := resourceLogs.AppendEmpty()
		sl := rl.ScopeLogs().AppendEmpty()
		logSlice := sl.LogRecords()
		dstLogRecord := logSlice.AppendEmpty()
		srcLogRecord := v.(plog.LogRecord)
		srcLogRecord.CopyTo(dstLogRecord)
	}

	return outld, nil
}
