package charger

import (
	"errors"
	"testing"

	"github.com/evcc-io/evcc/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrapper_AllMethodsReturnOriginalError(t *testing.T) {
	origErr := errors.New("connection refused")
	w := NewWrapper("test", "kebawh", map[string]any{"host": "1.2.3.4"}, origErr)

	t.Run("Status", func(t *testing.T) {
		s, err := w.Status()
		require.Error(t, err)
		assert.ErrorIs(t, err, origErr)
		assert.Equal(t, api.StatusNone, s)
	})

	t.Run("Enabled", func(t *testing.T) {
		e, err := w.Enabled()
		require.Error(t, err)
		assert.False(t, e)
	})

	t.Run("Enable", func(t *testing.T) {
		assert.ErrorIs(t, w.Enable(true), origErr)
	})

	t.Run("MaxCurrent", func(t *testing.T) {
		assert.ErrorIs(t, w.MaxCurrent(16), origErr)
	})

	t.Run("Features advertises Offline and Retryable", func(t *testing.T) {
		ww, ok := w.(api.FeatureDescriber)
		require.True(t, ok)
		feats := ww.Features()
		assert.Contains(t, feats, api.Offline)
		assert.Contains(t, feats, api.Retryable)
	})

	t.Run("WrappedConfig returns original type and config", func(t *testing.T) {
		ww, ok := w.(interface {
			WrappedConfig() (string, map[string]any)
		})
		require.True(t, ok)
		typ, cfg := ww.WrappedConfig()
		assert.Equal(t, "kebawh", typ)
		assert.Equal(t, "1.2.3.4", cfg["host"])
	})
}
