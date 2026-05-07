// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8seventsreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8seventsreceiver"

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
	"k8s.io/client-go/dynamic"
	k8s "k8s.io/client-go/kubernetes"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/k8sconfig"
)

type EventType string

const (
	EventTypeNormal  EventType = "Normal"
	EventTypeWarning EventType = "Warning"
)

// Config defines configuration for kubernetes events receiver.
type Config struct {
	k8sconfig.APIConfig `mapstructure:",squash"`

	// List of ‘namespaces’ to collect events from.
	Namespaces []string `mapstructure:"namespaces"`
	// List of namespaces to exclude from event collection.
	// When Namespaces is empty (watch all), events from these namespaces are dropped.
	// When Namespaces is non-empty, these are removed from the watch list (set difference).
	// Empty or nil means no exclusions.
	ExcludeNamespaces []string `mapstructure:"exclude_namespaces"`
	// List of ‘eventtypes’ to filter.
	EventTypes []EventType `mapstructure:"event_types,omitempty"`

	// Include only the specified involved objects. ObjectKind to List of Reasons.
	IncludeInvolvedObject map[string]InvolvedObjectProperties `mapstructure:"include_involved_objects,omitempty"`

	K8sLeaderElector *component.ID `mapstructure:"k8s_leader_elector"`

	// For mocking
	makeClient        func(apiConf k8sconfig.APIConfig) (k8s.Interface, error)
	makeDynamicClient func(apiConf k8sconfig.APIConfig) (dynamic.Interface, error)
}

type InvolvedObjectProperties struct {
	// Include only the specified reasons. If its empty, list events of all reasons.
	IncludeReasons []ReasonProperties `mapstructure:"include_reasons,omitempty"`
	// Exclude events with the specified reasons. Takes precedence over include_reasons.
	// If empty or absent, no events are excluded.
	// Note: Only the Name field of each entry is used for matching; Attributes are ignored.
	ExcludeReasons []ReasonProperties `mapstructure:"exclude_reasons,omitempty"`

	//Can be enhanced to take in object names with reg ex etc.
}

type ReasonProperties struct {
	Name       string     `mapstructure:"name"`
	Attributes []KeyValue `mapstructure:"attributes,omitempty"`
}

type KeyValue struct {
	// This is a required field.
	Key string `mapstructure:"key"`

	// This is a required field.
	Value any `mapstructure:"value"`
}

func (cfg *Config) Validate() error {
	for _, eventType := range cfg.EventTypes {
		switch eventType {
		case EventTypeNormal, EventTypeWarning:
		default:
			return fmt.Errorf("invalid event_type %s, must be one of %s or %s", eventType, EventTypeNormal, EventTypeWarning)
		}
	}
	return cfg.APIConfig.Validate()
}

func (cfg *Config) getK8sClient() (k8s.Interface, error) {
	if cfg.makeClient == nil {
		cfg.makeClient = k8sconfig.MakeClient
	}
	return cfg.makeClient(cfg.APIConfig)
}

func (cfg *Config) getDynamicClient() (dynamic.Interface, error) {
	if cfg.makeDynamicClient == nil {
		cfg.makeDynamicClient = k8sconfig.MakeDynamicClient
	}
	return cfg.makeDynamicClient(cfg.APIConfig)
}
