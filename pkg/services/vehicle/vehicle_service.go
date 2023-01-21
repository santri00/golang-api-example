package vehicle

import "github.com/santri00/golang-api-example/pkg/models"

// Service defines the interface methods for vehicle service
type Service interface {
	ListVehicles() (*models.Vehicles, error)
	GetVehicleByStockId(stockId int) (*models.Vehicle, error)
}

type service struct {
}

// NewService constructs the vehicle service instance
func NewService() Service {
	return &service{}
}

// ListVehicles returns the list of vehicles
func (svc *service) ListVehicles() (*models.Vehicles, error) {
	return nil, nil
}

// GetVehicleByStockId returns the vehicle data for the given stock ID
func (svc *service) GetVehicleByStockId(stockId int) (*models.Vehicle, error) {
	return nil, nil
}
