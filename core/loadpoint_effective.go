package core

import (
	"time"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/core/keys"
	"github.com/evcc-io/evcc/core/vehicle"
)

// publishEffectiveValues publishes all effective values
func (lp *Loadpoint) publishEffectiveValues() {
	lp.publish(keys.EffectivePriority, lp.EffectivePriority())
	lp.publish(keys.EffectivePlanTime, lp.EffectivePlanTime())
	lp.publish(keys.EffectiveMinCurrent, lp.effectiveMinCurrent())
	lp.publish(keys.EffectiveMaxCurrent, lp.effectiveMaxCurrent())
	lp.publish(keys.EffectiveLimitSoc, lp.effectiveLimitSoc())
}

// EffectivePriority returns the effective priority
func (lp *Loadpoint) EffectivePriority() int {
	if v := lp.GetVehicle(); v != nil {
		if res, ok := v.OnIdentified().GetPriority(); ok {
			return res
		}
	}
	return lp.GetPriority()
}

// EffectivePlanTime returns the effective plan time
func (lp *Loadpoint) EffectivePlanTime() time.Time {
	if v := lp.GetVehicle(); v != nil {
		if res := vehicle.Settings(lp.log, v).GetPlanTime(); !res.IsZero() {
			lp.publish(keys.EffectivePlanTime, res)
			return res
		}
	}
	return lp.GetPlanTime()
}

// effectiveMinCurrent returns the effective min current
func (lp *Loadpoint) effectiveMinCurrent() float64 {
	if v := lp.GetVehicle(); v != nil {
		if res, ok := v.OnIdentified().GetMinCurrent(); ok {
			return res
		}
	}

	if c, ok := lp.charger.(api.CurrentLimiter); ok {
		if res, _, err := c.GetMinMaxCurrent(); err == nil {
			lp.publish(keys.EffectiveMinCurrent, res)
			return res
		}
	}

	return lp.GetMinCurrent()
}

// effectiveMaxCurrent returns the effective max current
func (lp *Loadpoint) effectiveMaxCurrent() float64 {
	if v := lp.GetVehicle(); v != nil {
		if res, ok := v.OnIdentified().GetMaxCurrent(); ok {
			return res
		}
	}

	if c, ok := lp.charger.(api.CurrentLimiter); ok {
		if _, res, err := c.GetMinMaxCurrent(); err == nil {
			lp.publish(keys.EffectiveMaxCurrent, res)
			return res
		}
	}

	return lp.GetMaxCurrent()
}

// effectiveLimitSoc returns the effective session limit soc
// TODO take vehicle api limits into account
func (lp *Loadpoint) effectiveLimitSoc() int {
	lp.RLock()
	defer lp.RUnlock()

	if lp.limitSoc > 0 {
		return lp.limitSoc
	}

	if v := lp.GetVehicle(); v != nil {
		if soc := vehicle.Settings(lp.log, v).GetLimitSoc(); soc > 0 {
			return soc
		}
	}

	return 100
}

// EffectiveMinPower returns the effective min power for a single phase
func (lp *Loadpoint) EffectiveMinPower() float64 {
	// TODO check if 1p available
	return Voltage * lp.effectiveMinCurrent()
}

// EffectiveMaxPower returns the effective max power taking vehicle capabilities and phase scaling into account
func (lp *Loadpoint) EffectiveMaxPower() float64 {
	return Voltage * lp.effectiveMaxCurrent() * float64(lp.maxActivePhases())
}
