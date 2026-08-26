// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampmetricsfilterprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestParseCategoryMask(t *testing.T) {
	tests := []struct {
		name       string
		categories []string
		expected   categorySet
		wantErr    bool
	}{
		{
			name:       "empty selects everything",
			categories: nil,
			expected:   categoryAll,
		},
		{
			name:       "pod only",
			categories: []string{"podMetric"},
			expected:   categoryPod,
		},
		{
			name:       "cluster only",
			categories: []string{"clusterMetric"},
			expected:   categoryCluster,
		},
		{
			name:       "both listed",
			categories: []string{"podMetric", "clusterMetric"},
			expected:   categoryAll,
		},
		{
			name:       "case and space insensitive",
			categories: []string{" PodMetric "},
			expected:   categoryPod,
		},
		{
			name:       "unknown value rejected",
			categories: []string{"podMetrics"},
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mask, err := parseCategoryMask(tt.categories)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.expected, mask)
		})
	}
}

func TestCategoryForResourceType(t *testing.T) {
	tests := []struct {
		resourceType string
		expected     categorySet
	}{
		{"k8s_pod", categoryPod},
		{"K8S_POD", categoryPod},
		{"k8s-pod", categoryPod},
		{"Pod", categoryPod},
		{"k8s_cluster", categoryCluster},
		{"k8s_node", categoryCluster},
		{"k8s_deployment", categoryCluster},
		{"", categoryCluster},
	}

	for _, tt := range tests {
		t.Run(tt.resourceType, func(t *testing.T) {
			assert.Equal(t, tt.expected, categoryForResourceType(tt.resourceType))
		})
	}
}

func TestLoaderCategorizesMetrics(t *testing.T) {
	// container_cpu_usage is referenced by both a k8s_pod rule and a k8s_node
	// rule, so it must end up in both categories.
	yamlData := `
alertDefinitions:
  - resourceType: "k8s_pod"
    rules:
      - name: "pod cpu"
        expr: "container_cpu_usage > 80"
      - name: "pod memory"
        expr: "pod_memory_working_set > 90"
  - resourceType: "k8s_node"
    rules:
      - name: "node cpu"
        expr: "container_cpu_usage / node_cpu_total > 0.9"
  - resourceType: "k8s_cluster"
    rules:
      - name: "apiserver errors"
        expr: "rate(apiserver_request_total[5m]) > 1"
`

	got := processIntoSnapshot(t, []byte(yamlData))

	assert.Equal(t, map[string]categorySet{
		"container_cpu_usage":     categoryPod | categoryCluster,
		"pod_memory_working_set":  categoryPod,
		"node_cpu_total":          categoryCluster,
		"apiserver_request_total": categoryCluster,
	}, got)
}

func TestLoaderFlatFormatIsClusterScoped(t *testing.T) {
	yamlData := `
alertDefinitions:
  - name: "flat rule"
    expr: "some_metric > 1"
`

	got := processIntoSnapshot(t, []byte(yamlData))

	assert.Equal(t, map[string]categorySet{"some_metric": categoryCluster}, got)
}

func TestApplyDefinitionsProjectsCategories(t *testing.T) {
	snapshot := map[string]categorySet{
		"both_metric":    categoryPod | categoryCluster,
		"pod_metric":     categoryPod,
		"cluster_metric": categoryCluster,
	}

	tests := []struct {
		name     string
		mask     categorySet
		expected map[string]bool
	}{
		{
			name:     "pod instance",
			mask:     categoryPod,
			expected: map[string]bool{"both_metric": true, "pod_metric": true},
		},
		{
			name:     "cluster instance",
			mask:     categoryCluster,
			expected: map[string]bool{"both_metric": true, "cluster_metric": true},
		},
		{
			name:     "unrestricted instance",
			mask:     categoryAll,
			expected: map[string]bool{"both_metric": true, "pod_metric": true, "cluster_metric": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &filterProcessor{
				logger:       zap.NewNop(),
				categoryMask: tt.mask,
				metricsMap:   make(map[string]bool),
			}
			fp.applyDefinitions(snapshot)
			assert.Equal(t, tt.expected, fp.metricsMap)
		})
	}
}

func TestLoaderPublishesToAllSubscribers(t *testing.T) {
	loader := newTestLoader()

	var podMap, clusterMap map[string]bool
	pod := &filterProcessor{logger: zap.NewNop(), categoryMask: categoryPod, metricsMap: map[string]bool{}}
	cluster := &filterProcessor{logger: zap.NewNop(), categoryMask: categoryCluster, metricsMap: map[string]bool{}}

	loader.subscribe(pod.applyDefinitions)
	loader.subscribe(cluster.applyDefinitions)

	require.NoError(t, loader.process([]byte(`
alertDefinitions:
  - resourceType: "k8s_pod"
    rules:
      - name: "shared"
        expr: "shared_metric > 1"
  - resourceType: "k8s_node"
    rules:
      - name: "shared node"
        expr: "shared_metric + node_only_metric > 1"
`)))

	podMap = pod.metricsMap
	clusterMap = cluster.metricsMap

	assert.Equal(t, map[string]bool{"shared_metric": true}, podMap)
	assert.Equal(t, map[string]bool{"shared_metric": true, "node_only_metric": true}, clusterMap)
}

func newTestLoader() *definitionsLoader {
	return &definitionsLoader{
		logger:      zap.NewNop(),
		subscribers: make(map[int]func(map[string]categorySet)),
		latest:      make(map[string]categorySet),
	}
}

func processIntoSnapshot(t *testing.T, data []byte) map[string]categorySet {
	t.Helper()

	loader := newTestLoader()
	var got map[string]categorySet
	loader.subscribe(func(snapshot map[string]categorySet) { got = snapshot })
	require.NoError(t, loader.process(data))
	return got
}
