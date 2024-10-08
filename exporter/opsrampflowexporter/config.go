package opsrampflowexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/opsrampflowexporter"

import (
	"errors"

	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for opsramp flow exporter.
type Config struct {
	configretry.BackOffConfig      `mapstructure:"retry_on_failure"`
	exporterhelper.TimeoutSettings `mapstructure:",squash"`

	Endpoint  string `mapstructure:"endpoint"`
	AuthToken string `mapstructure:"authToken"`
}

var (
	errConfigNoEndpoint = errors.New("endpoint must be specified")
	errConfigNoToken    = errors.New("token must be specified")
)s

// Validate validates the opensearch server configuration.
func (cfg *Config) Validate() error {
	var multiErr []error

	if cfg.Endpoint == "" {
		multiErr = append(multiErr, errConfigNoEndpoint)
	}

	if cfg.AuthToken == "" {
		multiErr = append(multiErr, errConfigNoToken)
	}

	return errors.Join(multiErr...)
}
