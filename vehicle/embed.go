package vehicle

import (
	"github.com/evcc-io/evcc/api"
)

// TODO align phases with OnIdentify
type embed struct {
	Title_       string           `mapstructure:"title"`
	Icon_        string           `mapstructure:"icon"`
	Capacity_    float64          `mapstructure:"capacity"`
	Phases_      int              `mapstructure:"phases"`
	FixedPhases_ int              `mapstructure:"fixedphases"` // 0=auto, 1=force 1p, 3=force 3p (issue #30705)
	Identifiers_ []string         `mapstructure:"identifiers"`
	Features_    []api.Feature    `mapstructure:"features"`
	OnIdentify   api.ActionConfig `mapstructure:"onIdentify"`
}

// Title implements the api.Vehicle interface
func (v *embed) fromVehicle(title string, capacity float64) {
	if v.Title_ == "" {
		v.Title_ = title
	}
	if v.Capacity_ == 0 {
		v.Capacity_ = capacity
	}
}

// GetTitle implements the api.Vehicle interface
func (v *embed) GetTitle() string {
	return v.Title_
}

// SetTitle implements the api.TitleSetter interface
func (v *embed) SetTitle(title string) {
	v.Title_ = title
}

// Capacity implements the api.Vehicle interface
func (v *embed) Capacity() float64 {
	return v.Capacity_
}

var _ api.PhaseDescriber = (*embed)(nil)

// Phases returns the phases used by the vehicle
func (v *embed) Phases() int {
	return v.Phases_
}

var _ api.PhaseConfigurer = (*embed)(nil)

// PhasesConfigured returns a fixed phase-count override for charging this
// vehicle (0 = auto/dynamic switching, the default). Implements api.PhaseConfigurer.
func (v *embed) PhasesConfigured() int {
	return v.FixedPhases_
}

// Identifiers implements the api.Identifier interface
func (v *embed) Identifiers() []string {
	return v.Identifiers_
}

// OnIdentified returns the identify action
func (v *embed) OnIdentified() api.ActionConfig {
	return v.OnIdentify
}

var _ api.IconDescriber = (*embed)(nil)

// Icon implements the api.IconDescriber interface
func (v *embed) Icon() string {
	return v.Icon_
}

var _ api.FeatureDescriber = (*embed)(nil)

// Features implements the api.FeatureDescriber interface
func (v *embed) Features() []api.Feature {
	return v.Features_
}
