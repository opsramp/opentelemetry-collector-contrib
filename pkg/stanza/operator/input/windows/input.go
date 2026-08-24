// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:build windows

package windows // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/input/windows"

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sys/windows"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/internal/metadata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
)

// Input is an operator that creates entries using the windows event log api.
type Input struct {
	helper.InputOperator
	bookmark                 Bookmark
	buffer                   *Buffer
	channel                  string
	ignoreChannelErrors      bool
	query                    *string
	maxReads                 int
	currentMaxReads          int
	startAt                  string
	raw                      bool
	eventDataFormat          EventDataFormat
	includeLogRecordOriginal bool
	excludeProviders         map[string]struct{}
	pollInterval             time.Duration
	waitTimeout              time.Duration
	// cancelEvent is a manual-reset Windows event handle signaled by Stop() to unblock
	// WaitForMultipleObjects in awaitAndReadEvents. A plain context cancellation cannot
	// interrupt a blocking Windows syscall, so this handle bridges Go's cancellation model
	// to the Windows API layer.
	cancelEvent           windows.Handle
	persister             operator.Persister
	publisherCache        publisherCache
	cancel                context.CancelFunc
	wg                    sync.WaitGroup
	subscription          Subscription
	maxEventsPerPollCycle int
	eventsReadInPollCycle int
	remote                RemoteConfig
	remoteSessionHandle   windows.Handle
	startRemoteSession    func() error
	processEvent          func(context.Context, Event) error
	ReqOrgAttr            *bool
}

// newInput creates a new Input operator.
func newInput(settings component.TelemetrySettings) *Input {
	basicConfig := helper.NewBasicConfig("windowseventlog", "input")
	basicOperator, _ := basicConfig.Build(settings)

	input := &Input{
		InputOperator: helper.InputOperator{
			WriterOperator: helper.WriterOperator{
				BasicOperator: basicOperator,
			},
		},
	}
	input.startRemoteSession = input.defaultStartRemoteSession
	return input
}

// defaultStartRemoteSession starts a remote session for reading event logs from a remote server.
func (i *Input) defaultStartRemoteSession() error {
	if i.remote.Server == "" {
		return nil
	}

	login := EvtRPCLogin{
		Server:   windows.StringToUTF16Ptr(i.remote.Server),
		User:     windows.StringToUTF16Ptr(i.remote.Username),
		Password: windows.StringToUTF16Ptr(i.remote.Password),
	}
	if i.remote.Domain != "" {
		login.Domain = windows.StringToUTF16Ptr(i.remote.Domain)
	}

	sessionHandle, err := evtOpenSession(EvtRPCLoginClass, &login, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to open session for server %s: %w", i.remote.Server, err)
	}
	i.remoteSessionHandle = sessionHandle
	return nil
}

// stopRemoteSession stops the remote session if it is active.
func (i *Input) stopRemoteSession() error {
	if i.remoteSessionHandle != 0 {
		if err := evtClose(uintptr(i.remoteSessionHandle)); err != nil {
			return fmt.Errorf("failed to close remote session handle for server %s: %w", i.remote.Server, err)
		}
		i.remoteSessionHandle = 0
	}
	return nil
}

// isRemote checks if the input is configured for remote access.
func (i *Input) isRemote() bool {
	return i.remote.Server != ""
}

// isNonTransientError checks if the error is likely non-transient.
func isNonTransientError(err error) bool {
	return errors.Is(err, windows.ERROR_EVT_CHANNEL_NOT_FOUND) || errors.Is(err, windows.ERROR_ACCESS_DENIED)
}

