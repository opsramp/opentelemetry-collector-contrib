// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampmetricsfilterprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/opsrampmetricsfilterprocessor"

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/prometheus/prometheus/promql/parser"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// newKubeClient is a package variable so tests can substitute a fake client.
var newKubeClient = func() (kubernetes.Interface, error) {
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes config: %w", err)
	}
	client, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}
	return client, nil
}

// sourceKey identifies an alert definitions source. Processor instances that
// resolve to the same key share a single loader.
type sourceKey struct {
	filePath      string
	namespace     string
	configMapName string
	configMapKey  string
}

func sourceKeyFor(cfg *Config) sourceKey {
	if cfg.AlertDefinitionsFilePath != "" {
		return sourceKey{filePath: cfg.AlertDefinitionsFilePath}
	}
	return sourceKey{
		namespace:     cfg.Namespace,
		configMapName: cfg.AlertConfigMapName,
		configMapKey:  cfg.AlertConfigMapKey,
	}
}

// definitionsLoader owns a single watch and a single parse of one alert
// definitions source, fanning the categorized metric set out to every processor
// instance configured against that same source. This keeps the number of
// ConfigMap watches and PromQL parses at one regardless of how many category
// specific processor instances exist, and removes reload skew between them.
type definitionsLoader struct {
	key           sourceKey
	logger        *zap.Logger
	client        kubernetes.Interface
	watchFile     bool
	watchInterval time.Duration

	ctx    context.Context
	cancel context.CancelFunc

	mu          sync.Mutex
	refs        int
	started     bool
	nextSubID   int
	subscribers map[int]func(map[string]categorySet)
	latest      map[string]categorySet
	reloadCount int64
	lastReload  time.Time
	fileWatcher *FileWatcher
}

var (
	loadersMu sync.Mutex
	loaders   = make(map[sourceKey]*definitionsLoader)
)

// acquireLoader returns the shared loader for cfg's source, creating it on first
// use. Every successful call must be paired with a release.
func acquireLoader(cfg *Config, logger *zap.Logger) (*definitionsLoader, error) {
	key := sourceKeyFor(cfg)

	loadersMu.Lock()
	defer loadersMu.Unlock()

	if existing, ok := loaders[key]; ok {
		existing.mu.Lock()
		existing.refs++
		existing.mu.Unlock()
		return existing, nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	loader := &definitionsLoader{
		key:         key,
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
		refs:        1,
		subscribers: make(map[int]func(map[string]categorySet)),
		latest:      make(map[string]categorySet),
	}

	if cfg.AlertDefinitionsFilePath != "" {
		interval, err := time.ParseDuration(cfg.FileWatchInterval)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("invalid file_watch_interval: %w", err)
		}
		loader.watchInterval = interval
		loader.watchFile = cfg.WatchFileChanges
	} else {
		client, err := newKubeClient()
		if err != nil {
			cancel()
			return nil, err
		}
		loader.client = client
	}

	loaders[key] = loader
	return loader, nil
}

// start performs the initial load and begins watching. It is idempotent so that
// every subscriber can call it without coordinating.
func (l *definitionsLoader) start() {
	l.mu.Lock()
	if l.started {
		l.mu.Unlock()
		return
	}
	l.started = true
	l.mu.Unlock()

	if err := l.reload(); err != nil {
		l.logger.Error("Failed to load initial alert definitions", zap.Error(err))
	}

	if l.key.filePath != "" {
		if !l.watchFile {
			return
		}
		watcher, err := l.createFileWatcher()
		if err != nil {
			l.logger.Error("Failed to create file watcher", zap.Error(err))
			return
		}
		l.mu.Lock()
		l.fileWatcher = watcher
		l.mu.Unlock()

		l.logger.Info("Starting alert definitions file watcher",
			zap.String("file_path", l.key.filePath),
			zap.Bool("using_fsnotify", watcher.useFsNotify),
			zap.Duration("poll_interval", watcher.pollInterval))
		go watcher.Start()
		return
	}

	go l.watchConfigMap()
}

// release drops one reference, tearing the loader down once the last processor
// instance using it has shut down.
func (l *definitionsLoader) release() {
	loadersMu.Lock()
	defer loadersMu.Unlock()

	l.mu.Lock()
	l.refs--
	last := l.refs <= 0
	watcher := l.fileWatcher
	l.mu.Unlock()

	if !last {
		return
	}

	l.cancel()
	if watcher != nil {
		if err := watcher.Close(); err != nil {
			l.logger.Error("Failed to close file watcher", zap.Error(err))
		}
	}
	delete(loaders, l.key)
}

