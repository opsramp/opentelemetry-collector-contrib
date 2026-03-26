// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8seventsreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver"

import (
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	conventions "go.opentelemetry.io/otel/semconv/v1.38.0"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sinventory"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver/internal/metadata"
)

const (
	// Number of log attributes to add to the plog.LogRecordSlice.
	totalLogAttributes = 7

	// Number of resource attributes to add to the plog.ResourceLogs.
	totalResourceAttributes = 6
)

// By default k8s event has only two types of events (Normal, Warning), here are we allowing other types as well.
// For more info: https://github.com/kubernetes/api/blob/release-1.34/events/v1/types_swagger_doc_generated.go#L42
var severityMap = map[string]plog.SeverityNumber{
	"normal":   plog.SeverityNumberInfo,
	"warning":  plog.SeverityNumberWarn,
	"error":    plog.SeverityNumberError,
	"critical": plog.SeverityNumberFatal,
}

// k8sEventToLogRecord converts Kubernetes event to plog.LogRecordSlice and adds the resource attributes.
func k8sEventToLogData(logger *zap.Logger, ev *corev1.Event, version string, attributes []KeyValue) plog.Logs {
	ld := plog.NewLogs()
	rl := ld.ResourceLogs().AppendEmpty()
	sl := rl.ScopeLogs().AppendEmpty()
	sl.Scope().SetName(metadata.ScopeName)
	sl.Scope().SetVersion(version)
	lr := sl.LogRecords().AppendEmpty()

	resourceAttrs := rl.Resource().Attributes()
	resourceAttrs.EnsureCapacity(totalResourceAttributes)

	eventHost := ""
	if ev.ReportingInstance != "" {
		eventHost = ev.ReportingInstance
	} else if ev.Source.Host != "" {
		eventHost = ev.Source.Host
	}
	if eventHost != "" {
		resourceAttrs.PutStr(string(semconv.K8SNodeNameKey), eventHost)
	}
	// Attributes related to the object causing the event.
	resourceAttrs.PutStr("k8s.object.kind", ev.InvolvedObject.Kind)
	resourceAttrs.PutStr("k8s.object.name", ev.InvolvedObject.Name)
	resourceAttrs.PutStr("k8s.object.uid", string(ev.InvolvedObject.UID))
	resourceAttrs.PutStr("k8s.object.fieldpath", ev.InvolvedObject.FieldPath)
	resourceAttrs.PutStr("k8s.object.api_version", ev.InvolvedObject.APIVersion)
	resourceAttrs.PutStr("k8s.object.resource_version", ev.InvolvedObject.ResourceVersion)
	if ev.InvolvedObject.Namespace != "" {
		resourceAttrs.PutStr(string(semconv.K8SNamespaceNameKey), ev.InvolvedObject.Namespace)
	}

	//adding resource name and k8s.pod.name/ k8s.node.name/ k8s.daemonset.name etc accrording to kind
	object_key := "k8s." + strings.ToLower(ev.InvolvedObject.Kind) + ".name"
	object_value := ev.InvolvedObject.Name
	resourceAttrs.PutStr(object_key, object_value)
	resourceAttrs.PutStr("resourceName", object_value)

	resourceAttrs.PutStr("type", "event") // This should come from config. To be enhanced.

	lr.SetTimestamp(pcommon.NewTimestampFromTime(k8sinventory.GetEventTimestamp(ev)))

	// The Message field contains description about the event,
	// which is best suited for the "Body" of the LogRecordSlice.
	lr.Body().SetStr(ev.Message)

	// Set the "SeverityNumber" and "SeverityText" if a known type of
	// severity is found.
	if severityNumber, ok := severityMap[strings.ToLower(ev.Type)]; ok {
		lr.SetSeverityNumber(severityNumber)
		lr.SetSeverityText(ev.Type)
	} else {
		logger.Debug("unknown severity type", zap.String("type", ev.Type))
	}

	attrs := lr.Attributes()
	attrs.EnsureCapacity(totalLogAttributes)

	attrs.PutStr("k8s.event.type", ev.Type)
	attrs.PutStr("k8s.event.sourceComponent", ev.Source.Component)
	attrs.PutStr("k8s.event.reason", ev.Reason)
	attrs.PutStr("k8s.event.action", ev.Action)
	attrs.PutStr("k8s.event.start_time", ev.CreationTimestamp.String())
	attrs.PutStr("k8s.event.name", ev.Name)
	attrs.PutStr("k8s.event.uid", string(ev.UID))
	attrs.PutStr("level", ev.Type)
	attrs.PutStr(string(conventions.K8SNamespaceNameKey), ev.InvolvedObject.Namespace)


	// "Count" field of k8s event will be '0' in case it is
	// not present in the collected event from k8s.
	if ev.Count != 0 {
		attrs.PutInt("k8s.event.count", int64(ev.Count))
	}

	for _, kv := range attributes {
		val := pcommon.NewValueEmpty()
		err := val.FromRaw(kv.Value)
		if err != nil {
			continue
		}
		val.CopyTo(attrs.PutEmpty(kv.Key))
	}

	return ld
}
