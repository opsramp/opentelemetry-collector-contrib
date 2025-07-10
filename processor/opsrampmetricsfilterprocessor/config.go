// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampmetricsfilterprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/opsrampmetricsfilterprocessor"

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
)

// Config defines configuration for the Alert Metrics Extractor processor.
type Config struct {
	// AlertConfigMapName is the name of the ConfigMap containing alert definitions
	AlertConfigMapName string `mapstructure:"alert_configmap_name"`

	// AlertConfigMapKey is the key in the ConfigMap containing alert definitions YAML
	AlertConfigMapKey string `mapstructure:"alert_definitions_configmap_key"`

	// Namespace is the Kubernetes namespace where ConfigMaps are located
	Namespace string `mapstructure:"namespace"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	if cfg.AlertConfigMapName == "" {
		return fmt.Errorf("alert_configmap_name is required")
	}

	if cfg.AlertConfigMapKey == "" {
		return fmt.Errorf("alert_definitions_configmap_key is required")
	}

	if cfg.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	return nil
}
