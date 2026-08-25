// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package k8sattributesprocessor

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
)

type fakeSharedComponent struct {
	startErr  error
	starts    atomic.Int32
	shutdowns atomic.Int32
}

func (f *fakeSharedComponent) Start(context.Context, component.Host) error {
	f.starts.Add(1)
	return f.startErr
}

func (f *fakeSharedComponent) Shutdown(context.Context) error {
	f.shutdowns.Add(1)
	return nil
}

// The registry is reachable from several collectors running in the same
// process, so concurrent GetOrAdd must not race on the underlying map.
func TestSharedProcessorsConcurrentGetOrAdd(t *testing.T) {
	reg := newSharedProcessors()
	keys := []any{&struct{ a int }{1}, &struct{ b int }{2}, &struct{ c int }{3}}

	var created atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			sp := reg.GetOrAdd(keys[i%len(keys)], func() component.Component {
				created.Add(1)
				return &fakeSharedComponent{}
			})
			require.NotNil(t, sp)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, int32(len(keys)), created.Load(), "one instance per distinct config")
	assert.Len(t, reg.comps, len(keys))
}

func TestSharedProcessorSharesPerConfigIdentity(t *testing.T) {
	reg := newSharedProcessors()
	cfgA, cfgB := &Config{}, &Config{}

	first := reg.GetOrAdd(cfgA, func() component.Component { return &fakeSharedComponent{} })
	second := reg.GetOrAdd(cfgA, func() component.Component { return &fakeSharedComponent{} })
	other := reg.GetOrAdd(cfgB, func() component.Component { return &fakeSharedComponent{} })

	assert.Same(t, first, second, "same config is shared")
	assert.NotSame(t, first, other, "distinct configs are not shared")
	assert.Equal(t, 2, first.refs)
}

// A pipeline shutting down must not stop the processor while other pipelines
// are still using it.
func TestSharedProcessorShutdownIsRefCounted(t *testing.T) {
	reg := newSharedProcessors()
	cfg := &Config{}
	inner := &fakeSharedComponent{}

	create := func() component.Component { return inner }
	first := reg.GetOrAdd(cfg, create)
	second := reg.GetOrAdd(cfg, create)
	third := reg.GetOrAdd(cfg, create)

	ctx := context.Background()
	require.NoError(t, first.Shutdown(ctx))
	assert.Equal(t, int32(0), inner.shutdowns.Load(), "still referenced")

	require.NoError(t, second.Shutdown(ctx))
	assert.Equal(t, int32(0), inner.shutdowns.Load(), "still referenced")

	require.NoError(t, third.Shutdown(ctx))
	assert.Equal(t, int32(1), inner.shutdowns.Load(), "stopped on last release")
	assert.Empty(t, reg.comps, "entry dropped so it cannot leak")

	require.NoError(t, third.Shutdown(ctx))
	assert.Equal(t, int32(1), inner.shutdowns.Load(), "repeat shutdown is a no-op")
}

func TestSharedProcessorStartsOnce(t *testing.T) {
	reg := newSharedProcessors()
	cfg := &Config{}
	inner := &fakeSharedComponent{}

	create := func() component.Component { return inner }
	first := reg.GetOrAdd(cfg, create)
	second := reg.GetOrAdd(cfg, create)

	host := componenttest.NewNopHost()
	require.NoError(t, first.Start(context.Background(), host))
	require.NoError(t, second.Start(context.Background(), host))
	assert.Equal(t, int32(1), inner.starts.Load())
}

// Without a sticky error, only the first pipeline learns that the Kubernetes
// client failed and the rest run silently unenriched.
func TestSharedProcessorStartErrorIsReplayed(t *testing.T) {
	reg := newSharedProcessors()
	cfg := &Config{}
	wantErr := errors.New("boom")
	inner := &fakeSharedComponent{startErr: wantErr}

	create := func() component.Component { return inner }
	first := reg.GetOrAdd(cfg, create)
	second := reg.GetOrAdd(cfg, create)

	host := componenttest.NewNopHost()
	assert.ErrorIs(t, first.Start(context.Background(), host), wantErr)
	assert.ErrorIs(t, second.Start(context.Background(), host), wantErr)
	assert.Equal(t, int32(1), inner.starts.Load())
}

func TestSharedProcessorUnwrap(t *testing.T) {
	reg := newSharedProcessors()
	inner := &fakeSharedComponent{}
	sp := reg.GetOrAdd(&Config{}, func() component.Component { return inner })
	assert.Same(t, inner, sp.Unwrap())
}
