// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampmetricsfilterprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/opsrampmetricsfilterprocessor"

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/prometheus/promql/parser"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	client       kubernetes.Interface

	// Thread-safe map of metric names to extract
	metricsMutex sync.RWMutex
	metricsMap   map[string]bool

	// File watching
	fileWatcher      *FileWatcher
	lastModTime      time.Time
	watchIntervalDur time.Duration

	// Performance tracking
	reloadCount      int64
	lastReloadTime   time.Time
	processedMetrics int64
	filteredMetrics  int64

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

// FileWatcher handles file system watching for alert definitions
type FileWatcher struct {
	filePath     string
	logger       *zap.Logger
	callback     func() error
	watcher      *fsnotify.Watcher
	ctx          context.Context
	useFsNotify  bool
	pollInterval time.Duration
}

// Ensure filterProcessor implements processor.Metrics interface
var _ processor.Metrics = (*filterProcessor)(nil)

// newFilterProcessor creates a new instance of the filterProcessor
func newFilterProcessor(settings processor.Settings, config *Config, nextConsumer consumer.Metrics) (*filterProcessor, error) {
	ctx, cancel := context.WithCancel(context.Background())

	fp := &filterProcessor{
		config:       config,
		nextConsumer: nextConsumer,
		logger:       settings.Logger,
		metricsMap:   make(map[string]bool),
		ctx:          ctx,
		cancel:       cancel,
	}

	// Parse watch interval for file watching
	if config.AlertDefinitionsFilePath != "" {
		watchInterval, err := time.ParseDuration(config.FileWatchInterval)
		if err != nil {
			return nil, fmt.Errorf("invalid file_watch_interval: %w", err)
		}
		fp.watchIntervalDur = watchInterval
	} else {
		// Create Kubernetes client for ConfigMap mode
		k8sConfig, err := rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create Kubernetes config: %w", err)
		}
		client, err := kubernetes.NewForConfig(k8sConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
		}
		fp.client = client
	}

	// Initial load of alert definitions
	if err := fp.loadAlertDefinitions(); err != nil {
		fp.logger.Error("Failed to load initial alert definitions", zap.Error(err))
		// Don't return error here, just log it and continue
	}

	// Start watching for changes
	if config.AlertDefinitionsFilePath != "" {
		if config.WatchFileChanges {
			go fp.watchFile()
		}
	} else {
		go fp.watchConfigMap()
	}

	return fp, nil
}

