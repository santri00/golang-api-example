package vehicle

//TODO add unit tests

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
	//Create creates the vehicle record and adds to the storage.
	// The method returns the newly created stock ID for the given vehicle record
	Create(ctx context.Context, vehicleRecord *entities.Vehicle) (string, error)
	// GetAll fetches list of vehicle records as requested with the total record count
	// offset determines the starting record from which the list will be returned
	// limit determines the number of records that will be returned
	GetAll(ctx context.Context, offset int, limit int) (*[]entities.VehicleBasicRecord, int, error)
	// Get fetches the vehicle detailed record for the given stock ID
	Get(ctx context.Context, stockId string) (*entities.Vehicle, error)
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

	//generate new stock id (using uuid)
	//using uuid to uniquely identify each vehicle record is ok for our example use case
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
	//add vehicle record
	pipe.HSet(ctx, storageKey, cachekeys.BasicRecordCacheChildKey, basicRecord)
	pipe.HSet(ctx, storageKey, cachekeys.DetailedRecordCacheChildKey, string(detailedRecordBytes))
	//add to the vehicle list
	//this list is used to retrieve all vehicles with the basic record
	pipe.RPush(ctx, cachekeys.VehicleRecordListCacheKey, basicRecord)
	//execute all the commands
	_, e = pipe.Exec(ctx)
	if e != nil {
		return "", fmt.Errorf("failed to create the vehicle record. error: %s", e.Error())
	}

	return stockId, nil
}

// GetAll fetches list of vehicle (basic/header) records
func (repo *repository) GetAll(ctx context.Context, offset int, limit int) (*[]entities.VehicleBasicRecord, int, error) {
	if limit == 0 {
		return nil, 0, errors.New("invalid argument;limit value (should be greater than 0) must be provided")
	}

	//return list from cache
	data, e := repo.redisClient.LRange(ctx, cachekeys.VehicleRecordListCacheKey, int64(offset), int64(limit)).Result()
	if e != nil {
		return nil, 0, fmt.Errorf("failed to retrieve vehicle records. error: %s", e.Error())
	}

	//check if empty was returned
	if len(data) == 0 {
		return nil, 0, nil
	}

	//read data
	var result []entities.VehicleBasicRecord
	for _, record := range data {
		vehicle := entities.VehicleBasicRecord{}
		e = json.Unmarshal([]byte(record), &vehicle)
		if e != nil {
			return nil, 0, fmt.Errorf("failed to retrieve vehicle records. unmarshal error: %s", e.Error())
		}

		//add
		result = append(result, vehicle)
	}

	return &result, 0, nil
}

// Get fetches the vehicle detailed record for the given stock ID
func (repo *repository) Get(ctx context.Context, stockId string) (*entities.Vehicle, error) {
	if len(stockId) == 0 {
		return nil, errors.New("invalid argument;stock id must be provided")
	}

	//formats the vehicle hash key
	storageKey := fmt.Sprintf(cachekeys.VehicleRecordCacheKey, stockId)
	//fetch data from the cache
	data, e := repo.redisClient.HGetAll(ctx, storageKey).Result()
	if e != nil {
		return nil, fmt.Errorf("failed to retrieve the vehicle data for stock id %s; error: %s", stockId, e.Error())
	}

	//read data
	if len(data) == 0 {
		return nil, nil
	}

	basicRecord := entities.VehicleBasicRecord{}
	detailedRecord := entities.VehicleDetailedRecord{}
	for key, record := range data {
		switch key {
		case cachekeys.BasicRecordCacheChildKey:
			e = json.Unmarshal([]byte(record), &basicRecord)
			if e != nil {
				return nil, fmt.Errorf("failed to retrieve vehicle basic record. unmarshal error: %s", e.Error())
			}
		case cachekeys.DetailedRecordCacheChildKey:
			e = json.Unmarshal([]byte(record), &detailedRecord)
			if e != nil {
				return nil, fmt.Errorf("failed to retrieve vehicle detailed record. unmarshal error: %s", e.Error())
			}
		}
	}

	return &entities.Vehicle{BasicRecord: &basicRecord, DetailedRecord: &detailedRecord}, nil
}
