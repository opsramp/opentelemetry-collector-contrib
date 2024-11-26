// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package groupprocessscraper // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/groupprocessscraper"

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/common"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/host"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scrapererror"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/filter/filterset"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/groupprocessscraper/internal/handlecount"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/groupprocessscraper/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/hostmetricsreceiver/internal/scraper/groupprocessscraper/ucal"
)

const (
	cpuMetricsLen               = 1
	memoryMetricsLen            = 2
	memoryUtilizationMetricsLen = 1
	diskMetricsLen              = 1
	pagingMetricsLen            = 1
	threadMetricsLen            = 1
	contextSwitchMetricsLen     = 1
	fileDescriptorMetricsLen    = 1
	handleMetricsLen            = 1
	signalMetricsLen            = 1

	metricsLen = cpuMetricsLen + memoryMetricsLen + diskMetricsLen + memoryUtilizationMetricsLen + pagingMetricsLen + threadMetricsLen + contextSwitchMetricsLen + fileDescriptorMetricsLen + signalMetricsLen
)

// scraper for Process Metrics
type scraper struct {
	settings           receiver.Settings
	config             *Config
	mb                 *metadata.MetricsBuilder
	scrapeProcessDelay time.Duration
	ucals              map[int32]*ucal.CPUUtilizationCalculator
	logicalCores       int

	matchGroupFS map[string]filterset.FilterSet

	// for mocking
	getProcessCreateTime func(p processHandle, ctx context.Context) (int64, error)
	getProcessHandles    func(context.Context) (processHandles, error)

	handleCountManager handlecount.Manager
}

// newGroupProcessScraper creates a Process Scraper
func newGroupProcessScraper(settings receiver.Settings, cfg *Config) (*scraper, error) {
	scraper := &scraper{
		settings:             settings,
		config:               cfg,
		getProcessCreateTime: processHandle.CreateTimeWithContext,
		getProcessHandles:    getProcessHandlesInternal,
		scrapeProcessDelay:   cfg.ScrapeProcessDelay,
		ucals:                make(map[int32]*ucal.CPUUtilizationCalculator),
		handleCountManager:   handlecount.NewManager(),
	}

	var err error

	for _, gc := range cfg.GroupConfig {
		fs, err := filterset.CreateFilterSet(gc.Names, &gc.Config)
		if err != nil {
			return nil, fmt.Errorf("error creating process exclude filters: %w", err)
		}
		scraper.matchGroupFS[gc.GroupName] = fs
	}

	logicalCores, err := cpu.Counts(true)
	if err != nil {
		return nil, fmt.Errorf("error getting number of logical cores: %w", err)
	}

	scraper.logicalCores = logicalCores

	return scraper, nil
}

func (s *scraper) start(context.Context, component.Host) error {
	s.mb = metadata.NewMetricsBuilder(s.config.MetricsBuilderConfig, s.settings)
	return nil
}

func (s *scraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	var errs scrapererror.ScrapeErrors

	// If the boot time cache featuregate is disabled, this will refresh the
	// cached boot time value for use in the current scrape. This functionally
	// replicates the previous functionality in all but the most extreme
	// cases of boot time changing in the middle of a scrape.
	if !bootTimeCacheFeaturegate.IsEnabled() {
		host.EnableBootTimeCache(false)
		_, err := host.BootTimeWithContext(ctx)
		if err != nil {
			errs.AddPartial(1, fmt.Errorf(`retrieving boot time failed with error "%w", using cached boot time`, err))
		}
		host.EnableBootTimeCache(true)
	}

	data, err := s.getProcessMetadata()
	if err != nil {
		var partialErr scrapererror.PartialScrapeError
		if !errors.As(err, &partialErr) {
			return pmetric.NewMetrics(), err
		}

		errs.AddPartial(partialErr.Failed, partialErr)
	}

	presentPIDs := make(map[int32]struct{}, len(data))
	ctx = context.WithValue(ctx, common.EnvKey, s.config.EnvMap)

	for groupName, processMetadataList := range data {
		var totalCPUTime float64
		var totalMemoryUsage int64
		var totalProcessCount int64
		var totalThreadCount int64
		var totalOpenFDCount int64

		for _, md := range processMetadataList {
			presentPIDs[md.pid] = struct{}{}
			times, err := md.handle.TimesWithContext(ctx)
			if err != nil {
				errs.AddPartial(cpuMetricsLen, fmt.Errorf("error reading cpu times for process %q (pid %v): %w", md.executable.name, md.pid, err))
				continue
			}
			totalCPUTime += times.User + times.System + times.Idle + times.Nice + times.Iowait + times.Irq + times.Softirq + times.Steal + times.Guest + times.GuestNice

			mem, err := md.handle.MemoryInfoWithContext(ctx)
			if err != nil {
				errs.AddPartial(memoryMetricsLen, fmt.Errorf("error reading memory info for process %q (pid %v): %w", md.executable.name, md.pid, err))
				continue
			}
			totalMemoryUsage += int64(mem.RSS)

			threads, err := md.handle.NumThreadsWithContext(ctx)
			if err != nil {
				errs.AddPartial(threadMetricsLen, fmt.Errorf("error reading thread info for process %q (pid %v): %w", md.executable.name, md.pid, err))
				continue
			}
			totalThreadCount += int64(threads)

			fds, err := md.handle.NumFDsWithContext(ctx)
			if err != nil {
				errs.AddPartial(fileDescriptorMetricsLen, fmt.Errorf("error reading open file descriptor count for process %q (pid %v): %w", md.executable.name, md.pid, err))
				continue
			}
			totalOpenFDCount += int64(fds)

			totalProcessCount++
		}

		now := pcommon.NewTimestampFromTime(time.Now())
		s.mb.RecordGroupProcessCPUTimeDataPoint(now, totalCPUTime, metadata.AttributeStateTotal)
		s.mb.RecordGroupProcessMemoryUsageDataPoint(now, totalMemoryUsage)
		s.mb.RecordGroupProcessCountDataPoint(now, totalProcessCount)
		s.mb.RecordGroupProcessThreadsDataPoint(now, totalThreadCount)
		s.mb.RecordGroupProcessOpenFileDescriptorsDataPoint(now, totalOpenFDCount)

		s.mb.EmitForResource(metadata.WithResource(s.buildGroupResource(s.mb.NewResourceBuilder(), groupName)))
	}

	if s.config.MuteProcessAllErrors {
		return s.mb.Emit(), nil
	}

	return s.mb.Emit(), errs.Combine()
}

