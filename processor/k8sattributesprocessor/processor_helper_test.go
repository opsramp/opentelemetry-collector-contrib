// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sattributesprocessor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/cache"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/moid"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/redis"
)

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

// newRedisClientWithCache returns a Redis client whose host is empty (so it never
// dials a real server) but whose shared cache is pre-populated with the given
// moid->uuid entries. GetUuidValueInString then resolves purely from the cache.
func newRedisClientWithCache(entries map[string]string) *redis.Client {
	cacheObj := cache.GetCacheInstance(1000, time.Hour, 1000, time.Hour)
	for k, v := range entries {
		cacheObj.AddToPrimaryWithTTL(k, v, time.Hour)
	}
	return redis.NewClient(zap.NewNop(), cacheObj, "", "", "", time.Hour, time.Hour)
}

// newTestProcessor builds a minimal kubernetesprocessor for helper-level tests.
// passthroughMode is enabled so processResource returns before touching the kube client.
func newTestProcessor(redisClient *redis.Client, cluster string, enableRouting bool) *kubernetesprocessor {
	return &kubernetesprocessor{
		logger:          zap.NewNop(),
		passthroughMode: true,
		redisClient:     redisClient,
		redisConfig: redis.OpsrampRedisConfig{
			ClusterName:         cluster,
			EnableGpuNicRouting: enableRouting,
		},
	}
}

func gpuMoidKey(cluster, gpuUUID string) string {
	return moid.NewMoid(cluster).WithGPUUUID(gpuUUID).GPUMoid()
}

func nicMoidKey(cluster, node, device string) string {
	return moid.NewMoid(cluster).WithNICName(node + "-" + device).NICMoid()
}

// addGaugeMetric appends a gauge metric with one data-point per attribute set.
func addGaugeMetric(sm pmetric.ScopeMetrics, name string, dps []map[string]string) {
	m := sm.Metrics().AppendEmpty()
	m.SetName(name)
	g := m.SetEmptyGauge()
	for _, attrs := range dps {
		dp := g.DataPoints().AppendEmpty()
		dp.SetDoubleValue(1)
		for k, v := range attrs {
			dp.Attributes().PutStr(k, v)
		}
	}
}

func resAttrStr(rm pmetric.ResourceMetrics, key string) string {
	if v, ok := rm.Resource().Attributes().Get(key); ok {
		return v.Str()
	}
	return ""
}

func hasResAttr(rm pmetric.ResourceMetrics, key string) bool {
	_, ok := rm.Resource().Attributes().Get(key)
	return ok
}

func findMetric(rm pmetric.ResourceMetrics, name string) (pmetric.Metric, bool) {
	for i := 0; i < rm.ScopeMetrics().Len(); i++ {
		sm := rm.ScopeMetrics().At(i)
		for j := 0; j < sm.Metrics().Len(); j++ {
			if sm.Metrics().At(j).Name() == name {
				return sm.Metrics().At(j), true
			}
		}
	}
	return pmetric.Metric{}, false
}

// ---------------------------------------------------------------------------
// MoID format
// ---------------------------------------------------------------------------

func TestGPUAndNICMoidFormat(t *testing.T) {
	assert.Equal(t, "c1_GPU_GPU-xyz", moid.NewMoid("c1").WithGPUUUID("GPU-xyz").GPUMoid())
	assert.Equal(t, "c1_NIC_node1-mlx5_0", moid.NewMoid("c1").WithNICName("node1-mlx5_0").NICMoid())
}

// ---------------------------------------------------------------------------
// collectDCGMGPUUUIDs / collectRDMANICDevices
// ---------------------------------------------------------------------------

func TestCollectDCGMGPUUUIDs(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	// DCGM metrics with GPU- prefixed UUIDs (one duplicate) plus a non GPU- value.
	addGaugeMetric(sm, "DCGM_FI_DEV_GPU_TEMP", []map[string]string{
		{"UUID": "GPU-a"},
		{"UUID": "GPU-b"},
		{"UUID": "GPU-a"},
		{"UUID": "not-a-gpu"},
	})
	// Non-DCGM metric carrying a GPU- UUID must be ignored.
	addGaugeMetric(sm, "node_metric", []map[string]string{{"UUID": "GPU-c"}})

	assert.ElementsMatch(t, []string{"GPU-a", "GPU-b"}, collectDCGMGPUUUIDs(rm))
}

func TestCollectDCGMGPUUUIDs_None(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	addGaugeMetric(sm, "node_cpu_seconds_total", []map[string]string{{"cpu": "0"}})
	assert.Empty(t, collectDCGMGPUUUIDs(rm))
}

