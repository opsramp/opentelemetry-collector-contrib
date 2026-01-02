// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package groupprocessscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/groupprocessscraper"

import (
	"context"
	"errors"
	"runtime"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/featuregate"
	"go.opentelemetry.io/collector/scraper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/groupprocessscraper/internal/metadata"
)

// This file implements Factory for Group Process scraper.

var (
	bootTimeCacheFeaturegateID = "hostmetrics.groupprocess.bootTimeCache"
	bootTimeCacheFeaturegate   = featuregate.GlobalRegistry().MustRegister(
		bootTimeCacheFeaturegateID,
		featuregate.StageBeta,
		featuregate.WithRegisterDescription("When enabled, all groupprocess scrapes will use the boot time value that is cached at the start of the process."),
		featuregate.WithRegisterReferenceURL("https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/28849"),
		featuregate.WithRegisterFromVersion("v0.98.0"),
	)
)

// NewFactory creates a new factory for Group Process scraper.
func NewFactory() scraper.Factory {
	return scraper.NewFactory(metadata.Type, createDefaultConfig, scraper.WithMetrics(createMetricsScraper, metadata.MetricsStability))
}

// createDefaultConfig creates the default configuration for the Scraper.
func createDefaultConfig() component.Config {
	return &Config{
		MetricsBuilderConfig: metadata.DefaultMetricsBuilderConfig(),
	}
}

// createMetricsScraper creates a resource scraper based on provided config.
func createMetricsScraper(
	_ context.Context,
	settings scraper.Settings,
	cfg component.Config,
) (scraper.Metrics, error) {
	if runtime.GOOS != "linux" && runtime.GOOS != "windows" && runtime.GOOS != "darwin" {
		return nil, errors.New("groupprocess scraper only available on Linux, Windows, or MacOS")
	}

	s, err := newGroupProcessScraper(settings, cfg.(*Config))
	if err != nil {
		return nil, err
	}

	return scraper.NewMetrics(
		s.scrape,
		scraper.WithStart(s.start),
	)
}
