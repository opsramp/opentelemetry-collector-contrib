// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sclusterreceiver

import (
	"strconv"
	"testing"

	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver/internal/testutils"
)

func newBenchResourceWatcher(groups MetricsGroupsConfig) *resourceWatcher {
	return &resourceWatcher{
		client:        newFakeClientWithAllResources(),
		logger:        zap.NewNop(),
		metadataStore: metadata.NewStore(),
		config: &Config{
			MetricsBuilderConfig: metadata.NewDefaultMetricsBuilderConfig(),
			MetricsGroups:        groups,
		},
	}
}

// BenchmarkInformerCacheObject measures the bytes retained by a single object once it has been put
// through the informer transform, i.e. the per-object cost of keeping an informer running. Source
// objects are built outside the timed loop so B/op reflects only the cached copy.
func BenchmarkInformerCacheObject(b *testing.B) {
	cases := []struct {
		name string
		obj  any
	}{
		{
			"Pod", testutils.NewPodWithContainer("0",
				testutils.NewPodSpecWithContainer("container-name"),
				testutils.NewPodStatusWithContainer("container-name", "container-id")),
		},
		{"Node", testutils.NewNode("0")},
		{"Deployment", testutils.NewDeployment("0")},
		{"ReplicaSet", testutils.NewReplicaSet("0")},
		{"Service", testutils.NewService("0")},
		{"Job", testutils.NewJob("0")},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				out, err := transformObject(c.obj)
				if err != nil || out == nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkPrepareSharedInformerFactory measures receiver startup with every metrics group enabled
// versus an opt-out configuration that keeps only nodes and deployments.
func BenchmarkPrepareSharedInformerFactory(b *testing.B) {
	enabled, disabled := true, false

	cases := []struct {
		name   string
		groups MetricsGroupsConfig
	}{
		{"all_groups", MetricsGroupsConfig{}},
		{"nodes_and_deployments_only", MetricsGroupsConfig{
			EnabledByDefault: &disabled,
			Nodes:            ResourceGroupConfig{Enabled: &enabled},
			Deployments:      ResourceGroupConfig{Enabled: &enabled},
		}},
	}

	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				rw := newBenchResourceWatcher(c.groups)
				if err := rw.prepareSharedInformerFactory(); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkInformerCacheGrowth reports the cost of the objects an informer would hold for a cluster
// of the given pod count, which is what is reclaimed when the pods group is disabled.
func BenchmarkInformerCacheGrowth(b *testing.B) {
	for _, pods := range []int{500, 2000, 5000} {
		b.Run("pods_"+strconv.Itoa(pods), func(b *testing.B) {
			src := testutils.NewPodWithContainer("0",
				testutils.NewPodSpecWithContainer("container-name"),
				testutils.NewPodStatusWithContainer("container-name", "container-id"))

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cache := make([]any, 0, pods)
				for j := 0; j < pods; j++ {
					o, err := transformObject(src)
					if err != nil {
						b.Fatal(err)
					}
					cache = append(cache, o)
				}
				if len(cache) != pods {
					b.Fatal("unexpected cache size")
				}
			}
		})
	}
}
