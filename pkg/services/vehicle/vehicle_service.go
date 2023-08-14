package vehicle

import (
	"context"
	"errors"
	"github.com/santri00/golang-api-example/pkg/models"
	"github.com/santri00/golang-api-example/pkg/services/vehicle/mapper"
	"github.com/santri00/golang-api-example/pkg/storage/vehicle"
)

// Service defines the interface methods for vehicle service
type Service interface {
	//AddVehicle adds the vehicle data
	AddVehicle(ctx context.Context, vehicle *models.VehicleDetailedRecord) (*models.AddVehicleResponse, error)
	// ListVehicles returns the list of vehicles
	ListVehicles(ctx context.Context, offset, limit int) (*models.ListVehicleResponse, error)
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

// AddVehicle adds the vehicle
func (svc *service) AddVehicle(ctx context.Context, vehicle *models.VehicleDetailedRecord) (*models.AddVehicleResponse, error) {
	if vehicle == nil {
		return nil, errors.New("vehicle data is empty. The record can't be created")
	}

	//maps to entity
	entity := mapper.MapToVehicleEntity(vehicle)
	stockId, e := svc.vehicleRepo.Create(ctx, entity)
	return &models.AddVehicleResponse{StockId: stockId}, e
}

// ListVehicles returns the list of vehicles
func (svc *service) ListVehicles(ctx context.Context, offset, limit int) (*models.ListVehicleResponse, error) {
	records, totalRecordsCount, e := svc.vehicleRepo.GetAll(ctx, offset, limit)
	if e != nil {
		return nil, e
	}
	vehicles := mapper.MapFromVehicleBasicRecordEntity(records)
	res := &models.ListVehicleResponse{
		Metadata: models.Metadata{
			Offset: offset, Limit: limit, TotalCount: totalRecordsCount,
		},
		Records: vehicles,
	}
	return res, nil
}

// GetVehicleByStockId returns the vehicle data for the given stock ID
func (svc *service) GetVehicleByStockId(ctx context.Context, stockId string) (*models.VehicleDetailedRecord, error) {
	record, e := svc.vehicleRepo.Get(ctx, stockId)
	if e != nil {
		return nil, e
	}

	vehicle := mapper.MapFromVehicleEntity(record)
	return vehicle, nil
}
