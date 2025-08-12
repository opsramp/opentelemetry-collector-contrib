// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package recombine // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/transformer/recombine"

import (
	"bytes"
	"context"
	"errors"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
	"time"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/entry"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
)

const DefaultSourceIdentifier = "DefaultSourceIdentifier"

// Transformer is an operator that combines a field from consecutive log entries into a single
type Transformer struct {
	helper.TransformerOperator
	matchFirstLine        bool
	prog                  *vm.Program
	maxBatchSize          int
	maxUnmatchedBatchSize int
	maxSources            int
	overwriteWithNewest   bool
	combineField          entry.Field
	combineWith           string
	ticker                *time.Ticker
	forceFlushTimeout     time.Duration
	chClose               chan struct{}
	sourceIdentifier      entry.Field

	sync.Mutex
	batchPool  sync.Pool
	batchMap   map[string]*sourceBatch
	maxLogSize int64
}

// sourceBatch contains the status info of a batch
type sourceBatch struct {
	baseEntry              *entry.Entry
	numEntries             int
	recombined             *bytes.Buffer
	firstEntryObservedTime time.Time
	matchDetected          bool
}

func (t *Transformer) Start(_ operator.Persister) error {
	go t.flushLoop()
	return nil
}

func (t *Transformer) flushLoop() {
	for {
		select {
		case <-t.ticker.C:
			t.Lock()
			timeNow := time.Now()
			for source, batch := range t.batchMap {
				customLogger := initLogger()
				customLogger.Debug("Checking batch for force flush", zap.String("source", source), zap.Time("firstEntryObservedTime", batch.firstEntryObservedTime))
				timeSinceFirstEntry := timeNow.Sub(batch.firstEntryObservedTime)
				customLogger.Debug("Time since first entry", zap.String("source", source), zap.Duration("timeSinceFirstEntry", timeSinceFirstEntry), zap.Duration("forceFlushTimeout", t.forceFlushTimeout))
				if timeSinceFirstEntry < t.forceFlushTimeout {
					customLogger.Debug("Batch not ready for force flush", zap.String("source", source))
					continue
				}
				customLogger.Debug("Force flush condition met, flushing source", zap.String("source", source))
				if err := t.flushSource(context.Background(), source); err != nil {
					t.Logger().Error("there was error flushing combined logs", zap.Error(err))
					customLogger.Error("Error flushing combined logs during force flush", zap.String("source", source), zap.Error(err))
				}
			}
			// check every 1/5 forceFlushTimeout
			initLogger().Debug("Resetting ticker interval", zap.Duration("interval", t.forceFlushTimeout/5))
			t.ticker.Reset(t.forceFlushTimeout / 5)
			t.Unlock()
		case <-t.chClose:
			t.ticker.Stop()
			return
		}
	}
}

func (t *Transformer) Stop() error {
	t.Lock()
	defer t.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	t.flushAllSources(ctx)

	close(t.chClose)
	return nil
}