// Start starts the processor
func (fp *filterProcessor) Start(ctx context.Context, host component.Host) error {
	if fp.config.AlertDefinitionsFilePath != "" {
		fp.logger.Info("Starting alert metrics extractor processor with file path",
			zap.String("file_path", fp.config.AlertDefinitionsFilePath),
			zap.Bool("watch_changes", fp.config.WatchFileChanges),
			zap.String("watch_interval", fp.config.FileWatchInterval))
	} else {
		fp.logger.Info("Starting alert metrics extractor processor with ConfigMap",
			zap.String("configmap_name", fp.config.AlertConfigMapName),
			zap.String("configmap_key", fp.config.AlertConfigMapKey),
			zap.String("namespace", fp.config.Namespace))
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
func (fp *filterProcessor) Shutdown(ctx context.Context) error {
	fp.logger.Info("Shutting down alert metrics extractor processor")

	// Close file watcher if it exists
	if fp.fileWatcher != nil {
		if err := fp.fileWatcher.Close(); err != nil {
			fp.logger.Error("Failed to close file watcher", zap.Error(err))
		}
	}

	fp.cancel()
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

// loadAlertDefinitions loads alert definitions from either file or ConfigMap
func (fp *filterProcessor) loadAlertDefinitions() error {
	if fp.config.AlertDefinitionsFilePath != "" {
		return fp.loadAlertDefinitionsFromFile()
	}
	return fp.loadAlertDefinitionsFromConfigMap()
}

// loadAlertDefinitionsFromFile loads alert definitions from a file
func (fp *filterProcessor) loadAlertDefinitionsFromFile() error {
	fp.logger.Error("Loading alert definitions from file",
		zap.String("file_path", fp.config.AlertDefinitionsFilePath))

	// Get file modification time for change detection
	fileInfo, err := os.Stat(fp.config.AlertDefinitionsFilePath)
	if err != nil {
		fp.logger.Error("Failed to get file info",
			zap.String("file_path", fp.config.AlertDefinitionsFilePath),
			zap.Error(err))
		return fmt.Errorf("failed to get file info for %s: %w", fp.config.AlertDefinitionsFilePath, err)
	}
	fp.lastModTime = fileInfo.ModTime()

	// Read file content
	alertDefData, err := os.ReadFile(fp.config.AlertDefinitionsFilePath)
	if err != nil {
		fp.logger.Error("Failed to read alert definitions file",
			zap.String("file_path", fp.config.AlertDefinitionsFilePath),
			zap.Error(err))
		return fmt.Errorf("failed to read file %s: %w", fp.config.AlertDefinitionsFilePath, err)
	}

	return fp.processAlertDefinitionsData(alertDefData)
}

// loadAlertDefinitionsFromConfigMap loads alert definitions from the ConfigMap
func (fp *filterProcessor) loadAlertDefinitionsFromConfigMap() error {
	fp.logger.Info("Loading alert definitions from ConfigMap",
		zap.String("configmap", fp.config.AlertConfigMapName),
		zap.String("namespace", fp.config.Namespace),
		zap.String("key", fp.config.AlertConfigMapKey))

	configMap, err := fp.client.CoreV1().ConfigMaps(fp.config.Namespace).Get(
		context.TODO(),
		fp.config.AlertConfigMapName,
		metav1.GetOptions{},
	)
	if err != nil {
		fp.logger.Error("Failed to get ConfigMap",
			zap.String("configmap", fp.config.AlertConfigMapName),
			zap.String("namespace", fp.config.Namespace),
			zap.Error(err))
		return fmt.Errorf("failed to get ConfigMap %s/%s: %w", fp.config.Namespace, fp.config.AlertConfigMapName, err)
	}

	alertDefData, exists := configMap.Data[fp.config.AlertConfigMapKey]
	if !exists {
		fp.logger.Error("Alert definitions key not found in ConfigMap",
			zap.String("key", fp.config.AlertConfigMapKey),
			zap.Strings("available_keys", getKeys(configMap.Data)))
		return fmt.Errorf("key %s not found in ConfigMap %s/%s", fp.config.AlertConfigMapKey, fp.config.Namespace, fp.config.AlertConfigMapName)
	}

	return fp.processAlertDefinitionsData([]byte(alertDefData))
}

// processAlertDefinitionsData processes the alert definitions YAML data
func (fp *filterProcessor) processAlertDefinitionsData(data []byte) error {
	startTime := time.Now()

	var alertDefs AlertDefinitions

	if fp.config.AlertDefinitionsFilePath == "" {
		if err := yaml.Unmarshal(data, &alertDefs); err != nil {
			fp.logger.Error("Failed to unmarshal alert definitions", zap.Error(err))
			return fmt.Errorf("failed to unmarshal alert definitions: %w", err)
		}
	} else {
		// Flat format
		var flat struct {
			AlertDefinitions []AlertRule `yaml:"alertDefinitions"`
		}
		if err := yaml.Unmarshal(data, &flat); err != nil {
			fmt.Println("Failed to unmarshal flat format:", err)
			return fmt.Errorf("failed to unmarshal alert definitions: %w", err)
		}
		alertDefs.AlertDefinitions = []AlertDefinition{
			{
				ResourceType: "file", // placeholder
				Rules:        flat.AlertDefinitions,
			},
		}
	}

	// Validate that we have alert definitions
	if len(alertDefs.AlertDefinitions) == 0 {
		fp.logger.Warn("No alert definitions found in configuration")
		return nil
	}

	fp.logger.Warn("Loaded alert definitions",
		zap.Int("alert_definitions_count", len(alertDefs.AlertDefinitions)))

	// Extract metrics from alert expressions
	newMetricsMap := make(map[string]bool)
	totalRules := len(alertDefs.AlertDefinitions) // Now this is the direct count of rules
	invalidExpressions := 0

	// Iterate directly over the alert rules (no more nested structure)
	for _, alertDef := range alertDefs.AlertDefinitions {
		for _, rule := range alertDef.Rules {
			if rule.Expr == "" {
				fp.logger.Warn("Empty expression found in alert rule",
					zap.String("rule_name", rule.Name))
				invalidExpressions++
				continue
			}

			metrics := fp.extractMetricsFromExpression(rule.Expr)
			if len(metrics) == 0 {
				fp.logger.Warn("No metrics extracted from expression",
					zap.String("expression", rule.Expr),
					zap.String("rule_name", rule.Name))
			}

			for _, metric := range metrics {
				newMetricsMap[metric] = true
				fp.logger.Warn("Extracted metric from rule",
					zap.String("metric", metric),
					zap.String("rule_name", rule.Name),
					zap.String("expression", rule.Expr))
			}
		}
	}

	// Update the global metrics map and tracking
	fp.metricsMutex.Lock()
	previousCount := len(fp.metricsMap)
	fp.metricsMap = newMetricsMap
	fp.reloadCount++
	fp.lastReloadTime = time.Now()
	fp.metricsMutex.Unlock()

	processingDuration := time.Since(startTime)

	fp.logger.Warn("Successfully loaded alert definitions",
		zap.Int("alert_definitions", len(alertDefs.AlertDefinitions)),
		zap.Int("total_rules", totalRules),
		zap.Int("invalid_expressions", invalidExpressions),
		zap.Int("unique_metrics", len(newMetricsMap)),
		zap.Int("previous_metrics_count", previousCount),
		zap.Int64("reload_count", fp.reloadCount),
		zap.Duration("processing_duration", processingDuration))

	// Log the extracted metrics for debugging
	if len(newMetricsMap) > 0 {
		var metricsList []string
		for metric := range newMetricsMap {
			metricsList = append(metricsList, metric)
		}
		fp.logger.Warn("Extracted metrics from alert definitions",
			zap.Strings("metrics", metricsList))
	}

	return nil
}

// extractMetricNames recursively walks the AST to extract metric names
func extractMetricNames(node parser.Node, metrics map[string]struct{}) {
	switch n := node.(type) {
	case *parser.VectorSelector:
		metrics[n.Name] = struct{}{}
	case *parser.MatrixSelector:
		if vs, ok := n.VectorSelector.(*parser.VectorSelector); ok {
			metrics[vs.Name] = struct{}{}
		}
	}

	for _, child := range parser.Children(node) {
		extractMetricNames(child, metrics)
	}
}

// extractMetricsFromExpression extracts metric names from a PromQL expression
func (fp *filterProcessor) extractMetricsFromExpression(expr string) []string {
	parsedExpr, err := parser.ParseExpr(expr)
	if err != nil {
		fp.logger.Warn("Failed to parse PromQL expression", zap.String("expr", expr), zap.Error(err))
		return nil
	}

	metrics := make(map[string]struct{})
	extractMetricNames(parsedExpr, metrics)

	var result []string
	for metric := range metrics {
		result = append(result, metric)
	}

	return result
}

// watchFile watches for changes to the alert definitions file
func (fp *filterProcessor) watchFile() {
	// Create file watcher
	watcher, err := fp.createFileWatcher()
	if err != nil {
		fp.logger.Error("Failed to create file watcher", zap.Error(err))
		return
	}
	defer watcher.Close()

	fp.logger.Info("Starting enhanced file watcher",
		zap.String("file_path", fp.config.AlertDefinitionsFilePath),
		zap.Bool("using_fsnotify", watcher.useFsNotify),
		zap.Duration("poll_interval", watcher.pollInterval))

	watcher.Start()
}

// createFileWatcher creates a new file watcher with fsnotify support and polling fallback
func (fp *filterProcessor) createFileWatcher() (*FileWatcher, error) {
	fw := &FileWatcher{
		filePath:     fp.config.AlertDefinitionsFilePath,
		logger:       fp.logger,
		callback:     fp.loadAlertDefinitionsFromFile,
		ctx:          fp.ctx,
		pollInterval: fp.watchIntervalDur,
		useFsNotify:  true,
	}

	// Try to create fsnotify watcher
	var err error
	fw.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		fw.logger.Warn("Failed to create fsnotify watcher, falling back to polling",
			zap.Error(err))
		fw.useFsNotify = false
	} else {
		// Add the file to the watcher
		// Watch the directory since file might be replaced (common with editors)
		dir := filepath.Dir(fw.filePath)
		if err := fw.watcher.Add(dir); err != nil {
			fw.logger.Warn("Failed to add directory to watcher, falling back to polling",
				zap.String("directory", dir),
				zap.Error(err))
			fw.watcher.Close()
			fw.useFsNotify = false
		}
	}

	return fw, nil
}

// Start begins watching for file changes
func (fw *FileWatcher) Start() {
	if fw.useFsNotify {
		fw.startFsNotifyWatch()
	} else {
		fw.startPollingWatch()
	}
}

// Close stops the file watcher and cleans up resources
func (fw *FileWatcher) Close() error {
	if fw.watcher != nil {
		return fw.watcher.Close()
	}
	return nil
}

// startFsNotifyWatch uses fsnotify for real-time file change detection
func (fw *FileWatcher) startFsNotifyWatch() {
	fw.logger.Info("Starting fsnotify file watcher", zap.String("file_path", fw.filePath))

	for {
		select {
		case <-fw.ctx.Done():
			return
		case event, ok := <-fw.watcher.Events:
			if !ok {
				fw.logger.Error("File watcher events channel closed")
				return
			}

			// Check if this event is for our target file
			if filepath.Base(event.Name) == filepath.Base(fw.filePath) {
				fw.logger.Debug("File event received",
					zap.String("file", event.Name),
					zap.String("operation", event.Op.String()))

				// Handle file write, create, or rename events
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Rename) {
					fw.logger.Info("File modified, reloading",
						zap.String("file_path", fw.filePath),
						zap.String("event", event.Op.String()))

					if err := fw.callback(); err != nil {
						fw.logger.Error("Failed to reload file after change", zap.Error(err))
					} else {
						fw.logger.Info("Successfully reloaded file after change")
					}
				}
			}
		case err, ok := <-fw.watcher.Errors:
			if !ok {
				fw.logger.Error("File watcher errors channel closed")
				return
			}
			fw.logger.Error("File watcher error", zap.Error(err))
		}
	}
}

// startPollingWatch falls back to polling for file changes
func (fw *FileWatcher) startPollingWatch() {
	fw.logger.Info("Starting polling file watcher",
		zap.String("file_path", fw.filePath),
		zap.Duration("interval", fw.pollInterval))

	var lastModTime time.Time

	// Get initial modification time
	if fileInfo, err := os.Stat(fw.filePath); err == nil {
		lastModTime = fileInfo.ModTime()
	}

	ticker := time.NewTicker(fw.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-fw.ctx.Done():
			return
		case <-ticker.C:
			fileInfo, err := os.Stat(fw.filePath)
			if err != nil {
				fw.logger.Error("Failed to get file info during polling",
					zap.String("file_path", fw.filePath),
					zap.Error(err))
				continue
			}

			// Check if file has been modified
			if fileInfo.ModTime().After(lastModTime) {
				fw.logger.Info("File changed (polling), reloading",
					zap.String("file_path", fw.filePath),
					zap.Time("old_mod_time", lastModTime),
					zap.Time("new_mod_time", fileInfo.ModTime()))

				lastModTime = fileInfo.ModTime()
				if err := fw.callback(); err != nil {
					fw.logger.Error("Failed to reload file after polling change", zap.Error(err))
				} else {
					fw.logger.Info("Successfully reloaded file after polling change")
				}
			}
		}
	}
}

