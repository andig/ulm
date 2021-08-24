package meter

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/meter/lgpcs"
	"github.com/evcc-io/evcc/provider"
	"github.com/evcc-io/evcc/util"
)

/**
 * This meter supports the LGESS HOME 8 and LGESS HOME 10 systems from LG with / without battery.
 *
 ** Usages **
 * The following usages are supported:
 * - grid    ... for reading the power imported or exported to the grid
 * - pv      ... for reading the power produced by the Photovoltaik
 * - battery ... for reading the power imported or exported to the battery
 *  *
 ** Example configuration **
 *
 * meters:
 * - name: GridMeter
 *   type: lgess
 *   usage: grid
 *   uri: https://192.168.1.23
 *   password: "DE200....."
 * - name: PvMeter
 *   type: lgess
 *   usage: pv
 * - name: BatteryMeter
 *   type: lgess
 *   usage: battery
 *
 *
 ** Limitations **
 * It is not allowed to provide different URIs or passwords for different lgess meters since always the
 * same hardware instance is accessed with the different usages.
 *
 * */

// LgEss implements the api.Meter interface
type LgEss struct {
	usage string // grid, pv, battery
	essG  func() (interface{}, error)
}

func init() {
	registry.Add("lgess", NewLgEssFromConfig)
}

//go:generate go run ../cmd/tools/decorate.go -f decorateLgEss -b api.Meter -t "api.MeterEnergy,TotalEnergy,func() (float64, error)" -t "api.Battery,SoC,func() (float64, error)"

// NewLgEssFromConfig creates an LgEss Meter from generic config
func NewLgEssFromConfig(other map[string]interface{}) (api.Meter, error) {
	cc := struct {
		URI, Usage, Password string
		Cache                time.Duration
	}{
		Cache: 2 * time.Second,
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	if cc.Usage == "" {
		return nil, errors.New("missing usage")
	}

	return NewLgEss(cc.URI, cc.Usage, cc.Password, cc.Cache)
}

// NewLgEss creates an LgEss Meter
func NewLgEss(uri, usage, password string, cache time.Duration) (api.Meter, error) {
	lp, err := lgpcs.GetInstance(uri, password)
	if err != nil {
		return nil, err
	}

	m := &LgEss{
		usage: strings.ToLower(usage),
	}

	m.essG = provider.NewCached(func() (interface{}, error) {
		return lp.Data()
	}, cache).InterfaceGetter()

	// decorate api.MeterEnergy
	var totalEnergy func() (float64, error)
	if m.usage == "grid" {
		totalEnergy = m.totalEnergy
	}

	// decorate api.BatterySoC
	var batterySoC func() (float64, error)
	if usage == "battery" {
		batterySoC = m.batterySoC
	}

	return decorateLgEss(m, totalEnergy, batterySoC), nil
}

// CurrentPower implements the api.Meter interface
func (m *LgEss) CurrentPower() (float64, error) {
	res, err := m.essG()
	if err != nil {
		return 0, err
	}

	data := res.(lgpcs.EssData)

	switch m.usage {
	case "grid":
		return data.GridPower, nil
	case "pv":
		return data.PvTotalPower, nil
	case "battery":
		return data.BatConvPower, nil
	}

	return 0, fmt.Errorf("invalid usage: %s", m.usage)
}

// totalEnergy implements the api.MeterEnergy interface
func (m *LgEss) totalEnergy() (float64, error) {
	res, err := m.essG()
	if err != nil {
		return 0, err
	}

	data := res.(lgpcs.EssData)

	switch m.usage {
	case "grid":
		return data.CurrentGridFeedInEnergy / 1e3, nil
	case "pv":
		return data.CurrentPvGenerationSum / 1e3, nil
	}

	return 0, fmt.Errorf("invalid usage: %s", m.usage)
}

// batterySoC implements the api.Battery interface
func (m *LgEss) batterySoC() (float64, error) {
	res, err := m.essG()
	if err != nil {
		return 0, err
	}

	data := res.(lgpcs.EssData)

	return data.BatUserSoc, nil
}
