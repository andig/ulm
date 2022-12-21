package meter

import (
	"errors"
	"fmt"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/provider"
	"github.com/evcc-io/evcc/util"
)

func init() {
	registry.Add(api.Custom, NewConfigurableFromConfig)
}

//go:generate go run ../cmd/tools/decorate.go -f decorateMeter -b api.Meter -t "api.MeterEnergy,TotalEnergy,func() (float64, error)" -t "api.MeterCurrent,Currents,func() (float64, float64, float64, error)" -t "api.Battery,SoC,func() (float64, error)" -t "api.BatteryCapacity,Capacity,func() (float64, error)"

// NewConfigurableFromConfig creates api.Meter from config
func NewConfigurableFromConfig(other map[string]interface{}) (api.Meter, error) {
	var cc struct {
		Capacity_ *float64 `mapstructure:"capacity"`
		Power     provider.Config
		Energy    *provider.Config  // optional
		SoC       *provider.Config  // optional
		Currents  []provider.Config // optional
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	power, err := provider.NewFloatGetterFromConfig(cc.Power)
	if err != nil {
		return nil, fmt.Errorf("power: %w", err)
	}

	m, _ := NewConfigurable(cc.Capacity_, power)

	// decorate Meter with MeterEnergy
	var totalEnergyG func() (float64, error)
	if cc.Energy != nil {
		totalEnergyG, err = provider.NewFloatGetterFromConfig(*cc.Energy)
		if err != nil {
			return nil, fmt.Errorf("energy: %w", err)
		}
	}

	// decorate Meter with MeterCurrent
	var currentsG func() (float64, float64, float64, error)
	if len(cc.Currents) > 0 {
		if len(cc.Currents) != 3 {
			return nil, errors.New("need 3 currents")
		}

		var curr []func() (float64, error)
		for idx, cc := range cc.Currents {
			c, err := provider.NewFloatGetterFromConfig(cc)
			if err != nil {
				return nil, fmt.Errorf("currents[%d]: %w", idx, err)
			}

			curr = append(curr, c)
		}

		currentsG = collectCurrentProviders(curr)
	}

	// decorate Meter with BatterySoC
	var batterySoCG func() (float64, error)
	if cc.SoC != nil {
		batterySoCG, err = provider.NewFloatGetterFromConfig(*cc.SoC)
		if err != nil {
			return nil, fmt.Errorf("battery: %w", err)
		}
	}

	res := m.Decorate(totalEnergyG, currentsG, batterySoCG)

	return res, nil
}

// collectCurrentProviders combines phase getters into currents api function
func collectCurrentProviders(g []func() (float64, error)) func() (float64, float64, float64, error) {
	return func() (float64, float64, float64, error) {
		var currents []float64
		for _, currentG := range g {
			c, err := currentG()
			if err != nil {
				return 0, 0, 0, err
			}

			currents = append(currents, c)
		}

		return currents[0], currents[1], currents[2], nil
	}
}

// NewConfigurable creates a new meter
func NewConfigurable(capacity *float64, currentPowerG func() (float64, error)) (*Meter, error) {

	m := &Meter{
		capacity:      capacity,
		currentPowerG: currentPowerG,
	}
	return m, nil
}

// Meter is an api.Meter implementation with configurable getters and setters.
type Meter struct {
	capacity      *float64
	currentPowerG func() (float64, error)
}

// Decorate attaches additional capabilities to the base meter
func (m *Meter) Decorate(
	totalEnergy func() (float64, error),
	currents func() (float64, float64, float64, error),
	batterySoC func() (float64, error),
) api.Meter {

	var capacityG func() (float64, error)

	if _, ok := m.Capacity(); ok == nil {
		capacityG = func() (res float64, err error) {
			return m.Capacity()
		}
	} else {
		capacityG = nil
	}

	return decorateMeter(m, totalEnergy, currents, batterySoC, capacityG)
}

// CurrentPower implements the api.Meter interface
func (m *Meter) CurrentPower() (float64, error) {
	return m.currentPowerG()
}

var _ api.BatteryCapacity = (*Meter)(nil)

// Capacity implements the api.BatteryCapacity interface
func (m *Meter) Capacity() (float64, error) {
	if m.capacity != nil {
		return *m.capacity, nil
	} else {
		return 0, api.ErrNotAvailable
	}
}
