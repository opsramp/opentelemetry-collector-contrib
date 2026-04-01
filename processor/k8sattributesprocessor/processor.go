// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sattributesprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/cache"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/kube"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/moid"
	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/redis"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componentstatus"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/pprofile"
	"go.opentelemetry.io/collector/pdata/ptrace"
	semconv "go.opentelemetry.io/collector/semconv/v1.5.0"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	conventions "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/metadata"
)

const (
	clientIPLabelName string = "ip"
)

type kubernetesprocessor struct {
	cfg                    component.Config
	options                []option
	telemetrySettings      component.TelemetrySettings
	telemetry              *metadata.TelemetryBuilder
	logger                 *zap.Logger
	apiConfig              k8sconfig.APIConfig
	kc                     kube.Client
	passthroughMode        bool
	rules                  kube.ExtractionRules
	filters                kube.Filters
	addons                 []kube.AddOnMetadata
	redisConfig            redis.OpsrampRedisConfig
	podAssociations        []kube.Association
	podIgnore              kube.Excludes
	waitForMetadata        bool
	waitForMetadataTimeout time.Duration
	redisClient            *redis.Client
}

func (kp *kubernetesprocessor) initKubeClient(set component.TelemetrySettings, kubeClient kube.ClientProvider) error {
	if kubeClient == nil {
		kubeClient = kube.New
	}
	if !kp.passthroughMode {
		kc, err := kubeClient(set, kp.apiConfig, kp.rules, kp.filters, kp.podAssociations, kp.podIgnore, nil, kube.InformersFactoryList{}, kp.waitForMetadata, kp.waitForMetadataTimeout)
		if err != nil {
			return err
		}
		kp.kc = kc
	}
	return nil
}

func (kp *kubernetesprocessor) Start(_ context.Context, host component.Host) error {
	if metadata.ProcessorK8sattributesDontEmitV0K8sConventionsFeatureGate.IsEnabled() && !metadata.ProcessorK8sattributesEmitV1K8sConventionsFeatureGate.IsEnabled() {
		err := errors.New("processor.k8sattributes.DontEmitV0K8sConventions cannot be enabled without enabling processor.k8sattributes.EmitV1K8sConventions")
		kp.logger.Error("Invalid feature gate combination", zap.Error(err))
		componentstatus.ReportStatus(host, componentstatus.NewFatalErrorEvent(err))
		return err
	}

	if kp.rules.ContainerImageTag {
		kp.logger.Warn(
			"[WARNING] container.image.tag is being renamed to container.image.tags per Semantic Conventions. " +
				"Consider switching to the new schema by enabling the processor.k8sattributes.EmitV1K8sConventions and " +
				"processor.k8sattributes.DontEmitV0K8sConventions feature gates. " +
				"See processor README section 'Semantic Conventions Compatibility' for details.",
		)
	}
	if len(kp.rules.Labels) > 0 {
		kp.logger.Warn(
			"[WARNING] Pod label extraction attributes are being renamed (e.g. k8s.pod.labels.<key> -> k8s.pod.label.<key>) per Semantic Conventions. " +
				"Consider switching to the new schema by enabling the processor.k8sattributes.EmitV1K8sConventions and " +
				"processor.k8sattributes.DontEmitV0K8sConventions feature gates. " +
				"See processor README section 'Semantic Conventions Compatibility' for details.",
		)
	}
	if len(kp.rules.Annotations) > 0 {
		kp.logger.Warn(
			"[WARNING] Pod annotation extraction attributes are being renamed (e.g. k8s.pod.annotations.<key> -> k8s.pod.annotation.<key>) per Semantic Conventions. " +
				"Consider switching to the new schema by enabling the processor.k8sattributes.EmitV1K8sConventions and " +
				"processor.k8sattributes.DontEmitV0K8sConventions feature gates. " +
				"See processor README section 'Semantic Conventions Compatibility' for details.",
		)
	}

	allOptions := append(createProcessorOpts(kp.cfg), kp.options...)

	for _, opt := range allOptions {
		if err := opt(kp); err != nil {
			kp.logger.Error("Could not apply option", zap.Error(err))
			componentstatus.ReportStatus(host, componentstatus.NewFatalErrorEvent(err))
			return err
		}
	}

	kp.logger.Info("ops k8s attr processor start", zap.Any("redisHost", kp.redisConfig.RedisHost),
		zap.Any("redisPort", kp.redisConfig.RedisPort),
		zap.Any("redisPass", kp.redisConfig.RedisPass))

	// cache := lru.GetInstance(kp.redisConfig.LruCacheSize, kp.redisConfig.LruExpirationTime)

	// if cache == nil {
	// 	kp.logger.Error("Failed to initilize the cache with GetInstance()")
	// 	return nil
	// }

	// kp.redisClient = redis.NewClient(kp.logger, cache, kp.redisConfig.RedisHost, kp.redisConfig.RedisPort, kp.redisConfig.RedisPass)
	cacheObj := cache.GetCacheInstance(kp.redisConfig.PrimaryCacheSize, kp.redisConfig.PrimaryCacheEvictionTime, kp.redisConfig.SecondaryCacheSize, kp.redisConfig.SecondaryCacheEvictionTime)

	if cacheObj == nil {
		kp.logger.Error("Failed to initilize the cache with GetInstance()")
		return fmt.Errorf("failed to initialize cache")
	}

	kp.redisClient = redis.NewClient(kp.logger, cacheObj, kp.redisConfig.RedisHost, kp.redisConfig.RedisPort, kp.redisConfig.RedisPass, kp.redisConfig.PrimaryCacheEvictionTime, kp.redisConfig.SecondaryCacheEvictionTime)

	// This might have been set by an option already
	if kp.kc == nil {
		err := kp.initKubeClient(kp.telemetrySettings, kubeClientProvider)
		if err != nil {
			kp.logger.Error("Could not initialize kube client", zap.Error(err))
			componentstatus.ReportStatus(host, componentstatus.NewFatalErrorEvent(err))
			return err
		}
	}
	if !kp.passthroughMode {
		err := kp.kc.Start()
		if err != nil {
			componentstatus.ReportStatus(host, componentstatus.NewFatalErrorEvent(err))
			return err
		}
	}
	return nil
}

