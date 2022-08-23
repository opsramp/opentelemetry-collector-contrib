package stream_processor

import (
	"fmt"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
	"go.opentelemetry.io/collector/config"
)

// Config defines configuration for stream processor.
type Config struct {
	config.ProcessorSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
	adapter.BaseConfig       `mapstructure:",squash"`
	Query                    string `mapstructure:"query"`
}

var _ config.Processor = (*Config)(nil)

// Validate checks if the processor configuration is valid
func (cfg *Config) Validate() error {
	if len(cfg.Query) == 0 {
		return fmt.Errorf("correct sql query missed")
	}

	return nil
}
