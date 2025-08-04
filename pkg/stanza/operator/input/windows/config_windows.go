// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package windows // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/input/windows"

import (
	"errors"
	"fmt"
	"go.opentelemetry.io/collector/component"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
)

const (
	DefaultPersisterPath = "c:\\ProgramData\\OpenTelemetry\\Persister"
	DefaultFilePersister = "persister_data.json"
	DefaultLevelDBPath   = "leveldb_data"
)

func init() {
	operator.Register(operatorType, func() operator.Builder { return NewConfig() })
}

// Build will build a windows event log operator.
func (c *Config) Build(set component.TelemetrySettings) (operator.Operator, error) {
	inputOperator, err := c.InputConfig.Build(set)
	if err != nil {
		return nil, err
	}

	if c.Channel == "" {
		return nil, errors.New("missing required `channel` field")
	}

	if c.MaxReads < 1 {
		return nil, errors.New("the `max_reads` field must be greater than zero")
	}

	if c.StartAt != "end" && c.StartAt != "beginning" {
		return nil, errors.New("the `start_at` field must be set to `beginning` or `end`")
	}

	// Create the persister using the config
	persister, err := NewPersister(adapter.BaseConfig{
		PersisterType: c.PersisterType,
		PersisterPath: c.PersisterPath,
	}, c.Channel)

	fmt.Printf("************ error persister: %v", persister, "error: %v", err)
	if err != nil {
		return nil, fmt.Errorf("failed to create persister: %w", err)
	}

	return &Input{
		InputOperator:       inputOperator,
		buffer:              NewBuffer(),
		channel:             c.Channel,
		maxReads:            c.MaxReads,
		startAt:             c.StartAt,
		pollInterval:        c.PollInterval,
		raw:                 c.Raw,
		isReqAdditionalAttr: c.IsReqAdditionalAttr,
		excludeProviders:    c.ExcludeProviders,
		persister:           persister,
	}, nil
}
