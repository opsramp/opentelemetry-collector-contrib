// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package kubelet // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/kubelet"

import (
	"time"

	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
	stats "k8s.io/kubelet/pkg/apis/stats/v1alpha1"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kubeletstatsreceiver/internal/metadata"
)

func MetricsData(
	logger *zap.Logger, summary *stats.Summary,
	metadata Metadata,
	metricGroupsToCollect map[MetricGroup]bool,
	allNetworkInterfaces map[MetricGroup]bool,
	mbs *metadata.MetricsBuilders,
	cpuUsageCalculator *CPUUsageCalculator,
	podStatusEnabled bool,
) []pmetric.Metrics {
	cpuUsageCalculator.startScrape()
	defer cpuUsageCalculator.endScrape()

	acc := &metricDataAccumulator{
		metadata:              metadata,
		logger:                logger,
		metricGroupsToCollect: metricGroupsToCollect,
		allNetworkInterfaces:  allNetworkInterfaces,
		time:                  time.Now(),
		mbs:                   mbs,
		cpuUsageCalculator:    cpuUsageCalculator,
	}
	acc.nodeStats(summary.Node)
	for i := range summary.Pods {
		pod := &summary.Pods[i]
		acc.podStats(pod)
		for j := range pod.Containers {
			containerStats := &pod.Containers[j]
			// propagate the pod resource down to the container
			acc.containerStats(pod, containerStats)
		}

		for j := range pod.VolumeStats {
			volumeStats := &pod.VolumeStats[j]
			// propagate the pod resource down to the container
			acc.volumeStats(pod, volumeStats)
		}
	}

	// Pod state metrics come from the /pods endpoint rather than the stats summary so that pods and
	// containers without usage stats are still reported.
	if podStatusEnabled && acc.metadata.PodsMetadata != nil {
		for i := range acc.metadata.PodsMetadata.Items {
			acc.podStatus(&acc.metadata.PodsMetadata.Items[i])
		}
	}

	return acc.m
}
