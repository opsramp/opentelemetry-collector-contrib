// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package windows // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/input/windows"

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
)

// Input is an operator that creates entries using the windows event log api.
type (
	Input struct {
		helper.InputOperator
		bookmark            Bookmark
		subscription        Subscription
		buffer              Buffer
		channel             string
		maxReads            int
		startAt             string
		raw                 bool
		isReqAdditionalAttr *bool
		excludeProviders    []string
		pollInterval        time.Duration
		persister           Persister
		publisherCache      publisherCache
		cancel              context.CancelFunc
		wg                  sync.WaitGroup
	}
)

// Start will start reading events from a subscription.
func (i *Input) Start(persister operator.Persister) error {
	ctx, cancel := context.WithCancel(context.Background())
	i.cancel = cancel

	i.subscription = NewSubscription()
	if err := i.subscription.Open(i.channel, i.startAt, i.bookmark); err != nil {
		return fmt.Errorf("failed to open subscription: %w", err)
	}

	i.publisherCache = newPublisherCache()

	i.wg.Add(1)
	go i.readOnInterval(ctx)
	return nil
}

// Stop will stop reading events from a subscription.
func (i *Input) Stop() error {
	i.cancel()
	i.wg.Wait()

	if err := i.subscription.Close(); err != nil {
		return fmt.Errorf("failed to close subscription: %w", err)
	}

	if err := i.bookmark.Close(); err != nil {
		return fmt.Errorf("failed to close bookmark: %w", err)
	}

	if err := i.publisherCache.evictAll(); err != nil {
		return fmt.Errorf("failed to close publishers: %w", err)
	}

	return nil
}

// readOnInterval will read events with respect to the polling interval.
func (i *Input) readOnInterval(ctx context.Context) {
	defer i.wg.Done()

	ticker := time.NewTicker(i.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			i.readToEnd(ctx)
		}
	}
}

// readToEnd will read events from the subscription until it reaches the end of the channel.
func (i *Input) readToEnd(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if count := i.read(ctx); count == 0 {
				return
			}
		}
	}
}

// read will read events from the subscription.
func (i *Input) read(ctx context.Context) int {
	events, err := i.subscription.Read(i.maxReads)
	if err != nil {
		i.Logger().Error("Failed to read events from subscription", zap.Error(err))
		return 0
	}

	pstartedAtObj := time.Now().Local()
	pstartedAt := pstartedAtObj.Format("2006-01-02 15:04:05.000")

	for eI, event := range events {
		simpleEvent, err := event.RenderSimple(i.buffer)
		if err != nil {
			i.Logger().Error("Failed to render simple event", zap.Error(err))
			event.Close()
			continue
		}

		recordID := simpleEvent.RecordID
		dedupKey := fmt.Sprintf("dedup_%d", recordID)
		i.Logger().Debug("@@@ Suresh - stage 1 --- Processing event", zap.Any("pstartedAt", pstartedAt), zap.Int("current loop id:", eI), zap.Any("process started at ", time.Now().Local().Format("2006-01-02 15:04:05.000")), zap.Uint64("record_id", recordID), zap.String("channel", i.channel))

		// Deduplication: check if this record was already processed
		exists, dCheckErr := i.persister.Get(ctx, dedupKey)
		if dCheckErr != nil {
			i.Logger().Error("Failed to check deduplication key", zap.Error(dCheckErr))
			event.Close()
			continue
		}
		i.Logger().Debug("@@@ Suresh - stage 2 --- Checking deduplication", zap.Any("pstartedAt", pstartedAt), zap.Int("current loop id:", eI), zap.Any("process started at ", time.Now().Local().Format("2006-01-02 15:04:05.000")), zap.Uint64("record_id", recordID), zap.String("dedup_key", dedupKey), zap.Any("exists", exists))

		if exists != nil {
			i.Logger().Debug("Duplicate event, skipping", zap.Uint64("record_id", recordID))
			event.Close()
			continue
		}

		i.Logger().Debug("@@@ Suresh - stage 3 --- Processing event with simple", zap.Any("pstartedAt", pstartedAt), zap.Int("current loop id:", eI), zap.Any("process started at ", time.Now().Local().Format("2006-01-02 15:04:05.000")), zap.Uint64("record_id", recordID), zap.String("channel", i.channel))

		// Skip empty events
		if recordID == 0 && simpleEvent.Provider.Name == "" {
			i.Logger().Debug("Skipping empty event")
			event.Close()
			continue
		}

		i.Logger().Debug("@@Suresh - Reading event", zap.Any("pstartedAt", pstartedAt), zap.Int("current loop id:", eI), zap.Any("process started at ", time.Now().Local().Format("2006-01-02 15:04:05.000")), zap.Uint64("record_id", recordID), zap.String("channel", i.channel))

		err1 := i.processEventWithSimple(ctx, event, &simpleEvent)
		if err1 == nil {
			i.updateBookmarkOffset(ctx, event, recordID)
			// Persist processed RecordID for deduplication
			if err := i.persister.Set(ctx, dedupKey, []byte("1")); err != nil {
				i.Logger().Error("Failed to persist deduplication key", zap.Error(err))
			}
			i.Logger().Debug("@@@ Suresh - Processed successfully, updating bookmark", zap.Any("pstartedAt", pstartedAt), zap.Int("current loop id:", eI), zap.Any("process started at ", time.Now().Local().Format("2006-01-02 15:04:05.000")), zap.Uint64("record_id", recordID))
		} else {
			i.Logger().Error("Failed to process event", zap.Any("pstartedAt", pstartedAt), zap.Int("current loop id:", eI), zap.Any("process started at ", time.Now().Local().Format("2006-01-02 15:04:05.000")), zap.Uint64("record_id", recordID), zap.Error(err))
		}
		event.Close()
	}

	processDuration := time.Now().Local().Sub(pstartedAtObj).Milliseconds()
	i.Logger().Debug(
		"@@@ Suresh - Finished processing events",
		zap.Any("pstartedAt", pstartedAt),
		zap.Any("process finished at ", time.Now().Local().Format("2006-01-02 15:04:05.000")),
		zap.Any("process_duration_ms", processDuration),
	)

	return len(events)
}

