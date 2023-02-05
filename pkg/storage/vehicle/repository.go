package vehicle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/santri00/golang-api-example/pkg/storage/cachekeys"
	"github.com/santri00/golang-api-example/pkg/storage/entities"
	"github.com/santri00/golang-api-example/pkg/utilities"
)

// Repository defines the vehicle repository interface
type Repository interface {
	Create(ctx context.Context, vehicleRecord *entities.Vehicle) (string, error)
	GetAll(ctx context.Context, offset int, limit int) (*[]entities.VehicleBasicRecord, error)
	Get(ctx context.Context, stockId string) (*[]entities.Vehicle, error)
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

// Create creates the vehicle record and adds to the storage.
// The method returns the newly created stock ID for the given vehicle record
func (repo *repository) Create(ctx context.Context, vehicleRecord *entities.Vehicle) (string, error) {
	if vehicleRecord == nil {
		return "", errors.New("vehicle data is missing. vehicle record can't be created")
	}

	if vehicleRecord.BasicRecord == nil {
		return "", errors.New("vehicle basic record is missing. vehicle record can't be created")
	}

	if vehicleRecord.DetailedRecord == nil {
		return "", errors.New("vehicle detailed record is missing. vehicle record can't be created")
	}

	//generate new stock id
	stockId, e := utilities.NewUUIDString()
	if e != nil {
		return "", fmt.Errorf("failed to generate stock Id. error: %s", e.Error())
	}

	//updates stock id in the vehicle data
	vehicleRecord.BasicRecord.StockID = stockId

	//serialize vehicle data
	//Basic record
	basicRecordBytes, e := json.Marshal(vehicleRecord.BasicRecord)
	if e != nil {
		return "", fmt.Errorf("failed to serialize vehicle basic record. error: %s", e.Error())
	}

	basicRecord := string(basicRecordBytes)

	//detailed record
	detailedRecordBytes, e := json.Marshal(vehicleRecord.DetailedRecord)
	if e != nil {
		return "", fmt.Errorf("failed to serialize vehicle detailed record. error: %s", e.Error())
	}

	//formats the vehicle hash key
	storageKey := fmt.Sprintf(cachekeys.VehicleRecordCacheKey, stockId)

	//create a transaction pipeline
	pipe := repo.redisClient.TxPipeline()
	defer pipe.Close()

	pipe.HSet(ctx, storageKey, cachekeys.BasicRecordCacheChildKey, basicRecord)
	pipe.HSet(ctx, storageKey, cachekeys.DetailedRecordCacheChildKey, string(detailedRecordBytes))
	//add to the vehicle list
	//use this list for retrieval all vehicles with the basic record
	pipe.RPush(ctx, cachekeys.VehicleRecordListCacheKey, basicRecord)
	//execute all the commands
	_, e = pipe.Exec(ctx)
	if e != nil {
		return "", fmt.Errorf("failed to create the vehicle record. error: %s", e.Error())
	}

	return stockId, nil
}

// GetAll fetches list of vehicle (basic/header) records
func (repo *repository) GetAll(ctx context.Context, offset int, limit int) (*[]entities.VehicleBasicRecord, error) {
	return nil, nil
}

// Get fetches the vehicle detailed record for the given stock ID
func (repo *repository) Get(ctx context.Context, stockId string) (*[]entities.Vehicle, error) {
	return nil, nil
}
