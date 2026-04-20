package scrubbingprocessor

import (
	"context"
	"crypto/sha256"
	"fmt"
	"regexp"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

// compiledRule holds a pre-compiled regex alongside its settings so we never
// compile per-record.
type compiledRule struct {
	MaskingSettings
	re *regexp.Regexp
}

type scrubbingProcessor struct {
	logger *zap.Logger
	config *Config
	rules  []compiledRule
}

func newScrubbingProcessorProcessor(logger *zap.Logger, cfg *Config) (*scrubbingProcessor, error) {
	rules := make([]compiledRule, 0, len(cfg.Masking))
	for _, m := range cfg.Masking {
		re, err := regexp.Compile(m.Regexp)
		if err != nil {
			return nil, fmt.Errorf("invalid masking regexp %q: %w", m.Regexp, err)
		}
		// default mode
		if m.Mode == "" {
			m.Mode = ModeReplace
		}
		rules = append(rules, compiledRule{MaskingSettings: m, re: re})
	}
	return &scrubbingProcessor{
		logger: logger,
		config: cfg,
		rules:  rules,
	}, nil
}

func (sp *scrubbingProcessor) ProcessLogs(_ context.Context, logs plog.Logs) (plog.Logs, error) {
	if len(sp.rules) > 0 {
		sp.applyMasking(logs)
	}
	return logs, nil
}

// maskValue applies the masking mode to a single string and returns the result.
func (r *compiledRule) maskValue(input string) string {
	switch r.Mode {
	case ModePartial:
		// Requires at least one capture group. Keeps the parts outside the
		// first group intact and replaces only the group match.
		loc := r.re.FindStringSubmatchIndex(input)
		if len(loc) < 4 {
			return input // no match or no capture group
		}
		// loc[2]:loc[3] is the first capture group
		return input[:loc[2]] + r.Placeholder + input[loc[3]:]

	case ModeHash:
		return r.re.ReplaceAllStringFunc(input, func(match string) string {
			h := sha256.Sum256([]byte(match))
			return fmt.Sprintf("%x", h[:8]) // 16 hex chars
		})

	case ModeRedactKey:
		// For string values, replace entirely if it matches.
		if r.re.MatchString(input) {
			return r.Placeholder
		}
		return input

	default: // ModeReplace
		return r.re.ReplaceAllString(input, r.Placeholder)
	}
}

func (sp *scrubbingProcessor) applyMasking(ld plog.Logs) {
	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resourceAttributes := ld.ResourceLogs().At(i).Resource().Attributes()
		for idx := range sp.rules {
			rule := &sp.rules[idx]
			sp.maskAttributes(rule, resourceAttributes, ResourceAttribute)
		}
	}

	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		resLogs := ld.ResourceLogs().At(i)
		for k := 0; k < resLogs.ScopeLogs().Len(); k++ {
			scopedLog := resLogs.ScopeLogs().At(k)
			for z := 0; z < scopedLog.LogRecords().Len(); z++ {
				log := scopedLog.LogRecords().At(z)
				for idx := range sp.rules {
					rule := &sp.rules[idx]

					// record attributes
					sp.maskAttributes(rule, log.Attributes(), RecordAttribute)

					// body
					sp.maskBody(rule, log)
				}
			}
		}
	}
}

// maskAttributes masks values in an attribute map according to the rule.
func (sp *scrubbingProcessor) maskAttributes(rule *compiledRule, attrs pcommon.Map, scope AttributeType) {
	if rule.AttributeType != scope && rule.AttributeType != EmptyAttribute {
		return
	}

	if rule.AttributeKey == "" {
		// apply to all attributes
		if rule.Mode == ModeRedactKey {
			// collect keys to delete (can't modify map during Range)
			var toDelete []string
			attrs.Range(func(key string, val pcommon.Value) bool {
				if rule.re.MatchString(val.AsString()) {
					toDelete = append(toDelete, key)
				}
				return true
			})
			for _, key := range toDelete {
				attrs.Remove(key)
			}
		} else {
			attrs.Range(func(key string, val pcommon.Value) bool {
				val.SetStr(rule.maskValue(val.AsString()))
				return true
			})
		}
	} else {
		if rule.Mode == ModeRedactKey {
			if val, ok := attrs.Get(rule.AttributeKey); ok {
				if rule.re.MatchString(val.AsString()) {
					attrs.Remove(rule.AttributeKey)
				}
			}
		} else {
			if val, ok := attrs.Get(rule.AttributeKey); ok {
				val.SetStr(rule.maskValue(val.AsString()))
			}
		}
	}
}

// maskBody masks body content according to the rule.
// Body is always masked regardless of AttributeType — the attribute type only
// controls which attribute scope (resource vs record) is targeted.
func (sp *scrubbingProcessor) maskBody(rule *compiledRule, log plog.LogRecord) {

	switch log.Body().Type() {
	case pcommon.ValueTypeMap:
		if rule.Mode == ModeRedactKey {
			var toDelete []string
			log.Body().Map().Range(func(k string, v pcommon.Value) bool {
				if rule.re.MatchString(v.AsString()) {
					toDelete = append(toDelete, k)
				}
				return true
			})
			for _, key := range toDelete {
				log.Body().Map().Remove(key)
			}
		} else {
			log.Body().Map().Range(func(k string, v pcommon.Value) bool {
				v.SetStr(rule.maskValue(v.AsString()))
				return true
			})
		}
	case pcommon.ValueTypeStr:
		log.Body().SetStr(rule.maskValue(log.Body().AsString()))
	}
}
