package vehicle

import (
	"fmt"
	"strings"
	"time"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/vehicle/fiat"
)

// https://github.com/TA2k/ioBroker.fiat

// Fiat is an api.Vehicle implementation for Fiat cars
type Fiat struct {
	*embed
	*fiat.Provider
}

func init() {
	registry.Add("fiat", NewFiatFromConfig)
}

// NewFiatFromConfig creates a new vehicle
func NewFiatFromConfig(other map[string]interface{}) (api.Vehicle, error) {
	cc := struct {
		embed                    `mapstructure:",squash"`
		User, Password, VIN, PIN string
		Expiry                   time.Duration
		Cache                    time.Duration
	}{
		Expiry: expiry,
		Cache:  interval,
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	if cc.User == "" || cc.Password == "" {
		return nil, api.ErrMissingCredentials
	}

	v := &Fiat{
		embed: &cc.embed,
	}

	log := util.NewLogger("fiat").Redact(cc.User, cc.Password, cc.VIN)
	identity := fiat.NewIdentity(log, cc.User, cc.Password)

	err := identity.Login()
	if err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	api := fiat.NewAPI(log, identity)

	if cc.VIN == "" {
		cc.VIN, err = findVehicle(api.Vehicles())
		if err == nil {
			log.DEBUG.Printf("found vehicle: %v", cc.VIN)
		}
	}

	v.Provider = fiat.NewProvider(api, strings.ToUpper(cc.VIN), cc.PIN, cc.Expiry, cc.Cache)

	return v, err
}
