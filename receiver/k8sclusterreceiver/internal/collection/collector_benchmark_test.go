// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package collection

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"go.opentelemetry.io/collector/receiver/receivertest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver/internal/gvk"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver/internal/testutils"
)

// workloadCount is the number of objects created for each non-pod resource type. It is kept small
// relative to podCount to mirror a real cluster, where pods dominate.
const workloadCount = 50

func podStore(pods int) *testutils.MockStore {
	cache := make(map[string]any, pods)
	for i := 0; i < pods; i++ {
		id := strconv.Itoa(i)
		cache["pod"+id+"-uid"] = testutils.NewPodWithContainer(
			id,
			testutils.NewPodSpecWithContainer("container-name"),
			testutils.NewPodStatusWithContainer("container-name", "container-id"),
		)
	}
	return &testutils.MockStore{Cache: cache}
}

func workloadStore(n int, new func(id string) any) *testutils.MockStore {
	cache := make(map[string]any, n)
	for i := 0; i < n; i++ {
		id := strconv.Itoa(i)
		cache[id+"-uid"] = new(id)
	}
	return &testutils.MockStore{Cache: cache}
}

// benchStore builds a metadata store representing a cluster. When pods is 0 the Pod store is never
// registered, which is exactly the state produced by `metrics_groups: {pods: {enabled: false}}`.
func benchStore(pods int) *metadata.Store {
	ms := metadata.NewStore()

	if pods > 0 {
		ms.Setup(gvk.Pod, metadata.ClusterWideInformerKey, podStore(pods))
	}

	ms.Setup(gvk.Node, metadata.ClusterWideInformerKey,
		workloadStore(workloadCount, func(id string) any { return testutils.NewNode(id) }))
	ms.Setup(gvk.Namespace, metadata.ClusterWideInformerKey,
		workloadStore(workloadCount, func(id string) any { return testutils.NewNamespace(id) }))
	ms.Setup(gvk.Deployment, metadata.ClusterWideInformerKey,
		workloadStore(workloadCount, func(id string) any { return testutils.NewDeployment(id) }))
	ms.Setup(gvk.ReplicaSet, metadata.ClusterWideInformerKey,
		workloadStore(workloadCount, func(id string) any { return testutils.NewReplicaSet(id) }))
	ms.Setup(gvk.DaemonSet, metadata.ClusterWideInformerKey,
		workloadStore(workloadCount, func(id string) any { return testutils.NewDaemonset(id) }))
	ms.Setup(gvk.StatefulSet, metadata.ClusterWideInformerKey,
		workloadStore(workloadCount, func(id string) any { return testutils.NewStatefulset(id) }))
	ms.Setup(gvk.Job, metadata.ClusterWideInformerKey,
		workloadStore(workloadCount, func(id string) any { return testutils.NewJob(id) }))
	ms.Setup(gvk.CronJob, metadata.ClusterWideInformerKey,
		workloadStore(workloadCount, func(id string) any { return testutils.NewCronJob(id) }))
	ms.Setup(gvk.HorizontalPodAutoscaler, metadata.ClusterWideInformerKey,
		workloadStore(workloadCount, func(id string) any { return testutils.NewHPA(id) }))
	ms.Setup(gvk.ResourceQuota, metadata.ClusterWideInformerKey,
		workloadStore(workloadCount, func(id string) any { return testutils.NewResourceQuota(id) }))
	ms.Setup(gvk.ReplicationController, metadata.ClusterWideInformerKey,
		workloadStore(workloadCount, func(id string) any { return testutils.NewReplicationController(id) }))
	ms.Setup(gvk.Service, metadata.ClusterWideInformerKey,
		workloadStore(workloadCount, func(id string) any { return testutils.NewService(id) }))

	return ms
}

// BenchmarkCollectMetricData measures one scrape. The "pods_disabled" variants have no Pod store
// registered, which is what disabling the pods metrics group produces.
func BenchmarkCollectMetricData(b *testing.B) {
	for _, pods := range []int{0, 500, 2000, 5000} {
		name := fmt.Sprintf("pods_%d", pods)
		if pods == 0 {
			name = "pods_disabled"
		}
		b.Run(name, func(b *testing.B) {
			dc := NewDataCollector(receivertest.NewNopSettings(metadata.Type), benchStore(pods),
				metadata.NewDefaultMetricsBuilderConfig(), []string{"Ready"}, []string{"cpu", "memory"})
			ts := time.Now()

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				md := dc.CollectMetricData(ts)
				if md.ResourceMetrics().Len() == 0 {
					b.Fatal("no metrics produced")
				}
			}
		})
	}
}
