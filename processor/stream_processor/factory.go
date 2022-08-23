package stream_processor

import (
	"context"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
)

const (
	// The value of "type" key in configuration.
	typeStr = "sql_stream"
	// The stability level of the processor.
	stability = component.StabilityLevelInDevelopment
)

var processorCapabilities = consumer.Capabilities{MutatesData: true}

func NewFactory() component.ProcessorFactory {
	return component.NewProcessorFactory(typeStr, createDefaultConfig,
		component.WithLogsProcessorAndStabilityLevel(createSqlStreamProcessor,
			stability))
}

func createDefaultConfig() config.Processor {
	return &Config{
		ProcessorSettings: config.NewProcessorSettings(config.NewComponentID(typeStr)),
	}
}

func createSqlStreamProcessor(
	_ context.Context,
	set component.ProcessorCreateSettings,
	cfg config.Processor,
	nextConsumer consumer.Logs,
) (component.LogsProcessor, error) {
	return newSqlStreamProcessor(nextConsumer, set.Logger, cfg.(*Config)), nil
}
