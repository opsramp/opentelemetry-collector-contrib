// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package kubelet

import (
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	stats "k8s.io/kubelet/pkg/apis/stats/v1alpha1"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/metadata"
)

func podStatusFixture() *v1.PodList {
	return &v1.PodList{
		Items: []v1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "running-pod",
					Namespace: "default",
					UID:       "pod-uid-1",
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name: "app",
							Resources: v1.ResourceRequirements{
								Limits: v1.ResourceList{
									v1.ResourceCPU:     resource.MustParse("500m"),
									v1.ResourceMemory:  resource.MustParse("100M"),
									v1.ResourceStorage: resource.MustParse("2G"),
								},
								Requests: v1.ResourceList{
									v1.ResourceCPU:     resource.MustParse("250m"),
									v1.ResourceMemory:  resource.MustParse("50M"),
									v1.ResourceStorage: resource.MustParse("1G"),
								},
							},
						},
					},
				},
				Status: v1.PodStatus{
					Phase: v1.PodRunning,
					ContainerStatuses: []v1.ContainerStatus{
						{
							Name:         "app",
							ContainerID:  "containerd://abc123",
							RestartCount: 3,
							Ready:        true,
							State:        v1.ContainerState{Running: &v1.ContainerStateRunning{}},
						},
					},
				},
			},
			{
				// This pod has no entry in the stats summary, so it is only reachable via /pods.
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pending-pod",
					Namespace: "kube-system",
					UID:       "pod-uid-2",
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{{Name: "sidecar"}},
				},
				Status: v1.PodStatus{
					Phase:  v1.PodFailed,
					Reason: "Evicted",
					ContainerStatuses: []v1.ContainerStatus{
						{
							Name:         "sidecar",
							RestartCount: 0,
							Ready:        false,
							State: v1.ContainerState{
								Waiting: &v1.ContainerStateWaiting{Reason: "ImagePullBackOff"},
							},
						},
					},
				},
			},
		},
	}
}

func allPodStatusMetricsConfig() metadata.MetricsBuilderConfig {
	cfg := metadata.DefaultMetricsBuilderConfig()
	cfg.Metrics.K8sPodPhase.Enabled = true
	cfg.Metrics.K8sPodStatusReason.Enabled = true
	cfg.Metrics.K8sContainerRestarts.Enabled = true
	cfg.Metrics.K8sContainerReady.Enabled = true
	cfg.Metrics.K8sContainerStatusReason.Enabled = true
	cfg.Metrics.K8sContainerCPULimit.Enabled = true
	cfg.Metrics.K8sContainerCPURequest.Enabled = true
	cfg.Metrics.K8sContainerMemoryLimit.Enabled = true
	cfg.Metrics.K8sContainerMemoryRequest.Enabled = true
	cfg.Metrics.K8sContainerStorageLimit.Enabled = true
	cfg.Metrics.K8sContainerStorageRequest.Enabled = true
	return cfg
}

func newBuilders(cfg metadata.MetricsBuilderConfig) *metadata.MetricsBuilders {
	return &metadata.MetricsBuilders{
		NodeMetricsBuilder:      metadata.NewMetricsBuilder(cfg, receivertest.NewNopSettings(metadata.Type)),
		PodMetricsBuilder:       metadata.NewMetricsBuilder(cfg, receivertest.NewNopSettings(metadata.Type)),
		ContainerMetricsBuilder: metadata.NewMetricsBuilder(cfg, receivertest.NewNopSettings(metadata.Type)),
		OtherMetricsBuilder:     metadata.NewMetricsBuilder(cfg, receivertest.NewNopSettings(metadata.Type)),
	}
}

