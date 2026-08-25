// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sclusterreceiver

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/confmap/confmaptest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver/internal/metadata"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	disabled := false
	expectedMetricsGroups := MetricsGroupsConfig{
		Pods:                   ResourceGroupConfig{Enabled: &disabled},
		PersistentVolumeClaims: ResourceGroupConfig{Enabled: &disabled},
	}

	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	tests := []struct {
		id          component.ID
		expected    component.Config
		expectedErr error
	}{
		{
			id:       component.NewIDWithName(metadata.Type, ""),
			expected: createDefaultConfig(),
		},
		{
			id: component.NewIDWithName(metadata.Type, "all_settings"),
			expected: &Config{
				Distribution:               distributionKubernetes,
				CollectionInterval:         30 * time.Second,
				NodeConditionTypesToReport: []string{"Ready", "MemoryPressure"},
				AllocatableTypesToReport:   []string{"cpu", "memory"},
				MetadataExporters:          []string{"nop"},
				APIConfig: k8sconfig.APIConfig{
					AuthType: k8sconfig.AuthTypeServiceAccount,
				},
				MetadataCollectionInterval: 30 * time.Minute,
				MetricsBuilderConfig:       metadata.NewDefaultMetricsBuilderConfig(),
			},
		},
		{
			id: component.NewIDWithName(metadata.Type, "partial_settings"),
			expected: &Config{
				Distribution:               distributionOpenShift,
				CollectionInterval:         30 * time.Second,
				NodeConditionTypesToReport: []string{"Ready"},
				APIConfig: k8sconfig.APIConfig{
					AuthType: k8sconfig.AuthTypeServiceAccount,
				},
				MetadataCollectionInterval: 5 * time.Minute,
				MetricsBuilderConfig:       metadata.NewDefaultMetricsBuilderConfig(),
			},
		},
		{
			id: component.NewIDWithName(metadata.Type, "metrics_groups"),
			expected: &Config{
				Distribution:               distributionKubernetes,
				CollectionInterval:         defaultCollectionInterval,
				NodeConditionTypesToReport: []string{"Ready"},
				APIConfig: k8sconfig.APIConfig{
					AuthType: k8sconfig.AuthTypeServiceAccount,
				},
				MetadataCollectionInterval: defaultMetadataCollectionInterval,
				MetricsBuilderConfig:       metadata.NewDefaultMetricsBuilderConfig(),
				MetricsGroups:              expectedMetricsGroups,
			},
		},
		{
			// An empty metrics_groups key must behave exactly like omitting it.
			id:       component.NewIDWithName(metadata.Type, "empty_metrics_groups"),
			expected: createDefaultConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, sub.Unmarshal(cfg))

			assert.NoError(t, confmap.Validate(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestInvalidConfig(t *testing.T) {
	// No APIConfig
	cfg := &Config{
		Distribution:       distributionKubernetes,
		CollectionInterval: 30 * time.Second,
	}
	err := confmap.Validate(cfg)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "invalid authType for kubernetes: ")

	// Wrong distro
	cfg = &Config{
		APIConfig:          k8sconfig.APIConfig{AuthType: k8sconfig.AuthTypeNone},
		Distribution:       "wrong",
		CollectionInterval: 30 * time.Second,
	}
	expectedErr := "\"wrong\" is not a supported distribution. Must be one of: \"openshift\", \"kubernetes\""
	err = confmap.Validate(cfg)
	assert.Error(t, err)
	assert.ErrorContains(t, err, expectedErr)

	// Every metrics group turned off
	disabled := false
	cfg = &Config{
		APIConfig:          k8sconfig.APIConfig{AuthType: k8sconfig.AuthTypeNone},
		Distribution:       distributionKubernetes,
		CollectionInterval: 30 * time.Second,
		MetricsGroups:      MetricsGroupsConfig{EnabledByDefault: &disabled},
	}
	err = confmap.Validate(cfg)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "all metrics_groups are disabled")
}

func TestMetricsGroupsEnabledForKind(t *testing.T) {
	allKinds := []string{
		"Pod", "Node", "Namespace", "Deployment", "ReplicaSet", "ReplicationController", "DaemonSet",
		"StatefulSet", "Job", "CronJob", "HorizontalPodAutoscaler", "ResourceQuota", "Service",
		"EndpointSlice", "PersistentVolume", "PersistentVolumeClaim", "ClusterResourceQuota",
		// Unknown kinds follow enabled_by_default so newly added resources are not silently dropped.
		"SomeFutureKind",
	}

	var groups MetricsGroupsConfig
	for _, kind := range allKinds {
		assert.True(t, groups.enabledForKind(kind), "kind %s should be enabled by default", kind)
	}

	enabled, disabled := true, false
	groups.Pods = ResourceGroupConfig{Enabled: &disabled}
	groups.Services = ResourceGroupConfig{Enabled: &disabled}
	assert.False(t, groups.enabledForKind("Pod"))
	assert.False(t, groups.enabledForKind("Service"))
	// EndpointSlice only exists to back k8s.service.endpoint.count.
	assert.False(t, groups.enabledForKind("EndpointSlice"))
	assert.True(t, groups.enabledForKind("Node"))

	// enabled_by_default: false opts out of everything that is not set explicitly.
	optOut := MetricsGroupsConfig{
		EnabledByDefault: &disabled,
		Nodes:            ResourceGroupConfig{Enabled: &enabled},
	}
	for _, kind := range allKinds {
		if kind == "Node" {
			continue
		}
		assert.False(t, optOut.enabledForKind(kind), "kind %s should be disabled", kind)
	}
	assert.True(t, optOut.enabledForKind("Node"))
}