func (kp *kubernetesprocessor) Shutdown(context.Context) error {
	if kp.telemetry != nil {
		kp.telemetry.Shutdown()
	}
	if kp.kc == nil {
		return nil
	}
	if !kp.passthroughMode {
		kp.kc.Stop()
	}
	return nil
}

// processTraces process traces and add k8s metadata using resource IP or incoming IP as pod origin.
func (kp *kubernetesprocessor) processTraces(ctx context.Context, td ptrace.Traces) (ptrace.Traces, error) {
	rss := td.ResourceSpans()
	for i := 0; i < rss.Len(); i++ {
		kp.processResource(ctx, rss.At(i).Resource(), "traces")
		kp.processTraceResources(ctx, rss.At(i).Resource())

	}

	return td, nil
}

// processMetrics process metrics and add k8s metadata using resource IP, hostname or incoming IP as pod origin.
func (kp *kubernetesprocessor) processMetrics(ctx context.Context, md pmetric.Metrics) (pmetric.Metrics, error) {
	rm := md.ResourceMetrics()
	for i := 0; i < rm.Len(); i++ {
		kp.processResource(ctx, rm.At(i).Resource(), "metrics")
		kp.processopsrampResources(ctx, rm.At(i).Resource())
	}

	kp.filterOnlyOpsrampMetrics(md)

	return md, nil
}

// processLogs process logs and add k8s metadata using resource IP, hostname or incoming IP as pod origin.
func (kp *kubernetesprocessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	rl := ld.ResourceLogs()
	for i := 0; i < rl.Len(); i++ {
		kp.processResource(ctx, rl.At(i).Resource(), "logs")
		kp.processopsrampResources(ctx, rl.At(i).Resource())
		kp.addOpsrampEventResourceAttributes(ctx, rl.At(i).Resource())
		kp.processEventBody(rl.At(i))
	}

	return ld, nil
}

// processProfiles process profiles and add k8s metadata using resource IP, hostname or incoming IP as pod origin.
func (kp *kubernetesprocessor) processProfiles(ctx context.Context, pd pprofile.Profiles) (pprofile.Profiles, error) {
	rp := pd.ResourceProfiles()
	for i := 0; i < rp.Len(); i++ {
		kp.processResource(ctx, rp.At(i).Resource(), "profiles")
	}

	return pd, nil
}

