package loadpoint

import (
	"time"
)

// ThresholdsConfig defines pv mode hysteresis parameters
type ThresholdsConfig struct {
	Enable  ThresholdConfig `json:"enable"`
	Disable ThresholdConfig `json:"disable"`
}

// ThresholdConfig defines enable/disable hysteresis parameters
type ThresholdConfig struct {
	Delay     time.Duration `json:"delay"`
	Threshold float64       `json:"threshold"`
	// ImmediateThreshold, when > 0, bypasses Delay for the disable path: PV
	// charging stops right away once grid import reaches this (higher) power,
	// while smaller overshoots still ride out Delay (evcc-io/evcc#31061).
	ImmediateThreshold float64 `json:"immediateThreshold"`
}

// SocConfig defines soc settings, estimation and update behavior
type SocConfig struct {
	Poll     PollConfig `json:"poll"`
	Estimate *bool      `json:"estimate"`
}

// PollConfig defines the vehicle polling mode and interval
type PollConfig struct {
	Mode     PollMode      `json:"mode"`     // polling mode charging (default), connected, always
	Interval time.Duration `json:"interval"` // interval when not charging
}

//go:generate go tool enumer -type PollMode -trimprefix Poll -transform=lower -text
type PollMode int

// Poll modes
const (
	PollCharging PollMode = iota
	PollConnected
	PollAlways
)
