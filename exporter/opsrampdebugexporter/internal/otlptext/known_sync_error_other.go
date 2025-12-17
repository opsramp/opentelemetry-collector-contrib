// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !linux && !darwin && !windows

package otlptext // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opsrampdebugexporter/internal/otlptext"

// knownSyncError returns true if the given error is one of the known
// non-actionable errors returned by Sync on Plan 9.
func knownSyncError(err error) bool {
	return false
}