func TestCollectRDMANICDevices(t *testing.T) {
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	addGaugeMetric(sm, "rdma_rx_bytes", []map[string]string{
		{"device": "mlx5_0"},
		{"device": "mlx5_1"},
		{"device": "mlx5_0"},
		{"device": ""},
	})
	// Non-rdma metric carrying a device must be ignored.
	addGaugeMetric(sm, "node_metric", []map[string]string{{"device": "eth0"}})

	assert.ElementsMatch(t, []string{"mlx5_0", "mlx5_1"}, collectRDMANICDevices(rm))
}

// ---------------------------------------------------------------------------
// filterMetricDataPointsByAttr / dataPointCount
// ---------------------------------------------------------------------------

func TestFilterMetricDataPointsByAttr_Gauge(t *testing.T) {
	md := pmetric.NewMetrics()
	sm := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty()
	addGaugeMetric(sm, "x", []map[string]string{
		{"UUID": "a"}, {"UUID": "b"}, {"UUID": "a"},
	})
	m := sm.Metrics().At(0)
	filterMetricDataPointsByAttr(m, "UUID", "a")
	require.Equal(t, 2, m.Gauge().DataPoints().Len())
	for i := 0; i < m.Gauge().DataPoints().Len(); i++ {
		v, ok := m.Gauge().DataPoints().At(i).Attributes().Get("UUID")
		require.True(t, ok)
		assert.Equal(t, "a", v.Str())
	}
}

func TestFilterMetricDataPointsByAttr_Sum(t *testing.T) {
	md := pmetric.NewMetrics()
	sm := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty()
	m := sm.Metrics().AppendEmpty()
	m.SetName("s")
	s := m.SetEmptySum()
	for _, d := range []string{"mlx5_0", "mlx5_1", "mlx5_0"} {
		dp := s.DataPoints().AppendEmpty()
		dp.Attributes().PutStr("device", d)
	}
	filterMetricDataPointsByAttr(m, "device", "mlx5_0")
	assert.Equal(t, 2, m.Sum().DataPoints().Len())
}

func TestDataPointCount(t *testing.T) {
	md := pmetric.NewMetrics()
	sm := md.ResourceMetrics().AppendEmpty().ScopeMetrics().AppendEmpty()

	gauge := sm.Metrics().AppendEmpty()
	gauge.SetEmptyGauge().DataPoints().AppendEmpty()
	gauge.Gauge().DataPoints().AppendEmpty()
	assert.Equal(t, 2, dataPointCount(gauge))

	sum := sm.Metrics().AppendEmpty()
	sum.SetEmptySum().DataPoints().AppendEmpty()
	assert.Equal(t, 1, dataPointCount(sum))

	hist := sm.Metrics().AppendEmpty()
	hist.SetEmptyHistogram().DataPoints().AppendEmpty()
	assert.Equal(t, 1, dataPointCount(hist))

	summary := sm.Metrics().AppendEmpty()
	summary.SetEmptySummary().DataPoints().AppendEmpty()
	assert.Equal(t, 1, dataPointCount(summary))

	expHist := sm.Metrics().AppendEmpty()
	expHist.SetEmptyExponentialHistogram().DataPoints().AppendEmpty()
	assert.Equal(t, 1, dataPointCount(expHist))

	empty := sm.Metrics().AppendEmpty()
	assert.Equal(t, 0, dataPointCount(empty))
}

// ---------------------------------------------------------------------------
// processDcgmMetricByGpu
// ---------------------------------------------------------------------------

