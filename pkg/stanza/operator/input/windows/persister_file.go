// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package windows

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// FilePersister is a file-based implementation of the Persister interface.
type FilePersister struct {
	filePath string
	data     map[string][]byte
	mu       sync.Mutex
}

// NewFilePersister creates a new file-based persister.
func NewFilePersister(path string) (Persister, error) {
	// if path is empty, use a default path
	if path == "" {
		path = DefaultPersisterPath + "/" + DefaultFilePersister
	}

	// Ensure the parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	p := &FilePersister{
		filePath: path,
		data:     make(map[string][]byte),
	}

	// Load existing data from the file if it exists.
	file, err := os.Open(path)
	if err == nil {
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&p.data); err != nil {
			return nil, fmt.Errorf("failed to decode file data: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return p, nil
}

func (p *FilePersister) Get(ctx context.Context, key string) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.data[key], nil
}

func (p *FilePersister) Set(ctx context.Context, key string, value []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data[key] = value
	return p.saveToFile()
}

func (p *FilePersister) Delete(ctx context.Context, key string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.data, key)
	return p.saveToFile()
}

func (p *FilePersister) saveToFile() error {
	file, err := os.Create(p.filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	return json.NewEncoder(file).Encode(p.data)
}
