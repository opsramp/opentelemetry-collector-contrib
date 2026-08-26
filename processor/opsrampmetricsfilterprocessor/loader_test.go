// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampmetricsfilterprocessor

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/processor/processortest"
)

const twoScopeDefinitions = `
alertDefinitions:
  - resourceType: "k8s_pod"
    rules:
      - name: "pod cpu"
        expr: "container_cpu_usage > 80"
  - resourceType: "k8s_node"
    rules:
      - name: "node cpu"
        expr: "container_cpu_usage / node_cpu_total > 0.9"
`

func loaderCount() int {
	loadersMu.Lock()
	defer loadersMu.Unlock()
	return len(loaders)
}

func (fp *filterProcessor) snapshot() map[string]bool {
	fp.metricsMutex.RLock()
	defer fp.metricsMutex.RUnlock()
	out := make(map[string]bool, len(fp.metricsMap))
	for k, v := range fp.metricsMap {
		out[k] = v
	}
	return out
}

func newFileConfig(t *testing.T, path string, categories ...string) *Config {
	t.Helper()
	cfg := createDefaultConfig().(*Config)
	cfg.AlertDefinitionsFilePath = path
	cfg.WatchFileChanges = false
	cfg.MetricCategories = categories
	require.NoError(t, cfg.Validate())
	return cfg
}

// Exercises the real construction path: two category-scoped instances against
// one file source must share a single loader, receive correctly projected
// allow-lists, and tear the loader down only after the last one shuts down.
func TestSharedLoaderAcrossCategoryInstances(t *testing.T) {
	path := filepath.Join(t.TempDir(), "alert-definitions.yaml")
	require.NoError(t, os.WriteFile(path, []byte(twoScopeDefinitions), 0o600))

	before := loaderCount()
	set := processortest.NewNopSettings(component.MustNewType(typeStr))

	podProc, err := newFilterProcessor(set, newFileConfig(t, path, "podMetric"), consumertest.NewNop())
	require.NoError(t, err)
	clusterProc, err := newFilterProcessor(set, newFileConfig(t, path, "clusterMetric"), consumertest.NewNop())
	require.NoError(t, err)

	assert.Same(t, podProc.loader, clusterProc.loader, "instances on the same source must share a loader")
	assert.Equal(t, before+1, loaderCount(), "only one loader should be registered")
	assert.NotEqual(t, podProc.subID, clusterProc.subID)

	assert.Equal(t, map[string]bool{"container_cpu_usage": true}, podProc.snapshot())
	assert.Equal(t, map[string]bool{
		"container_cpu_usage": true,
		"node_cpu_total":      true,
	}, clusterProc.snapshot())

	require.NoError(t, podProc.Shutdown(context.Background()))
	assert.Equal(t, before+1, loaderCount(), "loader must survive while another instance holds it")

	require.NoError(t, clusterProc.Shutdown(context.Background()))
	assert.Equal(t, before, loaderCount(), "loader must be released by the last instance")
}

// A reload must reach every subscriber from a single publish.
func TestReloadUpdatesAllInstances(t *testing.T) {
	path := filepath.Join(t.TempDir(), "alert-definitions.yaml")
	require.NoError(t, os.WriteFile(path, []byte(twoScopeDefinitions), 0o600))

	set := processortest.NewNopSettings(component.MustNewType(typeStr))
	podProc, err := newFilterProcessor(set, newFileConfig(t, path, "podMetric"), consumertest.NewNop())
	require.NoError(t, err)
	clusterProc, err := newFilterProcessor(set, newFileConfig(t, path, "clusterMetric"), consumertest.NewNop())
	require.NoError(t, err)
	defer func() {
		require.NoError(t, podProc.Shutdown(context.Background()))
		require.NoError(t, clusterProc.Shutdown(context.Background()))
	}()

	require.NoError(t, os.WriteFile(path, []byte(`
alertDefinitions:
  - resourceType: "k8s_pod"
    rules:
      - name: "pod disk"
        expr: "pod_disk_usage > 50"
  - resourceType: "k8s_cluster"
    rules:
      - name: "apiserver"
        expr: "apiserver_request_total > 1"
`), 0o600))
	require.NoError(t, podProc.loader.reload())

	assert.Equal(t, map[string]bool{"pod_disk_usage": true}, podProc.snapshot())
	assert.Equal(t, map[string]bool{"apiserver_request_total": true}, clusterProc.snapshot())
}