// processEvent will process and send an event retrieved from windows event log.
func (i *Input) processEvent(ctx context.Context, event Event) error {
	i.Logger().Debug("@@@ Suresh - Processing event, stag 1:", zap.Any("i", i))

	if i.raw {
		rawEvent, err := event.RenderRaw(i.buffer)
		i.Logger().Debug("@@@ Suresh - Processing event, stag 1.1:", zap.Any("i", i))
		if err != nil {
			i.Logger().Error("Failed to render raw event", zap.Error(err))
			return err
		}
		i.Logger().Debug("@@@ Suresh - Processing event, stag 1.2:", zap.Any("i", i))
		i.sendEventRaw(ctx, rawEvent)
		return nil
	}

	i.Logger().Debug("@@@ Suresh - Processing event, stag 2:", zap.Any("i", i))

	isExcluded := func(providerName string) bool {
		for _, excludeProvider := range i.excludeProviders {
			if providerName == excludeProvider {
				return true
			}
		}
		return false
	}

	i.Logger().Debug("@@@ Suresh - Processing event, stag 3:", zap.Any("i", i))
	simpleEvent, err := event.RenderSimple(i.buffer)
	i.Logger().Debug("@@@ Suresh - Processing event --> inside raw, simpleEvent", zap.Any("simpleEvent", simpleEvent), zap.Error(err))
	if err != nil {
		i.Logger().Error("Failed to render simple event", zap.Error(err))
		return err
	}
	if isExcluded(simpleEvent.Provider.Name) {
		i.Logger().Debug("@@@ Suresh - Processing event --> inside isExcluded", zap.Any("simpleEvent.Provider.Name", simpleEvent.Provider.Name), zap.Any("check -=--- isExcluded", isExcluded(simpleEvent.Provider.Name)))
		return nil
	}

	publisher, openPublisherErr := i.publisherCache.get(simpleEvent.Provider.Name)
	if openPublisherErr != nil {
		i.Logger().Warn(
			"Failed to open event source, respective log entries cannot be formatted",
			zap.String("provider", simpleEvent.Provider.Name), zap.Error(openPublisherErr))
	}
	if !publisher.Valid() {
		return i.sendEvent(ctx, &simpleEvent)
	}

	formattedEvent, err := event.RenderFormatted(i.buffer, publisher)
	if err != nil {
		i.Logger().Error("Failed to render formatted event", zap.Error(err))
		return i.sendEvent(ctx, &simpleEvent)
	}
	return i.sendEvent(ctx, &formattedEvent)
}

