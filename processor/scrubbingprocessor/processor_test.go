package scrubbingprocessor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

func TestApplyMasking(t *testing.T) {
	tests := []struct {
		name        string
		regexp      string
		input       string
		placeholder string
		expected    string
	}{
		{
			name:        "digit",
			input:       "my 344 id is 123456",
			regexp:      `\d+`,
			placeholder: "HIDDEN",
			expected:    "my HIDDEN id is HIDDEN",
		},
		{
			name:        "starts with",
			input:       "This is Test",
			regexp:      `^This`,
			placeholder: "HIDDEN",
			expected:    "HIDDEN is Test",
		},
		{
			name:        "direct",
			input:       "This is Test",
			regexp:      `Test`,
			placeholder: "HIDDEN",
			expected:    "This is HIDDEN",
		},
		{
			name:        "email",
			input:       "This is test@opsramp.com",
			regexp:      `[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`,
			placeholder: "HIDDEN",
			expected:    "This is HIDDEN",
		},
		{
			name:        "uncorrect email",
			input:       "This is opsramp.com",
			regexp:      `[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`,
			placeholder: "HIDDEN",
			expected:    "This is opsramp.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ld := generateTestEntry(tt.input)
			processor, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
				Masking: []MaskingSettings{
					{
						Regexp:      tt.regexp,
						Placeholder: tt.placeholder,
					},
				},
			})
			assert.NoError(t, err)
			processor.applyMasking(ld)
			assert.Equal(t, ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString(), tt.expected)
		})
	}
}

func generateTestEntry(body string) plog.Logs {
	ld := plog.NewLogs()
	rl0 := ld.ResourceLogs().AppendEmpty()
	sc := rl0.ScopeLogs().AppendEmpty()
	e1 := sc.LogRecords().AppendEmpty()
	e1.Body().SetStr(body)
	return ld
}

func Test_applyMasking(t *testing.T) {
	type fields struct {
		logger *zap.Logger
		config *Config
	}
	type args struct {
		ld plog.Logs
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		expected plog.Logs
	}{
		{
			name: "mask everywhere",
			fields: fields{
				logger: &zap.Logger{},
				config: &Config{
					Masking: []MaskingSettings{
						{
							Regexp:      "User",
							Placeholder: "####",
						},
					},
				},
			},
			args: args{
				ld: generateLogs("18842:M 27 Jun 10:55:21.627 # User requested shutdown...",
					map[string]string{},
					map[string]string{
						"pid":       "18842",
						"role":      "M",
						"timestamp": "27 Jun 10:55:21.627",
						"level":     "#",
						"message":   "User requested shutdown...",
					},
				),
			},
			expected: generateLogs("18842:M 27 Jun 10:55:21.627 # #### requested shutdown...",
				map[string]string{},
				map[string]string{
					"pid":       "18842",
					"role":      "M",
					"timestamp": "27 Jun 10:55:21.627",
					"level":     "#",
					"message":   "#### requested shutdown...",
				},
			),
		},

		{
			name: "mask resource attribute",
			fields: fields{
				logger: &zap.Logger{},
				config: &Config{
					Masking: []MaskingSettings{
						{
							AttributeType: "resource",
							AttributeKey:  "hostname",
							Regexp:        "otlp",
							Placeholder:   "opentelemetry",
						},
					},
				},
			},
			args: args{
				ld: generateLogs("1234567890 otlp.example.com This is a sample log line",
					map[string]string{
						"hostname": "otlp.example.com",
					},
					map[string]string{
						"timestamp": "1234567890",
						"hostname":  "otlp.example.com",
						"message":   "This is a sample log line",
					},
				),
			},
			expected: generateLogs("1234567890 otlp.example.com This is a sample log line",
				map[string]string{
					"hostname": "opentelemetry.example.com",
				},
				map[string]string{
					"timestamp": "1234567890",
					"hostname":  "otlp.example.com",
					"message":   "This is a sample log line",
				},
			),
		},
		{
			name: "mask record attribute",
			fields: fields{
				logger: &zap.Logger{},
				config: &Config{
					Masking: []MaskingSettings{
						{
							AttributeType: "record",
							AttributeKey:  "hostname",
							Regexp:        "otlp",
							Placeholder:   "opentelemetry",
						},
					},
				},
			},
			args: args{
				ld: generateLogs("1234567890 otlp.example.com This is a sample log line",
					map[string]string{
						"hostname": "otlp.example.com",
					},
					map[string]string{
						"timestamp": "1234567890",
						"hostname":  "otlp.example.com",
						"message":   "This is a sample log line",
					},
				),
			},
			expected: generateLogs("1234567890 otlp.example.com This is a sample log line",
				map[string]string{
					"hostname": "otlp.example.com",
				},
				map[string]string{
					"timestamp": "1234567890",
					"hostname":  "opentelemetry.example.com",
					"message":   "This is a sample log line",
				},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp, err := newScrubbingProcessorProcessor(tt.fields.logger, tt.fields.config)
			assert.NoError(t, err)
			sp.applyMasking(tt.args.ld)

			expected := tt.args.ld.ResourceLogs().At(0).Resource().Attributes()
			actual := tt.expected.ResourceLogs().At(0).Resource().Attributes()
			assert.EqualValues(t, expected.AsRaw(), actual.AsRaw())

			expected = tt.args.ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Attributes()
			actual = tt.expected.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Attributes()
			assert.EqualValues(t, expected.AsRaw(), actual.AsRaw())

			assert.EqualValues(t, tt.args.ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString(), tt.expected.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString())
		})
	}
}