// An unparseable reload must leave the previous allow-list in place rather than
// silently dropping every metric.
func TestFailedReloadRetainsPreviousAllowList(t *testing.T) {
	path := filepath.Join(t.TempDir(), "alert-definitions.yaml")
	require.NoError(t, os.WriteFile(path, []byte(twoScopeDefinitions), 0o600))

	set := processortest.NewNopSettings(component.MustNewType(typeStr))
	proc, err := newFilterProcessor(set, newFileConfig(t, path, "podMetric"), consumertest.NewNop())
	require.NoError(t, err)
	defer func() { require.NoError(t, proc.Shutdown(context.Background())) }()

	require.NoError(t, os.WriteFile(path, []byte("alertDefinitions: [[[not yaml"), 0o600))
	require.Error(t, proc.loader.reload())

	assert.Equal(t, map[string]bool{"container_cpu_usage": true}, proc.snapshot())
}

// End-to-end: a category-scoped instance must drop metrics outside its category.
func TestConsumeMetricsRespectsCategory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "alert-definitions.yaml")
	require.NoError(t, os.WriteFile(path, []byte(twoScopeDefinitions), 0o600))

	sink := &consumertest.MetricsSink{}
	set := processortest.NewNopSettings(component.MustNewType(typeStr))
	proc, err := newFilterProcessor(set, newFileConfig(t, path, "podMetric"), sink)
	require.NoError(t, err)
	defer func() { require.NoError(t, proc.Shutdown(context.Background())) }()

	md := pmetric.NewMetrics()
	sm := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty()
	for _, name := range []string{"container_cpu_usage", "node_cpu_total", "unrelated_metric"} {
		m := sm.Metrics().AppendEmpty()
		m.SetName(name)
		m.SetEmptyGauge().DataPoints().AppendEmpty().SetDoubleValue(1)
	}

	require.NoError(t, proc.ConsumeMetrics(context.Background(), md))

	require.Len(t, sink.AllMetrics(), 1)
	got := sink.AllMetrics()[0]
	require.Equal(t, 1, got.ResourceMetrics().Len())
	forwarded := got.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics()
	require.Equal(t, 1, forwarded.Len())
	assert.Equal(t, "container_cpu_usage", forwarded.At(0).Name())
}

// Distinct sources must not share a loader.
func TestDistinctSourcesGetDistinctLoaders(t *testing.T) {
	dir := t.TempDir()
	pathA := filepath.Join(dir, "a.yaml")
	pathB := filepath.Join(dir, "b.yaml")
	require.NoError(t, os.WriteFile(pathA, []byte(twoScopeDefinitions), 0o600))
	require.NoError(t, os.WriteFile(pathB, []byte(twoScopeDefinitions), 0o600))

	set := processortest.NewNopSettings(component.MustNewType(typeStr))
	procA, err := newFilterProcessor(set, newFileConfig(t, pathA), consumertest.NewNop())
	require.NoError(t, err)
	procB, err := newFilterProcessor(set, newFileConfig(t, pathB), consumertest.NewNop())
	require.NoError(t, err)
	defer func() {
		require.NoError(t, procA.Shutdown(context.Background()))
		require.NoError(t, procB.Shutdown(context.Background()))
	}()

	assert.NotSame(t, procA.loader, procB.loader)
}

// Omitting metric_categories must reproduce the pre-change allow-list exactly.
func TestNoCategoriesKeepsEveryMetric(t *testing.T) {
	path := filepath.Join(t.TempDir(), "alert-definitions.yaml")
	require.NoError(t, os.WriteFile(path, []byte(twoScopeDefinitions), 0o600))

	set := processortest.NewNopSettings(component.MustNewType(typeStr))
	proc, err := newFilterProcessor(set, newFileConfig(t, path), consumertest.NewNop())
	require.NoError(t, err)
	defer func() { require.NoError(t, proc.Shutdown(context.Background())) }()

	assert.Equal(t, map[string]bool{
		"container_cpu_usage": true,
		"node_cpu_total":      true,
	}, proc.snapshot())
}

// A non-positive interval would panic time.NewTicker in the polling fallback.
func TestFileWatchIntervalMustBePositive(t *testing.T) {
	path := filepath.Join(t.TempDir(), "alert-definitions.yaml")
	require.NoError(t, os.WriteFile(path, []byte(twoScopeDefinitions), 0o600))

	cfg := createDefaultConfig().(*Config)
	cfg.AlertDefinitionsFilePath = path

	cfg.FileWatchInterval = "0s"
	require.ErrorContains(t, cfg.Validate(), "must be positive")

	cfg.FileWatchInterval = "-5s"
	require.ErrorContains(t, cfg.Validate(), "must be positive")

	cfg.FileWatchInterval = "not-a-duration"
	require.ErrorContains(t, cfg.Validate(), "invalid file_watch_interval format")

	cfg.FileWatchInterval = "15s"
	require.NoError(t, cfg.Validate())
}
