package opsrampflowexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opsrampflowexporter"

import (
	"context"
	"opsrampflowexporter/internal/metadata"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// NewFactory creates a factory for OpenSearch exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		metadata.Type,
		newDefaultConfig,
		exporter.WithLogs(createLogsExporter, metadata.LogsStability),
		exporter.WithMetrics(createMetricsExporter, metadata.MetricsStability),
	)
}

func newDefaultConfig() component.Config {
	return &Config{
		Endpoint:  "",
		AuthToken: "",
	}
}

func createLogsExporter(ctx context.Context,
	set exporter.Settings,
	cfg component.Config) (exporter.Logs, error) {
	c := cfg.(*Config)
	fe, e := newFlowExporter(c, set)
	if e != nil {
		return nil, e
	}

	return exporterhelper.NewLogsExporter(ctx, set, cfg,
		fe.consumeLogs,
		exporterhelper.WithStart(fe.Start),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}),
		exporterhelper.WithRetry(c.BackOffConfig))
}

func createMetricsExporter(ctx context.Context,
	set exporter.Settings,
	cfg component.Config) (exporter.Metrics, error) {
	c := cfg.(*Config)
	fe, e := newFlowExporter(c, set)
	if e != nil {
		return nil, e
	}

	return exporterhelper.NewMetricsExporter(ctx, set, cfg,
		fe.consumeMetrics,
		exporterhelper.WithStart(fe.Start),
		exporterhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}),
		exporterhelper.WithRetry(c.BackOffConfig))
}
