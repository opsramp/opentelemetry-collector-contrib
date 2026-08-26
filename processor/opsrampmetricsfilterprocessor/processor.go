// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampmetricsfilterprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/opsrampmetricsfilterprocessor"

import (
	"context"
	"strings"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"
)

// AlertRule represents a single alert rule
type AlertRule struct {
	Name              string `yaml:"name"`
	Interval          string `yaml:"interval"`
	Expr              string `yaml:"expr"`
	IsAvailability    bool   `yaml:"isAvailability"`
	WarnOperator      string `yaml:"warnOperator,omitempty"`
	WarnThreshold     string `yaml:"warnThreshold,omitempty"`
	CriticalOperator  string `yaml:"criticalOperator,omitempty"`
	CriticalThreshold string `yaml:"criticalThreshold,omitempty"`
	AlertSub          string `yaml:"alertSub,omitempty"`
	AlertBody         string `yaml:"alertBody,omitempty"`
}

// AlertDefinition represents a single alert definition group
type AlertDefinition struct {
	ResourceType string      `yaml:"resourceType"`
	Rules        []AlertRule `yaml:"rules"`
}

// AlertDefinitions represents the structure of alert definitions
type AlertDefinitions struct {
	AlertDefinitions []AlertDefinition `yaml:"alertDefinitions"`
}

// filterProcessor implements the alert metrics extractor processor
type filterProcessor struct {
	config       *Config
	nextConsumer consumer.Metrics
	logger       *zap.Logger

	// categoryMask selects which alert definition categories this instance keeps
	categoryMask categorySet

	// Thread-safe map of metric names to extract
	metricsMutex sync.RWMutex
	metricsMap   map[string]bool

	// Shared alert definitions source
	loader *definitionsLoader
	subID  int
}

// Ensure filterProcessor implements processor.Metrics interface
var _ processor.Metrics = (*filterProcessor)(nil)

// newFilterProcessor creates a new instance of the filterProcessor
func newFilterProcessor(settings processor.Settings, config *Config, nextConsumer consumer.Metrics) (*filterProcessor, error) {
	mask, err := parseCategoryMask(config.MetricCategories)
	if err != nil {
		return nil, err
	}

	fp := &filterProcessor{
		config:       config,
		nextConsumer: nextConsumer,
		logger:       settings.Logger,
		categoryMask: mask,
		metricsMap:   make(map[string]bool),
	}

	loader, err := acquireLoader(config, settings.Logger)
	if err != nil {
		return nil, err
	}
	fp.loader = loader
	fp.subID = loader.subscribe(fp.applyDefinitions)
	loader.start()

	return fp, nil
}

// applyDefinitions projects the shared, categorized metric set onto this
// instance's configured categories. The incoming map is shared with other
// subscribers and must not be modified.
func (fp *filterProcessor) applyDefinitions(all map[string]categorySet) {
	filtered := make(map[string]bool, len(all))
	for name, categories := range all {
		if categories&fp.categoryMask != 0 {
			filtered[name] = true
		}
	}

	fp.metricsMutex.Lock()
	previousCount := len(fp.metricsMap)
	fp.metricsMap = filtered
	fp.metricsMutex.Unlock()

	fp.logger.Info("Applied alert definitions",
		zap.String("categories", fp.categoryMask.String()),
		zap.Int("metrics_count", len(filtered)),
		zap.Int("previous_metrics_count", previousCount),
		zap.Int("source_metrics_count", len(all)))
}

// Start starts the processor
func (fp *filterProcessor) Start(context.Context, component.Host) error {
	if fp.config.AlertDefinitionsFilePath != "" {
		fp.logger.Info("Starting alert metrics extractor processor with file path",
			zap.String("file_path", fp.config.AlertDefinitionsFilePath),
			zap.Bool("watch_changes", fp.config.WatchFileChanges),
			zap.String("watch_interval", fp.config.FileWatchInterval),
			zap.String("categories", fp.categoryMask.String()))
	} else {
		fp.logger.Info("Starting alert metrics extractor processor with ConfigMap",
			zap.String("configmap_name", fp.config.AlertConfigMapName),
			zap.String("configmap_key", fp.config.AlertConfigMapKey),
			zap.String("namespace", fp.config.Namespace),
			zap.String("categories", fp.categoryMask.String()))
	}

	// Log the current state of the metrics map
	fp.metricsMutex.RLock()
	currentMetricsCount := len(fp.metricsMap)
	fp.metricsMutex.RUnlock()

	fp.logger.Info("Processor started with metrics configuration",
		zap.Int("metrics_count", currentMetricsCount))

	return nil
}

