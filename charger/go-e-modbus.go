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
	GoEModbusRegStatus     = 100 // Input
	GoEModbusRegMaxCurrent = 299 // Holding
	GoEModbusRegEnable     = 200 // Holding

	GoEModbusRegPower  = 120 // power reading in 0.01kW
	GoEModbusRegEnergy = 132 // energy loaded in deka watt
)

var GoEModbusRegCurrents = []uint16{114, 116, 118} // current readings

// GoEModbus is an api.ChargeController implementation for go-e wallboxes.
// It uses Modbus TCP to communicate with the wallbox at modbus client id 180.
type GoEModbus struct {
	conn *modbus.Connection
}

func init() {
	registry.Add("go-e-modbus", NewGoEModbusFromConfig)
}

// NewGoEModbusFromConfig creates a Go-E charger from generic config
func NewGoEModbusFromConfig(other map[string]interface{}) (api.Charger, error) {
	cc := struct {
		URI string
		ID  uint8
	}{
		URI: "192.168.0.8:502", // default
		ID:  180,               // default
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	return NewGoEModbus(cc.URI, cc.ID)
}

// NewGoEModbus creates a Phoenix charger
func NewGoEModbus(uri string, id uint8) (*GoEModbus, error) {
	conn, err := modbus.NewConnection(uri, "", "", 0, false, id)
	if err != nil {
		return nil, err
	}

	log := util.NewLogger("em-eth")
	conn.Logger(log.TRACE)

	wb := &GoEModbus{
		conn: conn,
	}

	return wb, nil
}

// Status implements the api.Charger interface
func (wb *GoEModbus) Status() (api.ChargeStatus, error) {
	b, err := wb.conn.ReadInputRegisters(GoEModbusRegStatus, 1)
	if err != nil {
		return api.StatusNone, err
	}

	switch b[1] {
	case 1: // WB ready, no car
		return api.StatusA, nil
	case 2: // car charging
		return api.StatusC, nil
	case 3, 4: // 3 -> WB ready, waiting for car, 4 -> charging stopped, car still connected
		return api.StatusB, nil
	default:
		return api.StatusF, fmt.Errorf("CAR_STATE unknown result: %d", b[1])
	}
}

// Enabled implements the api.Charger interface
func (wb *GoEModbus) Enabled() (bool, error) {
	b, err := wb.conn.ReadInputRegisters(GoEModbusRegEnable, 1)
	if err != nil {
		return false, err
	}

	return b[1] == 1, nil
}

// Enable implements the api.Charger interface
func (wb *GoEModbus) Enable(enable bool) error {
	var u uint16
	if enable {
		u = 0x0001
	}

	_, err := wb.conn.WriteSingleRegister(GoEModbusRegEnable, u)
	//let charger settle after update
	defer time.Sleep(2 * time.Second)

	return err
}

// MaxCurrent implements the api.Charger interface
func (wb *GoEModbus) MaxCurrent(current int64) error {
	if current < 6 {
		return fmt.Errorf("invalid current %d", current)
	}

	_, err := wb.conn.WriteSingleRegister(GoEModbusRegMaxCurrent, uint16(current))

	return err
}

// CurrentPower implements the api.Meter interface
func (wb *GoEModbus) CurrentPower() (float64, error) {
	b, err := wb.conn.ReadInputRegisters(GoEModbusRegPower, 2)
	if err != nil {
		return 0, err
	}

	return rs485.RTUUint32ToFloat64Swapped(b) / 100, err
}

// TotalEnergy implements the api.MeterEnergy interface
func (wb *GoEModbus) TotalEnergy() (float64, error) {
	b, err := wb.conn.ReadInputRegisters(GoEModbusRegEnergy, 2)
	if err != nil {
		return 0, err
	}

	return rs485.RTUUint32ToFloat64Swapped(b) / (60 * 60 * 100), err //Deka Watt -> kwh
}

// Currents implements the api.MeterCurrent interface
func (wb *GoEModbus) Currents() (float64, float64, float64, error) {
	var currents []float64
	for _, regCurrent := range GoEModbusRegCurrents {
		b, err := wb.conn.ReadInputRegisters(regCurrent, 2)
		if err != nil {
			return 0, 0, 0, err
		}

		currents = append(currents, rs485.RTUUint32ToFloat64Swapped(b)/10)
	}

	return currents[0], currents[1], currents[2], nil
}
