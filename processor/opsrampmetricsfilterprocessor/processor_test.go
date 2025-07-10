// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampmetricsfilterprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/processor/processortest"
)

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestCreateProcessor(t *testing.T) {
	cfg := &Config{
		AlertConfigMapName: "test-config",
		AlertConfigMapKey:  "alert-definitions.yaml",
		Namespace:          "test-namespace",
	}

	// This test will fail in CI/testing environment without Kubernetes cluster
	// But validates the configuration structure
	_, err := createMetricsProcessor(
		context.Background(),
		processortest.NewNopSettings(),
		cfg,
		consumertest.NewNop(),
	)

	// Expect error due to no Kubernetes cluster in test environment
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create")
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				AlertConfigMapName: "test-config",
				AlertConfigMapKey:  "alert-definitions.yaml",
				Namespace:          "test-namespace",
			},
			wantErr: false,
		},
		{
			name: "missing configmap name",
			config: &Config{
				AlertConfigMapKey: "alert-definitions.yaml",
				Namespace:         "test-namespace",
			},
			wantErr: true,
		},
		{
			name: "missing configmap key",
			config: &Config{
				AlertConfigMapName: "test-config",
				Namespace:          "test-namespace",
			},
			wantErr: true,
		},
		{
			name: "missing namespace",
			config: &Config{
				AlertConfigMapName: "test-config",
				AlertConfigMapKey:  "alert-definitions.yaml",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestExtractMetricsFromExpression(t *testing.T) {
	// This is a unit test that doesn't require Kubernetes
	processor := &filterProcessor{
		logger: processortest.NewNopSettings().Logger,
	}

	tests := []struct {
		name     string
		expr     string
		expected []string
	}{
		{
			name:     "simple metric",
			expr:     "up",
			expected: []string{"up"},
		},
		{
			name:     "metric with labels",
			expr:     "http_requests_total{job=\"api-server\"}",
			expected: []string{"http_requests_total"},
		},
		{
			name:     "binary operation",
			expr:     "rate(http_requests_total[5m]) > 0.5",
			expected: []string{"http_requests_total"},
		},
		{
			name:     "multiple metrics",
			expr:     "cpu_usage_percent + memory_usage_percent",
			expected: []string{"cpu_usage_percent", "memory_usage_percent"},
		},
		{
			name:     "function with metric",
			expr:     "increase(errors_total[1h])",
			expected: []string{"errors_total"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.extractMetricsFromExpression(tt.expr)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}