func (i *Input) sendEvent(ctx context.Context, eventXML *EventXML) error {
	var body any = eventXML.Original
	if !i.raw {
		body = formattedBody(eventXML)
	}

	e, err := i.NewEntry(body)
	i.Logger().Debug("@@@@ Suresh - stage 1-- Sending event", zap.Any("entry", e), zap.Any("body", body), zap.Error(err))
	if err != nil {
		i.Logger().Error("sendEvent -> Failed to create new entry", zap.Error(err))
		return err
	}

	e.Timestamp = parseTimestamp(eventXML.TimeCreated.SystemTime)
	e.Severity = parseSeverity(eventXML.RenderedLevel, eventXML.Level)

	eventData, er := i.ExtractEventData(eventXML.EventData)
	i.Logger().Debug("Suresh - stage 2- Debugging -- **** inside --- i.isAdditionalAttrReq - Extracted event data", zap.Any("event_data", eventData), zap.Any("err", er))
	if len(eventData) > 0 {
		for eK, eD := range eventData {
			eK = strings.ReplaceAll(strings.ToLower(eK), " ", "_")
			if str, ok := eD.(string); ok {
				e.AddAttribute(eK, str)
			} else {
				e.AddAttribute(eK, fmt.Sprintf("%v", eD))
			}
		}
	}
	if i.isAdditionalAttrReq() {
		e.AddAttribute("log_record_original", eventXML.Original)
		i.Logger().Debug("Suresh - stage 2- Debugging - Extracted event data", zap.Any("event_data", eventData), zap.Any("err", er))
	}
	return i.Write(ctx, e)
}

func (i *Input) sendEventRaw(ctx context.Context, eventRaw EventRaw) {
	body := eventRaw.parseBody()
	entry, err := i.NewEntry(body)
	if err != nil {
		i.Logger().Error("Failed to create entry", zap.Error(err))
		return
	}

	entry.Timestamp = eventRaw.parseTimestamp()
	entry.Severity = eventRaw.parseRenderedSeverity()
	i.Write(ctx, entry)
}

// getBookmarkXML will get the bookmark xml from the offsets database.
func (i *Input) getBookmarkOffset(ctx context.Context) (string, error) {
	bytes, err := i.persister.Get(ctx, i.channel)
	return string(bytes), err
}

// updateBookmarkOffset updates the bookmark xml and saves it in the offsets database.
func (i *Input) updateBookmarkOffset(ctx context.Context, event Event, recordID uint64) {
	if err := i.bookmark.Update(event); err != nil {
		i.Logger().Error("Failed to update bookmark", zap.Error(err), zap.Uint64("record_id", recordID))
		return
	}
	bookmarkXML, err := i.bookmark.Render(i.buffer)
	if err != nil {
		i.Logger().Error("Failed to render bookmark XML", zap.Error(err))
		return
	}
	if err = i.persister.Set(ctx, i.channel, []byte(bookmarkXML)); err != nil {
		i.Logger().Error("Failed to persist bookmark offset", zap.Error(err), zap.Uint64("record_id", recordID))
		return
	}
}

func (i *Input) isAdditionalAttrReq() bool {
	if i.isReqAdditionalAttr == nil {
		return true
	}
	return *i.isReqAdditionalAttr
}

func (i *Input) ExtractEventData(eventData EventData) (map[string]any, error) {
	result := make(map[string]any)
	for _, data := range eventData.Data {
		switch data.Name {
		case "LogonType":
			result["logon_type"] = data.Value
		case "AccountDomain":
			result["account_domain"] = data.Value
		case "AccountName":
			result["account_name"] = data.Value
		case "SecurityID", "TargetUserSid", "SubjectUserSid":
			result["security_id"] = data.Value
		case "LogonID", "TargetLogonId", "SubjectLogonId":
			result["logon_id"] = data.Value
		case "TargetUserName":
			result["user_name"] = data.Value
		case "TargetDomainName", "WorkstationName":
			result["domain_name"] = data.Value
		}
	}
	return result, nil
}

func (i *Input) processEventWithSimple(ctx context.Context, event Event, simpleEvent *EventXML) error {
	i.Logger().Debug("@@@ Suresh processEventWithSimple", zap.String("provider", simpleEvent.Provider.Name))

	if i.raw {
		rawEvent, err := event.RenderRaw(i.buffer)
		if err != nil {
			i.Logger().Error("Failed to render raw event", zap.Error(err))
			return err
		}
		i.sendEventRaw(ctx, rawEvent)
		return nil
	}

	isExcluded := func(providerName string) bool {
		for _, excludeProvider := range i.excludeProviders {
			if providerName == excludeProvider {
				return true
			}
		}
		return false
	}

	if isExcluded(simpleEvent.Provider.Name) {
		i.Logger().Debug("Event skipped due to excluded provider", zap.String("provider", simpleEvent.Provider.Name))
		return nil
	}

	publisher, err := i.publisherCache.get(simpleEvent.Provider.Name)
	if err != nil || !publisher.Valid() {
		i.Logger().Warn("Fallback to simple event due to invalid publisher", zap.Error(err))
		return i.sendEvent(ctx, simpleEvent)
	}

	formattedEvent, err := event.RenderFormatted(i.buffer, publisher)
	if err != nil {
		i.Logger().Error("Failed to render formatted event", zap.Error(err))
		return i.sendEvent(ctx, simpleEvent)
	}

	return i.sendEvent(ctx, &formattedEvent)
}
