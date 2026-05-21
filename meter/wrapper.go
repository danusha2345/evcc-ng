package meter

import (
	"fmt"

	"github.com/evcc-io/evcc/api"
)

// Wrapper holds a meter that failed to initialize, so the rest of the
// site can still come up. CurrentPower returns the original error wrapped
// with a recognisable prefix; Features() advertises Offline + Retryable
// so the rest of the system can skip / retry as appropriate
// (evcc-io/evcc#14496).
type Wrapper struct {
	name   string
	typ    string
	config map[string]any
	err    error
}

// NewWrapper creates an offline Meter wrapper that holds the init error.
func NewWrapper(name, typ string, other map[string]any, err error) api.Meter {
	return &Wrapper{
		name:   name,
		typ:    typ,
		config: other,
		err:    fmt.Errorf("meter not available: %w", err),
	}
}

// WrappedConfig exposes the underlying config so a retry loop / UI can
// re-instantiate the device when the user fixes it.
func (w *Wrapper) WrappedConfig() (string, map[string]any) {
	return w.typ, w.config
}

var _ api.Meter = (*Wrapper)(nil)

// CurrentPower implements the api.Meter interface
func (w *Wrapper) CurrentPower() (float64, error) {
	return 0, w.err
}

var _ api.FeatureDescriber = (*Wrapper)(nil)

// Features implements the api.FeatureDescriber interface
func (w *Wrapper) Features() []api.Feature {
	return []api.Feature{api.Offline, api.Retryable}
}
