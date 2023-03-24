package charger

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/modbus"
	"github.com/volkszaehler/mbmd/encoding"
)

const (
	// holding and input return same values
	charxRegName           = 100
	charxRegSwVersion      = 110
	charxRegNumControllers = 114

	// per-unit registers
	charxRegMeter = 112

	charxRegVoltages     = 232 // mV
	charxRegCurrents     = 238 // mA
	charxRegPower        = 244 // mW
	charxRegEnergy       = 250 // Wh
	charxRegSoc          = 264 // %
	charxRegEvid         = 265 // 10
	charxRegRfid         = 275 // 10
	charxRegChargeTime   = 287 // s
	charxRegChargeEnergy = 289 // Wh

	charxRegStatus     = 299 // IEC 61851-1
	charxRegEnable     = 300
	charxRegMaxCurrent = 301 // A

	charxOffset = 1000
)

// PhoenixCharx is an api.Charger implementation for Phoenix CHARX controller
type PhoenixCharx struct {
	conn      *modbus.Connection
	connector uint16
}

func init() {
	registry.Add("phoenix-charx", NewPhoenixCharxFromConfig)
}

//go:generate go run ../cmd/tools/decorate.go -f decoratePhoenixCharx -b *PhoenixCharx -r api.Charger -t "api.Meter,CurrentPower,func() (float64, error)" -t "api.MeterEnergy,TotalEnergy,func() (float64, error)" -t "api.PhaseCurrents,Currents,func() (float64, float64, float64, error)"

// NewPhoenixCharxFromConfig creates a Phoenix charger from generic config
func NewPhoenixCharxFromConfig(other map[string]interface{}) (api.Charger, error) {
	cc := struct {
		modbus.TcpSettings `mapstructure:",squash"`
		Connector          uint16
	}{
		TcpSettings: modbus.TcpSettings{
			ID: 1, // default
		},
		Connector: 1,
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	return NewPhoenixCharx(cc.URI, cc.ID, cc.Connector)
}

// NewPhoenixCharx creates a Phoenix charger
func NewPhoenixCharx(uri string, id uint8, connector uint16) (*PhoenixCharx, error) {
	conn, err := modbus.NewConnection(uri, "", "", 0, modbus.Tcp, id)
	if err != nil {
		return nil, err
	}

	log := util.NewLogger("charx")
	conn.Logger(log.TRACE)

	wb := &PhoenixCharx{
		conn:      conn,
		connector: connector,
	}

	controllers, err := wb.controllers()
	if err != nil {
		return nil, err
	}
	fmt.Println("controllers", controllers)

	if connector >= controllers {
		err = fmt.Errorf("invalid connector number: %d", connector)
	}

	return wb, err
}

func (wb *PhoenixCharx) controllers() (uint16, error) {
	b, err := wb.conn.ReadHoldingRegisters(charxRegNumControllers, 1)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint16(b), nil
}

func (wb *PhoenixCharx) register(reg uint16) uint16 {
	return wb.connector*charxOffset + reg
}

// Status implements the api.Charger interface
func (wb *PhoenixCharx) Status() (api.ChargeStatus, error) {
	b, err := wb.conn.ReadHoldingRegisters(wb.register(charxRegStatus), 1)
	if err != nil {
		return api.StatusNone, err
	}

	// TODO check IEC 61851-1 C1 state
	state := string(b[0])

	return api.ChargeStatus(state), nil
}

// Enabled implements the api.Charger interface
func (wb *PhoenixCharx) Enabled() (bool, error) {
	b, err := wb.conn.ReadHoldingRegisters(wb.register(charxRegEnable), 1)
	if err != nil {
		return false, err
	}

	return encoding.Uint16(b) == 1, nil
}

// Enable implements the api.Charger interface
func (wb *PhoenixCharx) Enable(enable bool) error {
	b := make([]byte, 2)
	if enable {
		binary.BigEndian.PutUint16(b, 1)
	}

	_, err := wb.conn.WriteMultipleRegisters(wb.register(charxRegEnable), 1, b)

	return err
}

// MaxCurrent implements the api.Charger interface
func (wb *PhoenixCharx) MaxCurrent(current int64) error {
	if current < 6 {
		return fmt.Errorf("invalid current %d", current)
	}

	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(current))

	_, err := wb.conn.WriteMultipleRegisters(wb.register(charxRegMaxCurrent), 1, b)

	return err
}

var _ api.ChargeTimer = (*PhoenixCharx)(nil)

// ChargingTime implements the api.ChargeTimer interface
func (wb *PhoenixCharx) ChargingTime() (time.Duration, error) {
	b, err := wb.conn.ReadHoldingRegisters(wb.register(charxRegChargeTime), 2)
	if err != nil {
		return 0, err
	}

	return time.Duration(encoding.Uint16(b)) * time.Second, nil
}

var _ api.Meter = (*PhoenixCharx)(nil)

// CurrentPower implements the api.Meter interface
func (wb *PhoenixCharx) CurrentPower() (float64, error) {
	b, err := wb.conn.ReadHoldingRegisters(wb.register(charxRegPower), 2)
	if err != nil {
		return 0, err
	}

	return float64(encoding.Int32(b)) / 1e3, nil
}

var _ api.MeterEnergy = (*PhoenixCharx)(nil)

// TotalEnergy implements the api.MeterEnergy interface
func (wb *PhoenixCharx) TotalEnergy() (float64, error) {
	b, err := wb.conn.ReadHoldingRegisters(wb.register(charxRegEnergy), 4)
	if err != nil {
		return 0, err
	}

	return float64(encoding.Int64(b)) / 1e3, nil
}

var _ api.PhaseCurrents = (*PhoenixCharx)(nil)

// Currents implements the api.PhaseCurrents interface
func (wb *PhoenixCharx) Currents() (float64, float64, float64, error) {
	b, err := wb.conn.ReadHoldingRegisters(wb.register(charxRegCurrents), 3*2)
	if err != nil {
		return 0, 0, 0, err
	}

	return float64(encoding.Int32(b)) / 1e3,
		float64(encoding.Int32(b[4:])) / 1e3,
		float64(encoding.Int32(b[8:])) / 1e3, nil
}

var _ api.PhaseVoltages = (*PhoenixCharx)(nil)

// Voltages implements the api.PhaseVoltages interface
func (wb *PhoenixCharx) Voltages() (float64, float64, float64, error) {
	b, err := wb.conn.ReadHoldingRegisters(wb.register(charxRegVoltages), 3*2)
	if err != nil {
		return 0, 0, 0, err
	}

	return float64(encoding.Int32(b)) / 1e3,
		float64(encoding.Int32(b[4:])) / 1e3,
		float64(encoding.Int32(b[8:])) / 1e3, nil
}

var _ api.Diagnosis = (*PhoenixCharx)(nil)

// Diagnose implements the api.Diagnosis interface
func (wb *PhoenixCharx) Diagnose() {
	if b, err := wb.conn.ReadHoldingRegisters(charxRegName, 10); err == nil {
		fmt.Printf("Name: %s\n", encoding.StringLsbFirst(b))
	}
	if b, err := wb.conn.ReadHoldingRegisters(charxRegSwVersion, 4); err == nil {
		fmt.Printf("Software version: %s\n", encoding.StringLsbFirst(b))
	}

	controllers, err := wb.controllers()
	if err == nil {
		fmt.Printf("Controllers: %d\n", controllers)
	}
}