// checkFileChanges checks if the file has been modified and reloads if necessary
func (fp *filterProcessor) checkFileChanges() {
	fileInfo, err := os.Stat(fp.config.AlertDefinitionsFilePath)
	if err != nil {
		fp.logger.Error("Failed to get file info during watch",
			zap.String("file_path", fp.config.AlertDefinitionsFilePath),
			zap.Error(err))
		return
	}

	// Check if file has been modified
	if fileInfo.ModTime().After(fp.lastModTime) {
		fp.logger.Info("File changed, reloading alert definitions",
			zap.String("file_path", fp.config.AlertDefinitionsFilePath),
			zap.Time("old_mod_time", fp.lastModTime),
			zap.Time("new_mod_time", fileInfo.ModTime()))

		if err := fp.loadAlertDefinitionsFromFile(); err != nil {
			fp.logger.Error("Failed to reload alert definitions from file", zap.Error(err))
		} else {
			fp.logger.Info("Successfully reloaded alert definitions from file")
		}
	}
}

// watchConfigMap watches for changes to the ConfigMap
func (fp *filterProcessor) watchConfigMap() {
	for {
		select {
		case <-fp.ctx.Done():
			return
		default:
			fp.doWatch()
		}
	}
}

// doWatch performs the actual watching with retry logic
func (fp *filterProcessor) doWatch() {
	defer func() {
		// Wait before retrying
		select {
		case <-fp.ctx.Done():
			return
		case <-time.After(30 * time.Second):
		}
	}()

	watcher, err := fp.client.CoreV1().ConfigMaps(fp.config.Namespace).Watch(
		context.TODO(),
		metav1.ListOptions{
			FieldSelector: fmt.Sprintf("metadata.name=%s", fp.config.AlertConfigMapName),
		},
	)
	if err != nil {
		fp.logger.Error("Failed to create ConfigMap watcher", zap.Error(err))
		return
	}
	defer watcher.Stop()

	fp.logger.Info("Started watching ConfigMap for changes",
		zap.String("configmap", fp.config.AlertConfigMapName),
		zap.String("namespace", fp.config.Namespace))

	for {
		select {
		case <-fp.ctx.Done():
			return
		case event, ok := <-watcher.ResultChan():
			if !ok {
				fp.logger.Warn("ConfigMap watcher channel closed, will retry")
				return
			}

			switch event.Type {
			case watch.Modified, watch.Added:
				fp.logger.Info("ConfigMap changed, reloading alert definitions")
				if err := fp.loadAlertDefinitionsFromConfigMap(); err != nil {
					fp.logger.Error("Failed to reload alert definitions", zap.Error(err))
				}
			case watch.Deleted:
				fp.logger.Warn("ConfigMap was deleted, clearing metrics map")
				fp.metricsMutex.Lock()
				fp.metricsMap = make(map[string]bool)
				fp.metricsMutex.Unlock()
			}
		}
	}
}

// getKeys returns a slice of keys from a map[string]string (for logging)
func getKeys(data map[string]string) []string {
	var keys []string
	for key := range data {
		keys = append(keys, key)
	}
	return keys
}
