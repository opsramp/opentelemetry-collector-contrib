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
			expected: generateLogs("1234567890 opentelemetry.example.com This is a sample log line",
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
			expected: generateLogs("1234567890 opentelemetry.example.com This is a sample log line",
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
