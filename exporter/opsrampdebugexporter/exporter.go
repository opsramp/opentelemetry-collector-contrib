// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

// NOTE: If you are making changes to this file, consider whether you want to make similar changes
// to the Logging exporter in /exporter/internal/common/logging_exporter.go, which has similar logic.
// This is especially important for security issues.

package opsrampdebugexporter // import "go.opentelemetry.io/collector/exporter/debugexporter"

import (
	"context"
	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opsrampdebugexporter/internal/otlptext"

	"go.uber.org/zap"

	"go.opentelemetry.io/collector/config/configtelemetry"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

var (
	//LogsOpsRampChannel = make(chan plog.Logs, 100)
	LogsOpsRampChannel = make(chan plog.Logs, 1000)
)

type debugExporter struct {
	verbosity        configtelemetry.Level
	logger           *zap.Logger
	logsMarshaler    plog.Marshaler
	metricsMarshaler pmetric.Marshaler
	tracesMarshaler  ptrace.Marshaler
}

func newDebugExporter(logger *zap.Logger, verbosity configtelemetry.Level) *debugExporter {
	return &debugExporter{
		verbosity:        verbosity,
		logger:           logger,
		logsMarshaler:    otlptext.NewTextLogsMarshaler(),
		metricsMarshaler: otlptext.NewTextMetricsMarshaler(),
		tracesMarshaler:  otlptext.NewTextTracesMarshaler(),
	}
}

func (s *debugExporter) pushTraces(_ context.Context, td ptrace.Traces) error {
	s.logger.Info("TracesExporter",
		zap.Int("resource spans", td.ResourceSpans().Len()),
		zap.Int("spans", td.SpanCount()))
	if s.verbosity != configtelemetry.LevelDetailed {
		return nil
	}

	buf, err := s.tracesMarshaler.MarshalTraces(td)
	if err != nil {
		return err
	}
	s.logger.Info(string(buf))
	return nil
}

func (s *debugExporter) pushMetrics(_ context.Context, md pmetric.Metrics) error {
	s.logger.Info("MetricsExporter",
		zap.Int("resource metrics", md.ResourceMetrics().Len()),
		zap.Int("metrics", md.MetricCount()),
		zap.Int("data points", md.DataPointCount()))
	if s.verbosity != configtelemetry.LevelDetailed {
		return nil
	}

	buf, err := s.metricsMarshaler.MarshalMetrics(md)
	if err != nil {
		return err
	}
	s.logger.Info(string(buf))
	return nil
}

func (s *debugExporter) pushLogs(_ context.Context, ld plog.Logs) error {
	s.logger.Info("LogsExporter",
		zap.Int("resource logs", ld.ResourceLogs().Len()),
		zap.Int("log records", ld.LogRecordCount()))
	select {
	case LogsOpsRampChannel <- ld:
		s.logger.Info("#######LogsExporter: Successfully sent to channel")
	default:
		s.logger.Info("#######LogsExporter: failed sent to channel")
	}
	if s.verbosity != configtelemetry.LevelDetailed {
		return nil
	}

	buf, err := s.logsMarshaler.MarshalLogs(ld)
	if err != nil {
		return err
	}
	s.logger.Info(string(buf))
	return nil
}
