package core

import (
	"time"

	"github.com/andig/evcc/api"
	"github.com/andig/evcc/core/wrapper"
	"github.com/andig/evcc/push"
	"github.com/andig/evcc/util"
	"github.com/pkg/errors"

	evbus "github.com/asaskevich/EventBus"
	"github.com/avast/retry-go"
	"github.com/benbjohnson/clock"
)

const (
	evChargeStart   = "start"   // update chargeTimer
	evChargeStop    = "stop"    // update chargeTimer
	evChargeCurrent = "current" // update fakeChargeMeter
	evChargePower   = "power"   // update chargeRater
)

// ThresholdConfig defines enable/disable hysteresis parameters
type ThresholdConfig struct {
	Delay     time.Duration
	Threshold float64
}

//go:generate mockgen -package mock -destination ../mock/mock_chargerhandler.go github.com/andig/evcc/core Handler

// Handler is the charger handler responsible for enabled state, target current and guard durations
type Handler interface {
	Prepare()
	SyncSettings()
	Enabled() bool
	Status() (api.ChargeStatus, error)
	TargetCurrent() int64
	Ramp(int64, ...bool) error
}

// LoadPoint is responsible for controlling charge depending on
// SoC needs and power availability.
type LoadPoint struct {
	clock    clock.Clock       // mockable time
	bus      evbus.Bus         // event bus
	pushChan chan<- push.Event // notifications
	uiChan   chan<- util.Param // client push messages
	log      *util.Logger

	// exposed public configuration
	Title      string `mapstructure:"title"`   // UI title
	Phases     int64  `mapstructure:"phases"`  // Phases- required for converting power and current
	ChargerRef string `mapstructure:"charger"` // Charger reference
	VehicleRef string `mapstructure:"vehicle"` // Vehicle reference
	Meters     struct {
		ChargeMeterRef string `mapstructure:"charge"` // Charge meter reference
	}
	Enable, Disable ThresholdConfig

	handler       Handler
	HandlerConfig `mapstructure:",squash"` // handle charger state and current

	chargeTimer api.ChargeTimer
	chargeRater api.ChargeRater

	chargeMeter api.Meter   // Charger usage meter
	vehicle     api.Vehicle // Vehicle

	// cached state
	status      api.ChargeStatus // Charger status
	charging    bool             // Charging cycle
	chargePower float64          // Charging power
	sitePower   float64          // Available power from site

	pvTimer time.Time
}

// NewLoadPointFromConfig creates a new loadpoint
func NewLoadPointFromConfig(log *util.Logger, cp configProvider, other map[string]interface{}) *LoadPoint {
	lp := NewLoadPoint(log)
	util.DecodeOther(log, other, &lp)

	if lp.Meters.ChargeMeterRef != "" {
		lp.chargeMeter = cp.Meter(lp.Meters.ChargeMeterRef)
	}
	if lp.VehicleRef != "" {
		lp.vehicle = cp.Vehicle(lp.VehicleRef)
	}

	if lp.ChargerRef == "" {
		lp.log.FATAL.Fatal("config: missing charger")
	}
	charger := cp.Charger(lp.ChargerRef)
	lp.configureChargerType(charger)

	lp.handler = &ChargerHandler{
		log:           lp.log,
		clock:         lp.clock,
		bus:           lp.bus,
		charger:       charger,
		HandlerConfig: lp.HandlerConfig,
	}

	return lp
}

// NewLoadPoint creates a LoadPoint with sane defaults
func NewLoadPoint(log *util.Logger) *LoadPoint {
	clock := clock.New()
	bus := evbus.New()

	lp := &LoadPoint{
		log:    log,   // logger
		clock:  clock, // mockable time
		bus:    bus,   // event bus
		Phases: 1,
		status: api.StatusNone,
		HandlerConfig: HandlerConfig{
			MinCurrent:    6,  // A
			MaxCurrent:    16, // A
			Sensitivity:   10, // A
			GuardDuration: 5 * time.Minute,
		},
	}

	return lp
}

