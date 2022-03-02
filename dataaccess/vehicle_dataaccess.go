package dataaccess

import (
	"encoding/json"

	"github.com/stripathy/demo-vehicle-api/dataaccess/entity"
	"github.com/stripathy/demo-vehicle-api/dataaccess/vehicledata"
)

//IVehicleDataAccess defines the interface for the vehicle dataaccess
type IVehicleDataAccess interface {
	GetAll() (*[]entity.Vehicle, error)
	Get(stockID int64) (*entity.Vehicle, error)
}

//VehicleDataAccess defines the vehicle dataaccess struct
type VehicleDataAccess struct {
	vehicleData vehicledata.IVehicleData
}

//New returns a new instance of Vehicle Dataaccess
func New() *VehicleDataAccess {
	return &VehicleDataAccess{vehicleData: vehicledata.New()}
}

//GetAll gets all vehicles from the data store
func (dataAccess *VehicleDataAccess) GetAll() (*[]entity.Vehicle, error) {

	data, e := dataAccess.vehicleData.ReadAll()
	if e != nil {
		return nil, e
	}

	if len(data) == 0 {
		return nil, nil
	}

	vehicles := []entity.Vehicle{}
	for _, record := range data {
		vehicle := entity.Vehicle{}

		e = json.Unmarshal([]byte(record), &vehicle)
		if e != nil {
			return nil, e
		}

		vehicles = append(vehicles, vehicle)
	}

	return &vehicles, nil
}

//Get gets the vehicle detail for the stockID from the data store
func (dataAccess *VehicleDataAccess) Get(stockID int64) (*entity.Vehicle, error) {

	data, e := dataAccess.vehicleData.Read(stockID)
	if e != nil {
		return nil, e
	}

	vehicle := entity.Vehicle{}
	e = json.Unmarshal([]byte(data), &vehicle)
	if e != nil {
		return nil, e
	}

	return &vehicle, nil
}