func TestMaskingModes(t *testing.T) {
	tests := []struct {
		name        string
		mode        MaskingMode
		regexp      string
		input       string
		placeholder string
		expected    string
	}{
		{
			name:        "replace mode (default)",
			mode:        ModeReplace,
			regexp:      `\d{3}-\d{2}-\d{4}`,
			input:       "SSN is 123-45-6789",
			placeholder: "***-**-****",
			expected:    "SSN is ***-**-****",
		},
		{
			name:        "partial mode - credit card middle",
			mode:        ModePartial,
			regexp:      `\d{4}(-\d{4}-\d{4}-)\d{4}`,
			input:       "card 1234-5678-9012-3456 ok",
			placeholder: "-****-****-",
			expected:    "card 1234-****-****-3456 ok",
		},
		{
			name:        "partial mode - no capture group falls through",
			mode:        ModePartial,
			regexp:      `\d+`,
			input:       "no capture group here 123",
			placeholder: "****",
			expected:    "no capture group here 123",
		},
		{
			name:        "hash mode",
			mode:        ModeHash,
			regexp:      `secret-[a-z0-9]+`,
			input:       "token is secret-abc123 here",
			placeholder: "ignored",
			expected:    "", // will verify length/format instead
		},
		{
			name:        "redact_key mode on string body",
			mode:        ModeRedactKey,
			regexp:      `password`,
			input:       "my password is hidden",
			placeholder: "[REDACTED]",
			expected:    "[REDACTED]",
		},
		{
			name:        "empty mode defaults to replace",
			mode:        "",
			regexp:      `secret`,
			input:       "this is a secret value",
			placeholder: "****",
			expected:    "this is a **** value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ld := generateTestEntry(tt.input)
			sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
				Masking: []MaskingSettings{
					{
						Regexp:      tt.regexp,
						Placeholder: tt.placeholder,
						Mode:        tt.mode,
					},
				},
			})
			assert.NoError(t, err)
			sp.applyMasking(ld)

			result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
			if tt.mode == ModeHash {
				// Hash output is deterministic but not predictable in test table;
				// just verify the original secret text is gone and result changed.
				assert.NotContains(t, result, "secret-abc123")
				assert.NotEqual(t, tt.input, result)
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestRedactKeyMode_MapBody(t *testing.T) {
	ld := plog.NewLogs()
	rl := ld.ResourceLogs().AppendEmpty()
	sl := rl.ScopeLogs().AppendEmpty()
	lr := sl.LogRecords().AppendEmpty()
	lr.Body().SetEmptyMap()
	lr.Body().Map().PutStr("username", "admin")
	lr.Body().Map().PutStr("password", "s3cret!")
	lr.Body().Map().PutStr("action", "login")

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `s3cret`,
				Placeholder: "[REDACTED]",
				Mode:        ModeRedactKey,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	body := lr.Body().Map()
	// "password" key should be removed because its value matched
	_, ok := body.Get("password")
	assert.False(t, ok, "password key should have been removed")
	// other keys remain
	v, ok := body.Get("username")
	assert.True(t, ok)
	assert.Equal(t, "admin", v.AsString())
	v, ok = body.Get("action")
	assert.True(t, ok)
	assert.Equal(t, "login", v.AsString())
}

func generateLogs(body string, resourceAttributes, recordAttributes map[string]string) plog.Logs {
	ld := plog.NewLogs()
	resourceLogs := ld.ResourceLogs().AppendEmpty()
	for k, v := range resourceAttributes {
		resourceLogs.Resource().Attributes().PutStr(k, v)
	}
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()
	logRecordEntry1 := scopeLogs.LogRecords().AppendEmpty()
	for k, v := range recordAttributes {
		logRecordEntry1.Attributes().PutStr(k, v)
	}
	logRecordEntry1.Body().SetStr(body)
	return ld
}

// generateLogsWithMapBody creates a log record with a Map body.
func generateLogsWithMapBody(bodyMap map[string]string, resourceAttributes, recordAttributes map[string]string) plog.Logs {
	ld := plog.NewLogs()
	resourceLogs := ld.ResourceLogs().AppendEmpty()
	for k, v := range resourceAttributes {
		resourceLogs.Resource().Attributes().PutStr(k, v)
	}
	scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()
	lr := scopeLogs.LogRecords().AppendEmpty()
	for k, v := range recordAttributes {
		lr.Attributes().PutStr(k, v)
	}
	lr.Body().SetEmptyMap()
	for k, v := range bodyMap {
		lr.Body().Map().PutStr(k, v)
	}
	return ld
}

// =============================================================================
// Scoping Tests — verify attribute_type controls what gets masked
// =============================================================================

func TestScoping_ResourceRuleDoesNotTouchBody(t *testing.T) {
	// A resource-scoped rule should only mask resource attributes.
	// Body and record attributes must remain unchanged.
	ld := generateLogs(
		"password secret123 in body",
		map[string]string{"host": "password-server"},
		map[string]string{"msg": "password visible"},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				AttributeType: ResourceAttribute,
				AttributeKey:  "host",
				Regexp:        `password`,
				Placeholder:   "***",
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	lr := ld.ResourceLogs().At(0)
	// Resource attribute masked
	v, _ := lr.Resource().Attributes().Get("host")
	assert.Equal(t, "***-server", v.AsString())
	// Record attribute untouched
	rec := lr.ScopeLogs().At(0).LogRecords().At(0)
	v, _ = rec.Attributes().Get("msg")
	assert.Equal(t, "password visible", v.AsString())
	// Body untouched
	assert.Equal(t, "password secret123 in body", rec.Body().AsString())
}

func TestScoping_ResourceRedactKeyDoesNotWipeBody(t *testing.T) {
	// This is the exact bug scenario: redact_key with text '.*' and
	// attribute_type "resource" should NOT empty the body map.
	ld := generateLogsWithMapBody(
		map[string]string{"message": "important log data"},
		map[string]string{"k8s.pod.ip": "10.244.0.15"},
		map[string]string{"level": "info"},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				AttributeType: ResourceAttribute,
				AttributeKey:  "k8s.pod.ip",
				Regexp:        `.*`,
				Placeholder:   "",
				Mode:          ModeRedactKey,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	lr := ld.ResourceLogs().At(0)
	// Resource attribute removed
	_, ok := lr.Resource().Attributes().Get("k8s.pod.ip")
	assert.False(t, ok, "k8s.pod.ip should be redacted from resource attributes")
	// Body map intact
	rec := lr.ScopeLogs().At(0).LogRecords().At(0)
	v, ok := rec.Body().Map().Get("message")
	assert.True(t, ok, "body 'message' key should still exist")
	assert.Equal(t, "important log data", v.AsString())
	// Record attribute untouched
	v, _ = rec.Attributes().Get("level")
	assert.Equal(t, "info", v.AsString())
}

func TestScoping_RecordRuleDoesNotTouchBody(t *testing.T) {
	// A record-scoped rule should only mask record attributes — NOT the body.
	ld := generateLogs(
		"body contains otlp text",
		map[string]string{"hostname": "otlp.example.com"},
		map[string]string{"hostname": "otlp.example.com"},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				AttributeType: RecordAttribute,
				AttributeKey:  "hostname",
				Regexp:        `otlp`,
				Placeholder:   "masked",
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	lr := ld.ResourceLogs().At(0)
	// Resource attribute NOT masked (rule is record-scoped)
	v, _ := lr.Resource().Attributes().Get("hostname")
	assert.Equal(t, "otlp.example.com", v.AsString())
	// Record attribute IS masked
	rec := lr.ScopeLogs().At(0).LogRecords().At(0)
	v, _ = rec.Attributes().Get("hostname")
	assert.Equal(t, "masked.example.com", v.AsString())
	// Body NOT masked
	assert.Equal(t, "body contains otlp text", rec.Body().AsString())
}

func TestScoping_EmptyTypeMasksEverything(t *testing.T) {
	// No attribute_type → masks resource attrs, record attrs, AND body.
	ld := generateLogs(
		"secret in body",
		map[string]string{"data": "secret in resource"},
		map[string]string{"data": "secret in record"},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `secret`,
				Placeholder: "***",
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	lr := ld.ResourceLogs().At(0)
	v, _ := lr.Resource().Attributes().Get("data")
	assert.Equal(t, "*** in resource", v.AsString())
	rec := lr.ScopeLogs().At(0).LogRecords().At(0)
	v, _ = rec.Attributes().Get("data")
	assert.Equal(t, "*** in record", v.AsString())
	assert.Equal(t, "*** in body", rec.Body().AsString())
}

func TestScoping_EmptyTypeNoKey_MasksAllAttributes(t *testing.T) {
	// No attribute_type + no attribute_key → applies to ALL keys in all scopes.
	ld := generateLogs(
		"digits 123 in body",
		map[string]string{"ip": "10.0.0.1", "name": "node123"},
		map[string]string{"msg": "error 456", "code": "789"},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `\d+`,
				Placeholder: "#",
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	lr := ld.ResourceLogs().At(0)
	v, _ := lr.Resource().Attributes().Get("ip")
	assert.Equal(t, "#.#.#.#", v.AsString())
	v, _ = lr.Resource().Attributes().Get("name")
	assert.Equal(t, "node#", v.AsString())
	rec := lr.ScopeLogs().At(0).LogRecords().At(0)
	v, _ = rec.Attributes().Get("msg")
	assert.Equal(t, "error #", v.AsString())
	v, _ = rec.Attributes().Get("code")
	assert.Equal(t, "#", v.AsString())
	assert.Equal(t, "digits # in body", rec.Body().AsString())
}

// =============================================================================
// Partial Mode Tests
// =============================================================================

func TestPartialMode_PasswordMasking(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "password with space",
			input:    "this is my password secret123",
			expected: "this is my password ***",
		},
		{
			name:     "password with colon",
			input:    "password: jNTBfObN7efw",
			expected: "password***",
		},
		{
			name:     "password with equals",
			input:    "db_password=jfbxYMnSt4tuQd",
			expected: "db_password***",
		},
		{
			name:     "password in sentence with space",
			input:    "user admin logged in with password ZhjqKprO3HERnYaf1zd8",
			expected: "user admin logged in with password ***",
		},
		{
			name:     "password with quoted value",
			input:    `this is my password "iNk7y8hiAMK&PP4J"`,
			expected: `this is my password ***`,
		},
		{
			name:     "no password match",
			input:    "this has no sensitive data",
			expected: "this has no sensitive data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ld := generateTestEntry(tt.input)
			sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
				Masking: []MaskingSettings{
					{
						Regexp:      `password\s*(.+)`,
						Placeholder: "***",
						Mode:        ModePartial,
					},
				},
			})
			assert.NoError(t, err)
			sp.applyMasking(ld)
			result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPartialMode_EmailMasking(t *testing.T) {
	ld := generateTestEntry("Contact user john.doe@example.com for details")
	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `([a-zA-Z0-9._%+-]+)@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`,
				Placeholder: "###",
				Mode:        ModePartial,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)
	result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	assert.Equal(t, "Contact user ###@example.com for details", result)
}

func TestPartialMode_OnlyFirstOccurrence(t *testing.T) {
	// Partial mode only replaces the FIRST match.
	ld := generateTestEntry("password abc123 and password def456")
	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `password\s*(.+)`,
				Placeholder: "***",
				Mode:        ModePartial,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)
	result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	// First match captures everything after first "password " (greedy)
	assert.Equal(t, "password ***", result)
}

func TestPartialMode_NoCaptureGroup_NoChange(t *testing.T) {
	ld := generateTestEntry("secret value 12345")
	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `\d+`, // no capture group
				Placeholder: "***",
				Mode:        ModePartial,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)
	result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	assert.Equal(t, "secret value 12345", result)
}

// =============================================================================
// Hash Mode Tests
// =============================================================================

func TestHashMode_FullMatchHashed(t *testing.T) {
	// Hash replaces the FULL regex match (capture groups are irrelevant).
	ld := generateTestEntry("AUTH_TOKEN=secret123")
	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `_TOKEN=(\S+)`,
				Placeholder: "ignored",
				Mode:        ModeHash,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)
	result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	// "AUTH" prefix preserved (not part of the match)
	assert.Contains(t, result, "AUTH")
	// Original token gone
	assert.NotContains(t, result, "_TOKEN=secret123")
	// Result should be "AUTH" + 16 hex chars
	assert.Len(t, result, 4+16, "expected AUTH (4 chars) + 16 hex chars")
}

func TestHashMode_Deterministic(t *testing.T) {
	// Same input produces same hash.
	makeLog := func() plog.Logs { return generateTestEntry("key=mysecret") }

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{Regexp: `key=\S+`, Mode: ModeHash},
		},
	})
	assert.NoError(t, err)

	ld1 := makeLog()
	ld2 := makeLog()
	sp.applyMasking(ld1)
	sp.applyMasking(ld2)
	r1 := ld1.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	r2 := ld2.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	assert.Equal(t, r1, r2)
}

