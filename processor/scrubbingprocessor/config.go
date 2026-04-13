package scrubbingprocessor

import (
	"fmt"

	"go.opentelemetry.io/collector/component"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
)

type AttributeType string

const (
	EmptyAttribute    AttributeType = ""
	ResourceAttribute AttributeType = "resource"
	RecordAttribute   AttributeType = "record"
)

// MaskingMode controls how a matched value is redacted.
type MaskingMode string

const (
	ModeReplace   MaskingMode = "replace" //nolint:gci
	ModePartial   MaskingMode = "partial"
	ModeHash      MaskingMode = "hash"
	ModeRedactKey MaskingMode = "redact_key"
)

type MaskingSettings struct {
	AttributeType AttributeType `mapstructure:"attribute_type"`
	AttributeKey  string        `mapstructure:"attribute_key"`
	Regexp        string        `mapstructure:"regexp"`
	Placeholder   string        `mapstructure:"placeholder"`
	Mode          MaskingMode   `mapstructure:"mode"` // defaults to "replace" if empty
}

// Config defines the configuration for the scrubbing processor.
type Config struct {
	adapter.BaseConfig `mapstructure:",squash"`
	Masking            []MaskingSettings `mapstructure:"masking"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	for i, m := range cfg.Masking {
		if m.Regexp == "" {
			return fmt.Errorf("masking[%d]: regexp must not be empty", i)
		}
		switch m.Mode {
		case "", ModeReplace, ModePartial, ModeHash, ModeRedactKey:
			// valid
		default:
			return fmt.Errorf("masking[%d]: unsupported mode %q (valid: replace, partial, hash, redact_key)", i, m.Mode)
		}
	}
	return nil
}
