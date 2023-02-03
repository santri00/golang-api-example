package vehicle

import (
	"context"
	"github.com/santri00/golang-api-example/pkg/models"
)

// Service defines the interface methods for vehicle service
type Service interface {
	ListVehicles(ctx context.Context) (*models.Vehicles, error)
	GetVehicleByStockId(ctx context.Context, stockId int) (*models.VehicleDetailedRecord, error)
}

type service struct {
}

// NewService constructs the vehicle service instance
func NewService() Service {
	return &service{}
}

// ListVehicles returns the list of vehicles
func (svc *service) ListVehicles(ctx context.Context) (*models.Vehicles, error) {
	return nil, nil
}

// GetVehicleByStockId returns the vehicle data for the given stock ID
func (svc *service) GetVehicleByStockId(ctx context.Context, stockId int) (*models.VehicleDetailedRecord, error) {
	return nil, nil
}
