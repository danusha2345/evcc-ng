package charger

import (
	"fmt"

	"github.com/evcc-io/evcc/api"
)

// Wrapper holds a charger that failed to initialize, so the rest of the
// site can still come up. Every method returns the original error wrapped
// with a recognisable prefix; loadpoint.AcceptableError treats
// api.ErrNotAvailable specially, but here we surface the *real* cause so
// users see why the device is offline. Features() advertises Offline +
// Retryable so the rest of the system can skip / retry as appropriate
// (evcc-io/evcc#14496).
type Wrapper struct {
	name     string
	typ      string
	config   map[string]any
	err      error
	disabled bool // intentionally disabled by the user — not retryable
}

// NewWrapper creates an offline Charger wrapper that holds the init error.
func NewWrapper(name, typ string, other map[string]any, err error) api.Charger {
	return &Wrapper{
		name:   name,
		typ:    typ,
		config: other,
		err:    fmt.Errorf("charger not available: %w", err),
	}
}

// NewDisabledWrapper creates a quiet offline stub for a device the user
// intentionally disabled. Unlike NewWrapper it is not Retryable, so the retry
// loop never revives it — re-enabling is a manual action (evcc-io/evcc#21144).
func NewDisabledWrapper(name, typ string, other map[string]any) api.Charger {
	return &Wrapper{
		name:     name,
		typ:      typ,
		config:   other,
		err:      api.ErrNotAvailable,
		disabled: true,
	}
}

// WrappedConfig exposes the underlying config so a retry loop / UI can
// re-instantiate the device when the user fixes it.
func (w *Wrapper) WrappedConfig() (string, map[string]any) {
	return w.typ, w.config
}

var _ api.Charger = (*Wrapper)(nil)

// Status implements the api.ChargeState interface
func (w *Wrapper) Status() (api.ChargeStatus, error) {
	return api.StatusNone, w.err
}

// Enabled implements the api.Charger interface
func (w *Wrapper) Enabled() (bool, error) {
	return false, w.err
}

// Enable implements the api.Charger interface
func (w *Wrapper) Enable(enable bool) error {
	return w.err
}

// MaxCurrent implements the api.CurrentController interface
func (w *Wrapper) MaxCurrent(current int64) error {
	return w.err
}

var _ api.FeatureDescriber = (*Wrapper)(nil)

// Features implements the api.FeatureDescriber interface
func (w *Wrapper) Features() []api.Feature {
	if w.disabled {
		return []api.Feature{api.Offline}
	}
	return []api.Feature{api.Offline, api.Retryable}
}
