// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package windows // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/input/windows"

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/entry"
)

// EventXML is the rendered xml of an event.
// See: https://learn.microsoft.com/en-us/windows/win32/wes/eventschema-schema
type EventXML struct {
	Original            string               `xml:"-"`
	EventID             EventID              `xml:"System>EventID"`
	Provider            Provider             `xml:"System>Provider"`
	Computer            string               `xml:"System>Computer"`
	Channel             string               `xml:"System>Channel"`
	RecordID            uint64               `xml:"System>EventRecordID"`
	TimeCreated         TimeCreated          `xml:"System>TimeCreated"`
	Level               string               `xml:"System>Level"`
	Task                string               `xml:"System>Task"`
	Opcode              string               `xml:"System>Opcode"`
	Keywords            string               `xml:"System>Keywords"`
	Security            *Security            `xml:"System>Security"`
	Execution           *Execution           `xml:"System>Execution"`
	EventData           EventData            `xml:"EventData"`
	UserData            *UserData            `xml:"UserData"`
	Correlation         *Correlation         `xml:"System>Correlation"`
	Version             uint8                `xml:"System>Version"`
	RenderingInfo       *RenderingInfo       `xml:"RenderingInfo"`
	ProcessingErrorData *ProcessingErrorData `xml:"ProcessingErrorData"`
	DebugData           *DebugData           `xml:"DebugData"`
	// BinaryEventData contains raw hex-encoded binary data logged by legacy providers.
	// https://learn.microsoft.com/en-us/windows/win32/wes/eventschema-binaryeventdata-eventtype-element
	BinaryEventData string `xml:"BinaryEventData"`
}

// parseTimestamp will parse the timestamp of the event.
func parseTimestamp(ts string) time.Time {
	if timestamp, err := time.Parse(time.RFC3339Nano, ts); err == nil {
		return timestamp
	}
	return time.Now()
}

