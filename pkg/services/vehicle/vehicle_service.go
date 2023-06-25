package vehicle

import (
	"context"
	"github.com/santri00/golang-api-example/pkg/models"
	"github.com/santri00/golang-api-example/pkg/storage/vehicle"
)

// Service defines the interface methods for vehicle service
type Service interface {
	AddVehicle(ctx context.Context, vehicle *models.Vehicle) (string, error)
	// ListVehicles returns the list of vehicles
	ListVehicles(ctx context.Context) (*models.Vehicles, error)
	// GetVehicleByStockId returns the vehicle data for the given stock ID
	GetVehicleByStockId(ctx context.Context, stockId string) (*models.VehicleDetailedRecord, error)
}

type service struct {
	vehicleRepo vehicle.Repository
}

// NewService constructs the vehicle service instance
func NewService(vehicleRepo vehicle.Repository) Service {
	return &service{vehicleRepo: vehicleRepo}
}

func (svc *service) AddVehicle(ctx context.Context, vehicle *models.Vehicle) (string, error) {
	return "", nil
}

// ListVehicles returns the list of vehicles
func (svc *service) ListVehicles(ctx context.Context) (*models.Vehicles, error) {
	return nil, nil
}

// GetVehicleByStockId returns the vehicle data for the given stock ID
func (svc *service) GetVehicleByStockId(ctx context.Context, stockId string) (*models.VehicleDetailedRecord, error) {
	return nil, nil
}
