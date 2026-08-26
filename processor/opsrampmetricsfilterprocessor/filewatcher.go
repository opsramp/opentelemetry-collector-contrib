// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package opsrampmetricsfilterprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/opsrampmetricsfilterprocessor"

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
)

// FileWatcher handles file system watching for alert definitions
type FileWatcher struct {
	filePath     string
	logger       *zap.Logger
	callback     func() error
	watcher      *fsnotify.Watcher
	ctx          context.Context
	useFsNotify  bool
	pollInterval time.Duration
}

// createFileWatcher creates a new file watcher with fsnotify support and polling fallback
func (l *definitionsLoader) createFileWatcher() (*FileWatcher, error) {
	fw := &FileWatcher{
		filePath:     l.key.filePath,
		logger:       l.logger,
		callback:     l.reload,
		ctx:          l.ctx,
		pollInterval: l.watchInterval,
		useFsNotify:  true,
	}

	var err error
	fw.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		fw.logger.Warn("Failed to create fsnotify watcher, falling back to polling",
			zap.Error(err))
		fw.useFsNotify = false
		return fw, nil
	}

	// Watch the directory since the file might be replaced rather than written
	// in place, which is what a ConfigMap volume update does.
	dir := filepath.Dir(fw.filePath)
	if err := fw.watcher.Add(dir); err != nil {
		fw.logger.Warn("Failed to add directory to watcher, falling back to polling",
			zap.String("directory", dir),
			zap.Error(err))
		fw.watcher.Close()
		fw.watcher = nil
		fw.useFsNotify = false
	}

	return fw, nil
}

// Start begins watching for file changes
func (fw *FileWatcher) Start() {
	if fw.useFsNotify {
		fw.startFsNotifyWatch()
	} else {
		fw.startPollingWatch()
	}
}

// Close stops the file watcher and cleans up resources
func (fw *FileWatcher) Close() error {
	if fw.watcher != nil {
		return fw.watcher.Close()
	}
	return nil
}

// startFsNotifyWatch uses fsnotify for real-time file change detection.
// Uses absolute path comparison to avoid false triggers from same-named files,
// and debounces rapid events (e.g. vim/VS Code multi-event saves).
func (fw *FileWatcher) startFsNotifyWatch() {
	fw.logger.Info("Starting fsnotify file watcher", zap.String("file_path", fw.filePath))

	// Pre-resolve the target path for accurate comparison
	absFilePath, err := filepath.Abs(fw.filePath)
	if err != nil {
		absFilePath = fw.filePath
	}

	const debounceDelay = 200 * time.Millisecond
	var (
		debounceTimer *time.Timer
		debounceMu    sync.Mutex
	)

	for {
		select {
		case <-fw.ctx.Done():
			debounceMu.Lock()
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceMu.Unlock()
			return
		case event, ok := <-fw.watcher.Events:
			if !ok {
				fw.logger.Error("File watcher events channel closed")
				return
			}

			// Use absolute path comparison to avoid false triggers
			absEventPath, _ := filepath.Abs(event.Name)
			if absEventPath != absFilePath {
				continue
			}

			fw.logger.Debug("File event received",
				zap.String("file", event.Name),
				zap.String("operation", event.Op.String()))

			// Debounce write/create/rename events to avoid redundant rapid reloads
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Rename) {
				debounceMu.Lock()
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				eventOp := event.Op.String()
				debounceTimer = time.AfterFunc(debounceDelay, func() {
					select {
					case <-fw.ctx.Done():
						return
					default:
					}
					fw.logger.Info("File modified, reloading",
						zap.String("file_path", fw.filePath),
						zap.String("event", eventOp))
					if err := fw.callback(); err != nil {
						fw.logger.Error("Failed to reload file after change", zap.Error(err))
					} else {
						fw.logger.Info("Successfully reloaded file after change")
					}
				})
				debounceMu.Unlock()
			}
		case err, ok := <-fw.watcher.Errors:
			if !ok {
				fw.logger.Error("File watcher errors channel closed")
				return
			}
			fw.logger.Error("File watcher error", zap.Error(err))
		}
	}
}

// startPollingWatch falls back to polling for file changes
func (fw *FileWatcher) startPollingWatch() {
	fw.logger.Info("Starting polling file watcher",
		zap.String("file_path", fw.filePath),
		zap.Duration("interval", fw.pollInterval))

	var lastModTime time.Time

	// Get initial modification time
	if fileInfo, err := os.Stat(fw.filePath); err == nil {
		lastModTime = fileInfo.ModTime()
	}

	ticker := time.NewTicker(fw.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-fw.ctx.Done():
			return
		case <-ticker.C:
			fileInfo, err := os.Stat(fw.filePath)
			if err != nil {
				fw.logger.Error("Failed to get file info during polling",
					zap.String("file_path", fw.filePath),
					zap.Error(err))
				continue
			}

			if fileInfo.ModTime().After(lastModTime) {
				fw.logger.Info("File changed (polling), reloading",
					zap.String("file_path", fw.filePath),
					zap.Time("old_mod_time", lastModTime),
					zap.Time("new_mod_time", fileInfo.ModTime()))

				lastModTime = fileInfo.ModTime()
				if err := fw.callback(); err != nil {
					fw.logger.Error("Failed to reload file after polling change", zap.Error(err))
				} else {
					fw.logger.Info("Successfully reloaded file after polling change")
				}
			}
		}
	}
}
