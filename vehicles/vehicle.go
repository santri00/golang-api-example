package vehicles

import (
	"fmt"
	"time"

	"github.com/stripathy/demo-vehicle-api/dataaccess"
	"github.com/stripathy/demo-vehicle-api/models"
)

//IVehicleResource defines the vehicle interface
type IVehicleResource interface {
	GetVehicles() (*models.Vehicles, error)
}

//VehicleResource defines a vehicle resource
type VehicleResource struct {
	vehicleDataAccess dataaccess.IVehicleDataAccess
}

//New instantiates a new Vehicle
func New(dataAccess dataaccess.IVehicleDataAccess) *VehicleResource {
	return &VehicleResource{
		vehicleDataAccess: dataAccess,
	}
}

//GetVehicles Returns all vehicles
func (vehicle *VehicleResource) GetVehicles() (*models.Vehicles, error) {
	var vehicles models.Vehicles

	records, e := vehicle.vehicleDataAccess.GetAll()
	if e != nil {
		return nil, e
	}

	if records == nil ||
		len(*records) == 0 {
		return nil, nil
	}

	//map the data to model
	for _, v := range *records {
		invDate, e := time.Parse("01/02/2006", v.InventoryDate)
		if e != nil {
			return nil, fmt.Errorf("Failed to parse the vehicle inventory date; error: %s", e.Error())
		}

		vehicle := models.Vehicle{
			StockID:         v.StockID,
			Name:            v.Name,
			Vin:             v.Vin,
			Model:           v.Model,
			Year:            v.Year,
			Price:           v.Price,
			DaysInInventory: int(time.Now().Sub(invDate).Hours() / 24),
		}

		vehicles = append(vehicles, vehicle)
	}

	return &vehicles, nil
}
