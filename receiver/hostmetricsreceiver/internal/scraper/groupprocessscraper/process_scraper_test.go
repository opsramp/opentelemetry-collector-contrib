// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package groupprocessscraper

import (
	"context"
	"errors"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/scraper/scrapererror"
	"go.opentelemetry.io/collector/scraper/scrapertest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/groupprocessscraper/internal/metadata"
)

func skipTestOnUnsupportedOS(t *testing.T) {
	if runtime.GOOS != "linux" && runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
		t.Skipf("skipping test on %v", runtime.GOOS)
	}
}

func TestScrape(t *testing.T) {
	t.Skip("TODO: This test needs to be rewritten for groupprocessscraper which has different metrics than processscraper")
}

func TestScrapeMetrics_NewError(t *testing.T) {
	t.Skip("TODO: This test needs to be rewritten for groupprocessscraper")
}

func TestScrapeMetrics_GetProcessesError(t *testing.T) {
	skipTestOnUnsupportedOS(t)

	scraper, err := newGroupProcessScraper(scrapertest.NewNopSettings(metadata.Type), &Config{MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig()})
	require.NoError(t, err, "Failed to create process scraper: %v", err)

	scraper.getProcessHandles = func(context.Context) (processHandles, error) { return nil, errors.New("err1") }

	err = scraper.start(context.Background(), componenttest.NewNopHost())
	require.NoError(t, err, "Failed to initialize process scraper: %v", err)

	md, err := scraper.scrape(context.Background())
	assert.EqualError(t, err, "err1")
	assert.Equal(t, 0, md.ResourceMetrics().Len())
	assert.False(t, scrapererror.IsPartialScrapeError(err))
}