// parseSeverity will parse the severity of the event.
// Prefers numeric level (more reliable across localized Windows installations)
// over rendered text, falling back to text when numeric is unavailable.
func parseSeverity(renderedLevel, level string) entry.Severity {
	// Prefer numeric level first (more reliable than localized text)
	switch level {
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
	}
	// Fallback to rendered text when numeric level is empty/unknown
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

// levelName maps a Windows ETW numeric level string to its canonical English name.
// Level 0 (LogAlways) returns "" because the Security channel and some other
// providers use keywords (Audit Success/Failure) instead of a severity level.
// Non-numeric values (e.g. already-rendered text) are returned unchanged.
func levelName(level string) string {
	switch level {
	case "1":
		return "Critical"
	case "2":
		return "Error"
	case "3":
		return "Warning"
	case "0", "4":
		return "Information"
	case "5":
		return "Verbose"
	}
	return level
}

// formattedBody will parse a body from the event.
func formattedBody(e *EventXML, eventDataFormat EventDataFormat) map[string]any {
	var rawMessage string
	level := levelName(e.Level)
	task := e.Task
	opcode := e.Opcode
	keywords := e.Keywords // System Keywords is a hex string
	var renderedKeywords []string

	if e.RenderingInfo != nil {
		rawMessage = e.RenderingInfo.Message
		if e.RenderingInfo.Level != "" {
			level = e.RenderingInfo.Level
		}
		if e.RenderingInfo.Task != "" {
			task = e.RenderingInfo.Task
		}
		if e.RenderingInfo.Opcode != "" {
			opcode = e.RenderingInfo.Opcode
		}
		if e.RenderingInfo.Keywords != nil {
			renderedKeywords = e.RenderingInfo.Keywords
		}
	}

	message, details := parseMessage(e.Channel, rawMessage)

	// When the publisher DLL is unavailable (simple render path), RenderingInfo.Message
	// is empty. Fall back to serializing named EventData fields (e.g. TargetUserName,
	// LogonType) as "Key: Value" pairs so the body message is never completely empty.
	if message == "" && len(e.EventData.Data) > 0 {
		var parts []string
		for _, d := range e.EventData.Data {
			if d.Name != "" && d.Value != "" {
				parts = append(parts, d.Name+": "+d.Value)
			}
		}
		message = strings.Join(parts, ", ")
	}

	// Use rendered keywords if available, otherwise use hex string
	var keywordsValue any = keywords
	if len(renderedKeywords) > 0 {
		keywordsValue = renderedKeywords
	}

	body := map[string]any{
		"event_id": map[string]any{
			"qualifiers": e.EventID.Qualifiers,
			"id":         e.EventID.ID,
		},
		"provider": map[string]any{
			"name":         e.Provider.Name,
			"guid":         e.Provider.GUID,
			"event_source": e.Provider.EventSourceName,
		},
		"system_time": e.TimeCreated.SystemTime,
		"computer":    e.Computer,
		"channel":     e.Channel,
		"record_id":   e.RecordID,
		"level":       level,
		"message":     message,
		"task":        task,
		"opcode":      opcode,
		"keywords":    keywordsValue,
		"event_data":  parseEventData(e.EventData, eventDataFormat),
		"version":     e.Version,
	}

	if len(details) > 0 {
		body["details"] = details
	}

	if e.Security != nil && e.Security.UserID != "" {
		body["security"] = map[string]any{
			"user_id": e.Security.UserID,
		}
	}

	if e.Execution != nil {
		body["execution"] = e.Execution.asMap()
	}

	if e.Correlation != nil {
		body["correlation"] = e.Correlation.asMap()
	}

	if e.RenderingInfo != nil {
		body["rendering_info"] = e.RenderingInfo.asMap()
	}

	if e.UserData != nil {
		body["user_data"] = e.UserData.asMap()
	}

	if e.ProcessingErrorData != nil {
		body["processing_error_data"] = e.ProcessingErrorData.asMap()
	}

	if e.DebugData != nil {
		body["debug_data"] = e.DebugData.asMap()
	}

	if e.BinaryEventData != "" {
		body["binary_event_data"] = e.BinaryEventData
	}

	return body
}

// parseMessage will attempt to parse a message into a message and details
func parseMessage(channel, message string) (string, map[string]any) {
	switch channel {
	case "Security":
		return parseSecurity(message)
	default:
		return message, nil
	}
}

// parseEventData converts EventData XML elements into a map.
// When format is EventDataFormatMap, named Data elements become direct keys and
// anonymous elements use numbered keys (param1, param2, …).
// When format is EventDataFormatArray, data is stored as a "data" slice of
// single-key maps, preserving the original collector format.
// see: https://learn.microsoft.com/en-us/windows/win32/wes/eventschema-datafieldtype-complextype
func parseEventData(eventData EventData, format EventDataFormat) map[string]any {
	outputMap := make(map[string]any, len(eventData.Data)+2)

	if eventData.Name != "" {
		outputMap["name"] = eventData.Name
	}
	if eventData.Binary != "" {
		outputMap["binary"] = eventData.Binary
	}

	if len(eventData.Data) == 0 {
		return outputMap
	}

	switch format {
	case EventDataFormatArray:
		dataMaps := make([]any, len(eventData.Data))
		for i, data := range eventData.Data {
			dataMaps[i] = map[string]any{
				data.Name: data.Value,
			}
		}
		outputMap["data"] = dataMaps
	default:
		anonymousCounter := 1
		for _, data := range eventData.Data {
			if data.Name != "" {
				outputMap[data.Name] = data.Value
			} else {
				key := fmt.Sprintf("param%d", anonymousCounter)
				outputMap[key] = data.Value
				anonymousCounter++
			}
		}
	}

	return outputMap
}

// EventID is the identifier of the event.
type EventID struct {
	Qualifiers uint16 `xml:"Qualifiers,attr"`
	ID         uint32 `xml:",chardata"`
}

// TimeCreated is the creation time of the event.
type TimeCreated struct {
	SystemTime string `xml:"SystemTime,attr"`
}

// Provider is the provider of the event.
type Provider struct {
	Name            string `xml:"Name,attr"`
	GUID            string `xml:"Guid,attr"`
	EventSourceName string `xml:"EventSourceName,attr"`
}

type EventData struct {
	// https://learn.microsoft.com/en-us/windows/win32/wes/eventschema-eventdatatype-complextype
	// ComplexData is not supported.
	Name   string `xml:"Name,attr"`
	Data   []Data `xml:"Data"`
	Binary string `xml:"Binary"`
}

type Data struct {
	// https://learn.microsoft.com/en-us/windows/win32/wes/eventschema-datafieldtype-complextype
	Name  string `xml:"Name,attr"`
	Value string `xml:",chardata"`
}

// Security contains info pertaining to the user triggering the event.
type Security struct {
	UserID string `xml:"UserID,attr"`
}

// Execution contains info pertaining to the process that triggered the event.
type Execution struct {
	// ProcessID and ThreadID are required on execution info
	ProcessID uint `xml:"ProcessID,attr"`
	ThreadID  uint `xml:"ThreadID,attr"`
	// These remaining fields are all optional for execution info
	ProcessorID   *uint `xml:"ProcessorID,attr"`
	SessionID     *uint `xml:"SessionID,attr"`
	KernelTime    *uint `xml:"KernelTime,attr"`
	UserTime      *uint `xml:"UserTime,attr"`
	ProcessorTime *uint `xml:"ProcessorTime,attr"`
}

func (e Execution) asMap() map[string]any {
	result := map[string]any{
		"process_id": e.ProcessID,
		"thread_id":  e.ThreadID,
	}

	if e.ProcessorID != nil {
		result["processor_id"] = *e.ProcessorID
	}

	if e.SessionID != nil {
		result["session_id"] = *e.SessionID
	}

	if e.KernelTime != nil {
		result["kernel_time"] = *e.KernelTime
	}

	if e.UserTime != nil {
		result["user_time"] = *e.UserTime
	}

	if e.ProcessorTime != nil {
		result["processor_time"] = *e.ProcessorTime
	}

	return result
}

// Correlation contains the activity identifiers that consumers can use to group related events together.
type Correlation struct {
	// ActivityID and RelatedActivityID are optional fields
	// https://learn.microsoft.com/en-us/windows/win32/wes/eventschema-correlation-systempropertiestype-element
	ActivityID        *string `xml:"ActivityID,attr"`
	RelatedActivityID *string `xml:"RelatedActivityID,attr"`
}

func (e Correlation) asMap() map[string]any {
	result := map[string]any{}

	if e.ActivityID != nil {
		result["activity_id"] = *e.ActivityID
	}

	if e.RelatedActivityID != nil {
		result["related_activity_id"] = *e.RelatedActivityID
	}

	return result
}

// RenderingInfo contains human-readable strings for event fields, populated
// when the event is rendered with a publisher metadata (RenderDeep path).
// https://learn.microsoft.com/en-us/windows/win32/wes/eventschema-renderinginfotype-complextype
type RenderingInfo struct {
	Culture  string   `xml:"Culture,attr"`
	Message  string   `xml:"Message"`
	Level    string   `xml:"Level"`
	Task     string   `xml:"Task"`
	Opcode   string   `xml:"Opcode"`
	Channel  string   `xml:"Channel"`
	Provider string   `xml:"Provider"`
	Keywords []string `xml:"Keywords>Keyword"`
}

func (r RenderingInfo) asMap() map[string]any {
	return map[string]any{
		"culture":  r.Culture,
		"message":  r.Message,
		"level":    r.Level,
		"task":     r.Task,
		"opcode":   r.Opcode,
		"channel":  r.Channel,
		"provider": r.Provider,
		"keywords": r.Keywords,
	}
}

// UserData contains provider-defined event data as an alternative to EventData.
// The structure is arbitrary and defined by each provider's XML manifest.
// https://learn.microsoft.com/en-us/windows/win32/wes/eventschema-userdatatype-complextype
type UserData struct {
	// Name is the local name of the first child element, which identifies the event type.
	Name string
	// Data holds the key-value pairs parsed from the first child element's children.
	Data map[string]string
}

// UnmarshalXML implements xml.Unmarshaler for UserData.
// It reads the first child element and collects its direct children as key-value pairs.
func (u *UserData) UnmarshalXML(d *xml.Decoder, _ xml.StartElement) error {
	// Find the first child element of <UserData>, which names the event type.
	var innerStart xml.StartElement
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			innerStart = t
			goto parseChildren
		case xml.EndElement:
			return nil // empty <UserData>
		}
	}

