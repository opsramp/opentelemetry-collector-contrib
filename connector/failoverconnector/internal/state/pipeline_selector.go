// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package state // import "github.com/open-telemetry/opentelemetry-collector-contrib/connector/failoverconnector/internal/state"

import (
	"context"
	"reflect"
	"sync/atomic"
	"time"
)

// PipelineSelector is meant to serve as the source of truth for the target priority level
type PipelineSelector struct {
	currentIndex    atomic.Int32
	stableIndex     atomic.Int32
	pipelineRetries []atomic.Int32
	constants       PSConstants
	RS              *RetryState

	errTryLock     *TryLock
	stableTryLock  *TryLock
	chans          []chan bool
	nextIndex      chan int
	newStableIndex chan int
}

func (p *PipelineSelector) handlePipelineError(idx int) {
	if idx != p.loadCurrent() {
		return
	}
	doRetry := p.indexIsStable(idx)
	p.updatePipelineIndex(idx)
	if !doRetry {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	p.RS.InvokeCancel()
	p.RS.UpdateCancelFunc(cancel)
	p.enableRetry(ctx)
}

func (p *PipelineSelector) enableRetry(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(p.constants.RetryInterval)
		defer ticker.Stop()

		stableIndex := p.loadStable()
		var cancelFunc context.CancelFunc
		for p.checkContinueRetry(stableIndex) {
			select {
			case <-ticker.C:
				if cancelFunc != nil {
					cancelFunc()
				}
				cancelFunc = p.handleRetry(ctx, stableIndex)
			case <-ctx.Done():
				return
			}
		}
		p.RS.InvokeCancel()
	}()
}

// handleRetry is responsible for launching goroutine and returning cancelFunc
func (p *PipelineSelector) handleRetry(parentCtx context.Context, stableIndex int) context.CancelFunc {
	retryCtx, cancelFunc := context.WithCancel(parentCtx)
	go p.retryHighPriorityPipelines(retryCtx, stableIndex, p.constants.RetryGap)
	return cancelFunc
}

// UpdatePipelineIndex is the main function that updates the pipeline indexes due to an error
func (p *PipelineSelector) updatePipelineIndex(idx int) {
	if p.indexIsStable(idx) {
		p.setToNextPriorityPipeline(idx)
		return
	}
	p.setToStableIndex(idx)
}

// NextPipeline skips through any lower priority pipelines that have exceeded their maxRetries
func (p *PipelineSelector) setToNextPriorityPipeline(idx int) {
	for ok := true; ok; ok = p.exceededMaxRetries(idx) {
		idx++
	}
	p.newStableIndex <- idx
}

// RetryHighPriorityPipelines responsible for single iteration through all higher priority pipelines
func (p *PipelineSelector) retryHighPriorityPipelines(ctx context.Context, stableIndex int, retryGap time.Duration) {
	ticker := time.NewTicker(retryGap)

	defer ticker.Stop()

	for i := 0; i < stableIndex; i++ {
		if stableIndex > p.loadStable() {
			return
		}
		if p.maxRetriesUsed(i) {
			continue
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.nextIndex <- i
		}
	}
}

// checkStopRetry checks if retry should be suspended if all higher priority levels have exceeded their max retries
func (p *PipelineSelector) checkContinueRetry(index int) bool {
	for i := 0; i < index; i++ {
		if p.loadRetryCount(i) < p.constants.MaxRetries {
			return true
		}
	}
	return false
}

func (p *PipelineSelector) exceededMaxRetries(idx int) bool {
	return idx < len(p.pipelineRetries) && (p.loadRetryCount(idx) >= p.constants.MaxRetries)
}

// SetToStableIndex returns the CurrentIndex to the known Stable Index
func (p *PipelineSelector) setToStableIndex(idx int) {
	p.incrementRetryCount(idx)
	p.nextIndex <- p.loadStable()
}

// MaxRetriesUsed exported access to maxRetriesUsed
func (p *PipelineSelector) maxRetriesUsed(idx int) bool {
	return p.loadRetryCount(idx) >= p.constants.MaxRetries
}

// SetNewStableIndex Update stableIndex to the passed stable index
func (p *PipelineSelector) setNewStableIndex(idx int) {
	p.resetRetryCount(idx)
	p.newStableIndex <- idx
}

// IndexIsStable returns if index passed is the stable index
func (p *PipelineSelector) indexIsStable(idx int) bool {
	return p.loadStable() == idx
}

func (p *PipelineSelector) loadStable() int {
	return int(p.stableIndex.Load())
}

func (p *PipelineSelector) loadCurrent() int {
	return int(p.currentIndex.Load())
}

func (p *PipelineSelector) loadRetryCount(idx int) int {
	return int(p.pipelineRetries[idx].Load())
}

func (p *PipelineSelector) incrementRetryCount(idx int) {
	p.pipelineRetries[idx].Add(1)
}

func (p *PipelineSelector) resetRetryCount(idx int) {
	p.pipelineRetries[idx].Store(0)
}

// ReportStable reports back to the failoverRouter that the current priority was stable
func (p *PipelineSelector) reportStable(idx int) {
	if p.indexIsStable(idx) {
		return
	}
	p.setNewStableIndex(idx)
}

func NewPipelineSelector(lenPriority int, consts PSConstants) *PipelineSelector {
	chans := make([]chan bool, lenPriority)

	for i := 0; i < lenPriority; i++ {
		chans[i] = make(chan bool)
	}

	ps := &PipelineSelector{
		pipelineRetries: make([]atomic.Int32, lenPriority),
		constants:       consts,
		RS:              &RetryState{},
		errTryLock:      NewTryLock(),
		stableTryLock:   NewTryLock(),
		chans:           chans,
		nextIndex:       make(chan int),
		newStableIndex:  make(chan int),
	}
	return ps
}

func (p *PipelineSelector) Start(done chan struct{}) {
	go p.ManageIndexes(done)
	go p.ListenToChannels(done)
}

func (p *PipelineSelector) ManageIndexes(done chan struct{}) {
	for {
		select {
		case <-done:
			return
		case idx := <-p.nextIndex:
			p.currentIndex.Store(int32(idx))
		case idx := <-p.newStableIndex:
			p.stableIndex.Store(int32(idx))
		}
	}
}

func (p *PipelineSelector) ListenToChannels(done chan struct{}) {
	cases := make([]reflect.SelectCase, len(p.chans)+1)
	for i, ch := range p.chans {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
	}
	cases[len(p.chans)] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(done)}

	for {
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			return
		}
		if value.Bool() {
			p.stableTryLock.TryExecute(p.reportStable, chosen)
		} else {
			p.errTryLock.TryExecute(p.handlePipelineError, chosen)
		}
	}
}

func (p *PipelineSelector) SelectedPipeline() (int, chan bool) {
	idx := p.loadCurrent()
	if idx < len(p.chans) {
		return idx, p.chans[idx]
	}
	return idx, nil
}

// For Testing
func (p *PipelineSelector) ChannelIndex(ch chan bool) int {
	for i, ch1 := range p.chans {
		if ch == ch1 {
			return i
		}
	}
	return -1
}
