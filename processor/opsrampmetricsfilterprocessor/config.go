// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampmetricsfilterprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/opsrampmetricsfilterprocessor"

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.opentelemetry.io/collector/component"
)

// Config defines configuration for the Alert Metrics Extractor processor.
type Config struct {
	// AlertConfigMapName is the name of the ConfigMap containing alert definitions
	// Optional: defaults to "opsramp-alert-user-config"
	AlertConfigMapName string `mapstructure:"alert_definitions_configmap_name"`

	// AlertConfigMapKey is the key in the ConfigMap containing alert definitions YAML
	// Optional: defaults to "alert-definitions.yaml"
	AlertConfigMapKey string `mapstructure:"alert_definitions_key"`

	// Namespace is the Kubernetes namespace where ConfigMaps are located
	// This is automatically populated from the NAMESPACE environment variable
	// Users should not set this field directly
	Namespace string `mapstructure:"namespace"`

	// AlertDefinitionsFilePath is the file path containing alert definitions
	// Optional: if provided, this takes precedence over ConfigMap
	AlertDefinitionsFilePath string `mapstructure:"alert_definitions_file_path"`

	// WatchFileChanges enables file watching for dynamic updates
	// Optional: defaults to true when using file path
	WatchFileChanges bool `mapstructure:"watch_file_changes"`

	// FileWatchInterval is the interval for checking file changes
	// Optional: defaults to 30 seconds
	FileWatchInterval string `mapstructure:"file_watch_interval"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	// If file path is provided, it takes precedence. Clear any ConfigMap fields
	// that may have been populated by createDefaultConfig() to avoid false conflicts.
	if cfg.AlertDefinitionsFilePath != "" {
		cfg.AlertConfigMapName = ""
		cfg.AlertConfigMapKey = ""
		cfg.Namespace = ""

		// Validate file path format
		if !filepath.IsAbs(cfg.AlertDefinitionsFilePath) {
			return fmt.Errorf("alert_definitions_file_path must be an absolute path: %s", cfg.AlertDefinitionsFilePath)
		}

		// Validate file extension
		ext := filepath.Ext(cfg.AlertDefinitionsFilePath)
		if ext != ".yaml" && ext != ".yml" {
			return fmt.Errorf("alert definitions file must have .yaml or .yml extension: %s", cfg.AlertDefinitionsFilePath)
		}

		// Check if file exists and is readable
		if _, err := os.Stat(cfg.AlertDefinitionsFilePath); os.IsNotExist(err) {
			return fmt.Errorf("alert definitions file does not exist: %s", cfg.AlertDefinitionsFilePath)
		}

		// Validate and set default watch interval
		if cfg.FileWatchInterval == "" {
			cfg.FileWatchInterval = "30s"
		} else if _, err := time.ParseDuration(cfg.FileWatchInterval); err != nil {
			return fmt.Errorf("invalid file_watch_interval format: %s (must be valid duration like '30s', '1m')", cfg.FileWatchInterval)
		}
	} else {
		// ConfigMap mode — apply defaults for name/key if unset,
		// then resolve namespace from env.
		if cfg.AlertConfigMapName == "" {
			cfg.AlertConfigMapName = "opsramp-alert-user-config"
		}
		if cfg.AlertConfigMapKey == "" {
			cfg.AlertConfigMapKey = "alert-definitions.yaml"
		}
		if cfg.Namespace == "" {
			cfg.Namespace = os.Getenv("NAMESPACE")
			if cfg.Namespace == "" {
				cfg.Namespace = "opsramp-agent"
			}
		}
	}

	return nil
}