parseChildren:
	u.Name = innerStart.Name.Local
	u.Data = make(map[string]string)

	// Collect direct children of the inner element as key-value pairs.
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch t := tok.(type) {
		case xml.StartElement:
			var value string
			if err := d.DecodeElement(&value, &t); err != nil {
				return err
			}
			u.Data[t.Name.Local] = value
		case xml.EndElement:
			// Consumed the inner element; skip remaining tokens up to </UserData>.
			return d.Skip()
		}
	}
}

func (u UserData) asMap() map[string]any {
	result := map[string]any{
		"name": u.Name,
	}
	if len(u.Data) > 0 {
		result["data"] = u.Data
	}
	return result
}

// ProcessingErrorData contains error information when an event cannot be rendered.
// https://learn.microsoft.com/en-us/windows/win32/wes/eventschema-processingerrordata-eventtype-element
type ProcessingErrorData struct {
	ErrorCode    uint32 `xml:"ErrorCode"`
	DataItemName string `xml:"DataItemName"`
	EventPayload string `xml:"EventPayload"`
}

func (p ProcessingErrorData) asMap() map[string]any {
	return map[string]any{
		"error_code":     p.ErrorCode,
		"data_item_name": p.DataItemName,
		"event_payload":  p.EventPayload,
	}
}

