package meter

import (
	"errors"
	"testing"

	"github.com/evcc-io/evcc/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapper(t *testing.T) {
	origErr := errors.New("modbus timeout")
	m := NewWrapper("grid", "sma", map[string]any{"uri": "tcp://1.2.3.4"}, origErr)

	t.Run("CurrentPower surfaces error and zero value", func(t *testing.T) {
		p, err := m.CurrentPower()
		require.Error(t, err)
		assert.ErrorIs(t, err, origErr)
		assert.Equal(t, 0.0, p)
	})

	t.Run("Features advertises Offline and Retryable", func(t *testing.T) {
		fd, ok := m.(api.FeatureDescriber)
		require.True(t, ok)
		feats := fd.Features()
		assert.Contains(t, feats, api.Offline)
		assert.Contains(t, feats, api.Retryable)
	})

	t.Run("WrappedConfig round-trips type and config", func(t *testing.T) {
		ww, ok := m.(interface {
			WrappedConfig() (string, map[string]any)
		})
		require.True(t, ok)
		typ, cfg := ww.WrappedConfig()
		assert.Equal(t, "sma", typ)
		assert.Equal(t, "tcp://1.2.3.4", cfg["uri"])
	})
}
