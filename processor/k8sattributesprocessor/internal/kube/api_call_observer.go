// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package kube // import "github.com/open-telemetry/opentelemetry-collector-contrib/processor/k8sattributesprocessor/internal/kube"

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

const defaultK8sAPICallProofLogPath = "/var/log/opsramp/k8s-api-invocations.log"

type k8sAPICallObserver struct {
	logger *zap.Logger

	proofLogPath string
	proofLogFile *os.File

	mu     sync.Mutex
	total  int64
	counts map[string]int64

	proofLogWarned bool
}

func newK8sAPICallObserver(logger *zap.Logger) *k8sAPICallObserver {
	return &k8sAPICallObserver{
		logger:       logger,
		proofLogPath: defaultK8sAPICallProofLogPath,
		counts:       make(map[string]int64),
	}
}

func (o *k8sAPICallObserver) Record(verb, resource, namespace string) {
	if o == nil || o.logger == nil {
		return
	}

	key := verb + " " + resource + " ns=" + namespace

	o.mu.Lock()
	o.total++
	o.counts[key]++
	totalCalls := o.total
	resourceCalls := o.counts[key]
	ts := time.Now().UTC().Format(time.RFC3339Nano)
	proofLine := fmt.Sprintf("ts=%s verb=%s resource=%s namespace=%s resource_call_count=%d total_call_count=%d\n", ts, verb, resource, namespace, resourceCalls, totalCalls)
	proofErr := o.appendProofLineLocked(proofLine)
	shouldWarn := proofErr != nil && !o.proofLogWarned
	if shouldWarn {
		o.proofLogWarned = true
	}
	if proofErr == nil {
		o.proofLogWarned = false
	}
	o.mu.Unlock()

	if shouldWarn {
		o.logger.Warn("failed writing k8s api proof log", zap.String("path", o.proofLogPath), zap.Error(proofErr))
	}

	o.logger.Info(
		"k8s api invoked",
		zap.String("verb", verb),
		zap.String("resource", resource),
		zap.String("namespace", namespace),
		zap.Int64("resource_call_count", resourceCalls),
		zap.Int64("total_call_count", totalCalls),
	)
}

func (o *k8sAPICallObserver) LogSummary() {
	if o == nil || o.logger == nil {
		return
	}

	o.mu.Lock()
	totalCalls := o.total
	entries := make([]string, 0, len(o.counts))
	for key, count := range o.counts {
		entries = append(entries, key+"="+strconv.FormatInt(count, 10))
	}
	ts := time.Now().UTC().Format(time.RFC3339Nano)
	proofLine := fmt.Sprintf("ts=%s summary total_call_count=%d breakdown=%v\n", ts, totalCalls, entries)
	proofErr := o.appendProofLineLocked(proofLine)
	shouldWarn := proofErr != nil && !o.proofLogWarned
	if shouldWarn {
		o.proofLogWarned = true
	}
	if proofErr == nil {
		o.proofLogWarned = false
	}
	o.mu.Unlock()

	if shouldWarn {
		o.logger.Warn("failed writing k8s api proof log", zap.String("path", o.proofLogPath), zap.Error(proofErr))
	}

	sort.Strings(entries)
	o.logger.Info(
		"k8s api invocation summary",
		zap.Int64("total_call_count", totalCalls),
		zap.Strings("breakdown", entries),
	)
}

func (o *k8sAPICallObserver) Close() {
	if o == nil {
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()
	if o.proofLogFile != nil {
		_ = o.proofLogFile.Close()
		o.proofLogFile = nil
	}
}

func (o *k8sAPICallObserver) appendProofLineLocked(line string) error {
	if err := o.ensureProofLogFileLocked(); err != nil {
		return err
	}
	_, err := o.proofLogFile.WriteString(line)
	return err
}

func (o *k8sAPICallObserver) ensureProofLogFileLocked() error {
	if o.proofLogFile != nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(o.proofLogPath), 0o755); err != nil {
		return err
	}

	f, err := os.OpenFile(o.proofLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	o.proofLogFile = f
	return nil
}