// processResource adds Pod metadata tags to resource based on pod association configuration
func (kp *kubernetesprocessor) processResource(ctx context.Context, resource pcommon.Resource, signalType string) {
	resource.Attributes().Range(func(k string, v pcommon.Value) bool {
		kp.logger.Debug("res-attributes", zap.Any(k, v.Str()))
		return true
	})

	if val, found := resource.Attributes().Get("type"); found && val.Str() == "event" {
		if kind, found := resource.Attributes().Get("k8s.object.kind"); found {
			if objectuid, found := resource.Attributes().Get("k8s.object.uid"); found {
				if kind.Str() == "Pod" {
					resource.Attributes().PutStr("k8s.pod.uid", objectuid.Str())
				} else if kind.Str() == "Node" {
					resource.Attributes().PutStr("k8s.node.uid", objectuid.Str())
				}
			}
		}
	}

	podIdentifierValue := extractPodID(ctx, resource.Attributes(), kp.podAssociations)
	kp.logger.Debug("evaluating pod identifier", zap.Any("value", podIdentifierValue))

	for i := range podIdentifierValue {
		if podIdentifierValue[i].Source.From == kube.ConnectionSource && podIdentifierValue[i].Value != "" {
			if kp.passthroughMode || kp.rules.PodIP {
				setResourceAttribute(resource.Attributes(), string(conventions.K8SPodIPKey), podIdentifierValue[i].Value)
			}
			break
		}
	}
	if kp.passthroughMode {
		return
	}

	var pod *kube.Pod
	var podFound bool
	podIdentifierStr := buildPodIdentifierString(podIdentifierValue)
	if podIdentifierValue.IsNotEmpty() {
		if pod, podFound = kp.kc.GetPod(podIdentifierValue); podFound {
			kp.logger.Debug("getting the pod", zap.Any("pod", pod))

			// Record successful pod association
			if kp.telemetry != nil {
				successAttr := metric.WithAttributes(
					attribute.String("status", "success"),
					attribute.String("pod_identifier", podIdentifierStr),
					attribute.String("otelcol.signal", signalType),
				)
				kp.telemetry.K8sPodAssociation.Add(ctx, 1, successAttr)
			}

			for key, val := range pod.Attributes {
				setResourceAttribute(resource.Attributes(), key, val)
			}
			kp.addContainerAttributes(resource.Attributes(), pod)
		} else {
			// Record failed pod association
			kp.logger.Debug("pod not found", zap.Any("podIdentifier", podIdentifierValue))
			if kp.telemetry != nil {
				errorAttr := metric.WithAttributes(
					attribute.String("status", "error"),
					attribute.String("pod_identifier", podIdentifierStr),
					attribute.String("otelcol.signal", signalType),
				)
				kp.telemetry.K8sPodAssociation.Add(ctx, 1, errorAttr)
			}
		}
	} else {
		// Record failed pod association when no identifier found
		kp.logger.Debug("no pod identifier found")
		if kp.telemetry != nil {
			errorAttr := metric.WithAttributes(
				attribute.String("status", "error"),
				attribute.String("pod_identifier", podIdentifierStr),
				attribute.String("otelcol.signal", signalType),
			)
			kp.telemetry.K8sPodAssociation.Add(ctx, 1, errorAttr)
		}
	}

	namespace := getNamespace(pod, resource.Attributes())
	if namespace != "" {
		attrsToAdd := kp.getAttributesForPodsNamespace(namespace)
		for key, val := range attrsToAdd {
			setResourceAttribute(resource.Attributes(), key, val)
		}

		if kp.rules.ServiceNamespace {
			setResourceAttribute(resource.Attributes(), string(conventions.ServiceNamespaceKey), namespace)
		}
	}

	nodeName := getNodeName(pod, resource.Attributes())
	if nodeName != "" {
		attrsToAdd := kp.getAttributesForPodsNode(nodeName)
		for key, val := range attrsToAdd {
			setResourceAttribute(resource.Attributes(), key, val)
		}
		nodeUID := kp.getUIDForPodsNode(nodeName)
		if nodeUID != "" {
			setResourceAttribute(resource.Attributes(), string(conventions.K8SNodeUIDKey), nodeUID)
		}
	}

	deployment := getDeploymentUID(pod, resource.Attributes())
	if deployment != "" {
		attrsToAdd := kp.getAttributesForPodsDeployment(deployment)
		for key, val := range attrsToAdd {
			setResourceAttribute(resource.Attributes(), key, val)
		}
	}

	statefulset := getStatefulSetUID(pod, resource.Attributes())
	if statefulset != "" {
		attrsToAdd := kp.getAttributesForPodsStatefulSet(statefulset)
		for key, val := range attrsToAdd {
			setResourceAttribute(resource.Attributes(), key, val)
		}
	}

	daemonset := getDaemonSetUID(pod, resource.Attributes())
	if daemonset != "" {
		attrsToAdd := kp.getAttributesForPodsDaemonSet(daemonset)
		for key, val := range attrsToAdd {
			setResourceAttribute(resource.Attributes(), key, val)
		}
	}

	job := getJobUID(pod, resource.Attributes())
	if job != "" {
		attrsToAdd := kp.getAttributesForPodsJob(job)
		for key, val := range attrsToAdd {
			setResourceAttribute(resource.Attributes(), key, val)
		}
	}
}

func setResourceAttribute(attributes pcommon.Map, key, val string) {
	attr, found := attributes.Get(key)
	if !found || attr.AsString() == "" {
		attributes.PutStr(key, val)
	}
}

func getNamespace(pod *kube.Pod, resAttrs pcommon.Map) string {
	if pod != nil && pod.Namespace != "" {
		return pod.Namespace
	}
	return stringAttributeFromMap(resAttrs, string(conventions.K8SNamespaceNameKey))
}

func getNodeName(pod *kube.Pod, resAttrs pcommon.Map) string {
	if pod != nil && pod.NodeName != "" {
		return pod.NodeName
	}
	return stringAttributeFromMap(resAttrs, string(conventions.K8SNodeNameKey))
}