func (t *Transformer) Process(ctx context.Context, e *entry.Entry) error {
	// Lock the recombine operator because process can't run concurrently
	t.Lock()
	defer t.Unlock()

	customLogger := initLogger()
	customLogger.Debug("Processing entry in recombine transformer was", zap.Any("entry", e)) // Get the environment for executing the expression.
	// In the future, we may want to provide access to the currently
	// batched entries so users can do comparisons to other entries
	// rather than just use absolute rules.
	env := helper.GetExprEnv(e)
	defer helper.PutExprEnv(env)

	m, err := expr.Run(t.prog, env)
	if err != nil {
		return t.HandleEntryError(ctx, e, err)
	}

	// this is guaranteed to be a boolean because of expr.AsBool
	matches := m.(bool)
	var s string
	err = e.Read(t.sourceIdentifier, &s)
	if err != nil {
		customLogger.Debug("Error in getting source identifier from entry", zap.Error(err))
		t.Logger().Warn("entry does not contain the source_identifier, so it may be pooled with other sources")
		s = DefaultSourceIdentifier
	}

	if s == "" {
		s = DefaultSourceIdentifier
	}

	switch {
	// This is the first entry in the next batch
	case matches && t.matchFirstLine:
		customLogger.Debug("First entry in new batch detected", zap.String("source", s), zap.Bool("matches", matches), zap.Bool("matchFirstLine", t.matchFirstLine))
		// Flush the existing batch
		if err := t.flushSource(ctx, s); err != nil {
			customLogger.Error("Error flushing existing batch before starting new batch", zap.String("source", s), zap.Error(err))
			return err
		}

		// Add the current log to the new batch
		t.addToBatch(ctx, e, s, matches)
		customLogger.Debug("Added entry to new batch", zap.String("source", s))
		return nil
	// This is the last entry in a complete batch
	case matches && !t.matchFirstLine:
		customLogger.Debug("Last entry in batch detected", zap.String("source", s), zap.Bool("matches", matches), zap.Bool("matchFirstLine", t.matchFirstLine))
		t.addToBatch(ctx, e, s, matches)
		customLogger.Debug("Added entry to batch before flush", zap.String("source", s))
		return t.flushSource(ctx, s)
	}

	// This is neither the first entry of a new log,
	// nor the last entry of a log, so just add it to the batch
	customLogger.Debug("Intermediate entry, adding to batch", zap.String("source", s), zap.Bool("matches", matches), zap.Bool("matchFirstLine", t.matchFirstLine))
	t.addToBatch(ctx, e, s, matches)
	return nil
}

// addToBatch adds the current entry to the current batch of entries that will be combined
func (t *Transformer) addToBatch(ctx context.Context, e *entry.Entry, source string, matches bool) {
	customLogger := initLogger()
	customLogger.Debug("Adding entry to batch", zap.String("source", source), zap.Any("entry", e))
	batch, ok := t.batchMap[source]
	if !ok {
		customLogger.Debug("No existing batch found for source, creating new batch", zap.String("source", source))
		if len(t.batchMap) >= t.maxSources {
			t.Logger().Error("Too many sources. Flushing all batched logs. Consider increasing max_sources parameter")
			t.flushAllSources(ctx)
		}
		batch = t.addNewBatch(source, e)
		customLogger.Debug("New batch created for source", zap.String("source", source))
	} else {
		batch.numEntries++
		customLogger.Debug("Incremented numEntries for batch", zap.String("source", source), zap.Int("numEntries", batch.numEntries))
		if t.overwriteWithNewest {
			batch.baseEntry = e
			customLogger.Debug("Overwriting baseEntry with newest entry", zap.String("source", source))
		}
	}

	// mark that match occurred to use max_unmatched_batch_size only when match didn't occur
	if matches && !batch.matchDetected {
		batch.matchDetected = true
		customLogger.Debug("Match detected for batch", zap.String("source", source))
	}

	// Combine the combineField of each entry in the batch,
	// separated by newlines
	var s string
	err := e.Read(t.combineField, &s)
	if err != nil {
		t.Logger().Error("entry does not contain the combine_field")
		customLogger.Error("Entry does not contain the combine_field", zap.String("source", source), zap.Error(err))
		return
	}
	if batch.recombined.Len() > 0 {
		batch.recombined.WriteString(t.combineWith)
		customLogger.Debug("Appended combineWith to recombined buffer", zap.String("source", source), zap.String("combineWith", t.combineWith))
	}
	batch.recombined.WriteString(s)
	customLogger.Debug("Appended entry field to recombined buffer", zap.String("source", source), zap.String("fieldValue", s))

	if (t.maxLogSize > 0 && int64(batch.recombined.Len()) > t.maxLogSize) ||
		batch.numEntries >= t.maxBatchSize ||
		(!batch.matchDetected && t.maxUnmatchedBatchSize > 0 && batch.numEntries >= t.maxUnmatchedBatchSize) {
		customLogger.Debug("Batch flush condition met", zap.String("source", source), zap.Int("numEntries", batch.numEntries), zap.Int("recombinedLen", batch.recombined.Len()))
		if err := t.flushSource(ctx, source); err != nil {
			t.Logger().Error("there was error flushing combined logs", zap.Error(err))
			customLogger.Error("Error flushing combined logs", zap.String("source", source), zap.Error(err))
		}
	}
}