// subscribe registers fn and immediately delivers the current snapshot so a late
// joining processor does not have to wait for the next reload.
func (l *definitionsLoader) subscribe(fn func(map[string]categorySet)) int {
	l.mu.Lock()
	id := l.nextSubID
	l.nextSubID++
	l.subscribers[id] = fn
	snapshot := l.latest
	l.mu.Unlock()

	fn(snapshot)
	return id
}

func (l *definitionsLoader) unsubscribe(id int) {
	l.mu.Lock()
	delete(l.subscribers, id)
	l.mu.Unlock()
}

// publish stores the new snapshot and notifies subscribers. The map is shared
// read only with every subscriber and must not be mutated after this point.
func (l *definitionsLoader) publish(metrics map[string]categorySet) {
	l.mu.Lock()
	l.latest = metrics
	l.reloadCount++
	l.lastReload = time.Now()
	reloadCount := l.reloadCount
	subscribers := make([]func(map[string]categorySet), 0, len(l.subscribers))
	for _, fn := range l.subscribers {
		subscribers = append(subscribers, fn)
	}
	l.mu.Unlock()

	l.logger.Info("Published alert definitions snapshot",
		zap.Int("unique_metrics", len(metrics)),
		zap.Int("subscribers", len(subscribers)),
		zap.Int64("reload_count", reloadCount))

	for _, fn := range subscribers {
		fn(metrics)
	}
}

func (l *definitionsLoader) reload() error {
	data, err := l.read()
	if err != nil {
		return err
	}
	return l.process(data)
}

func (l *definitionsLoader) read() ([]byte, error) {
	if l.key.filePath != "" {
		return l.readFile()
	}
	return l.readConfigMap()
}

