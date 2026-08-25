// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sclusterreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver"

import (
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver/internal/metadata"
)

// Config defines configuration for kubernetes cluster receiver.
type Config struct {
	APIConfig k8sconfig.APIConfig `mapstructure:",squash"`

	// Collection interval for metrics.
	CollectionInterval time.Duration `mapstructure:"collection_interval"`

	// Node condition types to report. See all condition types, see
	// here: https://kubernetes.io/docs/concepts/architecture/nodes/#condition.
	NodeConditionTypesToReport []string `mapstructure:"node_conditions_to_report"`
	// Allocate resource types to report. See all resource types, see
	// here: https://kubernetes.io/docs/tasks/administer-cluster/reserve-compute-resources/#node-allocatable
	AllocatableTypesToReport []string `mapstructure:"allocatable_types_to_report"`
	// List of exporters to which metadata from this receiver should be forwarded to.
	MetadataExporters []string `mapstructure:"metadata_exporters"`

	// Whether OpenShift support should be enabled or not.
	Distribution string `mapstructure:"distribution"`

	// Collection interval for metadata.
	// Metadata of the particular entity in the cluster is collected when the entity changes.
	// In addition metadata of all entities is collected periodically even if no changes happen.
	// Setting the duration to 0 will disable periodic collection (however will not impact
	// metadata collection on changes).
	MetadataCollectionInterval time.Duration `mapstructure:"metadata_collection_interval"`

	// MetricsBuilderConfig allows customizing scraped metrics/attributes representation.
	MetricsBuilderConfig metadata.MetricsBuilderConfig `mapstructure:",squash"`

	// MetricsGroups enables/disables collection per Kubernetes resource type. Disabling a group
	// also stops the informer that watches that resource, so no API watch or cache is maintained
	// for it. This takes precedence over the individual toggles in `metrics`.
	MetricsGroups MetricsGroupsConfig `mapstructure:"metrics_groups"`

	// Deprecated: This field is no longer supported, use cfg.Namespaces instead.
	Namespace string `mapstructure:"namespace"`

	// Namespaces to fetch resources from. If this is set, certain cluster-wide resources such as Nodes or Namespaces
	// will not be able to be observed. Setting this option is recommended in environments where due to security restrictions
	// the collector cannot be granted cluster-wide permissions.
	Namespaces []string `mapstructure:"namespaces"`

	// K8sLeaderElector defines the reference to the k8s leader elector extension
	// use this when k8s cluster receiver needs to be deployed in HA mode
	K8sLeaderElector *component.ID `mapstructure:"k8s_leader_elector"`
}

// ResourceGroupConfig toggles collection of a single Kubernetes resource type.
type ResourceGroupConfig struct {
	// Enabled falls back to MetricsGroupsConfig.EnabledByDefault when it is not set.
	Enabled *bool `mapstructure:"enabled"`
}

// MetricsGroupsConfig groups the collected metrics by the Kubernetes resource type they are
// derived from. Every group maps 1:1 to an informer, so disabling a group removes the
// corresponding LIST/WATCH against the API server and its client-side cache. All groups are
// enabled by default.
type MetricsGroupsConfig struct {
	// EnabledByDefault (default = true) applies to every group that is not set explicitly. Set it
	// to false to opt out of all resource types and then enable only the ones needed.
	EnabledByDefault *bool `mapstructure:"enabled_by_default"`

	// Pods also controls all k8s.container.* metrics, which are derived from pod objects.
	Pods                    ResourceGroupConfig `mapstructure:"pods"`
	Nodes                   ResourceGroupConfig `mapstructure:"nodes"`
	Namespaces              ResourceGroupConfig `mapstructure:"namespaces"`
	Deployments             ResourceGroupConfig `mapstructure:"deployments"`
	ReplicaSets             ResourceGroupConfig `mapstructure:"replicasets"`
	ReplicationControllers  ResourceGroupConfig `mapstructure:"replicationcontrollers"`
	DaemonSets              ResourceGroupConfig `mapstructure:"daemonsets"`
	StatefulSets            ResourceGroupConfig `mapstructure:"statefulsets"`
	Jobs                    ResourceGroupConfig `mapstructure:"jobs"`
	CronJobs                ResourceGroupConfig `mapstructure:"cronjobs"`
	HorizontalPodAutoscaler ResourceGroupConfig `mapstructure:"horizontalpodautoscalers"`
	ResourceQuotas          ResourceGroupConfig `mapstructure:"resourcequotas"`
	// Services also controls the EndpointSlice informer backing k8s.service.endpoint.count.
	Services               ResourceGroupConfig `mapstructure:"services"`
	PersistentVolumes      ResourceGroupConfig `mapstructure:"persistentvolumes"`
	PersistentVolumeClaims ResourceGroupConfig `mapstructure:"persistentvolumeclaims"`
	// ClusterResourceQuotas is only applicable to the openshift distribution.
	ClusterResourceQuotas ResourceGroupConfig `mapstructure:"clusterresourcequotas"`
}

// groupForKind returns the group owning the given Kubernetes kind. Unknown kinds map to an unset
// group so that they follow EnabledByDefault.
func (c MetricsGroupsConfig) groupForKind(kind string) ResourceGroupConfig {
	switch kind {
	case "Pod":
		return c.Pods
	case "Node":
		return c.Nodes
	case "Namespace":
		return c.Namespaces
	case "Deployment":
		return c.Deployments
	case "ReplicaSet":
		return c.ReplicaSets
	case "ReplicationController":
		return c.ReplicationControllers
	case "DaemonSet":
		return c.DaemonSets
	case "StatefulSet":
		return c.StatefulSets
	case "Job":
		return c.Jobs
	case "CronJob":
		return c.CronJobs
	case "HorizontalPodAutoscaler":
		return c.HorizontalPodAutoscaler
	case "ResourceQuota":
		return c.ResourceQuotas
	case "Service", "EndpointSlice":
		return c.Services
	case "PersistentVolume":
		return c.PersistentVolumes
	case "PersistentVolumeClaim":
		return c.PersistentVolumeClaims
	case "ClusterResourceQuota":
		return c.ClusterResourceQuotas
	default:
		return ResourceGroupConfig{}
	}
}

// metricsGroupKinds lists the Kubernetes kinds that own a metrics group.
var metricsGroupKinds = []string{
	"Pod", "Node", "Namespace", "Deployment", "ReplicaSet", "ReplicationController", "DaemonSet",
	"StatefulSet", "Job", "CronJob", "HorizontalPodAutoscaler", "ResourceQuota", "Service",
	"PersistentVolume", "PersistentVolumeClaim", "ClusterResourceQuota",
}

// enabledForKind reports whether the group owning the given Kubernetes kind is enabled.
func (c MetricsGroupsConfig) enabledForKind(kind string) bool {
	if enabled := c.groupForKind(kind).Enabled; enabled != nil {
		return *enabled
	}
	return c.EnabledByDefault == nil || *c.EnabledByDefault
}

func (c MetricsGroupsConfig) anyEnabled() bool {
	for _, kind := range metricsGroupKinds {
		if c.enabledForKind(kind) {
			return true
		}
	}
	return false
}

func (cfg *Config) Validate() error {
	switch cfg.Distribution {
	case distributionOpenShift:
	case distributionKubernetes:
	default:
		return fmt.Errorf("\"%s\" is not a supported distribution. Must be one of: \"openshift\", \"kubernetes\"", cfg.Distribution)
	}

	if !cfg.MetricsGroups.anyEnabled() {
		return errors.New("all metrics_groups are disabled, the receiver would not collect anything")
	}

	return nil
}
