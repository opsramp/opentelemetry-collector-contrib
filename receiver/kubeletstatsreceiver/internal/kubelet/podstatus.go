// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package kubelet // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/kubelet"

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	v1 "k8s.io/api/core/v1"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/metadata"
)

// allContainerStatusReasons is emitted in full on every scrape so that consumers never observe a
// missing series when a container transitions between reasons.
var allContainerStatusReasons = []metadata.AttributeK8sContainerStatusReason{
	metadata.AttributeK8sContainerStatusReasonContainerCreating,
	metadata.AttributeK8sContainerStatusReasonCrashLoopBackOff,
	metadata.AttributeK8sContainerStatusReasonCreateContainerConfigError,
	metadata.AttributeK8sContainerStatusReasonErrImagePull,
	metadata.AttributeK8sContainerStatusReasonImagePullBackOff,
	metadata.AttributeK8sContainerStatusReasonOOMKilled,
	metadata.AttributeK8sContainerStatusReasonCompleted,
	metadata.AttributeK8sContainerStatusReasonError,
	metadata.AttributeK8sContainerStatusReasonContainerCannotRun,
}

// podStatus emits state metrics sourced from the kubelet /pods endpoint. It is driven by the pod
// list rather than the stats summary so that pods and containers without usage stats (pending,
// image pull failures, crash loops) are still reported.
func (a *metricDataAccumulator) podStatus(pod *v1.Pod) {
	currentTime := pcommon.NewTimestampFromTime(a.time)

	if a.metricGroupsToCollect[PodMetricGroup] {
		mb := a.mbs.PodMetricsBuilder
		mb.RecordK8sPodPhaseDataPoint(currentTime, podPhaseToInt(pod.Status.Phase))
		mb.RecordK8sPodStatusReasonDataPoint(currentTime, podStatusReasonToInt(pod.Status.Reason))

		rb := mb.NewResourceBuilder()
		rb.SetK8sPodUID(string(pod.UID))
		rb.SetK8sPodName(pod.Name)
		rb.SetK8sNamespaceName(pod.Namespace)
		a.m = append(a.m, mb.Emit(metadata.WithResource(rb.Emit())))
	}

	if !a.metricGroupsToCollect[ContainerMetricGroup] {
		return
	}

	statuses := make(map[string]*v1.ContainerStatus, len(pod.Status.ContainerStatuses))
	for i := range pod.Status.ContainerStatuses {
		cs := &pod.Status.ContainerStatuses[i]
		statuses[cs.Name] = cs
	}

	for i := range pod.Spec.Containers {
		a.containerStatus(pod, &pod.Spec.Containers[i], statuses[pod.Spec.Containers[i].Name], currentTime)
	}
}

func (a *metricDataAccumulator) containerStatus(pod *v1.Pod, c *v1.Container, cs *v1.ContainerStatus, currentTime pcommon.Timestamp) {
	mb := a.mbs.ContainerMetricsBuilder

	if q, ok := c.Resources.Limits[v1.ResourceCPU]; ok {
		mb.RecordK8sContainerCPULimitDataPoint(currentTime, float64(q.MilliValue())/1000.0)
	}
	if q, ok := c.Resources.Requests[v1.ResourceCPU]; ok {
		mb.RecordK8sContainerCPURequestDataPoint(currentTime, float64(q.MilliValue())/1000.0)
	}
	if q, ok := c.Resources.Limits[v1.ResourceMemory]; ok {
		mb.RecordK8sContainerMemoryLimitDataPoint(currentTime, q.Value())
	}
	if q, ok := c.Resources.Requests[v1.ResourceMemory]; ok {
		mb.RecordK8sContainerMemoryRequestDataPoint(currentTime, q.Value())
	}
	if q, ok := c.Resources.Limits[v1.ResourceStorage]; ok {
		mb.RecordK8sContainerStorageLimitDataPoint(currentTime, q.Value())
	}
	if q, ok := c.Resources.Requests[v1.ResourceStorage]; ok {
		mb.RecordK8sContainerStorageRequestDataPoint(currentTime, q.Value())
	}

	if cs != nil {
		mb.RecordK8sContainerRestartsDataPoint(currentTime, int64(cs.RestartCount))
		mb.RecordK8sContainerReadyDataPoint(currentTime, boolToInt64(cs.Ready))

		reason := containerStatusReason(cs)
		for _, attrVal := range allContainerStatusReasons {
			var val int64
			if reason != "" && reason == attrVal.String() {
				val = 1
			}
			mb.RecordK8sContainerStatusReasonDataPoint(currentTime, val, attrVal)
		}
	}

	rb := mb.NewResourceBuilder()
	rb.SetK8sPodUID(string(pod.UID))
	rb.SetK8sPodName(pod.Name)
	rb.SetK8sNamespaceName(pod.Namespace)
	rb.SetK8sContainerName(c.Name)
	if a.metadata.Labels[MetadataLabelContainerID] && cs != nil && cs.ContainerID != "" {
		rb.SetContainerID(stripContainerID(cs.ContainerID))
	}
	a.m = append(a.m, mb.Emit(metadata.WithResource(rb.Emit())))
}

func containerStatusReason(cs *v1.ContainerStatus) string {
	switch {
	case cs.State.Terminated != nil:
		return cs.State.Terminated.Reason
	case cs.State.Waiting != nil:
		return cs.State.Waiting.Reason
	default:
		return ""
	}
}

func podPhaseToInt(phase v1.PodPhase) int64 {
	switch phase {
	case v1.PodPending:
		return 1
	case v1.PodRunning:
		return 2
	case v1.PodSucceeded:
		return 3
	case v1.PodFailed:
		return 4
	case v1.PodUnknown:
		return 5
	default:
		return 5
	}
}

func podStatusReasonToInt(reason string) int64 {
	switch reason {
	case "Evicted":
		return 1
	case "NodeAffinity":
		return 2
	case "NodeLost":
		return 3
	case "Shutdown":
		return 4
	case "UnexpectedAdmissionError":
		return 5
	default:
		return 6
	}
}

func boolToInt64(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