func getDeploymentUID(pod *kube.Pod, resAttrs pcommon.Map) string {
	if pod != nil && pod.DeploymentUID != "" {
		return pod.DeploymentUID
	}
	return stringAttributeFromMap(resAttrs, string(conventions.K8SDeploymentUIDKey))
}

func getStatefulSetUID(pod *kube.Pod, resAttrs pcommon.Map) string {
	if pod != nil && pod.StatefulSetUID != "" {
		return pod.StatefulSetUID
	}
	return stringAttributeFromMap(resAttrs, string(conventions.K8SStatefulSetUIDKey))
}

func getDaemonSetUID(pod *kube.Pod, resAttrs pcommon.Map) string {
	if pod != nil && pod.DaemonSetUID != "" {
		return pod.DaemonSetUID
	}
	return stringAttributeFromMap(resAttrs, string(conventions.K8SDaemonSetUIDKey))
}

func getJobUID(pod *kube.Pod, resAttrs pcommon.Map) string {
	if pod != nil && pod.JobUID != "" {
		return pod.JobUID
	}
	return stringAttributeFromMap(resAttrs, string(conventions.K8SJobUIDKey))
}

// addContainerAttributes looks if pod has any container identifiers and adds additional container attributes
func (kp *kubernetesprocessor) addContainerAttributes(attrs pcommon.Map, pod *kube.Pod) {
	containerName := stringAttributeFromMap(attrs, string(conventions.K8SContainerNameKey))
	containerID := stringAttributeFromMap(attrs, string(conventions.ContainerIDKey))
	var (
		containerSpec *kube.Container
		ok            bool
	)
	switch {
	case containerName != "":
		containerSpec, ok = pod.Containers.ByName[containerName]
		if !ok {
			return
		}
	case containerID != "":
		containerSpec, ok = pod.Containers.ByID[containerID]
		if !ok {
			return
		}
	// if there is only one container in the pod, we can fall back to that container
	case len(pod.Containers.ByID) == 1:
		for _, c := range pod.Containers.ByID {
			containerSpec = c
		}
	case len(pod.Containers.ByName) == 1:
		for _, c := range pod.Containers.ByName {
			containerSpec = c
		}
	default:
		return
	}
	if containerSpec.Name != "" {
		setResourceAttribute(attrs, string(conventions.K8SContainerNameKey), containerSpec.Name)
	}
	if containerSpec.ImageName != "" {
		setResourceAttribute(attrs, string(conventions.ContainerImageNameKey), containerSpec.ImageName)
	}
	enableStable := metadata.ProcessorK8sattributesEmitV1K8sConventionsFeatureGate.IsEnabled()
	disableLegacy := metadata.ProcessorK8sattributesDontEmitV0K8sConventionsFeatureGate.IsEnabled()
	if !disableLegacy && containerSpec.ImageTag != "" {
		setResourceAttribute(attrs, containerImageTag, containerSpec.ImageTag)
	}
	if enableStable && len(containerSpec.ImageTags) > 0 {
		sliceVal := attrs.PutEmptySlice(string(conventions.ContainerImageTagsKey))
		for _, tag := range containerSpec.ImageTags {
			sliceVal.AppendEmpty().SetStr(tag)
		}
	}
	if containerSpec.ServiceInstanceID != "" {
		setResourceAttribute(attrs, string(conventions.ServiceInstanceIDKey), containerSpec.ServiceInstanceID)
	}
	if containerSpec.ServiceVersion != "" {
		setResourceAttribute(attrs, string(conventions.ServiceVersionKey), containerSpec.ServiceVersion)
	}
	// attempt to get container ID from restart count
	runID := -1
	runIDAttr, ok := attrs.Get(string(conventions.K8SContainerRestartCountKey))
	if ok {
		containerRunID, err := intFromAttribute(runIDAttr)
		if err != nil {
			kp.logger.Debug(err.Error())
		} else {
			runID = containerRunID
		}
	} else {
		// take the highest runID (restart count) which represents the currently running container in most cases
		for containerRunID := range containerSpec.Statuses {
			if containerRunID > runID {
				runID = containerRunID
			}
		}
	}
	if runID != -1 {
		if containerStatus, ok := containerSpec.Statuses[runID]; ok {
			if _, found := attrs.Get(string(conventions.ContainerIDKey)); !found && containerStatus.ContainerID != "" {
				attrs.PutStr(string(conventions.ContainerIDKey), containerStatus.ContainerID)
			}
			if _, found := attrs.Get(string(conventions.ContainerImageRepoDigestsKey)); !found && containerStatus.ImageRepoDigest != "" {
				attrs.PutEmptySlice(string(conventions.ContainerImageRepoDigestsKey)).AppendEmpty().SetStr(containerStatus.ImageRepoDigest)
			}
		}
	}
}