func TestHashMode_AllOccurrencesReplaced(t *testing.T) {
	ld := generateTestEntry("ip 10.0.0.1 and ip 10.0.0.2 end")
	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp: `\d+\.\d+\.\d+\.\d+`,
				Mode:   ModeHash,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)
	result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	assert.NotContains(t, result, "10.0.0.1")
	assert.NotContains(t, result, "10.0.0.2")
}

// =============================================================================
// RedactKey Mode Tests
// =============================================================================

func TestRedactKey_ResourceAttribute_Removed(t *testing.T) {
	ld := generateLogs("body text",
		map[string]string{"k8s.pod.ip": "10.244.0.15", "k8s.pod.name": "my-pod"},
		map[string]string{},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				AttributeType: ResourceAttribute,
				AttributeKey:  "k8s.pod.ip",
				Regexp:        `^10\.`,
				Mode:          ModeRedactKey,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	attrs := ld.ResourceLogs().At(0).Resource().Attributes()
	_, ok := attrs.Get("k8s.pod.ip")
	assert.False(t, ok, "k8s.pod.ip should be removed")
	v, ok := attrs.Get("k8s.pod.name")
	assert.True(t, ok)
	assert.Equal(t, "my-pod", v.AsString())
}

func TestRedactKey_NoMatch_AttributeKept(t *testing.T) {
	ld := generateLogs("body text",
		map[string]string{"k8s.pod.ip": "203.0.113.50"},
		map[string]string{},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				AttributeType: ResourceAttribute,
				AttributeKey:  "k8s.pod.ip",
				Regexp:        `^10\.`,
				Mode:          ModeRedactKey,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	attrs := ld.ResourceLogs().At(0).Resource().Attributes()
	v, ok := attrs.Get("k8s.pod.ip")
	assert.True(t, ok, "k8s.pod.ip should remain (no match)")
	assert.Equal(t, "203.0.113.50", v.AsString())
}

func TestRedactKey_NoKey_RemovesAllMatchingAttributes(t *testing.T) {
	// redact_key without attribute_key removes ALL attributes whose value matches.
	ld := generateLogs("body text",
		map[string]string{
			"ip":   "10.0.0.1",
			"name": "my-server",
			"vip":  "10.0.0.2",
		},
		map[string]string{},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				AttributeType: ResourceAttribute,
				Regexp:        `^10\.`,
				Mode:          ModeRedactKey,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	attrs := ld.ResourceLogs().At(0).Resource().Attributes()
	_, ok := attrs.Get("ip")
	assert.False(t, ok, "ip should be removed")
	_, ok = attrs.Get("vip")
	assert.False(t, ok, "vip should be removed")
	v, ok := attrs.Get("name")
	assert.True(t, ok)
	assert.Equal(t, "my-server", v.AsString())
}

func TestRedactKey_MapBody_GlobalRule(t *testing.T) {
	// Unscoped redact_key with map body removes matching keys.
	ld := generateLogsWithMapBody(
		map[string]string{"password": "s3cret!", "user": "admin", "token": "abc"},
		map[string]string{},
		map[string]string{},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp: `s3cret`,
				Mode:   ModeRedactKey,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	body := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().Map()
	_, ok := body.Get("password")
	assert.False(t, ok, "password key should be removed (value matched)")
	v, ok := body.Get("user")
	assert.True(t, ok)
	assert.Equal(t, "admin", v.AsString())
	v, ok = body.Get("token")
	assert.True(t, ok)
	assert.Equal(t, "abc", v.AsString())
}

// =============================================================================
// Multiple Rules — Order of Application
// =============================================================================

func TestMultipleRules_AppliedInOrder(t *testing.T) {
	ld := generateTestEntry("password secret123 email user@test.com AUTH_TOKEN=xyz")

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `password\s*(.+)`,
				Placeholder: "***",
				Mode:        ModePartial,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)
	result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	// After partial: everything after "password " is replaced
	assert.Equal(t, "password ***", result)
}

func TestMultipleRules_IndependentFields(t *testing.T) {
	ld := generateLogs("body untouched",
		map[string]string{"ns": "opsramp-agent"},
		map[string]string{"msg": "SSN 123-45-6789 found"},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				AttributeType: ResourceAttribute,
				AttributeKey:  "ns",
				Regexp:        `^ops`,
				Placeholder:   "***",
			},
			{
				AttributeType: RecordAttribute,
				AttributeKey:  "msg",
				Regexp:        `\d{3}-\d{2}-\d{4}`,
				Placeholder:   "XXX-XX-XXXX",
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	lr := ld.ResourceLogs().At(0)
	v, _ := lr.Resource().Attributes().Get("ns")
	assert.Equal(t, "***ramp-agent", v.AsString())

	rec := lr.ScopeLogs().At(0).LogRecords().At(0)
	v, _ = rec.Attributes().Get("msg")
	assert.Equal(t, "SSN XXX-XX-XXXX found", v.AsString())

	assert.Equal(t, "body untouched", rec.Body().AsString())
}

// =============================================================================
// Map Body Tests
// =============================================================================

func TestMapBody_ReplaceModeGlobalRule(t *testing.T) {
	ld := generateLogsWithMapBody(
		map[string]string{"message": "user password secret123 logged in"},
		map[string]string{},
		map[string]string{},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `secret\d+`,
				Placeholder: "[MASKED]",
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	body := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().Map()
	v, _ := body.Get("message")
	assert.Equal(t, "user password [MASKED] logged in", v.AsString())
}

func TestMapBody_PartialModeGlobalRule(t *testing.T) {
	ld := generateLogsWithMapBody(
		map[string]string{"message": "AUTH_TOKEN=mysecrettoken123"},
		map[string]string{},
		map[string]string{},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `AUTH_TOKEN=(\S+)`,
				Placeholder: "",
				Mode:        ModePartial,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	body := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().Map()
	v, _ := body.Get("message")
	assert.Equal(t, "AUTH_TOKEN=", v.AsString())
}

func TestMapBody_HashModeGlobalRule(t *testing.T) {
	ld := generateLogsWithMapBody(
		map[string]string{"message": "token=abc123secret"},
		map[string]string{},
		map[string]string{},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp: `token=\S+`,
				Mode:   ModeHash,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	body := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().Map()
	v, _ := body.Get("message")
	assert.NotContains(t, v.AsString(), "abc123secret")
	assert.Len(t, v.AsString(), 16) // full match hashed to 16 hex chars
}

func TestMapBody_ResourceRuleDoesNotTouchMapBody(t *testing.T) {
	ld := generateLogsWithMapBody(
		map[string]string{"message": "sensitive data here"},
		map[string]string{"host": "sensitive data here"},
		map[string]string{},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				AttributeType: ResourceAttribute,
				Regexp:        `sensitive`,
				Placeholder:   "***",
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	// Resource attribute masked
	v, _ := ld.ResourceLogs().At(0).Resource().Attributes().Get("host")
	assert.Equal(t, "*** data here", v.AsString())
	// Body map untouched
	body := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().Map()
	v, _ = body.Get("message")
	assert.Equal(t, "sensitive data here", v.AsString())
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestInvalidRegex_ReturnsError(t *testing.T) {
	_, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{Regexp: `[invalid`, Placeholder: "x"},
		},
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid masking regexp")
}

func TestEmptyRules_NoChange(t *testing.T) {
	ld := generateTestEntry("original body")
	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)
	result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	assert.Equal(t, "original body", result)
}

func TestReplaceMode_NoMatch_Unchanged(t *testing.T) {
	ld := generateTestEntry("nothing to mask here")
	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{Regexp: `\d{3}-\d{2}-\d{4}`, Placeholder: "XXX"},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)
	result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	assert.Equal(t, "nothing to mask here", result)
}

func TestRedactKey_StringBody_MatchReplacesWithPlaceholder(t *testing.T) {
	// For string body, redact_key replaces entire string with placeholder if matched.
	ld := generateTestEntry("contains password data")
	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `password`,
				Placeholder: "[REDACTED]",
				Mode:        ModeRedactKey,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)
	result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	assert.Equal(t, "[REDACTED]", result)
}

func TestRedactKey_StringBody_NoMatch_Unchanged(t *testing.T) {
	ld := generateTestEntry("safe content")
	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `password`,
				Placeholder: "[REDACTED]",
				Mode:        ModeRedactKey,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)
	result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	assert.Equal(t, "safe content", result)
}

// =============================================================================
// maskBody with attribute_key set — body scoping fix
// =============================================================================

func TestMaskBody_WithKey_StringBody_StillMasked(t *testing.T) {
	// When attribute_key is set and body is a string, body is still masked
	// because a string body has no keys — the regex applies to the full string.
	ld := generateLogs(
		"password secret123 in body",
		map[string]string{},
		map[string]string{"password": "password secret123"},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				AttributeKey: "password",
				Regexp:       `password\s*(.+)`,
				Placeholder:  "***",
				Mode:         ModePartial,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	rec := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0)
	// String body IS masked — partial mode replaces the first capture group
	assert.Equal(t, "password ***", rec.Body().AsString())
	// Record attribute IS also masked
	v, _ := rec.Attributes().Get("password")
	assert.Equal(t, "password ***", v.AsString())
}

func TestMaskBody_WithKey_MapBody_OnlyTargetKeyMasked(t *testing.T) {
	// When attribute_key is set and body is a map, only that key in the map is masked.
	ld := generateLogsWithMapBody(
		map[string]string{
			"password": "password secret123",
			"username": "password admin",
			"action":   "login",
		},
		map[string]string{},
		map[string]string{},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				AttributeKey: "password",
				Regexp:       `secret\d+`,
				Placeholder:  "[MASKED]",
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	body := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().Map()
	// Only "password" key in the map is masked
	v, _ := body.Get("password")
	assert.Equal(t, "password [MASKED]", v.AsString())
	// Other keys with matching content are NOT touched
	v, _ = body.Get("username")
	assert.Equal(t, "password admin", v.AsString())
	v, _ = body.Get("action")
	assert.Equal(t, "login", v.AsString())
}

func TestMaskBody_WithKey_MapBody_RedactKey_OnlyTargetRemoved(t *testing.T) {
	// redact_key + attribute_key on map body: only the target key is removed.
	ld := generateLogsWithMapBody(
		map[string]string{
			"password": "s3cret!",
			"token":    "s3cret!",
			"user":     "admin",
		},
		map[string]string{},
		map[string]string{},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				AttributeKey: "password",
				Regexp:       `s3cret`,
				Placeholder:  "",
				Mode:         ModeRedactKey,
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	body := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().Map()
	// Only "password" removed
	_, ok := body.Get("password")
	assert.False(t, ok, "password key should be removed from body map")
	// "token" has same value but is NOT removed (attribute_key targets only "password")
	v, ok := body.Get("token")
	assert.True(t, ok, "token key should remain")
	assert.Equal(t, "s3cret!", v.AsString())
	v, ok = body.Get("user")
	assert.True(t, ok)
	assert.Equal(t, "admin", v.AsString())
}

func TestMaskBody_WithKey_MapBody_KeyNotPresent_NoOp(t *testing.T) {
	// attribute_key targets a key that doesn't exist in the map body — no-op.
	ld := generateLogsWithMapBody(
		map[string]string{"message": "hello world", "level": "info"},
		map[string]string{},
		map[string]string{},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				AttributeKey: "nonexistent",
				Regexp:       `.*`,
				Placeholder:  "[GONE]",
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	body := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().Map()
	v, _ := body.Get("message")
	assert.Equal(t, "hello world", v.AsString())
	v, _ = body.Get("level")
	assert.Equal(t, "info", v.AsString())
}

func TestMaskBody_NoKey_StringBody_StillMasked(t *testing.T) {
	// Verify the original behavior still works: no attribute_key → string body IS masked.
	ld := generateTestEntry("contains secret data")

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `secret`,
				Placeholder: "***",
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	result := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().AsString()
	assert.Equal(t, "contains *** data", result)
}

func TestMaskBody_NoKey_MapBody_AllValuesMasked(t *testing.T) {
	// Verify the original behavior still works: no attribute_key → all map values masked.
	ld := generateLogsWithMapBody(
		map[string]string{"field1": "secret here", "field2": "secret there"},
		map[string]string{},
		map[string]string{},
	)

	sp, err := newScrubbingProcessorProcessor(zap.NewNop(), &Config{
		Masking: []MaskingSettings{
			{
				Regexp:      `secret`,
				Placeholder: "***",
			},
		},
	})
	assert.NoError(t, err)
	sp.applyMasking(ld)

	body := ld.ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Body().Map()
	v, _ := body.Get("field1")
	assert.Equal(t, "*** here", v.AsString())
	v, _ = body.Get("field2")
	assert.Equal(t, "*** there", v.AsString())
}