// Shutdown stops the processor
func (fp *filterProcessor) Shutdown(context.Context) error {
	fp.logger.Info("Shutting down alert metrics extractor processor")

	if fp.loader != nil {
		fp.loader.unsubscribe(fp.subID)
		fp.loader.release()
	}

	return nil
}

// Capabilities returns the consumer capabilities
func (fp *filterProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

// ConsumeMetrics processes the metrics
func (fp *filterProcessor) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	fp.metricsMutex.RLock()
	defer fp.metricsMutex.RUnlock()

	// Count incoming metrics
	totalIncomingMetrics := 0
	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		rm := md.ResourceMetrics().At(i)
		for j := 0; j < rm.ScopeMetrics().Len(); j++ {
			sm := rm.ScopeMetrics().At(j)
			totalIncomingMetrics += sm.Metrics().Len()
		}
	}

	// Early return if no metrics to process
	if totalIncomingMetrics == 0 {
		fp.logger.Debug("No incoming metrics to process")
		return fp.nextConsumer.ConsumeMetrics(ctx, md)
	}

	// If no metrics are configured for filtering, drop all metrics
	if len(fp.metricsMap) == 0 {
		fp.logger.Info("No metrics filter configured, dropping all metrics",
			zap.Int("dropped_metrics", totalIncomingMetrics))
		// Return empty metrics to effectively drop all metrics
		return fp.nextConsumer.ConsumeMetrics(ctx, pmetric.NewMetrics())
	}

	// Create a new metrics object to hold filtered metrics
	filteredMetrics := pmetric.NewMetrics()
	totalFilteredMetrics := 0

	// Iterate through all resource metrics
	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		rm := md.ResourceMetrics().At(i)

		// Create a new resource metrics container only if we find matching metrics
		var filteredRM pmetric.ResourceMetrics
		resourceHasMatchingMetrics := false

		// Iterate through all scope metrics
		for j := 0; j < rm.ScopeMetrics().Len(); j++ {
			sm := rm.ScopeMetrics().At(j)

			// Create a new scope metrics container only if we find matching metrics
			var filteredSM pmetric.ScopeMetrics
			scopeHasMatchingMetrics := false

			// Iterate through all metrics
			for k := 0; k < sm.Metrics().Len(); k++ {
				metric := sm.Metrics().At(k)
				metricName := metric.Name()
				// Convert dots to underscores for Prometheus compatibility
				metricName = strings.ReplaceAll(metricName, ".", "_")

				// Check if this metric should be included
				if fp.metricsMap[metricName] {
					// Create containers only when we have a matching metric
					if !resourceHasMatchingMetrics {
						filteredRM = filteredMetrics.ResourceMetrics().AppendEmpty()
						rm.Resource().CopyTo(filteredRM.Resource())
						resourceHasMatchingMetrics = true
					}

					if !scopeHasMatchingMetrics {
						filteredSM = filteredRM.ScopeMetrics().AppendEmpty()
						sm.Scope().CopyTo(filteredSM.Scope())
						scopeHasMatchingMetrics = true
					}

					filteredMetric := filteredSM.Metrics().AppendEmpty()
					metric.CopyTo(filteredMetric)
					totalFilteredMetrics++
				}
			}
		}
	}

	fp.logger.Info("Metrics filtering completed",
		zap.Int("incoming_metrics", totalIncomingMetrics),
		zap.Int("filtered_metrics", totalFilteredMetrics),
		zap.Int("configured_filter_count", len(fp.metricsMap)))

	if totalFilteredMetrics == 0 {
		fp.logger.Warn("No metrics passed through the filter - check configuration")
	}

	return fp.nextConsumer.ConsumeMetrics(ctx, filteredMetrics)
}