// Start will start reading events from a subscription.
func (i *Input) Start(persister operator.Persister) error {
	ctx, cancel := context.WithCancel(context.Background())
	i.cancel = cancel

	i.persister = persister

	if i.isRemote() {
		if err := i.startRemoteSession(); err != nil {
			return fmt.Errorf("failed to start remote session for server %s: %w", i.remote.Server, err)
		}
	}

	i.bookmark = NewBookmark()
	offsetXML, err := i.getBookmarkOffset(ctx)
	if err != nil {
		_ = i.persister.Delete(ctx, i.getPersistKey())
	}

	if offsetXML != "" {
		if err := i.bookmark.Open(offsetXML); err != nil {
			return fmt.Errorf("failed to open bookmark: %w", err)
		}
	}

	i.publisherCache = newPublisherCache()

	subscriptionError := false
	subscription := NewLocalSubscription(i.Logger())
	if i.isRemote() {
		subscription = NewRemoteSubscription(i.remote.Server, i.Logger())
	}

	if err := subscription.Open(i.startAt, uintptr(i.remoteSessionHandle), i.channel, i.query, i.bookmark); err != nil {
		var errorString string
		if isNonTransientError(err) {
			if i.isRemote() {
				errorString = fmt.Sprintf("failed to open subscription for remote server: %s", i.remote.Server)
			} else {
				errorString = "failed to open local subscription"
			}
			if !i.ignoreChannelErrors {
				return fmt.Errorf("%s, error: %w", errorString, err)
			}
			subscriptionError = true
			i.Logger().Warn(errorString, zap.Error(err))
		} else {
			if i.isRemote() {
				i.Logger().Warn("Transient error opening subscription for remote server, continuing", zap.String("server", i.remote.Server), zap.Error(err))
			} else {
				i.Logger().Warn("Transient error opening local subscription, continuing", zap.Error(err))
			}
		}
	}

	if !subscriptionError {
		i.subscription = subscription
		if metadata.StanzaWindowsEventDrivenScrapingFeatureGate.IsEnabled() {
			cancelEvent, err := windows.CreateEvent(nil, 1, 0, nil) // manual-reset, initially non-signaled
			if err != nil {
				return fmt.Errorf("failed to create cancel event: %w", err)
			}
			i.cancelEvent = cancelEvent
			i.wg.Add(1)
			go i.awaitAndReadEvents(ctx)
		} else {
			i.wg.Add(1)
			go i.pollAndRead(ctx)
		}
	}

	return nil
}

// Stop will stop reading events from a subscription.
func (i *Input) Stop() error {
	// Warning: all calls made below must be safe to be done even if Start() was not called or failed.

	if i.cancel != nil {
		i.cancel()
	}

	if i.cancelEvent != 0 {
		// If this fails, wg.Wait() below will block forever since awaitAndReadEvents will never
		// return from WaitForMultipleObjects. Log loudly and continue.
		if err := windows.SetEvent(i.cancelEvent); err != nil {
			i.Logger().Error("Failed to signal cancel event during stop; shutdown may hang", zap.Error(err))
		}
	}

	i.wg.Wait()

	var errs error
	if i.cancelEvent != 0 {
		if err := windows.CloseHandle(i.cancelEvent); err != nil {
			errs = multierr.Append(errs, fmt.Errorf("failed to close cancel event: %w", err))
		}
		i.cancelEvent = 0
	}

	if err := i.subscription.Close(); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to close subscription: %w", err))
	}

	if err := i.bookmark.Close(); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to close bookmark: %w", err))
	}

	if err := i.publisherCache.evictAll(); err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to close publishers: %w", err))
	}

	return multierr.Append(errs, i.stopRemoteSession())
}

func (i *Input) pollAndRead(ctx context.Context) {
	defer i.wg.Done()

	for {
		i.eventsReadInPollCycle = 0

		select {
		case <-ctx.Done():
			return
		case <-time.After(i.pollInterval):
			i.read(ctx)
		}
	}
}

func (i *Input) read(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !i.readBatch(ctx) {
				return
			}
		}
	}
}

// readWithRetry reads events from the subscription, handling RPC_S_INVALID_BOUND by closing and
// reopening the subscription with a halved batch size until a read succeeds or a non-retryable
// error occurs.
func (i *Input) readWithRetry(maxReads int) ([]Event, error) {
	events, err := i.subscription.Read(maxReads)
	if !errors.Is(err, windows.RPC_S_INVALID_BOUND) {
		return events, err
	}

	// Error is RPC_S_INVALID_BOUND. Close the subscription and reopen it with a halved batch size.
	if closeErr := i.subscription.Close(); closeErr != nil {
		return nil, fmt.Errorf("failed to close subscription during RPC_S_INVALID_BOUND recovery: %w", closeErr)
	}
	if openErr := i.subscription.Open(i.startAt, uintptr(i.remoteSessionHandle), i.channel, i.query, i.bookmark); openErr != nil {
		return nil, fmt.Errorf("failed to reopen subscription during RPC_S_INVALID_BOUND recovery: %w", openErr)
	}
	newMaxReads := max(maxReads/2, 1)
	i.currentMaxReads = newMaxReads
	i.Logger().Debug("Encountered RPC_S_INVALID_BOUND, reduced batch size", zap.Int("current_batch_size", i.currentMaxReads), zap.Int("original_batch_size", i.maxReads))
	return i.readWithRetry(newMaxReads)
}

