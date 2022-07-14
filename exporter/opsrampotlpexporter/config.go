// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package opsrampotlpexporter // import "go.opentelemetry.io/collector/exporter/otlpexporter"

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

type MaskingSettings struct {
	Regexp      string `mapstructure:"regexp"`
	Placeholder string `mapstructure:"placeholder"`
}

type SecuritySettings struct {
	OAuthServiceURL string `mapstructure:"oauth_service_url"`
	ClientId        string `mapstructure:"client_id"`
	ClientSecret    string `mapstructure:"client_secret"`
	GrantType       string `mapstructure:"grant_type"`
}

func (s *SecuritySettings) Validate() error {
	if len(s.OAuthServiceURL) == 0 {
		return errors.New("oauth service url missed")
	}

	if len(s.ClientId) == 0 {
		return errors.New("client_id missed")
	}

	if len(s.ClientSecret) == 0 {
		return errors.New("client_secret missed")
	}

	if len(s.GrantType) == 0 {
		return errors.New("grant_type missed")
	}

	return nil
}

// Config defines configuration for OpenCensus exporter.
type Config struct {
	config.ExporterSettings        `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct
	exporterhelper.TimeoutSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`

	Security                      SecuritySettings         `mapstructure:"security"`
	configgrpc.GRPCClientSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	Masking                       []MaskingSettings        `mapstructure:"masking"`
}

var _ config.Exporter = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {
	if err := cfg.QueueSettings.Validate(); err != nil {
		return fmt.Errorf("queue settings has invalid configuration: %w", err)
	}

	if err := cfg.Security.Validate(); err != nil {
		return fmt.Errorf("security settings has invalid configuration: %w", err)
	}

	return nil
}
