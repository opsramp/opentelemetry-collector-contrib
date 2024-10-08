package scrubbingprocessor

import (
	"path/filepath"
	"testing"

	"github.com/open-telemetry/opentelemetry-collector-contrib/processor/scrubbingprocessor/internal/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/otelcol/otelcoltest"
)

func TestConfigLoad(t *testing.T) {
	factories, err := otelcoltest.NopFactories()
	assert.Nil(t, err)

	factory := NewFactory()
	factories.Processors[metadata.Type] = factory
	cfg, err := otelcoltest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)
	assert.Nil(t, err)
	require.NotNil(t, cfg)
}