func (kp *kubernetesprocessor) getAttributesForPodsNamespace(namespace string) map[string]string {
	ns, ok := kp.kc.GetNamespace(namespace)
	if !ok {
		return nil
	}
	return ns.Attributes
}

func (kp *kubernetesprocessor) getAttributesForPodsNode(nodeName string) map[string]string {
	node, ok := kp.kc.GetNode(nodeName)
	if !ok {
		return nil
	}
	return node.Attributes
}

func (kp *kubernetesprocessor) getAttributesForPodsDeployment(deploymentUID string) map[string]string {
	d, ok := kp.kc.GetDeployment(deploymentUID)
	if !ok {
		return nil
	}
	return d.Attributes
}

func (kp *kubernetesprocessor) getAttributesForPodsStatefulSet(statefulsetUID string) map[string]string {
	d, ok := kp.kc.GetStatefulSet(statefulsetUID)
	if !ok {
		return nil
	}
	return d.Attributes
}

func (kp *kubernetesprocessor) getAttributesForPodsDaemonSet(daemonsetUID string) map[string]string {
	d, ok := kp.kc.GetDaemonSet(daemonsetUID)
	if !ok {
		return nil
	}
	return d.Attributes
}

func (kp *kubernetesprocessor) getAttributesForPodsJob(jobUID string) map[string]string {
	j, ok := kp.kc.GetJob(jobUID)
	if !ok {
		return nil
	}
	return j.Attributes
}

func (kp *kubernetesprocessor) getUIDForPodsNode(nodeName string) string {
	node, ok := kp.kc.GetNode(nodeName)
	if !ok {
		return ""
	}
	return node.NodeUID
}

// buildPodIdentifierString combines all identifier values into a comma-separated string
func buildPodIdentifierString(podIdentifierValue kube.PodIdentifier) string {
	var identifiers []string
	for i := range podIdentifierValue {
		if podIdentifierValue[i].Value != "" {
			identifiers = append(identifiers, podIdentifierValue[i].Value)
		}
	}
	if len(identifiers) > 0 {
		return strings.Join(identifiers, ",")
	}
	return "unknown"
}

// intFromAttribute extracts int value from an attribute stored as string or int
func intFromAttribute(val pcommon.Value) (int, error) {
	switch val.Type() {
	case pcommon.ValueTypeInt:
		return int(val.Int()), nil
	case pcommon.ValueTypeStr:
		i, err := strconv.Atoi(val.Str())
		if err != nil {
			return 0, err
		}
		return i, nil
	default:
		return 0, fmt.Errorf("wrong attribute type %v, expected int", val.Type())
	}
}

func (kp *kubernetesprocessor) processEventBody(resourceLogs plog.ResourceLogs) {
	if val, found := resourceLogs.Resource().Attributes().Get("type"); found && val.Str() == "event" {

		bodyMap := map[string]string{}

		ilss := resourceLogs.ScopeLogs()
		for j := 0; j < ilss.Len(); j++ {
			ils := ilss.At(j)
			logs := ils.LogRecords()
			for k := 0; k < logs.Len(); k++ {
				lr := logs.At(k)

				bodyMap["message"] = lr.Body().AsString()

				body, err := json.Marshal(bodyMap)
				if err != nil {
					kp.logger.Error("Failed to marshal attributes as body ")
				}

				lr.Body().SetStr(string(body))
			}
		}
	}
}

func (kp *kubernetesprocessor) addOpsrampEventResourceAttributes(ctx context.Context, resource pcommon.Resource) {

	if val, found := resource.Attributes().Get("type"); found && val.Str() == "event" {
		resource.Attributes().PutStr("source", "kubernetes")

		host := ""
		if val, found := resource.Attributes().Get(semconv.AttributeK8SNodeName); found {
			host = val.Str()
			if host != "" {
				//overwrite node opsramp resource UUID in resourceUUID
				resource.Attributes().PutStr("k8s.node.name", host)

				if resourceUuid := kp.GetResourceUuidUsingResourceNodeMoid(ctx, resource); resourceUuid != "" {
					resource.Attributes().PutStr("k8s.node.resourceUUID", resourceUuid)
				}
			}
		}
	}
}

