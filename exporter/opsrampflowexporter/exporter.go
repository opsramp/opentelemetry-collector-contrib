// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampflowexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opsrampflowexporter"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
)

type flowExporter struct {
	endpoint string
	token    string
}

func newFlowExporter(cfg *Config, set exporter.Settings) (*flowExporter, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &flowExporter{
		endpoint: cfg.Endpoint,
		token:    cfg.AuthToken,
	}, nil
}

func (l *flowExporter) Start(ctx context.Context, host component.Host) error {

	return nil
}

func (l *flowExporter) pushLogData(ctx context.Context, ld plog.Logs) error {
	return nil
}