// readBatch will read events from the subscription
func (i *Input) readBatch(ctx context.Context) bool {
	maxBatchSize := i.getCurrentBatchSize()
	if maxBatchSize <= 0 {
		return false
	}

	events, err := i.readWithRetry(maxBatchSize)
	if err != nil {
		i.Logger().Error("Failed to read events from subscription", zap.Error(err))
		if i.isRemote() && (errors.Is(err, windows.ERROR_INVALID_HANDLE) || errors.Is(err, errSubscriptionHandleNotOpen)) {
			i.Logger().Info("Resubscribing, closing remote subscription")
			closeErr := i.subscription.Close()
			if closeErr != nil {
				i.Logger().Error("Failed to close remote subscription", zap.Error(closeErr))
				return false
			}
			if err := i.stopRemoteSession(); err != nil {
				i.Logger().Error("Failed to close remote session", zap.Error(err))
			}
			i.Logger().Info("Resubscribing, creating remote subscription")
			i.subscription = NewRemoteSubscription(i.remote.Server, i.Logger())
			if err := i.startRemoteSession(); err != nil {
				i.Logger().Error("Failed to re-establish remote session", zap.String("server", i.remote.Server), zap.Error(err))
				return false
			}
			if err := i.subscription.Open(i.startAt, uintptr(i.remoteSessionHandle), i.channel, i.query, i.bookmark); err != nil {
				i.Logger().Error("Failed to re-open subscription for remote server", zap.String("server", i.remote.Server), zap.Error(err))
				return false
			}
		}
		return false
	}

	pstartedAtObj := time.Now().Local()
	pstartedAt := pstartedAtObj.Format("2006-01-02 15:04:05.000")

	for eI, event := range events {
		// Always use RenderSimple to get full EventXML (needed for processEventWithSimple)
		// This ensures we have EventData, RecordID, etc. for both raw and non-raw modes
		parsedSimpleEvent, err := event.RenderSimple(i.buffer)
		if err != nil {
			i.Logger().Error("Failed to render simple event", zap.Error(err))
			event.Close()
			continue
		}

		simpleEvent := parsedSimpleEvent.toEventXML()
		recordID := simpleEvent.RecordID
		providerName := simpleEvent.Provider.Name

		if recordID > 0 {
			dedupKey := fmt.Sprintf("dedup_%d", recordID)
			// Deduplication: check if this record was already processed
			exists, dCheckErr := i.persister.Get(ctx, dedupKey)
			if dCheckErr != nil {
				i.Logger().Error("Failed to check deduplication key", zap.Error(dCheckErr))
				event.Close()
				continue
			}
			if exists != nil {
				i.Logger().Debug("Duplicate event, skipping", zap.Uint64("record_id", recordID))
				event.Close()
				continue
			}
		}

		// Skip empty events
		if recordID == 0 && providerName == "" {
			i.Logger().Debug("Skipping empty event")
			event.Close()
			continue
		}

		// Always go through processEventWithSimple - it handles both raw and non-raw modes
		// and attempts to get the formatted event (with RenderingInfo.Message) from the publisher
		err1 := i.processEventWithSimple(ctx, event, simpleEvent)
		if err1 == nil {
			i.updateBookmarkOffset(ctx, event)
			// Persist processed RecordID for deduplication (only when we have a valid recordID)
			if recordID > 0 {
				dedupKey := fmt.Sprintf("dedup_%d", recordID)
				if err := i.persister.Set(ctx, dedupKey, []byte("1")); err != nil {
					i.Logger().Error("Failed to persist deduplication key", zap.Error(err))
				}
			}
		} else {
			i.Logger().Error("Failed to process event", zap.Any("pstartedAt", pstartedAt), zap.Int("current loop id:", eI), zap.Any("process started at ", time.Now().Local().Format("2006-01-02 15:04:05.000")), zap.Uint64("record_id", recordID), zap.Error(err1))
		}
		event.Close()
	}

	i.eventsReadInPollCycle += len(events)
	return len(events) != 0
}

