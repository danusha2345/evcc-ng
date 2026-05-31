package vehicle

import (
	"fmt"
	"strings"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/util"
)

// Wrapper wraps an api.Vehicle to capture initialization errors
type Wrapper struct {
	embed
	typ    string
	config map[string]any
	err    error
}

// NewWrapper creates an offline Vehicle wrapper
func NewWrapper(name, typ string, other map[string]any, err error) api.Vehicle {
	return newWrapper(name, typ, other, fmt.Errorf("vehicle not available: %w", err), false)
}

// NewDisabledWrapper creates a quiet offline stub for a vehicle the user
// intentionally disabled. Unlike NewWrapper it is not Retryable, so the retry
// loop never revives it — re-enabling is a manual action (evcc-io/evcc#21144).
func NewDisabledWrapper(name, typ string, other map[string]any) api.Vehicle {
	return newWrapper(name, typ, other, api.ErrNotAvailable, true)
}

func newWrapper(name, typ string, other map[string]any, err error, disabled bool) api.Vehicle {
	var cc struct {
		embed `mapstructure:",squash"`
		Other map[string]any `mapstructure:",remain"`
	}

	// try to decode vehicle-specific config and look for title attribute
	_ = util.DecodeOther(other, &cc)

	if cc.Title_ == "" {
		//lint:ignore SA1019 as Title is safe on ascii
		cc.Title_ = strings.Title(name)
	}

	v := &Wrapper{
		embed:  cc.embed,
		typ:    typ,
		config: cc.Other,
		err:    err,
	}

	// disabled stubs are deliberately not Retryable so the retry loop leaves
	// them offline until the user re-enables them (evcc-io/evcc#21144)
	v.Features_ = append(v.Features_, api.Offline)
	if !disabled {
		v.Features_ = append(v.Features_, api.Retryable)
	}
	v.SetTitle(cc.Title_)

	return v
}

// WrappedConfig indicates a device with wrapped configuration
func (v *Wrapper) WrappedConfig() (string, map[string]any) {
	return v.typ, v.config
}

// SetTitle implements the api.TitleSetter interface
func (v *Wrapper) SetTitle(title string) {
	v.Title_ = title
}

var _ api.Battery = (*Wrapper)(nil)

// Soc implements the api.Battery interface
func (v *Wrapper) Soc() (float64, error) {
	return 0, v.err
}
