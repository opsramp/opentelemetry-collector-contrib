// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampk8sobjectsprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/opsrampk8sobjectsprocessor"

import (
	"context"

	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type logAttributesProcessor struct {
	logger *zap.Logger
}

// newOpsrampK8sObjectsProcessor returns a processor
func newOpsrampK8sObjectsProcessor(logger *zap.Logger) *logAttributesProcessor {
	return &logAttributesProcessor{
		logger: logger,
	}
}

func (a *logAttributesProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	/*rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		rs := rls.At(i)
		ilss := rs.ScopeLogs()
		resource := rs.Resource()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			logs := ils.LogRecords()
			library := ils.Scope()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)
				if a.skipExpr != nil {
					skip, err := a.skipExpr.Eval(ctx, ottllog.NewTransformContext(lr, library, resource))
					if err != nil {
						return ld, err
					}
					if skip {
						continue
					}
				}

				a.attrProc.Process(ctx, a.logger, lr.Attributes())
			}
		}
	}*/
	return ld, nil
}