// flushAllSources flushes all sources.
func (t *Transformer) flushAllSources(ctx context.Context) {
	customLogger := initLogger()
	customLogger.Debug("Flushing all sources in flushAllSources", zap.Int("numSources", len(t.batchMap)))
	var errs []error
	for source := range t.batchMap {
		customLogger.Debug("Flushing source in flushAllSources", zap.String("source", source))
		if err := t.flushSource(ctx, source); err != nil {
			customLogger.Error("Error flushing source in flushAllSources", zap.String("source", source), zap.Error(err))
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		customLogger.Error("There were errors flushing combined logs", zap.Int("numErrors", len(errs)), zap.Errors("errors", errs))
		t.Logger().Error("there was error flushing combined logs %s", zap.Error(errors.Join(errs...)))
	}
}

// flushSource combines the entries currently in the batch into a single entry,
// then forwards them to the next operator in the pipeline
func (t *Transformer) flushSource(ctx context.Context, source string) error {

	customLogger := initLogger()
	customLogger.Debug("Flushing source in recombine transformer", zap.String("source", source), zap.Any("batchMap", t.batchMap))
	batch := t.batchMap[source]
	// Skip flushing a combined log if the batch is empty
	if batch == nil {
		customLogger.Debug("No batch found for source, skipping flush", zap.String("source", source))
		return nil
	}

	if batch.baseEntry == nil {
		customLogger.Debug("Batch baseEntry is nil, removing batch", zap.String("source", source))
		t.removeBatch(source)
		return nil
	}

	// Set the recombined field on the entry
	recombinedStr := batch.recombined.String()
	customLogger.Debug("Setting recombined field on baseEntry", zap.String("combineField", t.combineField.String()), zap.String("recombined", recombinedStr))
	err := batch.baseEntry.Set(t.combineField, recombinedStr)
	if err != nil {
		customLogger.Error("Failed to set recombined field on baseEntry", zap.Error(err))
		return err
	}

	customLogger.Debug("Writing baseEntry to next operator", zap.Any("baseEntry", batch.baseEntry))
	err = t.Write(ctx, batch.baseEntry)
	if err != nil {
		customLogger.Error("Failed to write baseEntry to next operator", zap.Error(err))
	}
	t.removeBatch(source)
	customLogger.Debug("Batch removed after flush", zap.String("source", source))
	return err
}

// addNewBatch creates a new batch for the given source and adds the entry to it.
func (t *Transformer) addNewBatch(source string, e *entry.Entry) *sourceBatch {
	customLogger := initLogger()
	customLogger.Debug("Creating new batch in addNewBatch", zap.String("source", source), zap.Any("entry", e))

	batch := t.batchPool.Get().(*sourceBatch)
	batch.baseEntry = e
	batch.numEntries = 1
	batch.recombined.Reset()
	batch.firstEntryObservedTime = e.ObservedTimestamp
	batch.matchDetected = false
	t.batchMap[source] = batch
	customLogger.Debug("New batch created and added to batchMap", zap.String("source", source), zap.Any("batch", batch))
	return batch
}

// removeBatch removes the batch for the given source.
func (t *Transformer) removeBatch(source string) {
	batch := t.batchMap[source]
	delete(t.batchMap, source)
	t.batchPool.Put(batch)
}

func initLogger() *zap.Logger {
	writer := &lumberjack.Logger{
		Filename:   "/var/log/opsramp/debug.log", // or any path you prefer
		MaxSize:    10,                           // megabytes
		MaxBackups: 5,                            // number of old files to keep
		MaxAge:     30,                           // days to keep
		Compress:   true,                         // gzip
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(writer),
		zap.DebugLevel,
	)

	return zap.New(core)
}