// collect flattens the emitted metrics into a lookup of resource-attribute key -> metric name -> data points.
func collect(mds []pmetric.Metrics) map[string]map[string][]pmetric.NumberDataPoint {
	out := map[string]map[string][]pmetric.NumberDataPoint{}
	for _, md := range mds {
		for i := 0; i < md.ResourceMetrics().Len(); i++ {
			rm := md.ResourceMetrics().At(i)
			var attrs []string
			rm.Resource().Attributes().Range(func(k string, v pcommon.Value) bool {
				attrs = append(attrs, k+"="+v.AsString()+";")
				return true
			})
			sort.Strings(attrs)
			key := strings.Join(attrs, "")
			if _, ok := out[key]; !ok {
				out[key] = map[string][]pmetric.NumberDataPoint{}
			}
			for j := 0; j < rm.ScopeMetrics().Len(); j++ {
				ms := rm.ScopeMetrics().At(j).Metrics()
				for k := 0; k < ms.Len(); k++ {
					m := ms.At(k)
					var dps pmetric.NumberDataPointSlice
					switch m.Type() {
					case pmetric.MetricTypeGauge:
						dps = m.Gauge().DataPoints()
					case pmetric.MetricTypeSum:
						dps = m.Sum().DataPoints()
					default:
						continue
					}
					for d := 0; d < dps.Len(); d++ {
						out[key][m.Name()] = append(out[key][m.Name()], dps.At(d))
					}
				}
			}
		}
	}
	return out
}

func TestPodStatusMetrics(t *testing.T) {
	k8sMetadata := NewMetadata([]MetadataLabel{MetadataLabelContainerID}, podStatusFixture(), NodeInfo{}, nil)
	mbs := newBuilders(allPodStatusMetricsConfig())

	mds := MetricsData(zap.NewNop(), &stats.Summary{}, k8sMetadata, ValidMetricGroups, nil, mbs, NewCPUUsageCalculator(), true)
	byResource := collect(mds)

	runningPod := byResource["k8s.namespace.name=default;k8s.pod.name=running-pod;k8s.pod.uid=pod-uid-1;"]
	require.NotNil(t, runningPod)
	require.Len(t, runningPod["k8s.pod.phase"], 1)
	assert.Equal(t, int64(2), runningPod["k8s.pod.phase"][0].IntValue())
	require.Len(t, runningPod["k8s.pod.status_reason"], 1)
	assert.Equal(t, int64(6), runningPod["k8s.pod.status_reason"][0].IntValue())

	runningContainer := byResource["container.id=abc123;k8s.container.name=app;k8s.namespace.name=default;k8s.pod.name=running-pod;k8s.pod.uid=pod-uid-1;"]
	require.NotNil(t, runningContainer)
	require.Len(t, runningContainer["k8s.container.restarts"], 1)
	assert.Equal(t, int64(3), runningContainer["k8s.container.restarts"][0].IntValue())
	require.Len(t, runningContainer["k8s.container.ready"], 1)
	assert.Equal(t, int64(1), runningContainer["k8s.container.ready"][0].IntValue())
	require.Len(t, runningContainer["k8s.container.cpu_limit"], 1)
	assert.InDelta(t, 0.5, runningContainer["k8s.container.cpu_limit"][0].DoubleValue(), 0.0001)
	require.Len(t, runningContainer["k8s.container.cpu_request"], 1)
	assert.InDelta(t, 0.25, runningContainer["k8s.container.cpu_request"][0].DoubleValue(), 0.0001)
	require.Len(t, runningContainer["k8s.container.memory_limit"], 1)
	assert.Equal(t, int64(100000000), runningContainer["k8s.container.memory_limit"][0].IntValue())
	require.Len(t, runningContainer["k8s.container.memory_request"], 1)
	assert.Equal(t, int64(50000000), runningContainer["k8s.container.memory_request"][0].IntValue())
	require.Len(t, runningContainer["k8s.container.storage_limit"], 1)
	assert.Equal(t, int64(2000000000), runningContainer["k8s.container.storage_limit"][0].IntValue())
	require.Len(t, runningContainer["k8s.container.storage_request"], 1)
	assert.Equal(t, int64(1000000000), runningContainer["k8s.container.storage_request"][0].IntValue())
	// A running container matches none of the known reasons, so all series are zero.
	require.Len(t, runningContainer["k8s.container.status.reason"], len(allContainerStatusReasons))
	for _, dp := range runningContainer["k8s.container.status.reason"] {
		assert.Equal(t, int64(0), dp.IntValue())
	}

	pendingPod := byResource["k8s.namespace.name=kube-system;k8s.pod.name=pending-pod;k8s.pod.uid=pod-uid-2;"]
	require.NotNil(t, pendingPod)
	assert.Equal(t, int64(4), pendingPod["k8s.pod.phase"][0].IntValue())
	assert.Equal(t, int64(1), pendingPod["k8s.pod.status_reason"][0].IntValue())

	// No container.id resource attribute because the container has not been created yet.
	pendingContainer := byResource["k8s.container.name=sidecar;k8s.namespace.name=kube-system;k8s.pod.name=pending-pod;k8s.pod.uid=pod-uid-2;"]
	require.NotNil(t, pendingContainer)
	assert.Equal(t, int64(0), pendingContainer["k8s.container.ready"][0].IntValue())
	assert.Empty(t, pendingContainer["k8s.container.cpu_limit"])

	var reasonSet int
	for _, dp := range pendingContainer["k8s.container.status.reason"] {
		reason, ok := dp.Attributes().Get("k8s.container.status.reason")
		require.True(t, ok)
		if dp.IntValue() == 1 {
			reasonSet++
			assert.Equal(t, "ImagePullBackOff", reason.Str())
		}
	}
	assert.Equal(t, 1, reasonSet)
}

