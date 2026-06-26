// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sattributesprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"

import (
	"context"
	"strings"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/moid"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"
)

// ****************** Processing DCGM GPU metrics ******************

// GetResourceUuidUsingGPUMoid looks up the resource UUID for a GPU identified by its NVIDIA UUID (e.g. "GPU-0bbb5e8b-9833-a971-ef7f-e5ce00937b23").
func (op *kubernetesprocessor) GetResourceUuidUsingGPUMoid(ctx context.Context, gpuUUID string) (resourceUuid string) {
	gpuMoidKey := moid.NewMoid(op.redisConfig.ClusterName).WithGPUUUID(gpuUUID).GPUMoid()
	resourceUuid = op.redisClient.GetUuidValueInString(ctx, gpuMoidKey)
	op.logger.Debug("GPU MoID redis lookup", zap.String("moid", gpuMoidKey), zap.String("uuid", resourceUuid))
	return
}

// processDcgmMetricByGpu detects whether src contains DCGM_ metrics with per-data-point GPU UUID labels.
// If so it produces one ResourceMetrics per distinct GPU UUID, where each copy contains only the data-points for that GPU and has the correct resource uuid set.
// Returns (splits, true) when a split was performed; (nil, false) when src has no DCGM metrics.
func (kp *kubernetesprocessor) processDcgmMetricByGpu(ctx context.Context, src pmetric.ResourceMetrics) ([]pmetric.ResourceMetrics, bool) {
	// Collect all unique GPU UUIDs from DCGM_ data-point attributes.
	gpuUUIDs := collectDCGMGPUUUIDs(src)
	if len(gpuUUIDs) == 0 {
		return nil, false
	}

	kp.logger.Debug("Processing DCGM ResourceMetrics by GPU UUID", zap.Int("gpuCount", len(gpuUUIDs)))

	result := make([]pmetric.ResourceMetrics, 0, len(gpuUUIDs))
	for _, gpuUUID := range gpuUUIDs {
		gpuRM := pmetric.NewResourceMetrics()
		src.Resource().Attributes().CopyTo(gpuRM.Resource().Attributes())
		gpuRM.SetSchemaUrl(src.SchemaUrl())

		// Set the resource UUID for this specific GPU.
		gpuResourceUUID := kp.GetResourceUuidUsingGPUMoid(ctx, gpuUUID)
		if gpuResourceUUID != "" {
			gpuRM.Resource().Attributes().PutStr("uuid", gpuResourceUUID)
		}
		gpuRM.Resource().Attributes().PutStr("gpu.uuid", gpuUUID)

		// Copy ScopeMetrics, retaining only data-points that belong to this GPU.
		for i := 0; i < src.ScopeMetrics().Len(); i++ {
			srcSM := src.ScopeMetrics().At(i)
			dstSM := gpuRM.ScopeMetrics().AppendEmpty()
			srcSM.Scope().CopyTo(dstSM.Scope())
			dstSM.SetSchemaUrl(srcSM.SchemaUrl())

			for j := 0; j < srcSM.Metrics().Len(); j++ {
				srcM := srcSM.Metrics().At(j)
				dstM := dstSM.Metrics().AppendEmpty()
				srcM.CopyTo(dstM)
				if strings.HasPrefix(srcM.Name(), "DCGM_") {
					// Keep only data-points whose UUID attribute matches this GPU.
					filterMetricDataPointsByAttr(dstM, "UUID", gpuUUID)
				}
				// Non-DCGM metrics (if any) are included as-is in every split.
			}
		}

		result = append(result, gpuRM)
	}

	return result, true
}

