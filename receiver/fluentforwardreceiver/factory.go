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

package fluentforwardreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/fluentforwardreceiver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/consumer"
)

const (
	// The value of "type" key in configuration.
	typeStr = "fluentforward"
	// The stability level of the receiver.
	stability = component.StabilityLevelBeta
)

// NewFactory return a new component.ReceiverFactory for fluentd forwarder.
func NewFactory() component.ReceiverFactory {
	return component.NewReceiverFactory(
		typeStr,
		createDefaultConfig,
		component.WithLogsReceiver(createLogsReceiver, stability))
}

func createDefaultConfig() component.Config {
	return &Config{
		ReceiverSettings: config.NewReceiverSettings(component.NewID(typeStr)),
	}
}

func createLogsReceiver(
	_ context.Context,
	params component.ReceiverCreateSettings,
	cfg component.Config,
	consumer consumer.Logs,
) (component.LogsReceiver, error) {

	rCfg := cfg.(*Config)
	return newFluentReceiver(params.Logger, rCfg, consumer)
}