// DebugData contains data logged for Windows software tracing.
// https://learn.microsoft.com/en-us/windows/win32/wes/eventschema-debugdata-eventtype-element
type DebugData struct {
	SequenceNumber uint32 `xml:"SequenceNumber"`
	FlagName       string `xml:"FlagName"`
	LevelName      string `xml:"LevelName"`
	Component      string `xml:"Component"`
	SubComponent   string `xml:"SubComponent"`
	FileLine       string `xml:"FileLine"`
	Function       string `xml:"Function"`
	Message        string `xml:"Message"`
}

func (d DebugData) asMap() map[string]any {
	return map[string]any{
		"sequence_number": d.SequenceNumber,
		"flag_name":       d.FlagName,
		"level_name":      d.LevelName,
		"component":       d.Component,
		"sub_component":   d.SubComponent,
		"file_line":       d.FileLine,
		"function":        d.Function,
		"message":         d.Message,
	}
}

// UnmarshalEventXML will unmarshal EventXML from xml bytes.
// parsedEvent is the interface consumed by sendEvent. All fields accessed in
// sendEvent must go through this interface so that the compiler enforces that
// rawParsedEvent explicitly supports any new access added to the raw path.
type parsedEvent interface {
	getOriginal() string
	getSystemTime() string
	getLevel() string
	getRenderedLevel() string
	getProviderName() string
	getProviderGUID() string
	getEventSourceName() string
	getEventID() uint32
	getEventIDQualifiers() uint16
	getChannel() string
	getComputer() string
	getRecordID() uint64
	getUserID() string
	getProcessID() uint
	getThreadID() uint
	getTask() string
	getOpcode() string
	getKeywords() string
	getVersion() uint8
	getActivityID() string
	getRelatedActivityID() string
	getRenderedMessage() string
	getRenderedTask() string
	getRenderedOpcode() string
	getRenderedKeywords() []string
	getUserData() *UserData
	getBinaryEventData() string
	getEventData() EventData
	// formattedBody returns the structured body map for non-raw mode.
	// Panics if called on rawParsedEvent — only valid when raw=false.
	toEventXML() *EventXML
}