// awaitAndReadEvents is the event-driven alternative to pollAndRead. Instead of sleeping
// for a fixed interval it blocks on a Windows wait object that is signaled by the subscription
// when new events arrive. This reduces latency and avoids unnecessary wakeups.
func (i *Input) awaitAndReadEvents(ctx context.Context) {
	defer i.wg.Done()

	timeoutMs := uint32(i.waitTimeout.Milliseconds())
	for {
		ready, err := i.subscription.Wait(i.cancelEvent, timeoutMs)
		if err != nil {
			i.Logger().Error("Failed to wait for subscription signal", zap.Error(err))
			return
		}
		if !ready {
			// cancel event was signaled
			return
		}

		i.eventsReadInPollCycle = 0
		i.read(ctx)
	}
}

func (i *Input) getPublisherName(event Event) (name string, excluded bool) {
	providerName, err := event.GetPublisherName(i.buffer)
	if err != nil {
		i.Logger().Error("Failed to get provider name", zap.Error(err))
		return "", true
	}
	if _, exclude := i.excludeProviders[providerName]; exclude {
		return "", true
	}

	return providerName, false
}

func (i *Input) renderSimpleAndSend(ctx context.Context, event Event) error {
	render := event.RenderSimple
	if i.raw {
		render = event.RenderSimpleRaw
	}
	simpleEvent, err := render(i.buffer)
	if err != nil {
		return fmt.Errorf("render simple event: %w", err)
	}
	return i.sendEvent(ctx, simpleEvent)
}

func (i *Input) renderDeepAndSend(ctx context.Context, event Event, publisher Publisher) error {
	render := event.RenderDeep
	if i.raw {
		render = event.RenderDeepRaw
	}
	deepEvent, err := render(i.buffer, publisher)
	if err == nil {
		return i.sendEvent(ctx, deepEvent)
	}
	return multierr.Append(
		fmt.Errorf("render deep event: %w", err),
		i.renderSimpleAndSend(ctx, event),
	)
}

// processEvent will process and send an event retrieved from windows event log.
func (i *Input) processEventWithoutRenderingInfo(ctx context.Context, event Event) error {
	if len(i.excludeProviders) == 0 {
		return i.renderSimpleAndSend(ctx, event)
	}
	if _, exclude := i.getPublisherName(event); exclude {
		return nil
	}
	return i.renderSimpleAndSend(ctx, event)
}

func (i *Input) processEventWithRenderingInfo(ctx context.Context, event Event) error {
	providerName, exclude := i.getPublisherName(event)
	if exclude {
		return nil
	}

	publisher, err := i.publisherCache.get(providerName)
	if err != nil {
		return multierr.Append(
			fmt.Errorf("open event source for provider %q: %w", providerName, err),
			i.renderSimpleAndSend(ctx, event),
		)
	}

	if publisher.Valid() {
		return i.renderDeepAndSend(ctx, event, publisher)
	}
	return i.renderSimpleAndSend(ctx, event)
}

