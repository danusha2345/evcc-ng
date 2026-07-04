package core

import (
	"errors"
	"fmt"
	"slices"

	"github.com/cenkalti/backoff/v4"
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/core/keys"
	"github.com/evcc-io/evcc/tariff"
	"github.com/evcc-io/evcc/util/config"
	"github.com/evcc-io/evcc/util/modbus"
	"github.com/samber/lo"
)

type circuitStruct struct {
	Title      string   `json:"title,omitempty"`
	Icon       string   `json:"icon,omitempty"`
	Parent     string   `json:"parent,omitempty"`
	Power      float64  `json:"power"`
	Current    *float64 `json:"current,omitempty"`
	MaxPower   float64  `json:"maxPower,omitempty"`
	MaxCurrent float64  `json:"maxCurrent,omitempty"`
}

// publishCircuits returns a list of circuit titles
func (site *Site) publishCircuits() {
	cc := config.Circuits().Devices()
	res := make(map[string]circuitStruct, len(cc))

	names := make(map[api.Circuit]string, len(cc))
	for _, c := range cc {
		names[c.Instance()] = c.Config().Name
	}

	for _, c := range cc {
		instance := c.Instance()
		props := deviceProperties(c)

		data := circuitStruct{
			Title:      instance.GetTitle(),
			Icon:       props.Icon,
			Parent:     names[instance.GetParent()],
			Power:      instance.GetChargePower(),
			MaxPower:   instance.GetMaxPower(),
			MaxCurrent: instance.GetMaxCurrent(),
		}

		if instance.GetMaxCurrent() > 0 {
			data.Current = lo.EmptyableToPtr(instance.GetMaxPhaseCurrent())
		}

		res[c.Config().Name] = data
	}

	site.publish(keys.Circuits, res)
}

func (site *Site) dimMeters(dim *bool) error {
	if dim == nil {
		return nil
	}

	var errs error
	for _, dev := range slices.Concat(site.auxMeters, site.extMeters) {
		m, ok := api.Cap[api.Dimmer](dev.Instance())
		if !ok {
			continue
		}

		if dimmed, err := backoff.RetryWithData(m.Dimmed, modbus.Backoff()); err == nil {
			if *dim == dimmed {
				continue
			}
		} else {
			if !errors.Is(err, api.ErrNotAvailable) {
				errs = errors.Join(errs, fmt.Errorf("%s dimmed: %w", deviceTitleOrName(dev), err))
			}
			continue
		}

		if err := m.Dim(*dim); err == nil {
			site.log.DEBUG.Printf("%s dim: %t", deviceTitleOrName(dev), *dim)
		} else if !errors.Is(err, api.ErrNotAvailable) {
			errs = errors.Join(errs, fmt.Errorf("%s dim: %w", deviceTitleOrName(dev), err))
		}
	}

	return errs
}

// shouldFeedInCurtail reports whether PV should be curtailed based on the
// current feed-in tariff. Returns nil ("not managed" — leave registers
// alone, e.g. an external Huawei 70% limit) unless feed-in control is
// enabled and a feed-in tariff is configured. Then it returns true/false
// depending on whether the current price is at or below the threshold —
// the user-facing promise "stop exporting once it costs me money"
// (evcc-io/evcc#21747). Tristate matches the upstream curtail API (#30116),
// which also subsumes our earlier #30068 gating.
func (site *Site) shouldFeedInCurtail() *bool {
	if !site.GetFeedInControl() {
		return nil
	}
	t := site.GetTariff(api.TariffUsageFeedIn)
	if t == nil {
		return nil
	}
	price, err := tariff.Now(t)
	if err != nil {
		return nil
	}
	res := price <= site.GetFeedInControlThreshold()
	return &res
}

func (site *Site) curtailPV(percent *int) error {
	if percent == nil {
		return nil
	}

	curtail := *percent < 100

	var errs error
	for _, dev := range site.pvMeters {
		m, ok := api.Cap[api.Curtailer](dev.Instance())
		if !ok {
			continue
		}

		if curtailed, err := backoff.RetryWithData(m.Curtailed, modbus.Backoff()); err == nil {
			if curtail == curtailed {
				continue
			}
		} else {
			if !errors.Is(err, api.ErrNotAvailable) {
				errs = errors.Join(errs, fmt.Errorf("%s curtailed: %w", deviceTitleOrName(dev), err))
			}
			continue
		}

		if err := m.SetCurtailPercent(*percent); err == nil {
			site.log.DEBUG.Printf("%s curtail: %d%%", deviceTitleOrName(dev), *percent)
		} else if !errors.Is(err, api.ErrNotAvailable) {
			errs = errors.Join(errs, fmt.Errorf("%s curtail: %w", deviceTitleOrName(dev), err))
		}
	}

	return errs
}
