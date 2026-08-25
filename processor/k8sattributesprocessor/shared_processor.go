// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sattributesprocessor // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor"

import (
	"context"
	"sync"

	"go.opentelemetry.io/collector/component"
)

// sharedProcessors keeps one processor instance per config so that pipelines
// referencing the same processor share a single set of Kubernetes watchers.
//
// It replaces internal/sharedcomponent, which guards nothing and assumes a
// single collector per process. Hosts that build several collectors from
// independent goroutines can otherwise race on the registry map, and Go turns a
// concurrent map write into an unrecoverable fatal error. It also refcounts, so
// one pipeline shutting down no longer stops the processor for the others.
type sharedProcessors struct {
	mu    sync.Mutex
	comps map[any]*sharedProcessor
}

func newSharedProcessors() *sharedProcessors {
	return &sharedProcessors{comps: make(map[any]*sharedProcessor)}
}

// GetOrAdd returns the instance registered for key, creating it on first use.
// Every call takes a reference that must be released through Shutdown.
func (scs *sharedProcessors) GetOrAdd(key any, create func() component.Component) *sharedProcessor {
	scs.mu.Lock()
	defer scs.mu.Unlock()

	if c, ok := scs.comps[key]; ok {
		c.refs++
		return c
	}

	c := &sharedProcessor{Component: create(), reg: scs, key: key, refs: 1}
	scs.comps[key] = c
	return c
}

type sharedProcessor struct {
	component.Component

	reg *sharedProcessors
	key any

	// refs and stopped are guarded by reg.mu.
	refs    int
	stopped bool

	startMu  sync.Mutex
	started  bool
	startErr error
}

// Unwrap returns the wrapped component.
func (r *sharedProcessor) Unwrap() component.Component {
	return r.Component
}

// Start starts the wrapped component once. The first result is retained and
// replayed to every later caller, so a pipeline cannot mistake a failed start
// for a successful one and silently run without enrichment.
func (r *sharedProcessor) Start(ctx context.Context, host component.Host) error {
	// Deliberately not the registry lock: the underlying Start blocks while it
	// syncs metadata, which would stall unrelated collectors.
	r.startMu.Lock()
	defer r.startMu.Unlock()

	if !r.started {
		r.started = true
		r.startErr = r.Component.Start(ctx, host)
	}
	return r.startErr
}

// Shutdown releases one reference and stops the wrapped component once the last
// pipeline using it goes away.
func (r *sharedProcessor) Shutdown(ctx context.Context) error {
	r.reg.mu.Lock()
	if r.stopped {
		r.reg.mu.Unlock()
		return nil
	}
	r.refs--
	if r.refs > 0 {
		r.reg.mu.Unlock()
		return nil
	}
	r.stopped = true
	// Only drop the entry if it is still ours; a later GetOrAdd may have
	// replaced it after a previous shutdown.
	if cur, ok := r.reg.comps[r.key]; ok && cur == r {
		delete(r.reg.comps, r.key)
	}
	r.reg.mu.Unlock()

	return r.Component.Shutdown(ctx)
}
