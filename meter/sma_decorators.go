package meter

// This file has been generated - do not modify

import (
	"github.com/andig/evcc/api"
)

func decorateSMA(base api.Meter, meterEnergy func() (float64, error)) api.Meter {
	switch {
	case meterEnergy == nil:
		return base

	case meterEnergy != nil:
		return &struct{
			api.Meter
			api.MeterEnergy
		}{
			Meter: base,
			MeterEnergy: &decorateSMAMeterEnergyImpl{
				meterEnergy: meterEnergy,
			},
		}
	}

	return nil
}

type decorateSMAMeterEnergyImpl struct {
	meterEnergy func() (float64, error)
}

func (impl *decorateSMAMeterEnergyImpl) TotalEnergy() (float64, error) {
	return impl.meterEnergy()
}