func TestPodStatusMetricsDisabled(t *testing.T) {
	k8sMetadata := NewMetadata(nil, podStatusFixture(), NodeInfo{}, nil)
	mbs := newBuilders(metadata.DefaultMetricsBuilderConfig())

	mds := MetricsData(zap.NewNop(), &stats.Summary{}, k8sMetadata, ValidMetricGroups, nil, mbs, NewCPUUsageCalculator(), false)
	assert.Empty(t, collect(mds))
}

func TestPodStatusMetricGroupsFiltered(t *testing.T) {
	k8sMetadata := NewMetadata(nil, podStatusFixture(), NodeInfo{}, nil)
	mbs := newBuilders(allPodStatusMetricsConfig())

	mds := MetricsData(zap.NewNop(), &stats.Summary{}, k8sMetadata, map[MetricGroup]bool{PodMetricGroup: true}, nil, mbs, NewCPUUsageCalculator(), true)
	for _, metrics := range collect(mds) {
		assert.Empty(t, metrics["k8s.container.restarts"])
		assert.NotEmpty(t, metrics["k8s.pod.phase"])
	}
}

func TestPodPhaseToInt(t *testing.T) {
	assert.Equal(t, int64(1), podPhaseToInt(v1.PodPending))
	assert.Equal(t, int64(2), podPhaseToInt(v1.PodRunning))
	assert.Equal(t, int64(3), podPhaseToInt(v1.PodSucceeded))
	assert.Equal(t, int64(4), podPhaseToInt(v1.PodFailed))
	assert.Equal(t, int64(5), podPhaseToInt(v1.PodUnknown))
	assert.Equal(t, int64(5), podPhaseToInt(v1.PodPhase("Whatever")))
}

func TestPodStatusReasonToInt(t *testing.T) {
	assert.Equal(t, int64(1), podStatusReasonToInt("Evicted"))
	assert.Equal(t, int64(2), podStatusReasonToInt("NodeAffinity"))
	assert.Equal(t, int64(3), podStatusReasonToInt("NodeLost"))
	assert.Equal(t, int64(4), podStatusReasonToInt("Shutdown"))
	assert.Equal(t, int64(5), podStatusReasonToInt("UnexpectedAdmissionError"))
	assert.Equal(t, int64(6), podStatusReasonToInt(""))
}