func (e *EventXML) getOriginal() string          { return e.Original }
func (e *EventXML) getSystemTime() string        { return e.TimeCreated.SystemTime }
func (e *EventXML) getLevel() string             { return e.Level }
func (e *EventXML) getProviderName() string      { return e.Provider.Name }
func (e *EventXML) getProviderGUID() string      { return e.Provider.GUID }
func (e *EventXML) getEventSourceName() string   { return e.Provider.EventSourceName }
func (e *EventXML) getEventID() uint32           { return e.EventID.ID }
func (e *EventXML) getEventIDQualifiers() uint16 { return e.EventID.Qualifiers }
func (e *EventXML) getChannel() string           { return e.Channel }
func (e *EventXML) getComputer() string          { return e.Computer }
func (e *EventXML) getRecordID() uint64          { return e.RecordID }
func (e *EventXML) getTask() string              { return e.Task }
func (e *EventXML) getOpcode() string            { return e.Opcode }
func (e *EventXML) getKeywords() string          { return e.Keywords }
func (e *EventXML) getVersion() uint8            { return e.Version }
func (e *EventXML) getBinaryEventData() string   { return e.BinaryEventData }
func (e *EventXML) getEventData() EventData      { return e.EventData }
func (e *EventXML) toEventXML() *EventXML        { return e }

func (e *EventXML) getRenderedLevel() string {
	if e.RenderingInfo != nil {
		return e.RenderingInfo.Level
	}
	return ""
}
func (e *EventXML) getRenderedMessage() string {
	if e.RenderingInfo != nil {
		return e.RenderingInfo.Message
	}
	return ""
}
func (e *EventXML) getRenderedTask() string {
	if e.RenderingInfo != nil {
		return e.RenderingInfo.Task
	}
	return ""
}
func (e *EventXML) getRenderedOpcode() string {
	if e.RenderingInfo != nil {
		return e.RenderingInfo.Opcode
	}
	return ""
}
func (e *EventXML) getRenderedKeywords() []string {
	if e.RenderingInfo != nil {
		return e.RenderingInfo.Keywords
	}
	return nil
}
func (e *EventXML) getUserID() string {
	if e.Security != nil {
		return e.Security.UserID
	}
	return ""
}
func (e *EventXML) getProcessID() uint {
	if e.Execution != nil {
		return e.Execution.ProcessID
	}
	return 0
}
func (e *EventXML) getThreadID() uint {
	if e.Execution != nil {
		return e.Execution.ThreadID
	}
	return 0
}
func (e *EventXML) getActivityID() string {
	if e.Correlation != nil && e.Correlation.ActivityID != nil {
		return *e.Correlation.ActivityID
	}
	return ""
}
func (e *EventXML) getRelatedActivityID() string {
	if e.Correlation != nil && e.Correlation.RelatedActivityID != nil {
		return *e.Correlation.RelatedActivityID
	}
	return ""
}
func (e *EventXML) getUserData() *UserData { return e.UserData }

// unmarshalEventXML will unmarshal EventXML from xml bytes.
// Illegal XML 1.0 characters (e.g. U+0001 found in some Sysmon events) are
// stripped before parsing so that a single malformed event does not halt the
// entire receiver.
func UnmarshalEventXML(data []byte) (parsedEvent, error) {
	sanitized := sanitizeXMLBytes(data)
	var eventXML *EventXML
	if err := xml.Unmarshal(sanitized, &eventXML); err != nil {
		return nil, fmt.Errorf("failed to unmarshal xml bytes into event: %w (%s)", err, string(sanitized))
	}
	// The sanitized bytes are only required for XML unmarshalling - the original data is preserved.
	eventXML.Original = string(data)
	return eventXML, nil
}

