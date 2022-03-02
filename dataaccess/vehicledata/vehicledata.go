package vehicledata

import (
	"errors"
	"sync"
)

//IData Represents the interface
type IVehicleData interface {
	ReadAll() ([]string, error)
	Read(int64) (string, error)
}

//VehicleData Represents the struct
type VehicleData struct {
}

//New returns a new instance of Vehicle data
func New() *VehicleData {
	return &VehicleData{}
}

var vehicleData = map[int64]string{
	1001: `{ "stock_id": 1001, "name": "2021 MB GLC 300 4MATIC SUV", "vin": "WDC0G4KKXRF100001", "model": "GLC 300", "year": "2021", "price": "48000.99", "inventory_date": "11/10/2020"}`,
	1002: `{ "stock_id": 1002, "name": "2021 MB GLC 300 4MATIC SUV", "vin": "WDC0G4KKXRF100002", "model": "GLC 300", "year": "2021", "price": "49000.99", "inventory_date": "11/14/2020"}`,
	1003: `{ "stock_id": 1003, "name": "2021 MB GLC 300 4MATIC SUV", "vin": "WDC0G4KKXRF100003", "model": "GLC 300", "year": "2021", "price": "50000.99", "inventory_date": "11/16/2020"}`,
	1004: `{ "stock_id": 1004, "name": "2021 MB GLE 350 4MATIC SUV", "vin": "WDC0G4KKXRF100004", "model": "GLE 350", "year": "2021", "price": "58000.99", "inventory_date": "11/10/2020"}`,
	1005: `{ "stock_id": 1005, "name": "2021 MB GLE 350 4MATIC SUV", "vin": "WDC0G4KKXRF100006", "model": "GLE 350", "year": "2021", "price": "60000.99", "inventory_date": "11/14/2020"}`,
}

var lock = sync.RWMutex{}

//ReadAll Selects all vehicle records
func (vd *VehicleData) ReadAll() ([]string, error) {
	lock.RLock()
	defer lock.RUnlock()

	var records []string
	for _, record := range vehicleData {
		records = append(records, record)
	}

	return records, nil
}

//ReadAll Selects all vehicle records
func (vd *VehicleData) Read(stockId int64) (string, error) {
	lock.RLock()
	defer lock.RUnlock()

	var record string
	var exists bool
	if record, exists = vehicleData[stockId]; !exists {
		return "", errors.New("The vehicle record is not found")
	}

	return record, nil
}
