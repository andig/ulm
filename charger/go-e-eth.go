package charger

import (
	"fmt"
        "time"

	"github.com/andig/evcc/api"
	"github.com/andig/evcc/util"
	"github.com/andig/evcc/util/modbus"
	"github.com/volkszaehler/mbmd/meters/rs485"
)

const (
	goEEthRegStatus     = 100 // Input
	goEEthRegMaxCurrent = 299 // Holding
	goEEthRegEnable     = 200 // Holding

	goEEthRegPower  = 120 // power reading in 0.01kW
	goEEthRegEnergy = 132 // energy loaded in deka watt
)

var goEEthRegCurrents = []uint16{114, 116, 118} // current readings

// GoEEth is an api.ChargeController implementation for go-e wallboxes.
// It uses Modbus TCP to communicate with the wallbox at modbus client id 180.
type GoEEth struct {
	conn *modbus.Connection
}

func init() {
	registry.Add("go-e-eth", NewGoEEthFromConfig)
}

//go:generate go run ../cmd/tools/decorate.go -f decorateGoEEth -b *GoEEth -r api.Charger -t "api.Meter,CurrentPower,func() (float64, error)" -t "api.MeterEnergy,TotalEnergy,func() (float64, error)" -t "api.MeterCurrent,Currents,func() (float64, float64, float64, error)"

// NewGoEEthFromConfig creates a Go-E charger from generic config
func NewGoEEthFromConfig(other map[string]interface{}) (api.Charger, error) {
	cc := struct {
		URI   string
		ID    uint8
		Meter struct {
			Power, Energy, Currents bool
		}
	}{
		URI: "192.168.0.8:502", // default
		ID:  180,               // default
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	wb, err := NewGoEEth(cc.URI, cc.ID)

	var currentPower func() (float64, error)
	if cc.Meter.Power {
		currentPower = wb.currentPower
	}

	var totalEnergy func() (float64, error)
	if cc.Meter.Energy {
		totalEnergy = wb.totalEnergy
	}

	var currents func() (float64, float64, float64, error)
	if cc.Meter.Currents {
		currents = wb.currents
	}

	return decorateGoEEth(wb, currentPower, totalEnergy, currents), err
}

// NewGoEEth creates a Phoenix charger
func NewGoEEth(uri string, id uint8) (*GoEEth, error) {
	conn, err := modbus.NewConnection(uri, "", "", 0, false, id)
	if err != nil {
		return nil, err
	}

	log := util.NewLogger("em-eth")
	conn.Logger(log.TRACE)

	wb := &GoEEth{
		conn: conn,
	}

	return wb, nil
}

// Status implements the api.Charger interface
func (wb *GoEEth) Status() (api.ChargeStatus, error) {
	b, err := wb.conn.ReadInputRegisters(goEEthRegStatus, 1)
	if err != nil {
		return api.StatusNone, err
	}

	switch b[1] {
	case 0: // unkown, defect
            return api.StatusF, nil
        case 1: // WB ready, no car
	    return api.StatusA, nil
        case 2: // car charging
            return api.StatusC, nil
        case 3: // WB ready, waiting for car
            return api.StatusB, nil
        case 4: // charging stopped, car still connected
            return api.StatusB, nil
        default:
            return api.StatusF, fmt.Errorf("CAR_STATE unknown result: %d", b[1])
        }
}

// Enabled implements the api.Charger interface
func (wb *GoEEth) Enabled() (bool, error) {
	b, err := wb.conn.ReadInputRegisters(goEEthRegEnable, 1)
	if err != nil {
		return false, err
	}

	return b[1] == 1, nil
}

// Enable implements the api.Charger interface
func (wb *GoEEth) Enable(enable bool) error {
	var u uint16
	if enable {
		u = 0x0001
	}

	_, err := wb.conn.WriteSingleRegister(goEEthRegEnable, u)
        //let charger settle after update
        defer time.Sleep(2 * time.Second) 

	return err
}

// MaxCurrent implements the api.Charger interface
func (wb *GoEEth) MaxCurrent(current int64) error {
	if current < 6 {
		return fmt.Errorf("invalid current %d", current)
	}

	_, err := wb.conn.WriteSingleRegister(goEEthRegMaxCurrent, uint16(current))

	return err
}

// CurrentPower implements the api.Meter interface
func (wb *GoEEth) currentPower() (float64, error) {
	b, err := wb.conn.ReadInputRegisters(goEEthRegPower, 2)
	if err != nil {
		return 0, err
	}

	return rs485.RTUUint32ToFloat64Swapped(b) / 100, err
}

// totalEnergy implements the api.MeterEnergy interface
func (wb *GoEEth) totalEnergy() (float64, error) {
	b, err := wb.conn.ReadInputRegisters(goEEthRegEnergy, 2)
	if err != nil {
		return 0, err
	}

	return rs485.RTUUint32ToFloat64Swapped(b) / (60 * 60 * 100) , err //Deka Watt -> kwh
}

// currents implements the api.MeterCurrent interface
func (wb *GoEEth) currents() (float64, float64, float64, error) {
	var currents []float64
	for _, regCurrent := range goEEthRegCurrents {
		b, err := wb.conn.ReadInputRegisters(regCurrent, 2)
		if err != nil {
			return 0, 0, 0, err
		}

		currents = append(currents, rs485.RTUUint32ToFloat64Swapped(b)/10)
	}

	return currents[0], currents[1], currents[2], nil
}