// rawEventXML holds the fields needed when raw=true. Includes all System fields
// plus EventData, UserData, and RenderingInfo for comprehensive attribute extraction.
type rawEventXML struct {
	Original        string               `xml:"-"`
	TimeCreated     TimeCreated          `xml:"System>TimeCreated"`
	Level           string               `xml:"System>Level"`
	Provider        Provider             `xml:"System>Provider"`
	EventID         EventID              `xml:"System>EventID"`
	Channel         string               `xml:"System>Channel"`
	Computer        string               `xml:"System>Computer"`
	RecordID        uint64               `xml:"System>EventRecordID"`
	Task            string               `xml:"System>Task"`
	Opcode          string               `xml:"System>Opcode"`
	Keywords        string               `xml:"System>Keywords"`
	Version         uint8                `xml:"System>Version"`
	Security        *Security            `xml:"System>Security"`
	Execution       *Execution           `xml:"System>Execution"`
	Correlation     *Correlation         `xml:"System>Correlation"`
	EventData       EventData            `xml:"EventData"`
	UserData        *UserData            `xml:"UserData"`
	BinaryEventData string               `xml:"BinaryEventData"`
	RenderingInfo   *rawRenderingInfoXML `xml:"RenderingInfo"`
}

// rawRenderingInfoXML holds rendered fields from RenderingInfo for raw mode.
type rawRenderingInfoXML struct {
	Level    string   `xml:"Level"`
	Message  string   `xml:"Message"`
	Task     string   `xml:"Task"`
	Opcode   string   `xml:"Opcode"`
	Keywords []string `xml:"Keywords>Keyword"`
}

func (r *rawEventXML) getOriginal() string          { return r.Original }
func (r *rawEventXML) getSystemTime() string        { return r.TimeCreated.SystemTime }
func (r *rawEventXML) getLevel() string             { return r.Level }
func (r *rawEventXML) getProviderName() string      { return r.Provider.Name }
func (r *rawEventXML) getProviderGUID() string      { return r.Provider.GUID }
func (r *rawEventXML) getEventSourceName() string   { return r.Provider.EventSourceName }
func (r *rawEventXML) getEventID() uint32           { return r.EventID.ID }
func (r *rawEventXML) getEventIDQualifiers() uint16 { return r.EventID.Qualifiers }
func (r *rawEventXML) getChannel() string           { return r.Channel }
func (r *rawEventXML) getComputer() string          { return r.Computer }
func (r *rawEventXML) getRecordID() uint64          { return r.RecordID }
func (r *rawEventXML) getTask() string              { return r.Task }
func (r *rawEventXML) getOpcode() string            { return r.Opcode }
func (r *rawEventXML) getKeywords() string          { return r.Keywords }
func (r *rawEventXML) getVersion() uint8            { return r.Version }
func (r *rawEventXML) getBinaryEventData() string   { return r.BinaryEventData }
func (r *rawEventXML) getEventData() EventData      { return r.EventData }
func (r *rawEventXML) getUserData() *UserData       { return r.UserData }

func (r *rawEventXML) getRenderedLevel() string {
	if r.RenderingInfo != nil {
		return r.RenderingInfo.Level
	}
	return ""
}
func (r *rawEventXML) getRenderedMessage() string {
	if r.RenderingInfo != nil {
		return r.RenderingInfo.Message
	}
	return ""
}
func (r *rawEventXML) getRenderedTask() string {
	if r.RenderingInfo != nil {
		return r.RenderingInfo.Task
	}
	return ""
}
func (r *rawEventXML) getRenderedOpcode() string {
	if r.RenderingInfo != nil {
		return r.RenderingInfo.Opcode
	}
	return ""
}
func (r *rawEventXML) getRenderedKeywords() []string {
	if r.RenderingInfo != nil {
		return r.RenderingInfo.Keywords
	}
	return nil
}
func (r *rawEventXML) getUserID() string {
	if r.Security != nil {
		return r.Security.UserID
	}
	return ""
}
func (r *rawEventXML) getProcessID() uint {
	if r.Execution != nil {
		return r.Execution.ProcessID
	}
	return 0
}
func (r *rawEventXML) getThreadID() uint {
	if r.Execution != nil {
		return r.Execution.ThreadID
	}
	return 0
}
func (r *rawEventXML) getActivityID() string {
	if r.Correlation != nil && r.Correlation.ActivityID != nil {
		return *r.Correlation.ActivityID
	}
	return ""
}
func (r *rawEventXML) getRelatedActivityID() string {
	if r.Correlation != nil && r.Correlation.RelatedActivityID != nil {
		return *r.Correlation.RelatedActivityID
	}
	return ""
}

