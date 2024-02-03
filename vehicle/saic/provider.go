package saic

import (
	"time"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/provider"
	"github.com/evcc-io/evcc/vehicle/saic/requests"
)

// https://github.com/SAIC-iSmart-API/reverse-engineering

// Provider implements the vehicle api
type Provider struct {
	statusG func() (requests.ChargeStatus, error)
	purge   func()
	Vin     string
	Api     *API
}

// NewProvider creates a vehicle api provider
func NewProvider(api *API, vin string, cache time.Duration) *Provider {
	impl := &Provider{
		Vin: vin,
		Api: api,
	}
	c := provider.ResettableCached(func() (requests.ChargeStatus, error) {
		return api.Status(vin)
	}, cache)

	impl.statusG = c.Get
	impl.purge = c.Reset
	return impl
}

var _ api.Battery = (*Provider)(nil)

// Soc implements the api.Vehicle interface
func (v *Provider) Soc() (float64, error) {
	res, err := v.statusG()
	if err != nil {
		return 0, err
	}

	val := res.ChrgMgmtData.BmsPackSOCDsp
	if val > 1000 {
		v.Api.Logger.ERROR.Printf("Invalid raw SOC value: %d", val)
		v.purge()
		return float64(val), api.ErrMustRetry
	}

	return float64(val) / 10.0, nil
}

var _ api.ChargeState = (*Provider)(nil)

// Status implements the api.ChargeState interface
func (v *Provider) Status() (api.ChargeStatus, error) {
	res, err := v.statusG()
	if err != nil {
		return api.StatusNone, err
	}

	status := api.StatusA // disconnected
	if res.RvsChargeStatus.ChargingGunState != 0 {
		if (res.ChrgMgmtData.BmsChrgSts & 0x01) == 0 {
			status = api.StatusB
		} else {
			status = api.StatusC
		}
	}

	return status, nil
}

var _ api.VehicleFinishTimer = (*Provider)(nil)

// FinishTime implements the api.VehicleFinishTimer interface
func (v *Provider) FinishTime() (time.Time, error) {
	res, err := v.statusG()
	if err == nil {
		ctr := res.ChrgMgmtData.ChrgngRmnngTime
		return time.Now().Add(time.Duration(ctr) * time.Minute), err
	}

	return time.Time{}, err
}

var _ api.VehicleRange = (*Provider)(nil)

// Range implements the api.VehicleRange interface
func (v *Provider) Range() (int64, error) {
	res, err := v.statusG()
	if err != nil {
		return 0, err
	}
	val := res.RvsChargeStatus.FuelRangeElec
	if val < 10 {
		// Ok, 0 would be possible, but it's more likely that it's an invalid answer.
		v.Api.Logger.WARN.Printf("Suspicous raw RANGE value: %d", val)
		return val, api.ErrMustRetry
	}
	return val / 10, nil
}

var _ api.VehicleOdometer = (*Provider)(nil)

// Odometer implements the api.VehicleOdometer interface
func (v *Provider) Odometer() (float64, error) {
	res, err := v.statusG()
	if err != nil {
		return 0, err
	}

	return float64(res.RvsChargeStatus.Mileage), nil
}

var _ api.Resurrector = (*Provider)(nil)

func (v *Provider) WakeUp() error {
	err := v.Api.Wakeup(v.Vin)
	return err
}
