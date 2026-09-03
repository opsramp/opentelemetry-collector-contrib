// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sattributesprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/moid"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/redis"
)

const enrichCluster = "enrich-cluster"

func podMoidKey(cluster, namespace, pod string) string {
	return moid.NewMoid(cluster).WithNamespaceName(namespace).WithPodName(pod).PodMoid()
}

func newEnrichProcessor(t *testing.T, entries map[string]string, enrich bool) *kubernetesprocessor {
	t.Helper()
	kp := newTestProcessor(newRedisClientWithCache(entries), enrichCluster, false)
	kp.redisConfig.ClusterUid = "cluster-uuid"
	kp.redisConfig.EnrichAttributesFromRedis = enrich
	return kp
}

func podResource(t *testing.T, signalType, pod string) pcommon.Resource {
	t.Helper()
	res := pcommon.NewResource()
	res.Attributes().PutStr("type", signalType)
	res.Attributes().PutStr("k8s.namespace.name", "ns-1")
	res.Attributes().PutStr("k8s.pod.name", pod)
	return res
}

func attr(t *testing.T, res pcommon.Resource, key string) string {
	t.Helper()
	if v, ok := res.Attributes().Get(key); ok {
		return v.Str()
	}
	return ""
}

func TestApplyRedisAttributes(t *testing.T) {
	t.Run("nil data is a no-op", func(t *testing.T) {
		res := pcommon.NewResource()
		assert.NotPanics(t, func() { applyRedisAttributes(res, nil) })
		assert.Equal(t, 0, res.Attributes().Len())
	})

	t.Run("empty fields are skipped", func(t *testing.T) {
		res := pcommon.NewResource()
		applyRedisAttributes(res, &redis.RedisData{NodeName: "node-1"})

		assert.Equal(t, "node-1", attr(t, res, "k8s.node.name"))
		_, found := res.Attributes().Get("k8s.deployment.name")
		assert.False(t, found)
	})

	t.Run("existing non-empty attribute is preserved", func(t *testing.T) {
		res := pcommon.NewResource()
		res.Attributes().PutStr("k8s.node.name", "already-set")
		applyRedisAttributes(res, &redis.RedisData{NodeName: "node-1"})

		assert.Equal(t, "already-set", attr(t, res, "k8s.node.name"))
	})

	t.Run("existing empty attribute is overwritten", func(t *testing.T) {
		res := pcommon.NewResource()
		res.Attributes().PutStr("k8s.node.name", "")
		applyRedisAttributes(res, &redis.RedisData{NodeName: "node-1"})

		assert.Equal(t, "node-1", attr(t, res, "k8s.node.name"))
	})

	t.Run("all labels applied", func(t *testing.T) {
		res := pcommon.NewResource()
		applyRedisAttributes(res, &redis.RedisData{
			NodeName: "node-1", NamespaceName: "ns-1",
			DeploymentName: "dep-1", DaemonSetName: "ds-1",
			StatefulSetName: "ss-1", ReplicaSetName: "rs-1",
			PodName: "pod-1", PodUid: "pod-uid-1", PodIp: "10.0.0.1",
		})

		for key, want := range map[string]string{
			"k8s.node.name":        "node-1",
			"k8s.namespace.name":   "ns-1",
			"k8s.deployment.name":  "dep-1",
			"k8s.daemonset.name":   "ds-1",
			"k8s.statefulset.name": "ss-1",
			"k8s.replicaset.name":  "rs-1",
			"k8s.pod.name":         "pod-1",
			"k8s.pod.uid":          "pod-uid-1",
			"k8s.pod.ip":           "10.0.0.1",
		} {
			assert.Equal(t, want, attr(t, res, key), key)
		}
	})
}

