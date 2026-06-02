// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package windows // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/input/windows"

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/entry"
)

// EventRaw is the rendered xml of an event with all fields parsed for attribute extraction.
// The body is the original XML of the entire event.
type EventRaw struct {
	// System fields
	TimeCreated TimeCreated  `xml:"System>TimeCreated"`
	Level       string       `xml:"System>Level"`
	Provider    Provider     `xml:"System>Provider"`
	EventID     EventID      `xml:"System>EventID"`
	Channel     string       `xml:"System>Channel"`
	Computer    string       `xml:"System>Computer"`
	RecordID    uint64       `xml:"System>EventRecordID"`
	Task        string       `xml:"System>Task"`
	Opcode      string       `xml:"System>Opcode"`
	Keywords    string       `xml:"System>Keywords"`
	Version     uint8        `xml:"System>Version"`
	Security    *Security    `xml:"System>Security"`
	Execution   *Execution   `xml:"System>Execution"`
	Correlation *Correlation `xml:"System>Correlation"`

	// Data fields
	EventData       EventData `xml:"EventData"`
	UserData        *UserData `xml:"UserData"`
	BinaryEventData string    `xml:"BinaryEventData"`

	// RenderingInfo (populated when using RenderFormattedRaw with publisher)
	RenderingInfo *RenderingInfoRaw `xml:"RenderingInfo"`

	// Body stores the original XML
	Body string `xml:"-"`
}

// RenderingInfoRaw holds rendered fields from RenderingInfo for raw mode.
type RenderingInfoRaw struct {
	Level    string   `xml:"Level"`
	Message  string   `xml:"Message"`
	Task     string   `xml:"Task"`
	Opcode   string   `xml:"Opcode"`
	Keywords []string `xml:"Keywords>Keyword"`
}

// parseTimestamp will parse the timestamp of the event.
func (e *EventRaw) ParseTimestamp() time.Time {
	if timestamp, err := time.Parse(time.RFC3339Nano, e.TimeCreated.SystemTime); err == nil {
		return timestamp
	}
	return time.Now()
}

// ParseRenderedSeverity returns the severity of the event.
// Prefers numeric level (more reliable across localized Windows installations)
// over rendered text, falling back to text only when numeric is unavailable.
func (e *EventRaw) ParseRenderedSeverity() entry.Severity {
	// Prefer numeric level first (more reliable than localized text)
	switch e.Level {
	case "1":
		return entry.Fatal
	case "2":
		return entry.Error
	case "3":
		return entry.Warn
	case "0", "4":
		return entry.Info
	case "5":
		return entry.Debug // Verbose
	}
	// Fallback to rendered text when numeric level is empty/unknown
	renderedLevel := e.GetRenderedLevel()
	switch renderedLevel {
	case "Critical":
		return entry.Fatal
	case "Error":
		return entry.Error
	case "Warning":
		return entry.Warn
	case "Information":
		return entry.Info
	}
	return entry.Default
}

// parseSeverity will parse the severity of the event when RenderingInfo is not populated
func (e *EventRaw) ParseSeverity() entry.Severity {
	switch e.Level {
	case "0":
		return entry.Info // LogAlways - used by Security audit events
	case "1":
		return entry.Fatal
	case "2":
		return entry.Error
	case "3":
		return entry.Warn
	case "4":
		return entry.Info
	case "5":
		return entry.Debug // Verbose
	default:
		return entry.Default
	}
}

func (e *EventRaw) parseBody() string {
	return e.Body
}

// === Getter methods for attribute extraction ===

func (e *EventRaw) GetProviderName() string      { return e.Provider.Name }
func (e *EventRaw) GetProviderGUID() string      { return e.Provider.GUID }
func (e *EventRaw) GetEventSourceName() string   { return e.Provider.EventSourceName }
func (e *EventRaw) GetEventID() uint32           { return e.EventID.ID }
func (e *EventRaw) GetEventIDQualifiers() uint16 { return e.EventID.Qualifiers }
func (e *EventRaw) GetChannel() string           { return e.Channel }
func (e *EventRaw) GetComputer() string          { return e.Computer }
func (e *EventRaw) GetRecordID() uint64          { return e.RecordID }
func (e *EventRaw) GetTask() string              { return e.Task }
func (e *EventRaw) GetOpcode() string            { return e.Opcode }
func (e *EventRaw) GetKeywords() string          { return e.Keywords }
func (e *EventRaw) GetVersion() uint8            { return e.Version }
func (e *EventRaw) GetBinaryEventData() string   { return e.BinaryEventData }
func (e *EventRaw) GetEventData() EventData      { return e.EventData }
func (e *EventRaw) GetUserData() *UserData       { return e.UserData }

func (e *EventRaw) GetRenderedLevel() string {
	if e.RenderingInfo != nil {
		return e.RenderingInfo.Level
	}
	return ""
}

func (e *EventRaw) GetRenderedMessage() string {
	// First try RenderingInfo.Message (available when RenderFormattedRaw succeeds)
	if e.RenderingInfo != nil && e.RenderingInfo.Message != "" {
		return e.RenderingInfo.Message
	}

	// Fallback: construct message from EventData fields (locale-neutral)
	// This works for Security/Setup channels where EvtFormatMessage may fail due to permissions
	if len(e.EventData.Data) > 0 {
		var parts []string
		for _, d := range e.EventData.Data {
			if d.Name != "" && d.Value != "" {
				parts = append(parts, d.Name+": "+d.Value)
			}
		}
		if len(parts) > 0 {
			return strings.Join(parts, ", ")
		}
	}
	return ""
}

func (e *EventRaw) GetRenderedTask() string {
	if e.RenderingInfo != nil {
		return e.RenderingInfo.Task
	}
	return ""
}

func (e *EventRaw) GetRenderedOpcode() string {
	if e.RenderingInfo != nil {
		return e.RenderingInfo.Opcode
	}
	return ""
}

func (e *EventRaw) GetRenderedKeywords() []string {
	if e.RenderingInfo != nil {
		return e.RenderingInfo.Keywords
	}
	return nil
}

func (e *EventRaw) GetUserID() string {
	if e.Security != nil {
		return e.Security.UserID
	}
	return ""
}

func (e *EventRaw) GetProcessID() uint {
	if e.Execution != nil {
		return e.Execution.ProcessID
	}
	return 0
}

func (e *EventRaw) GetThreadID() uint {
	if e.Execution != nil {
		return e.Execution.ThreadID
	}
	return 0
}

func (e *EventRaw) GetActivityID() string {
	if e.Correlation != nil && e.Correlation.ActivityID != nil {
		return *e.Correlation.ActivityID
	}
	return ""
}

func (e *EventRaw) GetRelatedActivityID() string {
	if e.Correlation != nil && e.Correlation.RelatedActivityID != nil {
		return *e.Correlation.RelatedActivityID
	}
	return ""
}

// unmarshalEventRaw will unmarshal EventRaw from xml bytes.
// Uses sanitizeXMLBytes to strip the default namespace so that Go's xml package
// can correctly match XML paths like "System>Level".
func UnmarshalEventRaw(bytes []byte) (EventRaw, error) {
	// Sanitize XML to strip namespace and illegal characters
	sanitized := sanitizeXMLBytes(bytes)

	var eventRaw EventRaw
	if err := xml.Unmarshal(sanitized, &eventRaw); err != nil {
		return EventRaw{}, fmt.Errorf("failed to unmarshal xml bytes into event: %w (%s)", err, string(bytes))
	}
	// Keep original XML as the body
	eventRaw.Body = string(bytes)
	return eventRaw, nil
}