func (*rawEventXML) toEventXML() *EventXML {
	panic("toEventXML called on rawEventXML: only valid in non-raw mode")
}

// unmarshalRawEventXML parses only the fields needed when raw=true and returns
// a rawParsedEvent. Use this instead of unmarshalEventXML when raw=true to
// avoid populating fields that will not be used.
func unmarshalRawEventXML(data []byte) (parsedEvent, error) {
	sanitized := sanitizeXMLBytes(data)
	var raw rawEventXML
	if err := xml.Unmarshal(sanitized, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal xml bytes into event: %w (%s)", err, string(sanitized))
	}
	// The sanitized bytes are only required for XML unmarshalling - the original data is preserved.
	raw.Original = string(data)
	return &raw, nil
}

// sanitizeXMLBytes removes characters that are illegal in XML 1.0 documents
// and strips the default XML namespace to allow Go's xml package to parse
// Windows Event Log XML correctly.
// XML 1.0 permits: #x9 | #xA | #xD | [#x20-#xD7FF] | [#xE000-#xFFFD] | [#x10000-#x10FFFF]
// All other code points (e.g. U+0001–U+0008, U+000B, U+000C, U+000E–U+001F) are dropped.
func sanitizeXMLBytes(data []byte) []byte {
	// Strip default namespace to allow Go xml package to parse correctly
	// Pattern: xmlns='...' or xmlns="..."
	data = stripDefaultNamespace(data)

	if !hasIllegalXMLBytes(data) {
		return data
	}
	return bytes.Map(func(r rune) rune {
		if (r >= 0x20 && r <= 0xD7FF) || r == 0x09 || r == 0x0A || r == 0x0D ||
			(r >= 0xE000 && r <= 0xFFFD) ||
			r >= 0x10000 {
			return r
		}
		return -1
	}, data)
}

// stripDefaultNamespace removes xmlns='...' or xmlns="..." from the XML
// to allow Go's xml package to parse without namespace qualification issues.
func stripDefaultNamespace(data []byte) []byte {
	// Look for xmlns= patterns and remove them
	// Pattern 1: xmlns='...'
	// Pattern 2: xmlns="..."
	result := data

	// Handle xmlns='...'
	if idx := bytes.Index(result, []byte("xmlns='")); idx >= 0 {
		endIdx := bytes.Index(result[idx+7:], []byte("'"))
		if endIdx >= 0 {
			// Remove xmlns='...' including trailing space if present
			endPos := idx + 7 + endIdx + 1
			if endPos < len(result) && result[endPos] == ' ' {
				endPos++
			}
			result = append(result[:idx], result[endPos:]...)
		}
	}

	// Handle xmlns="..."
	if idx := bytes.Index(result, []byte("xmlns=\"")); idx >= 0 {
		endIdx := bytes.Index(result[idx+7:], []byte("\""))
		if endIdx >= 0 {
			// Remove xmlns="..." including trailing space if present
			endPos := idx + 7 + endIdx + 1
			if endPos < len(result) && result[endPos] == ' ' {
				endPos++
			}
			result = append(result[:idx], result[endPos:]...)
		}
	}

	return result
}

// hasIllegalXMLBytes reports whether data contains any character that is
// illegal in an XML 1.0 document. It operates on raw bytes to avoid
// allocation: single-byte control characters are caught by a simple range
// check, and the only illegal multi-byte sequences handled explicitly are
// U+FFFE (EF BF BE) and U+FFFF (EF BF BF).
func hasIllegalXMLBytes(data []byte) bool {
	for i := range len(data) {
		b := data[i]
		if b < 0x20 {
			// 0x09 (tab), 0x0A (LF), 0x0D (CR) are the only legal control chars.
			if b != 0x09 && b != 0x0A && b != 0x0D {
				return true
			}
		} else if b == 0xEF && i+2 < len(data) {
			// U+FFFE = EF BF BE, U+FFFF = EF BF BF — both illegal in XML 1.0.
			if data[i+1] == 0xBF && (data[i+2] == 0xBE || data[i+2] == 0xBF) {
				return true
			}
		}
	}
	return false
}
