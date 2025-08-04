// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package windows

import (
	"context"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"os"
	"path/filepath"
)

// LevelDBPersister is a LevelDB-based implementation of the Persister interface.
type LevelDBPersister struct {
	db *leveldb.DB
}

// NewLevelDBPersister creates a new LevelDB-based persister.
func NewLevelDBPersister(path string) (Persister, error) {
	// if path is empty, use a default path
	if path == "" {
		path = DefaultPersisterPath + "/" + DefaultLevelDBPath
	}

	// Ensure the parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("NewLevelDBPersister:failed to create directory: %w", err)
	}

	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open LevelDB at %s: %w", path, err)
	}
	return &LevelDBPersister{db: db}, nil
}

func (p *LevelDBPersister) Get(ctx context.Context, key string) ([]byte, error) {
	value, err := p.db.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		return nil, nil
	}
	return value, err
}

func (p *LevelDBPersister) Set(ctx context.Context, key string, value []byte) error {
	return p.db.Put([]byte(key), value, nil)
}

func (p *LevelDBPersister) Delete(ctx context.Context, key string) error {
	return p.db.Delete([]byte(key), nil)
}