func TestProcessDcgmMetricByGpu(t *testing.T) {
	const cluster = "dcgm-cluster"
	rc := newRedisClientWithCache(map[string]string{
		gpuMoidKey(cluster, "GPU-aaaa"): "uuid-aaaa",
		gpuMoidKey(cluster, "GPU-bbbb"): "uuid-bbbb",
	})
	kp := newTestProcessor(rc, cluster, true)

	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("k8s.node.name", "worker1")
	sm := rm.ScopeMetrics().AppendEmpty()
	addGaugeMetric(sm, "DCGM_FI_DEV_GPU_TEMP", []map[string]string{
		{"UUID": "GPU-aaaa"}, {"UUID": "GPU-bbbb"},
	})
	// Non-DCGM metric is copied as-is into every split.
	addGaugeMetric(sm, "up", []map[string]string{{}})

	splits, ok := kp.processDcgmMetricByGpu(context.Background(), rm)
	require.True(t, ok)
	require.Len(t, splits, 2)

	wantUUID := map[string]string{"GPU-aaaa": "uuid-aaaa", "GPU-bbbb": "uuid-bbbb"}
	seen := map[string]bool{}
	for _, split := range splits {
		gpu := resAttrStr(split, "gpu.uuid")
		require.Contains(t, wantUUID, gpu)
		seen[gpu] = true
		assert.Equal(t, wantUUID[gpu], resAttrStr(split, "uuid"))

		// DCGM metric retains only the data-point for this GPU.
		dcgm, found := findMetric(split, "DCGM_FI_DEV_GPU_TEMP")
		require.True(t, found)
		require.Equal(t, 1, dataPointCount(dcgm))
		v, _ := dcgm.Gauge().DataPoints().At(0).Attributes().Get("UUID")
		assert.Equal(t, gpu, v.Str())

		// Non-DCGM metric is present (parity with reference PR — kept in every split).
		_, hasUp := findMetric(split, "up")
		assert.True(t, hasUp)
	}
	assert.Len(t, seen, 2)
}

func TestProcessDcgmMetricByGpu_NotDCGM(t *testing.T) {
	kp := newTestProcessor(newRedisClientWithCache(nil), "c", true)
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	addGaugeMetric(sm, "node_cpu_seconds_total", []map[string]string{{"cpu": "0"}})

	splits, ok := kp.processDcgmMetricByGpu(context.Background(), rm)
	assert.False(t, ok)
	assert.Nil(t, splits)
}

// ---------------------------------------------------------------------------
// processRdmaMetricByNic
// ---------------------------------------------------------------------------

func TestProcessRdmaMetricByNic(t *testing.T) {
	const cluster = "rdma-cluster"
	const node = "worker2"
	rc := newRedisClientWithCache(map[string]string{
		nicMoidKey(cluster, node, "mlx5_0"): "uuid-nic0",
		nicMoidKey(cluster, node, "mlx5_1"): "uuid-nic1",
	})
	kp := newTestProcessor(rc, cluster, true)

	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("k8s.node.name", node)
	sm := rm.ScopeMetrics().AppendEmpty()
	addGaugeMetric(sm, "rdma_rx_bytes", []map[string]string{
		{"device": "mlx5_0"}, {"device": "mlx5_1"},
	})
	// rdma_ metric without a device label -> empty after filtering -> dropped.
	addGaugeMetric(sm, "rdma_scrape_errors_total", []map[string]string{{}})
	// Non-rdma metric -> dropped entirely from NIC splits.
	addGaugeMetric(sm, "go_goroutines", []map[string]string{{}})

	splits, ok := kp.processRdmaMetricByNic(context.Background(), rm)
	require.True(t, ok)
	require.Len(t, splits, 2)

	wantUUID := map[string]string{"mlx5_0": "uuid-nic0", "mlx5_1": "uuid-nic1"}
	seen := map[string]bool{}
	for _, split := range splits {
		dev := resAttrStr(split, "nic.device")
		require.Contains(t, wantUUID, dev)
		seen[dev] = true
		assert.Equal(t, wantUUID[dev], resAttrStr(split, "uuid"))

		rdma, found := findMetric(split, "rdma_rx_bytes")
		require.True(t, found)
		require.Equal(t, 1, dataPointCount(rdma))
		v, _ := rdma.Gauge().DataPoints().At(0).Attributes().Get("device")
		assert.Equal(t, dev, v.Str())

		// Empty rdma_ metric removed; non-rdma metric dropped.
		_, hasScrape := findMetric(split, "rdma_scrape_errors_total")
		assert.False(t, hasScrape)
		_, hasGo := findMetric(split, "go_goroutines")
		assert.False(t, hasGo)
	}
	assert.Len(t, seen, 2)
}

func TestProcessRdmaMetricByNic_NotRDMA(t *testing.T) {
	kp := newTestProcessor(newRedisClientWithCache(nil), "c", true)
	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	sm := rm.ScopeMetrics().AppendEmpty()
	addGaugeMetric(sm, "node_cpu_seconds_total", []map[string]string{{"cpu": "0"}})

	splits, ok := kp.processRdmaMetricByNic(context.Background(), rm)
	assert.False(t, ok)
	assert.Nil(t, splits)
}

// ---------------------------------------------------------------------------
// processMetrics — end-to-end routing behaviour
// ---------------------------------------------------------------------------