func TestProcessOpsrampResources_EnrichmentMatrix(t *testing.T) {
	withUUID := `{"resourceUuid":"uuid-1","k8s.node.name":"node-1","k8s.deployment.name":"dep-1"}`
	withoutUUID := `{"k8s.node.name":"node-2","k8s.deployment.name":"dep-2"}`

	// The cache is a process-wide singleton, so each case needs its own MoID key.
	setup := func(t *testing.T, pod, payload string, enrich bool) *kubernetesprocessor {
		t.Helper()
		entries := map[string]string{}
		if payload != "" {
			entries[podMoidKey(enrichCluster, "ns-1", pod)] = payload
		}
		return newEnrichProcessor(t, entries, enrich)
	}

	t.Run("metrics with uuid get labels and uuid", func(t *testing.T) {
		kp := setup(t, "pod-metrics-hit", withUUID, true)
		res := podResource(t, "RESOURCE", "pod-metrics-hit")

		kp.processopsrampResources(context.Background(), res, signalMetrics)

		assert.Equal(t, "uuid-1", attr(t, res, "uuid"))
		assert.Equal(t, "node-1", attr(t, res, "k8s.node.name"))
		assert.Equal(t, "dep-1", attr(t, res, "k8s.deployment.name"))
	})

	t.Run("metrics without uuid are not enriched", func(t *testing.T) {
		kp := setup(t, "pod-metrics-nouuid", withoutUUID, true)
		res := podResource(t, "RESOURCE", "pod-metrics-nouuid")

		kp.processopsrampResources(context.Background(), res, signalMetrics)

		_, found := res.Attributes().Get("uuid")
		assert.False(t, found)
		_, found = res.Attributes().Get("k8s.node.name")
		assert.False(t, found, "no label work on the metrics miss path")
	})

	t.Run("logs without uuid still get labels", func(t *testing.T) {
		kp := setup(t, "pod-logs-nouuid", withoutUUID, true)
		res := podResource(t, "log", "pod-logs-nouuid")

		kp.processopsrampResources(context.Background(), res, signalLogs)

		_, found := res.Attributes().Get("resourceUUID")
		assert.False(t, found)
		assert.Equal(t, "node-2", attr(t, res, "k8s.node.name"))
		assert.Equal(t, "dep-2", attr(t, res, "k8s.deployment.name"))
	})

	t.Run("logs with uuid get labels and resourceUUID", func(t *testing.T) {
		kp := setup(t, "pod-logs-hit", withUUID, true)
		res := podResource(t, "log", "pod-logs-hit")

		kp.processopsrampResources(context.Background(), res, signalLogs)

		assert.Equal(t, "uuid-1", attr(t, res, "resourceUUID"))
		assert.Equal(t, "uuid-1", attr(t, res, "k8s.pod.resourceUUID"))
		assert.Equal(t, "node-1", attr(t, res, "k8s.node.name"))
	})

	t.Run("flag off leaves labels untouched", func(t *testing.T) {
		kp := setup(t, "pod-flag-off", withUUID, false)
		res := podResource(t, "RESOURCE", "pod-flag-off")

		kp.processopsrampResources(context.Background(), res, signalMetrics)

		assert.Equal(t, "uuid-1", attr(t, res, "uuid"), "uuid resolution is unaffected by the flag")
		_, found := res.Attributes().Get("k8s.node.name")
		assert.False(t, found)
	})

	t.Run("key missing from redis", func(t *testing.T) {
		kp := setup(t, "pod-absent", "", true)
		res := podResource(t, "RESOURCE", "pod-absent")

		kp.processopsrampResources(context.Background(), res, signalMetrics)

		_, found := res.Attributes().Get("uuid")
		assert.False(t, found)
	})
}

func TestProcessResource_Passthrough(t *testing.T) {
	t.Run("keeps the k8s events object-kind fixup", func(t *testing.T) {
		kp := newEnrichProcessor(t, nil, true)

		res := pcommon.NewResource()
		res.Attributes().PutStr("type", "event")
		res.Attributes().PutStr("k8s.object.kind", "Pod")
		res.Attributes().PutStr("k8s.object.uid", "object-uid-1")

		kp.processResource(context.Background(), res)

		assert.Equal(t, "object-uid-1", attr(t, res, "k8s.pod.uid"))
	})

	t.Run("skips pod association", func(t *testing.T) {
		kp := newEnrichProcessor(t, nil, true)

		res := pcommon.NewResource()
		res.Attributes().PutStr("type", "RESOURCE")
		res.Attributes().PutStr("k8s.pod.name", "pod-1")
		before := res.Attributes().Len()

		kp.processResource(context.Background(), res)

		assert.Equal(t, before, res.Attributes().Len())
		_, found := res.Attributes().Get("k8s.pod.ip")
		assert.False(t, found)
	})
}

func TestProcessTraces_SkipsRedisWhenNotConfigured(t *testing.T) {
	kp := newTestProcessor(nil, enrichCluster, false)

	td := ptrace.NewTraces()
	rs := td.ResourceSpans().AppendEmpty()
	rs.Resource().Attributes().PutStr("k8s.pod.name", "pod-1")
	rs.Resource().Attributes().PutStr("k8s.namespace.name", "ns-1")

	assert.NotPanics(t, func() {
		_, err := kp.processTraces(context.Background(), td)
		require.NoError(t, err)
	})
}

