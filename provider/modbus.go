package provider

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/modbus"
	gridx "github.com/grid-x/modbus"
	"github.com/volkszaehler/mbmd/meters"
	"github.com/volkszaehler/mbmd/meters/sunspec"
)

// Modbus implements modbus RTU and TCP access
type Modbus struct {
	log    *util.Logger
	conn   *modbus.Connection
	device meters.Device
	op     modbus.Operation
	scale  float64
}

func init() {
	registry.Add("modbus", NewModbusFromConfig)
}

// NewModbusFromConfig creates Modbus plugin
func NewModbusFromConfig(other map[string]interface{}) (Provider, error) {
	cc := struct {
		modbus.Settings `mapstructure:",squash"`
		Register        modbus.Register
		Model           string
		Value           string
		Scale           float64
		Delay           time.Duration
		ConnectDelay    time.Duration
		Timeout         time.Duration
	}{
		Scale: 1,
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	// assume RTU if not set and this is a known RS485 meter model
	if cc.RTU == nil {
		b := modbus.IsRS485(cc.Model)
		cc.RTU = &b
	}

	conn, err := modbus.NewConnection(cc.URI, cc.Device, cc.Comset, cc.Baudrate, modbus.ProtocolFromRTU(cc.RTU), cc.ID)
	if err != nil {
		return nil, err
	}

	// set non-default timeout
	if cc.Timeout > 0 {
		conn.Timeout(cc.Timeout)
	}

	// set non-default delay
	if cc.Delay > 0 {
		conn.Delay(cc.Delay)
	}

	// set non-default connect delay
	if cc.ConnectDelay > 0 {
		conn.ConnectDelay(cc.ConnectDelay)
	}

	log := util.NewLogger("modbus")
	conn.Logger(log.TRACE)

	var device meters.Device
	var op modbus.Operation

	if cc.Value == "" && cc.Register.Address == 0 {
		return nil, errors.New("either value or register required")
	}

	if cc.Model == "" && cc.Value != "" {
		return nil, errors.New("need device model when value is configured")
	}

	// no register configured - need device
	if cc.Register.Address == 0 {
		device, err = modbus.NewDevice(cc.Model, cc.SubDevice)
		if err != nil {
			return nil, err
		}

		// prepare device; ignore KOSTAL implementation errors
		if err := device.Initialize(conn); err != nil && !errors.Is(err, meters.ErrPartiallyOpened) {
			return nil, err
		}

		// parse value
		if err := modbus.ParseOperation(device, cc.Value, &op); err != nil {
			return nil, fmt.Errorf("invalid value %s", cc.Value)
		}
	}

	// // register configured
	// if cc.Register.Address != 0 {
	// 	if err := modbus.RegisterOperation(cc.Register); err != nil {
	// 		return nil, err
	// 	}
	// }

	if op.Type() == modbus.OperationTypeUnknown {
		return nil, fmt.Errorf("invalid register: %v/%v/%v", cc.Device, cc.Value, cc.Register)
	}

	mb := &Modbus{
		log:    log,
		conn:   conn,
		device: device,
		op:     op,
		scale:  cc.Scale,
	}
	return mb, nil
}

func (m *Modbus) readBytes() ([]byte, error) {
	if m.op.Type() != modbus.OperationTypeRegister {
		return nil, errors.New("expected rtu reading")
	}

	switch reg := m.op.Register; reg.Type {
	case modbus.RegisterTypeHoldings:
		return m.conn.ReadHoldingRegisters(reg.Address, reg.Encoding.Len())

	case modbus.RegisterTypeInput:
		return m.conn.ReadInputRegisters(reg.Address, reg.Encoding.Len())

	// case modbus.RegisterTypeCoils:
	// 	return m.conn.ReadCoils(reg.Address, reg.Encoding.Len())

	default:
		return nil, fmt.Errorf("invalid read type: %s", reg.Type)
	}
}

func (m *Modbus) floatGetter() (f float64, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()

	defer func() {
		// silence NaN reading errors by assuming zero
		if err != nil && errors.Is(err, meters.ErrNaN) {
			f = 0
			err = nil
		}
	}()

	switch m.op.Type() {
	case modbus.OperationTypeRegister:
		bytes, err := m.readBytes()
		if err != nil {
			return 0, fmt.Errorf("read %d failed: %w", m.op.Register.Address, err)
		}

		// TODO transform
		_ = bytes
		// return m.scale * op.Transform(bytes), nil
		return 0, nil

	case modbus.OperationTypeMeasurement:
		// TODO sunspec catch-all?
		res, err := m.device.(*sunspec.SunSpec).QueryOp(m.conn, m.op.Measurement)
		return m.scale * res.Value, err

	case modbus.OperationTypeSunSpec:
		res, err := m.device.(*sunspec.SunSpec).QueryPoint(
			m.conn,
			m.op.SunSpec.Model,
			m.op.SunSpec.Block,
			m.op.SunSpec.Point,
		)
		return m.scale * res, err

	default:
		return 0, fmt.Errorf("invalid operation type: %s", m.op.Type())
	}
}

var _ FloatProvider = (*Modbus)(nil)

// FloatGetter executes configured modbus read operation and implements func() (float64, error)
func (m *Modbus) FloatGetter() (func() (f float64, err error), error) {
	return m.floatGetter, nil
}

var _ IntProvider = (*Modbus)(nil)

// IntGetter executes configured modbus read operation and implements IntProvider
func (m *Modbus) IntGetter() (func() (int64, error), error) {
	g, err := m.FloatGetter()

	return func() (int64, error) {
		res, err := g()
		return int64(math.Round(res)), err
	}, err
}

var _ StringProvider = (*Modbus)(nil)

// StringGetter implements StringProvider
func (m *Modbus) StringGetter() (func() (string, error), error) {
	return func() (string, error) {
		b, err := m.readBytes()
		if err != nil {
			return "", err
		}

		return strings.TrimSpace(string(bytes.TrimLeft(b, "\x00"))), nil
	}, nil
}

// UintFromBytes converts byte slice to bigendian uint value
func UintFromBytes(bytes []byte) (u uint64, err error) {
	switch l := len(bytes); l {
	case 1:
		u = uint64(bytes[0])
	case 2:
		u = uint64(binary.BigEndian.Uint16(bytes))
	case 4:
		u = uint64(binary.BigEndian.Uint32(bytes))
	case 8:
		u = binary.BigEndian.Uint64(bytes)
	default:
		err = fmt.Errorf("unexpected length: %d", l)
	}

	return u, err
}

var _ BoolProvider = (*Modbus)(nil)

// BoolGetter implements BoolProvider
func (m *Modbus) BoolGetter() (func() (bool, error), error) {
	return func() (bool, error) {
		bytes, err := m.readBytes()
		if err != nil {
			return false, err
		}

		u, err := UintFromBytes(bytes)
		return u > 0, err
	}, nil
}

var _ SetFloatProvider = (*Modbus)(nil)

// FloatSetter executes configured modbus write operation and implements SetFloatProvider
func (m *Modbus) FloatSetter(_ string) (func(float64) error, error) {
	op := m.op.MBMD
	if op.FuncCode == 0 {
		return nil, errors.New("modbus plugin does not support writing to sunspec")
	}

	// need multiple registers for float
	if op.FuncCode != gridx.FuncCodeWriteMultipleRegisters {
		return nil, fmt.Errorf("invalid write function code: %d", op.FuncCode)
	}

	return func(val float64) error {
		val = m.scale * val

		var err error
		switch op.ReadLen {
		case 2:
			var b [4]byte
			binary.BigEndian.PutUint32(b[:], math.Float32bits(float32(val)))
			_, err = m.conn.WriteMultipleRegisters(op.OpCode, 2, b[:])

		case 4:
			var b [8]byte
			binary.BigEndian.PutUint64(b[:], math.Float64bits(val))
			_, err = m.conn.WriteMultipleRegisters(op.OpCode, 4, b[:])

		default:
			err = fmt.Errorf("invalid write length: %d", op.ReadLen)
		}

		return err
	}, nil
}

var _ SetIntProvider = (*Modbus)(nil)

// IntSetter executes configured modbus write operation and implements SetIntProvider
func (m *Modbus) IntSetter(_ string) (func(int64) error, error) {
	op := m.op.MBMD
	if op.FuncCode == 0 {
		return nil, errors.New("modbus plugin does not support writing to sunspec")
	}

	return func(val int64) error {
		ival := int64(m.scale * float64(val))

		// if funccode is configured, execute the read directly
		var err error
		switch op.FuncCode {
		case gridx.FuncCodeWriteSingleRegister:
			_, err = m.conn.WriteSingleRegister(op.OpCode, uint16(ival))

		case gridx.FuncCodeWriteMultipleRegisters:
			switch op.ReadLen {
			case 1:
				var b [2]byte
				binary.BigEndian.PutUint16(b[:], uint16(ival))
				_, err = m.conn.WriteMultipleRegisters(op.OpCode, 1, b[:])

			case 2:
				var b [4]byte
				binary.BigEndian.PutUint32(b[:], uint32(ival))
				_, err = m.conn.WriteMultipleRegisters(op.OpCode, 2, b[:])

			case 4:
				var b [8]byte
				binary.BigEndian.PutUint64(b[:], uint64(ival))
				_, err = m.conn.WriteMultipleRegisters(op.OpCode, 4, b[:])

			default:
				err = fmt.Errorf("invalid write length: %d", op.ReadLen)
			}

		case gridx.FuncCodeWriteSingleCoil:
			if ival != 0 {
				// Modbus protocol requires 0xFF00 for ON
				// and 0x0000 for OFF
				ival = 0xFF00
			}
			_, err = m.conn.WriteSingleCoil(op.OpCode, uint16(ival))

		default:
			err = fmt.Errorf("invalid write function code: %d", op.FuncCode)
		}

		return err
	}, nil
}

var _ SetBoolProvider = (*Modbus)(nil)

// BoolSetter executes configured modbus write operation and implements SetBoolProvider
func (m *Modbus) BoolSetter(param string) (func(bool) error, error) {
	set, err := m.IntSetter(param)

	return func(val bool) error {
		var ival int64
		if val {
			ival = 1
		}

		return set(ival)
	}, err
}
