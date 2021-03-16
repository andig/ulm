package charger

// Code generated by github.com/andig/cmd/tools/decorate.go. DO NOT EDIT.

import (
	"github.com/mark-sch/evcc/api"
)

func decorateWarp(base *Warp, meter func() (float64, error), meterEnergy func() (float64, error)) api.Charger {
	switch {
	case meter == nil && meterEnergy == nil:
		return base

	case meter != nil && meterEnergy == nil:
		return &struct {
			*Warp
			api.Meter
		}{
			Warp: base,
			Meter: &decorateWarpMeterImpl{
				meter: meter,
			},
		}

	case meter == nil && meterEnergy != nil:
		return &struct {
			*Warp
			api.MeterEnergy
		}{
			Warp: base,
			MeterEnergy: &decorateWarpMeterEnergyImpl{
				meterEnergy: meterEnergy,
			},
		}

	case meter != nil && meterEnergy != nil:
		return &struct {
			*Warp
			api.Meter
			api.MeterEnergy
		}{
			Warp: base,
			Meter: &decorateWarpMeterImpl{
				meter: meter,
			},
			MeterEnergy: &decorateWarpMeterEnergyImpl{
				meterEnergy: meterEnergy,
			},
		}
	}

	return nil
}

type decorateWarpMeterImpl struct {
	meter func() (float64, error)
}

func (impl *decorateWarpMeterImpl) CurrentPower() (float64, error) {
	return impl.meter()
}

type decorateWarpMeterEnergyImpl struct {
	meterEnergy func() (float64, error)
}

func (impl *decorateWarpMeterEnergyImpl) TotalEnergy() (float64, error) {
	return impl.meterEnergy()
}