// processResource adds Pod metadata tags to resource based on pod association configuration
func (kp *kubernetesprocessor) processopsrampResources(ctx context.Context, resource pcommon.Resource) {

	var resourceUuid string

	for _, addon := range kp.addons {
		// If receiver has already added some attributes with some value, then we do not overwrite here.
		// For ex. type = event is already added for kube events. We should not overwrite it with type = RESOURCE.
		kp.logger.Debug("addon", zap.Any("key", addon.Key))

		if _, found := resource.Attributes().Get(addon.Key); !found {
			kp.logger.Debug("addon not found adding it", zap.Any("key", addon.Key))

			resource.Attributes().PutStr(addon.Key, addon.Value)
		}
	}
	var resourceType string

	if _, found := resource.Attributes().Get("map.to.namespace"); found {
		resourceUuid = kp.GetResourceUuidUsingNamespaceMoid(ctx, resource)
		resourceType = "namespace"
	} else if podname, found := resource.Attributes().Get("k8s.pod.name"); found {
		if resourceUuid = kp.GetResourceUuidUsingPodMoid(ctx, resource); resourceUuid == "" {
			kp.logger.Debug("opsramp resourceuuid not found in redis", zap.Any("podname", podname.Str()))
		}
		resourceType = "pod"
	} else if nodename, found := resource.Attributes().Get("k8s.node.name"); found {
		if resourceUuid = kp.GetResourceUuidUsingResourceNodeMoid(ctx, resource); resourceUuid == "" {
			kp.logger.Debug("opsramp resourceuuid not found in redis", zap.Any("nodename", nodename.Str()))
		}
		resourceType = "node"
	} else if dpname, found := resource.Attributes().Get("k8s.deployment.name"); found {
		if resourceUuid = kp.GetResourceUuidUsingWorkloadMoid(ctx, resource, dpname, "deployment"); resourceUuid == "" {
			kp.logger.Debug("opsramp resourceuuid not found in redis", zap.Any("deployment", dpname.Str()))
		}
		resourceType = "deployment"
	} else if rsname, found := resource.Attributes().Get("k8s.replicaset.name"); found {
		if resourceUuid = kp.GetResourceUuidUsingWorkloadMoid(ctx, resource, rsname, "replicaset"); resourceUuid == "" {
			kp.logger.Debug("opsramp resourceuuid not found in redis", zap.Any("replicaset", rsname.Str()))
		}
		resourceType = "replicaset"
	} else if ssname, found := resource.Attributes().Get("k8s.statefulset.name"); found {
		if resourceUuid = kp.GetResourceUuidUsingWorkloadMoid(ctx, resource, ssname, "statefulset"); resourceUuid == "" {
			kp.logger.Debug("opsramp resourceuuid not found in redis", zap.Any("statefulset", ssname.Str()))
		}
		resourceType = "statefulset"
	} else if dsname, found := resource.Attributes().Get("k8s.daemonset.name"); found {
		if resourceUuid = kp.GetResourceUuidUsingWorkloadMoid(ctx, resource, dsname, "daemonset"); resourceUuid == "" {
			kp.logger.Debug("opsramp resourceuuid not found in redis", zap.Any("daemonset", dsname.Str()))
		}
		resourceType = "daemonset"
	} else {
		if resourceUuid = kp.redisConfig.ClusterUid; resourceUuid == "" {
			kp.logger.Debug("opsramp resourceuuid not found", zap.Any("clustername", kp.redisConfig.ClusterName))
		}
		resourceType = "cluster"

		/*
			No need to get it from redis. As its directly available in config.
			if resourceUuid = kp.GetResourceUuidUsingClusterMoid(ctx, resource); resourceUuid == "" {
				kp.logger.Debug("opsramp resourceuuid not found in redis", zap.Any("clustername", kp.redisConfig.ClusterName))
			}
		*/
	}

	if resourceUuid != "" {
		if val, found := resource.Attributes().Get("type"); found {
			if val.Str() == "event" || val.Str() == "log" {
				resource.Attributes().PutStr("resourceUUID", resourceUuid)
				resource.Attributes().PutStr("k8s."+resourceType+".resourceUUID", resourceUuid)
			} else {
				resource.Attributes().PutStr("uuid", resourceUuid)
			}
		} else {
			kp.logger.Debug("type resource attribute not found hence not adding uuid/resourceuuid")

		}
	}

}

// processTraceResource adds Pod metadata tags to resource based on pod association configuration
func (kp *kubernetesprocessor) processTraceResources(ctx context.Context, resource pcommon.Resource) {
	var resourceUuid, resourceName string

	for _, addon := range kp.addons {
		// If receiver has already added some attributes with some value, then we do not overwrite here.
		// For ex. type = event is already added for kube events. We should not overwrite it with type = RESOURCE.
		kp.logger.Debug("addon", zap.Any("key", addon.Key))

		if _, found := resource.Attributes().Get(addon.Key); !found {
			kp.logger.Debug("addon not found adding it", zap.Any("key", addon.Key))

			resource.Attributes().PutStr(addon.Key, addon.Value)
		}
	}

	if _, found := resource.Attributes().Get("k8s.pod.name"); found {
		resourceUuid = kp.GetResourceUuidUsingPodMoid(ctx, resource)
		podname, _ := resource.Attributes().Get("k8s.pod.name")
		if resourceUuid != "" {
			resourceName = podname.Str()
			resource.Attributes().PutStr("k8s.pod.resourceUUID", resourceUuid)
			kp.addAdditionalResourceUuid(ctx, resource)
		} else {
			kp.logger.Debug("opsramp resourceuuid not found in redis", zap.Any("podname", podname.Str()))
		}
	}

	if resourceUuid != "" {
		resource.Attributes().PutStr("resourceUUID", resourceUuid)
	}
	if resourceName != "" {
		resource.Attributes().PutStr("resourceName", resourceName)
	}
}