// configureChargerType ensures that chargeMeter, Rate and Timer can use charger capabilities
func (lp *LoadPoint) configureChargerType(charger api.Charger) {
	// ensure charge meter exists
	if lp.chargeMeter == nil {
		if mt, ok := charger.(api.Meter); ok {
			lp.chargeMeter = mt
		} else {
			mt := &wrapper.ChargeMeter{}
			_ = lp.bus.Subscribe(evChargeCurrent, lp.evChargeCurrentHandler)
			_ = lp.bus.Subscribe(evChargeStop, func() {
				mt.SetPower(0)
			})
			lp.chargeMeter = mt
		}
	}

	// ensure charge rater exists
	if rt, ok := charger.(api.ChargeRater); ok {
		lp.chargeRater = rt
	} else {
		rt := wrapper.NewChargeRater(lp.log, lp.chargeMeter)
		_ = lp.bus.Subscribe(evChargePower, rt.SetChargePower)
		_ = lp.bus.Subscribe(evChargeStart, rt.StartCharge)
		_ = lp.bus.Subscribe(evChargeStop, rt.StopCharge)
		lp.chargeRater = rt
	}

	// ensure charge timer exists
	if ct, ok := charger.(api.ChargeTimer); ok {
		lp.chargeTimer = ct
	} else {
		ct := wrapper.NewChargeTimer()
		_ = lp.bus.Subscribe(evChargeStart, ct.StartCharge)
		_ = lp.bus.Subscribe(evChargeStop, ct.StopCharge)
		lp.chargeTimer = ct
	}
}

// notify sends push messages to clients
func (lp *LoadPoint) notify(event string) {
	lp.pushChan <- push.Event{Event: event}
}

// publish sends values to UI and databases
func (lp *LoadPoint) publish(key string, val interface{}) {
	lp.uiChan <- util.Param{Key: key, Val: val}
}

// evChargeStartHandler sends external start event
func (lp *LoadPoint) evChargeStartHandler() {
	lp.notify(evChargeStart)
}

// evChargeStopHandler sends external stop event
func (lp *LoadPoint) evChargeStopHandler() {
	lp.publishChargeProgress()
	lp.notify(evChargeStop)
}

// evChargeCurrentHandler updates the dummy charge meter's charge power. This simplifies the main flow
// where the charge meter can always be treated as present. It assumes that the charge meter cannot consume
// more than total household consumption. If physical charge meter is present this handler is not used.
func (lp *LoadPoint) evChargeCurrentHandler(current int64) {
	power := float64(current*lp.Phases) * Voltage

	if !lp.handler.Enabled() || lp.status != api.StatusC {
		// if disabled we cannot be charging
		power = 0
	}
	// TODO
	// else if power > 0 && lp.Site.pvMeter != nil {
	// 	// limit charge power to generation plus grid consumption/ minus grid delivery
	// 	// as the charger cannot have consumed more than that
	// 	// consumedPower := consumedPower(lp.pvPower, lp.batteryPower, lp.gridPower)
	// 	consumedPower := lp.Site.consumedPower()
	// 	power = math.Min(power, consumedPower)
	// }

	// handler only called if charge meter was replaced by dummy
	lp.chargeMeter.(*wrapper.ChargeMeter).SetPower(power)

	// expose for UI
	lp.publish("chargeCurrent", current)
}

// Name returns the human-readable loadpoint title
func (lp *LoadPoint) Name() string {
	return lp.Title
}

// Prepare loadpoint configuration by adding missing helper elements
func (lp *LoadPoint) Prepare(uiChan chan<- util.Param, pushChan chan<- push.Event) {
	lp.pushChan = pushChan
	lp.uiChan = uiChan

	// event handlers
	_ = lp.bus.Subscribe(evChargeStart, lp.evChargeStartHandler)
	_ = lp.bus.Subscribe(evChargeStop, lp.evChargeStopHandler)

	// prepare charger status
	lp.handler.Prepare()
}

// connected returns the EVs connection state
func (lp *LoadPoint) connected() bool {
	return lp.status == api.StatusB || lp.status == api.StatusC
}