// sendEvent will send a parsedEvent as an entry to the operator's output.
//
// raw=true path: All getXxx() methods in parsedEvent interface are called.
// If you add a field access here, add a corresponding method to parsedEvent
// and rawEventXML.
func (i *Input) sendEvent(ctx context.Context, event parsedEvent) error {
	var body any = event.getOriginal()
	if !i.raw {
		body = formattedBody(event.toEventXML(), i.eventDataFormat)
	}

	e, err := i.NewEntry(body)
	if err != nil {
		return fmt.Errorf("create entry: %w", err)
	}

	e.Timestamp = parseTimestamp(event.getSystemTime())
	e.Severity = parseSeverity(event.getRenderedLevel(), event.getLevel())

	// === Provider attributes ===
	if providerName := event.getProviderName(); providerName != "" {
		e.AddAttribute("event_name", providerName)
		e.AddAttribute("provider_name", providerName)
	}
	if eventSourceName := event.getEventSourceName(); eventSourceName != "" {
		e.AddAttribute("event_source_name", eventSourceName)
	}

	// === EventID attributes ===
	if eventID := event.getEventID(); eventID != 0 {
		e.AddAttribute("event_identifier", fmt.Sprintf("%d", eventID))
		e.AddAttribute("event_id", fmt.Sprintf("%d", eventID))
	}
	if qualifiers := event.getEventIDQualifiers(); qualifiers != 0 {
		e.AddAttribute("event_id_qualifiers", fmt.Sprintf("%d", qualifiers))
	}

	// === Channel ===
	if channel := event.getChannel(); channel != "" {
		e.AddAttribute("event_channel", channel)
	}

	// === Record and Version ===
	if recordID := event.getRecordID(); recordID != 0 {
		e.AddAttribute("record_id", fmt.Sprintf("%d", recordID))
	}
	if version := event.getVersion(); version != 0 {
		e.AddAttribute("version", fmt.Sprintf("%d", version))
	}

	// === Level: source_level is raw numeric, level is text ===
	if level := event.getLevel(); level != "" {
		e.AddAttribute("source_level", level)     // Raw numeric: "0", "1", "2", "3", "4", "5"
		e.AddAttribute("level", levelName(level)) // Text: "Information", "Error", "Warning", "Critical", "Verbose"
	}
	if renderedLevel := event.getRenderedLevel(); renderedLevel != "" {
		e.AddAttribute("rendered_level", renderedLevel)
	}

	// === Task and Opcode ===
	if task := event.getTask(); task != "" {
		e.AddAttribute("task", task)
	}
	if renderedTask := event.getRenderedTask(); renderedTask != "" {
		e.AddAttribute("rendered_task", renderedTask)
	}
	if opcode := event.getOpcode(); opcode != "" {
		e.AddAttribute("opcode", opcode)
	}
	if renderedOpcode := event.getRenderedOpcode(); renderedOpcode != "" {
		e.AddAttribute("rendered_opcode", renderedOpcode)
	}

	// === Rendered Keywords (human-readable keywords from RenderingInfo) ===
	if renderedKeywords := event.getRenderedKeywords(); len(renderedKeywords) > 0 {
		e.AddAttribute("rendered_keywords", strings.Join(renderedKeywords, ","))
	}

	// === Security: User ID (SID) ===
	if userID := event.getUserID(); userID != "" {
		e.AddAttribute("user_id", userID)
	}

	// === Execution: Process ID only ===
	if processID := event.getProcessID(); processID != 0 {
		e.AddAttribute("process_id", fmt.Sprintf("%d", processID))
	}

	// === Correlation: Activity IDs ===
	if activityID := event.getActivityID(); activityID != "" {
		e.AddAttribute("activity_id", activityID)
	}
	if relatedActivityID := event.getRelatedActivityID(); relatedActivityID != "" {
		e.AddAttribute("related_activity_id", relatedActivityID)
	}

	// === Rendered Message ===
	if message := event.getRenderedMessage(); message != "" {
		e.AddAttribute("rendered_message", message)
	}

	// === Binary Event Data ===
	if binaryData := event.getBinaryEventData(); binaryData != "" {
		e.AddAttribute("binary_event_data", binaryData)
	}

	// === EventData fields as individual attributes ===
	eventData, _ := i.ExtractEventData(event.getEventData())
	if len(eventData) > 0 {
		// Add event_data as JSON-encoded structured attribute
		if eventDataJSON, err := json.Marshal(eventData); err == nil {
			e.AddAttribute("event_data", string(eventDataJSON))
		}
		// Also add individual fields as flat attributes
		for eK, eD := range eventData {
			eK = strings.ReplaceAll(strings.ToLower(eK), " ", "_")
			if str, ok := eD.(string); ok {
				e.AddAttribute(eK, str)
			} else {
				e.AddAttribute(eK, fmt.Sprintf("%v", eD))
			}
		}
	}

	// === UserData fields as individual attributes ===
	if userData := event.getUserData(); userData != nil {
		if userData.Name != "" {
			e.AddAttribute("userdata_name", userData.Name)
		}
		for k, v := range userData.Data {
			attrKey := "userdata_" + strings.ReplaceAll(strings.ToLower(k), " ", "_")
			e.AddAttribute(attrKey, v)
		}
	}

	if i.remote.Server != "" {
		e.AddAttribute("server.address", i.remote.Server)
	}

	if i.isAdditionalAttrReq() {
		e.AddAttribute("log_record_original", event.getOriginal())
	}

	return i.Write(ctx, e)
}