// getProcessMetadata returns a slice of processMetadata, including handles,
// for all currently running processes. If errors occur obtaining information
// for some processes, an error will be returned, but any processes that were
// successfully obtained will still be returned.
func (s *scraper) getProcessMetadata() (map[string][]*processMetadata, error) {
	ctx := context.WithValue(context.Background(), common.EnvKey, s.config.EnvMap)
	handles, err := s.getProcessHandles(ctx)
	if err != nil {
		return nil, err
	}

	var errs scrapererror.ScrapeErrors

	//data := make([]*processMetadata, 0, handles.Len())
	data := make(map[string][]*processMetadata)
	for i := 0; i < handles.Len(); i++ {
		pid := handles.Pid(i)
		handle := handles.At(i)

		exe, err := getProcessExecutable(ctx, handle)
		if err != nil {
			if !s.config.MuteProcessExeError {
				errs.AddPartial(1, fmt.Errorf("error reading process executable for pid %v: %w", pid, err))
			}
		}

		name, err := getProcessName(ctx, handle, exe)
		if err != nil {
			if !s.config.MuteProcessNameError {
				errs.AddPartial(1, fmt.Errorf("error reading process name for pid %v: %w", pid, err))
			}
			continue
		}
		cgroup, err := getProcessCgroup(ctx, handle)
		if err != nil {
			if !s.config.MuteProcessCgroupError {
				errs.AddPartial(1, fmt.Errorf("error reading process cgroup for pid %v: %w", pid, err))
			}
			continue
		}

		executable := &executableMetadata{name: name, path: exe, cgroup: cgroup}

		groupConfigName := ""
		for groupName, fs := range s.matchGroupFS {
			if fs.Matches(executable.name) {
				groupConfigName = groupName
				break
			}
		}

		if groupConfigName == "" {
			continue
		}

		command, err := getProcessCommand(ctx, handle)
		if err != nil {
			errs.AddPartial(0, fmt.Errorf("error reading command for process %q (pid %v): %w", executable.name, pid, err))
		}

		username, err := handle.UsernameWithContext(ctx)
		if err != nil {
			if !s.config.MuteProcessUserError {
				errs.AddPartial(0, fmt.Errorf("error reading username for process %q (pid %v): %w", executable.name, pid, err))
			}
		}

		createTime, err := s.getProcessCreateTime(handle, ctx)
		if err != nil {
			errs.AddPartial(0, fmt.Errorf("error reading create time for process %q (pid %v): %w", executable.name, pid, err))
			// set the start time to now to avoid including this when a scrape_process_delay is set
			createTime = time.Now().UnixMilli()
		}
		if s.scrapeProcessDelay.Milliseconds() > (time.Now().UnixMilli() - createTime) {
			continue
		}

		parentPid, err := parentPid(ctx, handle, pid)
		if err != nil {
			errs.AddPartial(0, fmt.Errorf("error reading parent pid for process %q (pid %v): %w", executable.name, pid, err))
		}

		md := &processMetadata{
			pid:             pid,
			parentPid:       parentPid,
			executable:      executable,
			command:         command,
			username:        username,
			handle:          handle,
			createTime:      createTime,
			groupConfigName: groupConfigName,
		}

		data[groupConfigName] = append(data[groupConfigName], md)
	}

	return data, errs.Combine()
}
