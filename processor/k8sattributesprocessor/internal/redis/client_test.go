// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/cache"
)

const fullPayload = `{
	"resourceUuid": "uuid-1",
	"resourceHash": 42,
	"k8s.node.name": "node-1",
	"k8s.namespace.name": "ns-1",
	"k8s.deployment.name": "dep-1",
	"k8s.daemonset.name": "ds-1",
	"k8s.statefulset.name": "ss-1",
	"k8s.replicaset.name": "rs-1",
	"k8s.pod.name": "pod-1",
	"k8s.pod.uid": "pod-uid-1",
	"k8s.pod.ip": "10.0.0.1"
}`

// cacheObj is a process-wide singleton, so every test must use unique keys.
func testClient(t *testing.T, enabled bool) *Client {
	t.Helper()
	return &Client{
		logger:                     zap.NewNop(),
		CacheObject:                cache.GetCacheInstance(1000, time.Hour, 1000, time.Hour),
		PrimaryCacheEvictionTime:   time.Hour,
		SecondaryCacheEvictionTime: time.Hour,
		Enabled:                    enabled,
	}
}

func TestDecode(t *testing.T) {
	c := testClient(t, false)

	tests := []struct {
		name string
		val  string
		want *RedisData
	}{
		{
			name: "full json",
			val:  fullPayload,
			want: &RedisData{
				ResourceUuid: "uuid-1", ResourceHash: 42,
				NodeName: "node-1", NamespaceName: "ns-1",
				DeploymentName: "dep-1", DaemonSetName: "ds-1",
				StatefulSetName: "ss-1", ReplicaSetName: "rs-1",
				PodName: "pod-1", PodUid: "pod-uid-1", PodIp: "10.0.0.1",
			},
		},
		{
			name: "partial json keeps zero values",
			val:  `{"resourceUuid":"uuid-2","k8s.pod.name":"pod-2"}`,
			want: &RedisData{ResourceUuid: "uuid-2", PodName: "pod-2"},
		},
		{
			name: "json without uuid",
			val:  `{"k8s.pod.name":"pod-3"}`,
			want: &RedisData{PodName: "pod-3"},
		},
		{
			name: "legacy bare uuid",
			val:  "uuid-legacy",
			want: &RedisData{ResourceUuid: "uuid-legacy"},
		},
		{
			name: "empty is the negative-cache marker",
			val:  "",
			want: nil,
		},
		{
			name: "malformed json",
			val:  `{"resourceUuid":`,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, c.decode("k", tt.val))
		})
	}
}

func TestGetResourceUuidIsNilSafe(t *testing.T) {
	var nilData *RedisData
	assert.Empty(t, nilData.GetResourceUuid())
	assert.Equal(t, "uuid-1", (&RedisData{ResourceUuid: "uuid-1"}).GetResourceUuid())
}

// GoClient is left nil: any Redis round trip would panic, so completing the
// call proves the lookup was served entirely from cache.
func TestGetRedisData_PrimaryCacheHitSkipsRedis(t *testing.T) {
	c := testClient(t, true)
	key := "primary-hit-key"
	c.CacheObject.AddToPrimaryWithTTL(key, fullPayload, time.Hour)

	got := c.GetRedisData(context.Background(), key)

	require.NotNil(t, got)
	assert.Equal(t, "uuid-1", got.ResourceUuid)
	assert.Equal(t, "dep-1", got.DeploymentName)
	assert.Equal(t, "10.0.0.1", got.PodIp)
}

func TestGetRedisData_SecondaryNegativeHitSkipsRedis(t *testing.T) {
	c := testClient(t, true)
	key := "secondary-miss-key"
	c.CacheObject.AddToSecondaryWithTTL(key, "", time.Hour)

	assert.Nil(t, c.GetRedisData(context.Background(), key))
}

func TestGetRedisData_DisabledClientReturnsNil(t *testing.T) {
	c := testClient(t, false)

	assert.Nil(t, c.GetRedisData(context.Background(), "disabled-client-key"))
}

func TestGetRedisData_NilCacheObjectDoesNotPanic(t *testing.T) {
	c := &Client{logger: zap.NewNop(), Enabled: false}

	assert.NotPanics(t, func() {
		assert.Nil(t, c.GetRedisData(context.Background(), "nil-cache-key"))
		c.negativeCache("nil-cache-key")
	})
}

func TestWrappersShareTheSingleFetch(t *testing.T) {
	c := testClient(t, true)
	key := "wrapper-key"
	c.CacheObject.AddToPrimaryWithTTL(key, fullPayload, time.Hour)
	ctx := context.Background()

	assert.Equal(t, "uuid-1", c.GetUuidValueInString(ctx, key))

	withDeployment := c.GetRedisDataWithDeployment(ctx, key)
	require.NotNil(t, withDeployment)
	assert.Equal(t, "uuid-1", withDeployment.ResourceUuid)
	assert.Equal(t, "dep-1", withDeployment.DeploymentName)
	assert.Equal(t, uint64(42), withDeployment.ResourceHash)
}

func TestWrappersOnMiss(t *testing.T) {
	c := testClient(t, false)
	ctx := context.Background()

	assert.Empty(t, c.GetUuidValueInString(ctx, "wrapper-miss-key"))
	assert.Nil(t, c.GetRedisDataWithDeployment(ctx, "wrapper-miss-key"))
}