func TestProcessMetrics_DCGMRoutingSplitsPerGPU(t *testing.T) {
	const cluster = "pm-dcgm"
	rc := newRedisClientWithCache(map[string]string{
		gpuMoidKey(cluster, "GPU-1"): "uuid-g1",
		gpuMoidKey(cluster, "GPU-2"): "uuid-g2",
	})
	kp := newTestProcessor(rc, cluster, true)

	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("k8s.node.name", "n1")
	sm := rm.ScopeMetrics().AppendEmpty()
	addGaugeMetric(sm, "DCGM_FI_DEV_GPU_TEMP", []map[string]string{
		{"UUID": "GPU-1"}, {"UUID": "GPU-2"},
	})

	out, err := kp.processMetrics(context.Background(), md)
	require.NoError(t, err)
	require.Equal(t, 2, out.ResourceMetrics().Len())

	got := map[string]string{}
	for i := 0; i < out.ResourceMetrics().Len(); i++ {
		r := out.ResourceMetrics().At(i)
		// Internal marker must never leak downstream.
		assert.False(t, hasResAttr(r, "_dcgm_processed"))
		got[resAttrStr(r, "gpu.uuid")] = resAttrStr(r, "uuid")
	}
	assert.Equal(t, map[string]string{"GPU-1": "uuid-g1", "GPU-2": "uuid-g2"}, got)
}

func TestProcessMetrics_RDMARoutingSplitsPerNIC(t *testing.T) {
	const cluster = "pm-rdma"
	const node = "n2"
	rc := newRedisClientWithCache(map[string]string{
		nicMoidKey(cluster, node, "mlx5_0"): "uuid-n0",
		nicMoidKey(cluster, node, "mlx5_1"): "uuid-n1",
	})
	kp := newTestProcessor(rc, cluster, true)

	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("k8s.node.name", node)
	sm := rm.ScopeMetrics().AppendEmpty()
	addGaugeMetric(sm, "rdma_rx_bytes", []map[string]string{
		{"device": "mlx5_0"}, {"device": "mlx5_1"},
	})

	out, err := kp.processMetrics(context.Background(), md)
	require.NoError(t, err)
	require.Equal(t, 2, out.ResourceMetrics().Len())

	got := map[string]string{}
	for i := 0; i < out.ResourceMetrics().Len(); i++ {
		r := out.ResourceMetrics().At(i)
		assert.False(t, hasResAttr(r, "_rdma_processed"))
		got[resAttrStr(r, "nic.device")] = resAttrStr(r, "uuid")
	}
	assert.Equal(t, map[string]string{"mlx5_0": "uuid-n0", "mlx5_1": "uuid-n1"}, got)
}

// Routing disabled: DCGM batch must NOT be split (parity with original pipeline).
func TestProcessMetrics_RoutingDisabledNoSplit(t *testing.T) {
	kp := newTestProcessor(newRedisClientWithCache(nil), "pm-off", false)

	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("k8s.node.name", "n3")
	// Pre-set uuid so the resource survives filterOnlyOpsrampMetrics.
	rm.Resource().Attributes().PutStr("uuid", "preset")
	sm := rm.ScopeMetrics().AppendEmpty()
	addGaugeMetric(sm, "DCGM_FI_DEV_GPU_TEMP", []map[string]string{
		{"UUID": "GPU-1"}, {"UUID": "GPU-2"},
	})

	out, err := kp.processMetrics(context.Background(), md)
	require.NoError(t, err)
	require.Equal(t, 1, out.ResourceMetrics().Len())
	assert.False(t, hasResAttr(out.ResourceMetrics().At(0), "gpu.uuid"))
	assert.Equal(t, "preset", resAttrStr(out.ResourceMetrics().At(0), "uuid"))
}

// redisClient nil: no routing, no split, no panic.
func TestProcessMetrics_RedisNilNoSplit(t *testing.T) {
	kp := newTestProcessor(nil, "pm-nil", true)

	md := pmetric.NewMetrics()
	rm := md.ResourceMetrics().AppendEmpty()
	rm.Resource().Attributes().PutStr("k8s.node.name", "n4")
	rm.Resource().Attributes().PutStr("uuid", "preset")
	sm := rm.ScopeMetrics().AppendEmpty()
	addGaugeMetric(sm, "DCGM_FI_DEV_GPU_TEMP", []map[string]string{
		{"UUID": "GPU-1"}, {"UUID": "GPU-2"},
	})

	out, err := kp.processMetrics(context.Background(), md)
	require.NoError(t, err)
	require.Equal(t, 1, out.ResourceMetrics().Len())
	assert.False(t, hasResAttr(out.ResourceMetrics().At(0), "gpu.uuid"))
}
