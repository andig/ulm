package hems

import (
	"errors"
	"strings"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/core/site"
	"github.com/evcc-io/evcc/hems/eebus"
	"github.com/evcc-io/evcc/hems/semp"
	"github.com/evcc-io/evcc/server"
)

// NewFromConfig creates new HEMS from config
func NewFromConfig(typ string, other map[string]interface{}, site site.API, httpd *server.HTTPd) (api.HEMS, error) {
	switch strings.ToLower(typ) {
	case "sma", "shm", "semp":
		return semp.New(other, site, httpd)
	case "eebus":
		return eebus.New(other, site)
	default:
		return nil, errors.New("unknown hems: " + typ)
	}
}