// collectDCGMGPUUUIDs scans all DCGM_ metrics in rm and returns the list of distinct GPU UUIDs found in data-point "UUID" attributes.
func collectDCGMGPUUUIDs(rm pmetric.ResourceMetrics) []string {
	set := make(map[string]struct{})
	collect := func(attrs pcommon.Map) {
		if v, found := attrs.Get("UUID"); found && strings.HasPrefix(v.Str(), "GPU-") {
			set[v.Str()] = struct{}{}
		}
	}
	for i := 0; i < rm.ScopeMetrics().Len(); i++ {
		sm := rm.ScopeMetrics().At(i)
		for j := 0; j < sm.Metrics().Len(); j++ {
			m := sm.Metrics().At(j)
			if !strings.HasPrefix(m.Name(), "DCGM_") {
				continue
			}
			switch m.Type() {
			case pmetric.MetricTypeGauge:
				dps := m.Gauge().DataPoints()
				for k := 0; k < dps.Len(); k++ {
					collect(dps.At(k).Attributes())
				}
			case pmetric.MetricTypeSum:
				dps := m.Sum().DataPoints()
				for k := 0; k < dps.Len(); k++ {
					collect(dps.At(k).Attributes())
				}
			case pmetric.MetricTypeSummary:
				dps := m.Summary().DataPoints()
				for k := 0; k < dps.Len(); k++ {
					collect(dps.At(k).Attributes())
				}
			case pmetric.MetricTypeHistogram:
				dps := m.Histogram().DataPoints()
				for k := 0; k < dps.Len(); k++ {
					collect(dps.At(k).Attributes())
				}
			case pmetric.MetricTypeExponentialHistogram:
				dps := m.ExponentialHistogram().DataPoints()
				for k := 0; k < dps.Len(); k++ {
					collect(dps.At(k).Attributes())
				}
			}
		}
	}
	uuids := make([]string, 0, len(set))
	for u := range set {
		uuids = append(uuids, u)
	}
	return uuids
}

// ****************** Processing RDMA NIC metrics ******************

// GetResourceUuidUsingNICMoid looks up the resource UUID for a NIC identified by its node name and device name (e.g. nodeName="worker1.r06lef14.lan", device="mlx5_0").
func (op *kubernetesprocessor) GetResourceUuidUsingNICMoid(ctx context.Context, nodeName, device string) (resourceUuid string) {
	nicMoidKey := moid.NewMoid(op.redisConfig.ClusterName).WithNICName(nodeName + "-" + device).NICMoid()
	resourceUuid = op.redisClient.GetUuidValueInString(ctx, nicMoidKey)
	op.logger.Debug("NIC MoID redis lookup", zap.String("moid", nicMoidKey), zap.String("uuid", resourceUuid))
	return
}

// processRdmaMetricByNic detects whether src contains rdma_ metrics with per-data-point device labels.
// If so it produces one ResourceMetrics per distinct NIC device, where each copy contains only the data-points for that device and has the correct resource uuid set.
// Returns (perNicRMs, true) when processed; (nil, false) when src has no rdma_ metrics.
func (kp *kubernetesprocessor) processRdmaMetricByNic(ctx context.Context, src pmetric.ResourceMetrics) ([]pmetric.ResourceMetrics, bool) {
	devices := collectRDMANICDevices(src)
	if len(devices) == 0 {
		return nil, false
	}

	// Node name is a resource-level attribute set by processResource.
	nodeName := ""
	if v, found := src.Resource().Attributes().Get("k8s.node.name"); found {
		nodeName = v.Str()
	}

	kp.logger.Debug("Processing RDMA ResourceMetrics by NIC device",
		zap.Int("nicCount", len(devices)), zap.String("node", nodeName))

	result := make([]pmetric.ResourceMetrics, 0, len(devices))
	for _, device := range devices {
		nicRM := pmetric.NewResourceMetrics()
		src.Resource().Attributes().CopyTo(nicRM.Resource().Attributes())
		nicRM.SetSchemaUrl(src.SchemaUrl())

		// Set the resource UUID for this specific NIC.
		nicResourceUUID := kp.GetResourceUuidUsingNICMoid(ctx, nodeName, device)
		if nicResourceUUID != "" {
			nicRM.Resource().Attributes().PutStr("uuid", nicResourceUUID)
		}
		nicRM.Resource().Attributes().PutStr("nic.device", device)

		// Copy ScopeMetrics, retaining only data-points that belong to this NIC device.
		for i := 0; i < src.ScopeMetrics().Len(); i++ {
			srcSM := src.ScopeMetrics().At(i)
			dstSM := nicRM.ScopeMetrics().AppendEmpty()
			srcSM.Scope().CopyTo(dstSM.Scope())
			dstSM.SetSchemaUrl(srcSM.SchemaUrl())

			for j := 0; j < srcSM.Metrics().Len(); j++ {
				srcM := srcSM.Metrics().At(j)
				if !strings.HasPrefix(srcM.Name(), "rdma_") {
					// Only rdma_ metrics belong to NIC resources; skip go_*, process_*, scrape_*, etc.
					continue
				}
				dstM := dstSM.Metrics().AppendEmpty()
				srcM.CopyTo(dstM)
				// Keep only data-points whose device attribute matches this NIC.
				filterMetricDataPointsByAttr(dstM, "device", device)
			}
			// Drop rdma_ metrics that had no data-points for this device (e.g. rdma_scrape_errors_total which carries no device label).
			dstSM.Metrics().RemoveIf(func(m pmetric.Metric) bool {
				return dataPointCount(m) == 0
			})
		}

		result = append(result, nicRM)
	}

	return result, true
}