// updateChargeStatus updates car status and detects car connected/disconnected events
func (lp *LoadPoint) updateChargeStatus() error {
	status, err := lp.handler.Status()
	if err != nil {
		return err
	}

	lp.log.DEBUG.Printf("charger status: %s", status)

	if prevStatus := lp.status; status != prevStatus {
		lp.status = status

		// changed from A - connected
		if prevStatus == api.StatusA {
			lp.log.INFO.Printf("car connected (%s)", string(status))
		}

		// changed to A -  disconnected
		if status == api.StatusA {
			lp.log.INFO.Println("car disconnected")
		}

		// update whenever there is a state change
		lp.bus.Publish(evChargeCurrent, lp.handler.TargetCurrent())

		// start/stop charging cycle
		if lp.charging = status == api.StatusC; lp.charging {
			lp.log.INFO.Println("start charging ->")
			lp.bus.Publish(evChargeStart)
		} else {
			lp.log.INFO.Println("stop charging <-")
			lp.bus.Publish(evChargeStop)
		}
	}

	return nil
}

func (lp *LoadPoint) maxCurrent(mode api.ChargeMode) int64 {
	// grid meter will always be available, if as wrapped pv meter
	targetPower := lp.chargePower - lp.sitePower
	lp.log.DEBUG.Printf("target power: %.0fW = %.0fW charge - %.0fW available", targetPower, lp.chargePower, lp.sitePower)

	// get max charge current
	targetCurrent := clamp(powerToCurrent(targetPower, Voltage, lp.Phases), 0, lp.MaxCurrent)

	// in MinPV mode return at least minCurrent
	if mode == api.ModeMinPV && targetCurrent < lp.MinCurrent {
		return lp.MinCurrent
	}

	// in PV mode disable if not connected and minCurrent not possible
	if mode == api.ModePV && lp.status != api.StatusC {
		lp.pvTimer = time.Time{}

		if targetCurrent < lp.MinCurrent {
			return 0
		}

		return lp.MinCurrent
	}

	// read only once to simplify testing
	enabled := lp.handler.Enabled()

	if mode == api.ModePV && enabled && targetCurrent < lp.MinCurrent {
		// kick off disable sequence
		if lp.sitePower >= lp.Disable.Threshold {
			lp.log.DEBUG.Printf("site power %.0f >= disable threshold %.0f", lp.sitePower, lp.Disable.Threshold)

			if lp.pvTimer.IsZero() {
				lp.log.DEBUG.Println("start pv disable timer")
				lp.pvTimer = lp.clock.Now()
			}

			if lp.clock.Since(lp.pvTimer) >= lp.Disable.Delay {
				lp.log.DEBUG.Println("pv disable timer elapsed")
				return 0
			}
		} else {
			// reset timer
			lp.pvTimer = lp.clock.Now()
		}

		return lp.MinCurrent
	}

	if mode == api.ModePV && !enabled {
		// kick off enable sequence
		if targetCurrent >= lp.MinCurrent ||
			(lp.Enable.Threshold != 0 && lp.sitePower <= lp.Enable.Threshold) {
			lp.log.DEBUG.Printf("site power %.0f < enable threshold %.0f", lp.sitePower, lp.Enable.Threshold)

			if lp.pvTimer.IsZero() {
				lp.log.DEBUG.Println("start pv enable timer")
				lp.pvTimer = lp.clock.Now()
			}

			if lp.clock.Since(lp.pvTimer) >= lp.Enable.Delay {
				lp.log.DEBUG.Println("pv enable timer elapsed")
				return lp.MinCurrent
			}
		} else {
			// reset timer
			lp.pvTimer = lp.clock.Now()
		}

		return 0
	}

	// reset timer to disabled state
	lp.log.DEBUG.Printf("pv timer reset")
	lp.pvTimer = time.Time{}

	return targetCurrent
}

// updateMeter updates and publishes single meter
func (lp *LoadPoint) updateMeter(name string, meter api.Meter, power *float64) error {
	value, err := meter.CurrentPower()
	if err != nil {
		return err
	}

	*power = value // update value if no error

	lp.log.DEBUG.Printf("%s power: %.1fW", name, *power)
	lp.publish(name+"Power", *power)

	return nil
}