// getBookmarkXML will get the bookmark xml from the offsets database.
func (i *Input) getBookmarkOffset(ctx context.Context) (string, error) {
	bytes, err := i.persister.Get(ctx, i.getPersistKey())
	return string(bytes), err
}

// updateBookmark will update the bookmark xml and save it in the offsets database.
func (i *Input) updateBookmarkOffset(ctx context.Context, event Event) {
	if err := i.bookmark.Update(event); err != nil {
		i.Logger().Error("Failed to update bookmark from event", zap.Error(err))
		return
	}

	bookmarkXML, err := i.bookmark.Render(i.buffer)
	if err != nil {
		i.Logger().Error("Failed to render bookmark xml", zap.Error(err))
		return
	}

	if err := i.persister.Set(ctx, i.getPersistKey(), []byte(bookmarkXML)); err != nil {
		i.Logger().Error("failed to set offsets", zap.Error(err))
		return
	}
}

func (i *Input) getPersistKey() string {
	if i.query != nil {
		return *i.query
	}

	return i.channel
}

func (i *Input) getCurrentBatchSize() int {
	if i.maxEventsPerPollCycle == 0 {
		return i.currentMaxReads
	}

	return min(i.currentMaxReads, i.maxEventsPerPollCycle-i.eventsReadInPollCycle)
}

func (i *Input) isAdditionalAttrReq() bool {
	if i.ReqOrgAttr == nil {
		return false
	}
	return *i.ReqOrgAttr
}

func (i *Input) ExtractEventData(eventData EventData) (map[string]any, error) {
	result := make(map[string]any)

	// Include EventData name if present
	if eventData.Name != "" {
		result["name"] = eventData.Name
	}

	// Include binary data if present
	if eventData.Binary != "" {
		result["binary"] = eventData.Binary
	}

	// Collect anonymous data values into message
	var messageValues []string

	for _, data := range eventData.Data {
		// Handle anonymous data (no Name attribute) - aggregate into message
		if data.Name == "" {
			if data.Value != "" {
				messageValues = append(messageValues, data.Value)
			}
			continue
		}

		// Skip unwanted fields
		switch data.Name {
		case "CountOfCredentialsReturned", "ProcessCreationTime", "ReadOperation", "SubjectDomainName":
			// Skip these fields
			continue
		}

		// Extract known named fields with normalized keys
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
		case "TargetUserName", "SubjectUserName":
			result["user_name"] = data.Value
		case "TargetDomainName", "WorkstationName":
			result["domain_name"] = data.Value
		default:
			// Include all other named fields with normalized key (lowercase, spaces to underscores)
			key := strings.ReplaceAll(strings.ToLower(data.Name), " ", "_")
			result[key] = data.Value
		}
	}

	// Set message from anonymous data
	if len(messageValues) > 0 {
		result["message"] = strings.Join(messageValues, "\n")
	}

	return result, nil
}