// collectRDMANICDevices scans all rdma_ metrics in rm and returns the list of distinct NIC device names found in data-point "device" attributes.
func collectRDMANICDevices(rm pmetric.ResourceMetrics) []string {
	set := make(map[string]struct{})
	collect := func(attrs pcommon.Map) {
		if v, found := attrs.Get("device"); found && v.Str() != "" {
			set[v.Str()] = struct{}{}
		}
	}
	for i := 0; i < rm.ScopeMetrics().Len(); i++ {
		sm := rm.ScopeMetrics().At(i)
		for j := 0; j < sm.Metrics().Len(); j++ {
			m := sm.Metrics().At(j)
			if !strings.HasPrefix(m.Name(), "rdma_") {
				continue
			}
			switch m.Type() {
			case pmetric.MetricTypeGauge:
				dps := m.Gauge().DataPoints()
				for k := 0; k < dps.Len(); k++ {
					collect(dps.At(k).Attributes())
				}
			case pmetric.MetricTypeSum:
				dps := m.Sum().DataPoints()
				for k := 0; k < dps.Len(); k++ {
					collect(dps.At(k).Attributes())
				}
			case pmetric.MetricTypeSummary:
				dps := m.Summary().DataPoints()
				for k := 0; k < dps.Len(); k++ {
					collect(dps.At(k).Attributes())
				}
			case pmetric.MetricTypeHistogram:
				dps := m.Histogram().DataPoints()
				for k := 0; k < dps.Len(); k++ {
					collect(dps.At(k).Attributes())
				}
			case pmetric.MetricTypeExponentialHistogram:
				dps := m.ExponentialHistogram().DataPoints()
				for k := 0; k < dps.Len(); k++ {
					collect(dps.At(k).Attributes())
				}
			}
		}
	}
	devices := make([]string, 0, len(set))
	for d := range set {
		devices = append(devices, d)
	}
	return devices
}

// filterMetricDataPointsByAttr removes all data-points from m that do not have
// an attribute with the given key equal to the given value.
func filterMetricDataPointsByAttr(m pmetric.Metric, attrKey, attrValue string) {
	keep := func(attrs pcommon.Map) bool {
		v, found := attrs.Get(attrKey)
		return found && v.Str() == attrValue
	}
	switch m.Type() {
	case pmetric.MetricTypeGauge:
		m.Gauge().DataPoints().RemoveIf(func(dp pmetric.NumberDataPoint) bool {
			return !keep(dp.Attributes())
		})
	case pmetric.MetricTypeSum:
		m.Sum().DataPoints().RemoveIf(func(dp pmetric.NumberDataPoint) bool {
			return !keep(dp.Attributes())
		})
	case pmetric.MetricTypeSummary:
		m.Summary().DataPoints().RemoveIf(func(dp pmetric.SummaryDataPoint) bool {
			return !keep(dp.Attributes())
		})
	case pmetric.MetricTypeHistogram:
		m.Histogram().DataPoints().RemoveIf(func(dp pmetric.HistogramDataPoint) bool {
			return !keep(dp.Attributes())
		})
	case pmetric.MetricTypeExponentialHistogram:
		m.ExponentialHistogram().DataPoints().RemoveIf(func(dp pmetric.ExponentialHistogramDataPoint) bool {
			return !keep(dp.Attributes())
		})
	}
}

// dataPointCount returns the total number of data-points in m regardless of metric type.
func dataPointCount(m pmetric.Metric) int {
	switch m.Type() {
	case pmetric.MetricTypeGauge:
		return m.Gauge().DataPoints().Len()
	case pmetric.MetricTypeSum:
		return m.Sum().DataPoints().Len()
	case pmetric.MetricTypeSummary:
		return m.Summary().DataPoints().Len()
	case pmetric.MetricTypeHistogram:
		return m.Histogram().DataPoints().Len()
	case pmetric.MetricTypeExponentialHistogram:
		return m.ExponentialHistogram().DataPoints().Len()
	}
	return 0
}