// updateMeter updates and publishes single meter
func (lp *LoadPoint) updateMeters() (err error) {
	retryMeter := func(s string, m api.Meter, f *float64) {
		if m != nil {
			e := retry.Do(func() error {
				return lp.updateMeter(s, m, f)
			}, retry.Attempts(3))

			if e != nil {
				err = errors.Wrapf(e, "updating %s meter", s)
				lp.log.ERROR.Printf("%v", err)
			}
		}
	}

	// read PV meter before charge meter
	retryMeter("charge", lp.chargeMeter, &lp.chargePower)

	return err
}

// chargeDuration returns for how long the charge cycle has been running
func (lp *LoadPoint) chargeDuration() time.Duration {
	d, err := lp.chargeTimer.ChargingTime()
	if err != nil {
		lp.log.ERROR.Printf("charge timer error: %v", err)
		return 0
	}
	return d
}

// chargedEnergy returns energy consumption since charge start in kWh
func (lp *LoadPoint) chargedEnergy() float64 {
	f, err := lp.chargeRater.ChargedEnergy()
	if err != nil {
		lp.log.ERROR.Printf("charge rater error: %v", err)
		return 0
	}
	return f
}

// publish charged energy and duration
func (lp *LoadPoint) publishChargeProgress() {
	lp.publish("chargedEnergy", 1e3*lp.chargedEnergy()) // return Wh for UI
	lp.publish("chargeDuration", lp.chargeDuration())
}

// remainingChargeDuration returns the remaining charge time
func (lp *LoadPoint) remainingChargeDuration(chargePercent float64) time.Duration {
	if !lp.charging {
		return -1
	}

	if lp.chargePower > 0 && lp.vehicle != nil {
		whRemaining := (1 - chargePercent/100.0) * 1e3 * float64(lp.vehicle.Capacity())
		return time.Duration(float64(time.Hour) * whRemaining / lp.chargePower)
	}

	return -1
}

// publish state of charge and remaining charge duration
func (lp *LoadPoint) publishSoC() {
	if lp.vehicle == nil {
		return
	}

	if lp.connected() {
		f, err := lp.vehicle.ChargeState()
		if err == nil {
			lp.log.DEBUG.Printf("vehicle soc: %.1f%%", f)
			lp.publish("socCharge", f)
			lp.publish("chargeEstimate", lp.remainingChargeDuration(f))
			return
		}
		lp.log.ERROR.Printf("vehicle error: %v", err)
	}

	lp.publish("socCharge", -1)
	lp.publish("chargeEstimate", -1)
}

// Update is the main control function. It reevaluates meters and charger state
func (lp *LoadPoint) Update(mode api.ChargeMode, sitePower float64) float64 {
	// read and publish meters first
	_ = lp.updateMeters()

	lp.sitePower = sitePower

	// update ChargeRater here to make sure initial meter update is caught
	lp.bus.Publish(evChargeCurrent, lp.handler.TargetCurrent())
	lp.bus.Publish(evChargePower, lp.chargePower)

	// update progress and soc before status is updated
	lp.publishChargeProgress()
	lp.publishSoC()

	// read and publish status
	if err := retry.Do(lp.updateChargeStatus, retry.Attempts(3)); err != nil {
		lp.log.ERROR.Printf("charge controller error: %v", err)
		return lp.chargePower
	}

	lp.publish("connected", lp.connected())
	lp.publish("charging", lp.charging)

	// sync settings with charger
	if lp.status != api.StatusA {
		lp.handler.SyncSettings()
	}

	// check if car connected and ready for charging
	var err error

	// execute loading strategy
	switch mode {
	case api.ModeOff:
		err = lp.handler.Ramp(0, true)

	case api.ModeNow:
		// ensure that new connections happen at min current
		current := lp.MinCurrent
		if lp.connected() {
			current = lp.MaxCurrent
		}
		err = lp.handler.Ramp(current, true)

	case api.ModeMinPV, api.ModePV:
		targetCurrent := lp.maxCurrent(mode)
		if !lp.connected() {
			// ensure minimum current when not connected
			// https://github.com/andig/evcc/issues/105
			targetCurrent = min(lp.MinCurrent, targetCurrent)
		}
		lp.log.DEBUG.Printf("target charge current: %dA", targetCurrent)

		err = lp.handler.Ramp(targetCurrent)
	}

	if err != nil {
		lp.log.ERROR.Println(err)
	}

	return lp.chargePower
}