// function to add resourceuuid to the resource
func (kp *kubernetesprocessor) addAdditionalResourceUuid(ctx context.Context, resource pcommon.Resource) {
	var additionalResourceUuid string

	if dpName, found := resource.Attributes().Get("k8s.deployment.name"); found {
		additionalResourceUuid = kp.GetResourceUuidUsingWorkloadMoid(ctx, resource, dpName, "deployment")
		if additionalResourceUuid != "" {
			resource.Attributes().PutStr("k8s.deployment.resourceUUID", additionalResourceUuid)
			if rsName, found := resource.Attributes().Get("k8s.replicaset.name"); found {
				additionalResourceUuid = kp.GetResourceUuidUsingWorkloadMoid(ctx, resource, rsName, "replicaset")
				if additionalResourceUuid != "" {
					resource.Attributes().PutStr("k8s.replicaset.resourceUUID", additionalResourceUuid)
				}
			}
			return
		}
	}

	if rsName, found := resource.Attributes().Get("k8s.replicaset.name"); found {
		additionalResourceUuid = kp.GetResourceUuidUsingWorkloadMoid(ctx, resource, rsName, "replicaset")
		if additionalResourceUuid != "" {
			resource.Attributes().PutStr("k8s.replicaset.resourceUUID", additionalResourceUuid)
			return
		}
	}

	if ssName, found := resource.Attributes().Get("k8s.statefulset.name"); found {
		additionalResourceUuid = kp.GetResourceUuidUsingWorkloadMoid(ctx, resource, ssName, "statefulset")
		if additionalResourceUuid != "" {
			resource.Attributes().PutStr("k8s.statefulset.resourceUUID", additionalResourceUuid)
			return
		}
	}

	if dsName, found := resource.Attributes().Get("k8s.daemonset.name"); found {
		additionalResourceUuid = kp.GetResourceUuidUsingWorkloadMoid(ctx, resource, dsName, "daemonset")
		if additionalResourceUuid != "" {
			resource.Attributes().PutStr("k8s.daemonset.resourceUUID", additionalResourceUuid)
			return
		}
	}
}

func (kp *kubernetesprocessor) filterOnlyOpsrampMetrics(md pmetric.Metrics) {
	md.ResourceMetrics().RemoveIf(func(rmetrics pmetric.ResourceMetrics) bool {
		resource := rmetrics.Resource()
		if _, found := resource.Attributes().Get("uuid"); !found {
			return true
		}
		return false
	})
}

func (op *kubernetesprocessor) GetResourceUuidUsingPodMoid(ctx context.Context, resource pcommon.Resource) (resourceUuid string) {
	var namespace, podname pcommon.Value
	var found bool

	if namespace, found = resource.Attributes().Get("k8s.namespace.name"); !found {
		op.logger.Debug("k8s.namespace.name not found in resource attributes hence not able to get resource uuid using pod moid")
		return
	}
	if podname, found = resource.Attributes().Get("k8s.pod.name"); !found {
		return
	}

	podMoid := moid.NewMoid(op.redisConfig.ClusterName).WithNamespaceName(namespace.Str()).WithPodName(podname.Str())

	//removed workload name in POD MoId
	/*if rsname, found = resource.Attributes().Get("k8s.replicaset.name"); found {
		podMoid.WithReplicasetName(rsname.Str())
	} else if dsname, found = resource.Attributes().Get("k8s.daemonset.name"); found {
		podMoid.WithDaemonsetName(dsname.Str())
	} else if ssname, found = resource.Attributes().Get("k8s.statefulset.name"); found {
		podMoid.WithStatefulsetName(ssname.Str())
	}*/

	podMoidKey := podMoid.PodMoid()

	// Check if replicaset is found
	if _, rsFound := resource.Attributes().Get("k8s.replicaset.name"); rsFound {
		// If replicaset is found but deployment is not found, fetch deployment with redis data
		if _, deploymentFound := resource.Attributes().Get("k8s.deployment.name"); !deploymentFound {
			op.logger.Debug("Fetching deployment name from redis for pod moid", zap.Any("key", podMoidKey))
			redisData := op.redisClient.GetRedisDataWithDeployment(ctx, podMoidKey)
			op.logger.Debug("Fetched Redis data for pod moid", zap.Any("redisData", redisData))
			if redisData != nil && redisData.ResourceUuid != "" {
				resourceUuid = redisData.ResourceUuid
				if redisData.DeploymentName != "" {
					op.logger.Debug("Fetched deployment name from redis for pod moid", zap.Any("deploymentName", redisData.DeploymentName))
					resource.Attributes().PutStr("k8s.deployment.name", redisData.DeploymentName)
				}
				op.logger.Debug("redis KV ", zap.Any("key", podMoidKey), zap.Any("value", resourceUuid))
			}
		}
	}
	if resourceUuid == "" {
		resourceUuid = op.redisClient.GetUuidValueInString(ctx, podMoidKey)
		op.logger.Debug("redis KV ", zap.Any("key", podMoidKey), zap.Any("value", resourceUuid))
	}
	return
}