func TestHasPodIdentity(t *testing.T) {
	tests := []struct {
		name  string
		attrs map[string]string
		want  bool
	}{
		{name: "pod name", attrs: map[string]string{"k8s.pod.name": "pod-1"}, want: true},
		{name: "pod uid", attrs: map[string]string{"k8s.pod.uid": "pod-uid-1"}, want: true},
		{name: "empty pod name", attrs: map[string]string{"k8s.pod.name": ""}},
		{name: "pod ip only", attrs: map[string]string{"k8s.pod.ip": "10.0.0.1"}},
		{name: "node only", attrs: map[string]string{"k8s.node.name": "node-1"}},
		{name: "workload only", attrs: map[string]string{"k8s.deployment.name": "dep-1"}},
		{name: "no attributes"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := pcommon.NewResource()
			for k, v := range tt.attrs {
				res.Attributes().PutStr(k, v)
			}
			assert.Equal(t, tt.want, hasPodIdentity(res))
		})
	}
}

// The Redis labels describe a pod and its owners, so non-pod-scoped signals are
// left alone even when their MoID resolves.
func TestProcessOpsrampResources_EnrichmentRequiresPodIdentity(t *testing.T) {
	nodeKey := moid.NewMoid(enrichCluster).WithNodeName("node-gate").NodeMoid()
	kp := newEnrichProcessor(t, map[string]string{
		nodeKey: `{"resourceUuid":"uuid-node","k8s.namespace.name":"ns-1","k8s.deployment.name":"dep-1"}`,
	}, true)

	res := pcommon.NewResource()
	res.Attributes().PutStr("type", "RESOURCE")
	res.Attributes().PutStr("k8s.node.name", "node-gate")

	kp.processopsrampResources(context.Background(), res, signalMetrics)

	assert.Equal(t, "uuid-node", attr(t, res, "uuid"), "uuid resolution is unaffected by the gate")
	_, found := res.Attributes().Get("k8s.namespace.name")
	assert.False(t, found)
	_, found = res.Attributes().Get("k8s.deployment.name")
	assert.False(t, found)
}

func TestProcessOpsrampResources_ClusterFallbackUnchanged(t *testing.T) {
	kp := newEnrichProcessor(t, nil, true)
	res := pcommon.NewResource()
	res.Attributes().PutStr("type", "RESOURCE")

	kp.processopsrampResources(context.Background(), res, signalMetrics)

	assert.Equal(t, "cluster-uuid", attr(t, res, "uuid"))
}

func TestValidate_PassthroughRequiresRedis(t *testing.T) {
	tests := []struct {
		name        string
		passthrough bool
		redisHost   string
		wantErr     bool
	}{
		{name: "passthrough without redis", passthrough: true, wantErr: true},
		{name: "passthrough with redis", passthrough: true, redisHost: "redis-host"},
		{name: "no passthrough no redis"},
		{name: "no passthrough with redis", redisHost: "redis-host"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Passthrough: tt.passthrough}
			cfg.AuthType = k8sconfig.AuthTypeServiceAccount
			cfg.RedisConfig.RedisHost = tt.redisHost
			if tt.redisHost != "" {
				cfg.RedisConfig.RedisPort = "6379"
				cfg.RedisConfig.RedisPass = "pass"
				cfg.RedisConfig.ClusterName = "cluster"
				cfg.RedisConfig.ClusterUid = "cluster-uuid"
				cfg.RedisConfig.NodeName = "node-1"
			}

			err := cfg.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "redis_config.redisHost is required when passthrough is enabled")
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestGetRedisDataUsingPodMoid_BackfillsDeployment(t *testing.T) {
	key := podMoidKey(enrichCluster, "ns-1", "pod-backfill")
	kp := newEnrichProcessor(t, map[string]string{
		key: `{"resourceUuid":"uuid-1","k8s.deployment.name":"dep-1"}`,
	}, false)

	res := podResource(t, "RESOURCE", "pod-backfill")
	res.Attributes().PutStr("k8s.replicaset.name", "rs-1")

	got := kp.GetRedisDataUsingPodMoid(context.Background(), res)

	require.NotNil(t, got)
	assert.Equal(t, "uuid-1", got.ResourceUuid)
	assert.Equal(t, "dep-1", attr(t, res, "k8s.deployment.name"))
}
