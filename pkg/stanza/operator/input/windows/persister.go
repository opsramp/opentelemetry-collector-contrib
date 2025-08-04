// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package windows // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/input/windows"

import (
	"context"
	"fmt"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
)

// Persister is an interface used to persist data
type Persister interface {
	Get(context.Context, string) ([]byte, error)
	Set(context.Context, string, []byte) error
	Delete(context.Context, string) error
}

// scopedPersister wraps a Persister and adds a scope prefix to keys
type scopedPersister struct {
	Persister
	scope string
}

// NewPersister creates a new Persister based on the config
func NewPersister(cfg adapter.BaseConfig, scope string) (Persister, error) {
	var p Persister
	var err error

	fmt.Printf("NewPersister: PersisterType: %v, PersisterPath:%v", cfg.PersisterType, cfg.PersisterPath)

	switch cfg.PersisterType {
	case "leveldb", "level_db":
		p, err = NewLevelDBPersister(cfg.PersisterPath)
	default:
		p, err = NewFilePersister(cfg.PersisterPath)
	}
	if err != nil {
		_ = fmt.Errorf("NewPersister: Failed to create Persister: %v", err)
		return nil, err
	}
	return &scopedPersister{
		Persister: p,
		scope:     scope,
	}, nil
}

func (p *scopedPersister) Get(ctx context.Context, key string) ([]byte, error) {
	return p.Persister.Get(ctx, fmt.Sprintf("%s.%s", p.scope, key))
}
func (p *scopedPersister) Set(ctx context.Context, key string, value []byte) error {
	return p.Persister.Set(ctx, fmt.Sprintf("%s.%s", p.scope, key), value)
}
func (p *scopedPersister) Delete(ctx context.Context, key string) error {
	return p.Persister.Delete(ctx, fmt.Sprintf("%s.%s", p.scope, key))
}
