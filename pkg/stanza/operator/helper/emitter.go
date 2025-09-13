// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package helper // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"

import (
	"context"
	"go.uber.org/zap"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/entry"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
)

// LogEmitter is a stanza operator that emits log entries to a channel
type LogEmitter struct {
	OutputOperator
	logChan       chan []*entry.Entry
	closeChan     chan struct{}
	stopOnce      sync.Once
	batchMux      sync.Mutex
	batch         []*entry.Entry
	wg            sync.WaitGroup
	maxBatchSize  uint
	flushInterval time.Duration
}

var (
	defaultFlushInterval      = 100 * time.Millisecond
	defaultMaxBatchSize  uint = 100
)

type EmitterOption interface {
	apply(*LogEmitter)
}

func WithMaxBatchSize(maxBatchSize uint) EmitterOption {
	return maxBatchSizeOption{maxBatchSize}
}

type maxBatchSizeOption struct {
	maxBatchSize uint
}

func (o maxBatchSizeOption) apply(e *LogEmitter) {
	e.maxBatchSize = o.maxBatchSize
}

func WithFlushInterval(flushInterval time.Duration) EmitterOption {
	return flushIntervalOption{flushInterval}
}

type flushIntervalOption struct {
	flushInterval time.Duration
}

func (o flushIntervalOption) apply(e *LogEmitter) {
	e.flushInterval = o.flushInterval
}

// NewLogEmitter creates a new receiver output
func NewLogEmitter(set component.TelemetrySettings, opts ...EmitterOption) *LogEmitter {
	op, _ := NewOutputConfig("log_emitter", "log_emitter").Build(set)
	e := &LogEmitter{
		OutputOperator: op,
		logChan:        make(chan []*entry.Entry),
		closeChan:      make(chan struct{}),
		maxBatchSize:   defaultMaxBatchSize,
		batch:          make([]*entry.Entry, 0, defaultMaxBatchSize),
		flushInterval:  defaultFlushInterval,
	}
	for _, opt := range opts {
		opt.apply(e)
	}
	return e
}

// Start starts the goroutine(s) required for this operator
func (e *LogEmitter) Start(_ operator.Persister) error {
	e.wg.Add(1)
	go e.flusher()
	return nil
}

// Stop will close the log channel and stop running goroutines
func (e *LogEmitter) Stop() error {
	e.Logger().Debug("Stopping log emitter")
	e.stopOnce.Do(func() {
		e.Logger().Debug("Closing log emitter channels")
		close(e.closeChan)
		e.Logger().Debug("Waiting for log emitter goroutines to finish")
		e.wg.Wait()
		e.Logger().Debug("Closing log channel")

		close(e.logChan)
		e.Logger().Debug("Log emitter stopped")
	})

	e.Logger().Debug("Stopped log emitter")
	return nil
}

// OutChannel returns the channel on which entries will be sent to.
func (e *LogEmitter) OutChannel() <-chan []*entry.Entry {
	e.set.Logger.Debug("OutChannel called", zap.Int("channel_length", len(e.logChan)), zap.Int("channel_capacity", cap(e.logChan)))
	return e.logChan
}

// OutChannelForWrite returns the channel on which entries can be sent to.
func (e *LogEmitter) OutChannelForWrite() chan []*entry.Entry {
	e.Logger().Debug("OutChannelForWrite called")
	return e.logChan
}

// Process will emit an entry to the output channel
func (e *LogEmitter) Process(ctx context.Context, ent *entry.Entry) error {
	e.Logger().Debug("Emitting entry in emitter")
	if oldBatch := e.appendEntry(ent); len(oldBatch) > 0 {
		e.flush(ctx, oldBatch)
	}

	return nil
}

// appendEntry appends the entry to the current batch. If maxBatchSize is reached, a new batch will be made, and the old batch
// (which should be flushed) will be returned
func (e *LogEmitter) appendEntry(ent *entry.Entry) []*entry.Entry {
	e.batchMux.Lock()
	defer e.batchMux.Unlock()

	e.batch = append(e.batch, ent)
	if uint(len(e.batch)) >= e.maxBatchSize {
		var oldBatch []*entry.Entry
		oldBatch, e.batch = e.batch, make([]*entry.Entry, 0, e.maxBatchSize)
		e.Logger().Debug("Appended entry to batch", zap.Int("new_batch_size", len(oldBatch)))
		return oldBatch
	}

	return nil
}

// flusher flushes the current batch every flush interval. Intended to be run as a goroutine
func (e *LogEmitter) flusher() {
	defer e.wg.Done()

	ticker := time.NewTicker(e.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if oldBatch := e.makeNewBatch(); len(oldBatch) > 0 {
				e.flush(context.Background(), oldBatch)
			}
		case <-e.closeChan:
			// flush currently batched entries
			if oldBatch := e.makeNewBatch(); len(oldBatch) > 0 {
				e.flush(context.Background(), oldBatch)
			}
			return
		}
	}
}

// flush flushes the provided batch to the log channel.
func (e *LogEmitter) flush(ctx context.Context, batch []*entry.Entry) {
	e.Logger().Debug("Flushing batch", zap.Int("batch_size", len(batch)))
	select {
	case e.logChan <- batch:
		e.Logger().Debug("Successfully sent batch to log channel",
			zap.Int("batch_size", len(batch)))
	case <-ctx.Done():
		e.Logger().Debug("Context done before batch could be sent",
			zap.Int("batch_size", len(batch)),
			zap.Error(ctx.Err()))
	default:
		e.Logger().Warn("Blocking on sending batch to log channel",
			zap.Int("batch_size", len(batch)))
	}
}

// makeNewBatch replaces the current batch on the log emitter with a new batch, returning the old one
func (e *LogEmitter) makeNewBatch() []*entry.Entry {
	e.batchMux.Lock()
	defer e.batchMux.Unlock()

	if len(e.batch) == 0 {
		return nil
	}

	var oldBatch []*entry.Entry
	oldBatch, e.batch = e.batch, make([]*entry.Entry, 0, e.maxBatchSize)
	return oldBatch
}
