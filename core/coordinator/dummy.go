package coordinator

import (
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/core/vehicle"
)

type dummy struct{}

// NewDummy creates a dummy coordinator without vehicles
func NewDummy() API {
	return &dummy{}
}

func (a *dummy) GetVehicles() []api.Vehicle {
	return nil
}

func (a *dummy) Acquire(v api.Vehicle) {}

func (a *dummy) Release(v api.Vehicle) {}

func (a *dummy) IdentifyVehicleByStatus() api.Vehicle {
	return nil
}

func (a *dummy) GetVehicleIndex(v api.Vehicle) int {
	return -1
}

func (a *dummy) Settings(v api.Vehicle) vehicle.API {
	return nil
}