func (op *kubernetesprocessor) GetResourceUuidUsingWorkloadMoid(ctx context.Context, resource pcommon.Resource, workloadName pcommon.Value, workloadType string) (resourceUuid string) {
	var namespace pcommon.Value
	var found bool

	if namespace, found = resource.Attributes().Get("k8s.namespace.name"); !found {
		op.logger.Debug("k8s.namespace.name not found in resource attributes hence not able to get resource uuid using workload moid")
		return
	}

	workloadMoid := moid.NewMoid(op.redisConfig.ClusterName).WithNamespaceName(namespace.Str())

	var workloadMoidKey string
	if workloadType == "deployment" {
		workloadMoidKey = workloadMoid.WithDeploymentName(workloadName.Str()).DeploymentMoid()
	} else if workloadType == "replicaset" {
		workloadMoidKey = workloadMoid.WithReplicasetName(workloadName.Str()).ReplicaSetMoid()
	} else if workloadType == "daemonset" {
		workloadMoidKey = workloadMoid.WithDaemonsetName(workloadName.Str()).DaemonSetMoid()
	} else if workloadType == "statefulset" {
		workloadMoidKey = workloadMoid.WithStatefulsetName(workloadName.Str()).StatefulSetMoid()
	}
	if workloadMoidKey == "" {
		return
	}

	resourceUuid = op.redisClient.GetUuidValueInString(ctx, workloadMoidKey)
	op.logger.Debug("redis KV ", zap.Any("key", workloadMoidKey), zap.Any("value", resourceUuid))
	return
}

func (op *kubernetesprocessor) GetResourceUuidUsingResourceNodeMoid(ctx context.Context, resource pcommon.Resource) (resourceUuid string) {
	var nodename pcommon.Value
	var found bool
	if nodename, found = resource.Attributes().Get("k8s.node.name"); !found {
		op.logger.Debug("k8s.node.name not found in resource attributes hence not able to get resource uuid using node moid")
		return
	}

	nodeMoidKey := moid.NewMoid(op.redisConfig.ClusterName).WithNodeName(nodename.Str()).NodeMoid()

	resourceUuid = op.redisClient.GetUuidValueInString(ctx, nodeMoidKey)
	op.logger.Debug("redis KV ", zap.Any("key", nodeMoidKey), zap.Any("value", resourceUuid))
	return
}

func (op *kubernetesprocessor) GetResourceUuidUsingCurrentNodeMoid(ctx context.Context, resource pcommon.Resource) (resourceUuid string) {

	nodeMoidKey := moid.NewMoid(op.redisConfig.ClusterName).WithNodeName(op.redisConfig.NodeName).NodeMoid()

	resourceUuid = op.redisClient.GetUuidValueInString(ctx, nodeMoidKey)
	op.logger.Debug("redis KV ", zap.Any("key", nodeMoidKey), zap.Any("value", resourceUuid))
	return
}

func (op *kubernetesprocessor) GetResourceUuidUsingNamespaceMoid(ctx context.Context, resource pcommon.Resource) (resourceUuid string) {
	var namespace pcommon.Value
	var found bool

	if namespace, found = resource.Attributes().Get("k8s.namespace.name"); !found {
		op.logger.Debug("k8s.namespace.name not found in resource attributes hence not able to get resource uuid using namespace moid")
		return
	}

	namespaceMoidKey := moid.NewMoid(op.redisConfig.ClusterName).WithNamespaceName(namespace.Str()).NamespaceMoid()

	resourceUuid = op.redisClient.GetUuidValueInString(ctx, namespaceMoidKey)
	op.logger.Debug("redis KV ", zap.Any("key", namespaceMoidKey), zap.Any("value", resourceUuid))
	return
}

func (op *kubernetesprocessor) GetResourceUuidUsingClusterMoid(ctx context.Context, resource pcommon.Resource) (resourceUuid string) {

	nodeMoidKey := moid.NewMoid(op.redisConfig.ClusterName).WithClusterUuid(op.redisConfig.ClusterUid).ClusterMoid()

	resourceUuid = op.redisClient.GetUuidValueInString(ctx, nodeMoidKey)
	op.logger.Debug("redis KV ", zap.Any("key", nodeMoidKey), zap.Any("value", resourceUuid))
	return
}
