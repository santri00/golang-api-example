package vehicle

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/santri00/golang-api-example/pkg/storage/entities"
)

// Repository defines the vehicle repository interface
type Repository interface {
	Create(ctx context.Context, vehicleRecord *entities.Vehicle) (int, error)
	GetAll(ctx context.Context, offset int, limit int) (*[]entities.VehicleBasicRecord, error)
	Get(ctx context.Context, stockId int) (*[]entities.Vehicle, error)
}

type repository struct {
	redisClient *redis.Client
}

// New initializes the repository
func New(client *redis.Client) Repository {
	return &repository{
		redisClient: client,
	}
}

// Create adds the vehicle record data to the storage
func (repo *repository) Create(ctx context.Context, vehicleRecord *entities.Vehicle) (int, error) {
	return 0, nil
}

// GetAll fetches list of vehicle (basic/header) records
func (repo *repository) GetAll(ctx context.Context, offset int, limit int) (*[]entities.VehicleBasicRecord, error) {
	return nil, nil
}

// Get fetches the vehicle detailed record for the given stock ID
func (repo *repository) Get(ctx context.Context, stockId int) (*[]entities.Vehicle, error) {
	return nil, nil
}
