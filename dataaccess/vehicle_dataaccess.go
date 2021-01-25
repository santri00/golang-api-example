package dataaccess

import (
	"encoding/json"

	"github.com/stripathy/demo-vehicle-api/dataaccess/data"
	"github.com/stripathy/demo-vehicle-api/dataaccess/entity"
)

//IVehicleDataAccess defines the interface for the vehicle dataaccess
type IVehicleDataAccess interface {
	GetAll() (*[]entity.Vehicle, error)
	Get(stockID int64) (*entity.Vehicle, error)
}

//VehicleDataAccess defines the vehicle dataaccess struct
type VehicleDataAccess struct {
}

//New returns a new instance of Vehicle Dataaccess
func New() *VehicleDataAccess {
	return &VehicleDataAccess{}
}

//GetAll gets all vehicles from the data store
func (dataAccess *VehicleDataAccess) GetAll() (*[]entity.Vehicle, error) {
	dataJSON := data.VehiclesDataJSON
	vehicles := []entity.Vehicle{}

	e := json.Unmarshal([]byte(dataJSON), &vehicles)
	if e != nil {
		return nil, e
	}

	return &vehicles, nil
}

//Get gets the vehicle detail for the stockID from the data store
func (dataAccess *VehicleDataAccess) Get(stockID int64) (*entity.Vehicle, error) {
	return nil, nil
}