func (l *definitionsLoader) readFile() ([]byte, error) {
	l.logger.Info("Loading alert definitions from file",
		zap.String("file_path", l.key.filePath))

	if _, err := os.Stat(l.key.filePath); err != nil {
		l.logger.Error("Failed to get file info",
			zap.String("file_path", l.key.filePath),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get file info for %s: %w", l.key.filePath, err)
	}

	data, err := os.ReadFile(l.key.filePath)
	if err != nil {
		l.logger.Error("Failed to read alert definitions file",
			zap.String("file_path", l.key.filePath),
			zap.Error(err))
		return nil, fmt.Errorf("failed to read file %s: %w", l.key.filePath, err)
	}
	return data, nil
}

func (l *definitionsLoader) readConfigMap() ([]byte, error) {
	l.logger.Info("Loading alert definitions from ConfigMap",
		zap.String("configmap", l.key.configMapName),
		zap.String("namespace", l.key.namespace),
		zap.String("key", l.key.configMapKey))

	configMap, err := l.client.CoreV1().ConfigMaps(l.key.namespace).Get(
		l.ctx,
		l.key.configMapName,
		metav1.GetOptions{},
	)
	if err != nil {
		l.logger.Error("Failed to get ConfigMap",
			zap.String("configmap", l.key.configMapName),
			zap.String("namespace", l.key.namespace),
			zap.Error(err))
		return nil, fmt.Errorf("failed to get ConfigMap %s/%s: %w", l.key.namespace, l.key.configMapName, err)
	}

	data, exists := configMap.Data[l.key.configMapKey]
	if !exists {
		l.logger.Error("Alert definitions key not found in ConfigMap",
			zap.String("key", l.key.configMapKey),
			zap.Strings("available_keys", getKeys(configMap.Data)))
		return nil, fmt.Errorf("key %s not found in ConfigMap %s/%s", l.key.configMapKey, l.key.namespace, l.key.configMapName)
	}
	return []byte(data), nil
}

// categorizedRule pairs a rule with the category implied by the resourceType of
// the definition group it came from.
type categorizedRule struct {
	rule     AlertRule
	category categorySet
}

// process parses the alert definitions and builds the metric name to category
// mapping that is shared by all subscribers.
func (l *definitionsLoader) process(data []byte) error {
	startTime := time.Now()

	rules, err := parseCategorizedRules(data)
	if err != nil {
		l.logger.Error("Failed to unmarshal alert definitions", zap.Error(err))
		return err
	}

	if len(rules) == 0 {
		l.logger.Warn("No alert definitions found in configuration")
		return nil
	}

	metrics := make(map[string]categorySet)
	invalidExpressions := 0

	for _, cr := range rules {
		if cr.rule.Expr == "" {
			l.logger.Warn("Empty expression found in alert rule",
				zap.String("rule_name", cr.rule.Name))
			invalidExpressions++
			continue
		}

		names := extractMetricsFromExpression(l.logger, cr.rule.Expr)
		if len(names) == 0 {
			l.logger.Warn("No metrics extracted from expression",
				zap.String("expression", cr.rule.Expr),
				zap.String("rule_name", cr.rule.Name))
			continue
		}

		for _, name := range names {
			// OR rather than assign: a metric referenced by both a k8s_pod rule
			// and a non-pod rule must end up in both categories.
			metrics[name] |= cr.category
		}
	}

	l.logger.Info("Loaded alert definitions",
		zap.Int("total_rules", len(rules)),
		zap.Int("invalid_expressions", invalidExpressions),
		zap.Int("unique_metrics", len(metrics)),
		zap.Duration("processing_duration", time.Since(startTime)))

	l.publish(metrics)
	return nil
}

// parseCategorizedRules reads the nested resourceType format, falling back to
// the flat VM agent format which carries no resourceType.
func parseCategorizedRules(data []byte) ([]categorizedRule, error) {
	var nested AlertDefinitions
	if err := yaml.Unmarshal(data, &nested); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert definitions: %w", err)
	}

	var rules []categorizedRule
	for _, def := range nested.AlertDefinitions {
		category := categoryForResourceType(def.ResourceType)
		for _, rule := range def.Rules {
			rules = append(rules, categorizedRule{rule: rule, category: category})
		}
	}
	if len(rules) > 0 {
		return rules, nil
	}

	var flat struct {
		AlertDefinitions []AlertRule `yaml:"alertDefinitions"`
	}
	if err := yaml.Unmarshal(data, &flat); err != nil {
		return nil, fmt.Errorf("failed to unmarshal alert definitions: %w", err)
	}
	// The flat format has no resourceType, so those rules are cluster scoped.
	for _, rule := range flat.AlertDefinitions {
		rules = append(rules, categorizedRule{rule: rule, category: categoryCluster})
	}
	return rules, nil
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
func extractMetricsFromExpression(logger *zap.Logger, expr string) []string {
	parsedExpr, err := parser.NewParser(parser.Options{}).ParseExpr(expr)
	if err != nil {
		logger.Warn("Failed to parse PromQL expression", zap.String("expr", expr), zap.Error(err))
		return nil
	}

	metrics := make(map[string]struct{})
	extractMetricNames(parsedExpr, metrics)

	var result []string
	for metric := range metrics {
		// Selectors such as {__name__=~"..."} yield an empty name.
		if metric == "" {
			continue
		}
		result = append(result, metric)
	}

	return result
}

// watchConfigMap watches for changes to the ConfigMap
func (l *definitionsLoader) watchConfigMap() {
	for {
		select {
		case <-l.ctx.Done():
			return
		default:
			l.doWatch()
		}
	}
}

// doWatch performs the actual watching with retry logic
func (l *definitionsLoader) doWatch() {
	defer func() {
		// Wait before retrying
		select {
		case <-l.ctx.Done():
		case <-time.After(30 * time.Second):
		}
	}()

	watcher, err := l.client.CoreV1().ConfigMaps(l.key.namespace).Watch(
		l.ctx,
		metav1.ListOptions{
			FieldSelector: fmt.Sprintf("metadata.name=%s", l.key.configMapName),
		},
	)
	if err != nil {
		l.logger.Error("Failed to create ConfigMap watcher", zap.Error(err))
		return
	}
	defer watcher.Stop()

	l.logger.Info("Started watching ConfigMap for changes",
		zap.String("configmap", l.key.configMapName),
		zap.String("namespace", l.key.namespace))

	for {
		select {
		case <-l.ctx.Done():
			return
		case event, ok := <-watcher.ResultChan():
			if !ok {
				l.logger.Warn("ConfigMap watcher channel closed, will retry")
				return
			}

			switch event.Type {
			case watch.Modified, watch.Added:
				l.logger.Info("ConfigMap changed, reloading alert definitions")
				if err := l.reload(); err != nil {
					l.logger.Error("Failed to reload alert definitions", zap.Error(err))
				}
			case watch.Deleted:
				l.logger.Warn("ConfigMap was deleted, clearing metrics map")
				l.publish(make(map[string]categorySet))
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