func (i *Input) processEventWithSimple(ctx context.Context, event Event, simpleEvent *EventXML) error {
	isExcluded := func(providerName string) bool {
		for excludeProvider := range i.excludeProviders {
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

	// For both raw and non-raw modes, try to get formatted event with RenderingInfo
	// This ensures we get the localized message from Windows Event API
	publisher, err := i.publisherCache.get(simpleEvent.Provider.Name)
	if err != nil || !publisher.Valid() {
		i.Logger().Debug("Publisher unavailable, using simple event",
			zap.String("provider", simpleEvent.Provider.Name),
			zap.String("channel", simpleEvent.Channel),
			zap.Error(err))
		if i.raw {
			return i.sendEventXMLAsRaw(ctx, simpleEvent)
		}
		return i.sendEvent(ctx, simpleEvent)
	}

	formattedEvent, err := event.RenderFormatted(i.buffer, publisher)
	if err != nil {
		i.Logger().Debug("RenderFormatted failed, using simple event",
			zap.String("provider", simpleEvent.Provider.Name),
			zap.String("channel", simpleEvent.Channel),
			zap.Error(err))
		if i.raw {
			return i.sendEventXMLAsRaw(ctx, simpleEvent)
		}
		return i.sendEvent(ctx, simpleEvent)
	}

	if i.raw {
		return i.sendEventXMLAsRaw(ctx, formattedEvent)
	}
	return i.sendEvent(ctx, formattedEvent)
}

// sendEventXMLAsRaw sends EventXML in raw mode - body is the original XML,
// but attributes are extracted from the parsed EventXML (including RenderingInfo.Message)
func (i *Input) sendEventXMLAsRaw(ctx context.Context, eventXML *EventXML) error {
	// Body is the original XML string
	body := eventXML.Original
	e, err := i.NewEntry(body)
	if err != nil {
		i.Logger().Error("Failed to create entry", zap.Error(err))
		return err
	}

	// Get rendered values from RenderingInfo if available
	var renderedLevel, renderedTask, renderedOpcode, renderedMessage string
	var renderedKeywords []string
	if eventXML.RenderingInfo != nil {
		renderedLevel = eventXML.RenderingInfo.Level
		renderedTask = eventXML.RenderingInfo.Task
		renderedOpcode = eventXML.RenderingInfo.Opcode
		renderedMessage = eventXML.RenderingInfo.Message
		renderedKeywords = eventXML.RenderingInfo.Keywords
	}

	e.Timestamp = parseTimestamp(eventXML.TimeCreated.SystemTime)
	e.Severity = parseSeverity(renderedLevel, eventXML.Level)

	// === System Time ===
	if systemTime := eventXML.TimeCreated.SystemTime; systemTime != "" {
		e.AddAttribute("system_time", systemTime)
	}

	// === Provider attributes ===
	if providerName := eventXML.Provider.Name; providerName != "" {
		e.AddAttribute("provider_name", providerName)
	}

	// === EventID ===
	if eventID := eventXML.EventID.ID; eventID != 0 {
		e.AddAttribute("event_id", fmt.Sprintf("%d", eventID))
	}

	// === Channel ===
	if channel := eventXML.Channel; channel != "" {
		e.AddAttribute("channel", channel)
	}

	// === Computer ===
	if computer := eventXML.Computer; computer != "" {
		e.AddAttribute("computer", computer)
	}

	// === Record ID ===
	if recordID := eventXML.RecordID; recordID != 0 {
		e.AddAttribute("record_id", fmt.Sprintf("%d", recordID))
	}

	// === Level: prefer rendered (localized), fallback to numeric ===
	if renderedLevel != "" {
		e.AddAttribute("level", renderedLevel)
	} else if level := eventXML.Level; level != "" {
		e.AddAttribute("source_level", level)
		e.AddAttribute("level", levelName(level))
	}

	// === Task: prefer rendered (localized), fallback to numeric ===
	if renderedTask != "" {
		e.AddAttribute("task", renderedTask)
	} else if task := eventXML.Task; task != "" {
		e.AddAttribute("task", task)
	}

	// === Opcode: prefer rendered (localized), fallback to numeric ===
	if renderedOpcode != "" {
		e.AddAttribute("opcode", renderedOpcode)
	} else if opcode := eventXML.Opcode; opcode != "" {
		e.AddAttribute("opcode", opcode)
	}

	// === Keywords ===
	if len(renderedKeywords) > 0 {
		e.AddAttribute("keywords", strings.Join(renderedKeywords, ","))
	}

	// === Security: User ID (SID) ===
	if eventXML.Security != nil && eventXML.Security.UserID != "" {
		e.AddAttribute("user_id", eventXML.Security.UserID)
	}

	// === Message (from RenderingInfo - localized by Windows via EvtFormatMessage) ===
	if renderedMessage != "" {
		e.AddAttribute("message", renderedMessage)
	}

	// === EventData fields as individual attributes ===
	eventData, _ := i.ExtractEventData(eventXML.EventData)
	for eK, eD := range eventData {
		eK = strings.ReplaceAll(strings.ToLower(eK), " ", "_")
		if str, ok := eD.(string); ok {
			e.AddAttribute(eK, str)
		} else {
			e.AddAttribute(eK, fmt.Sprintf("%v", eD))
		}
	}

	if i.remote.Server != "" {
		e.AddAttribute("server.address", i.remote.Server)
	}

	if i.isAdditionalAttrReq() {
		e.AddAttribute("log_record_original", eventXML.Original)
	}

	i.Write(ctx, e)
	return nil
}

func (i *Input) sendEventRaw(ctx context.Context, eventRaw EventRaw) {
	body := eventRaw.parseBody()
	e, err := i.NewEntry(body)
	if err != nil {
		i.Logger().Error("Failed to create entry", zap.Error(err))
		return
	}

	e.Timestamp = eventRaw.ParseTimestamp()
	e.Severity = eventRaw.ParseRenderedSeverity()

	// === System Time ===
	if systemTime := eventRaw.TimeCreated.SystemTime; systemTime != "" {
		e.AddAttribute("system_time", systemTime)
	}

	// === Provider attributes ===
	if providerName := eventRaw.GetProviderName(); providerName != "" {
		e.AddAttribute("provider_name", providerName)
	}

	// === EventID ===
	if eventID := eventRaw.GetEventID(); eventID != 0 {
		e.AddAttribute("event_id", fmt.Sprintf("%d", eventID))
	}

	// === Channel ===
	if channel := eventRaw.GetChannel(); channel != "" {
		e.AddAttribute("channel", channel)
	}

	// === Computer ===
	if computer := eventRaw.GetComputer(); computer != "" {
		e.AddAttribute("computer", computer)
	}

	// === Record ID ===
	if recordID := eventRaw.GetRecordID(); recordID != 0 {
		e.AddAttribute("record_id", fmt.Sprintf("%d", recordID))
	}

	// === Level: source_level is raw numeric, level is text ===
	if level := eventRaw.Level; level != "" {
		e.AddAttribute("source_level", level)     // Raw numeric: "0", "1", "2", "3", "4", "5"
		e.AddAttribute("level", levelName(level)) // Text: "Information", "Error", "Warning", "Critical", "Verbose"
	}

	// === Task: prefer rendered text over numeric ===
	if renderedTask := eventRaw.GetRenderedTask(); renderedTask != "" {
		e.AddAttribute("task", renderedTask)
	} else if task := eventRaw.GetTask(); task != "" {
		e.AddAttribute("task", task)
	}

	// === Opcode: prefer rendered text over numeric ===
	if renderedOpcode := eventRaw.GetRenderedOpcode(); renderedOpcode != "" {
		e.AddAttribute("opcode", renderedOpcode)
	} else if opcode := eventRaw.GetOpcode(); opcode != "" {
		e.AddAttribute("opcode", opcode)
	}

	// === Keywords ===
	if renderedKeywords := eventRaw.GetRenderedKeywords(); len(renderedKeywords) > 0 {
		e.AddAttribute("keywords", strings.Join(renderedKeywords, ","))
	}

	// === Security: User ID (SID) ===
	if userID := eventRaw.GetUserID(); userID != "" {
		e.AddAttribute("user_id", userID)
	}

	// === Message ===
	if message := eventRaw.GetRenderedMessage(); message != "" {
		e.AddAttribute("message", message)
	}

	// === EventData fields as individual attributes ===
	eventData, _ := i.ExtractEventData(eventRaw.GetEventData())
	for eK, eD := range eventData {
		eK = strings.ReplaceAll(strings.ToLower(eK), " ", "_")
		if str, ok := eD.(string); ok {
			e.AddAttribute(eK, str)
		} else {
			e.AddAttribute(eK, fmt.Sprintf("%v", eD))
		}
	}

	if i.remote.Server != "" {
		e.AddAttribute("server.address", i.remote.Server)
	}

	i.Write(ctx, e)
}
