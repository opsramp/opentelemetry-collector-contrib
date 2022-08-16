package stream_processor

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/service/servicetest"
)

func TestConfigLoad(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.Nil(t, err)

	factory := NewFactory()
	factories.Processors[typeStr] = factory
	cfg, err := servicetest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)
	assert.Nil(t, err)
	require.NotNil(t, cfg)
	id := config.NewComponentIDWithName(typeStr, "")
	component := cfg.Processors[id]
	assert.Equal(t, component.(*Config).Query, "select * where price = 5;")

}
