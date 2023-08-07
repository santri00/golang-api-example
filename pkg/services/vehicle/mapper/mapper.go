package mapper

import (
	"github.com/santri00/golang-api-example/pkg/models"
	"github.com/santri00/golang-api-example/pkg/storage/entities"
	"time"
)

// MapToVehicleEntity maps model object to vehicle entity
func MapToVehicleEntity(vehicle *models.VehicleDetailedRecord) *entities.Vehicle {
	if vehicle == nil {
		return nil
	}

	vehicleBasicRecord := &entities.VehicleBasicRecord{
		Title:              vehicle.Title,
		Price:              vehicle.Price,
		Miles:              vehicle.Miles,
		CoverImageLink:     vehicle.CoverImageLink,
		InventoryAddedDate: vehicle.InventoryAddedDate,
	}
	vehicleDetailedRecord := &entities.VehicleDetailedRecord{
		Vin:                   vehicle.Vin,
		Model:                 vehicle.Model,
		Year:                  vehicle.Year,
		Brand:                 vehicle.Brand,
		Type:                  vehicle.Type,
		ExteriorColor:         vehicle.ExteriorColor,
		InteriorColor:         vehicle.InteriorColor,
		MilesPerGallonCity:    vehicle.MilesPerGallonCity,
		MilesPerGallonHighway: vehicle.MilesPerGallonHighway,
		TransmissionType:      vehicle.TransmissionType,
		Features:              vehicle.Features,
		ExteriorImageLinks:    vehicle.ExteriorImageLinks,
		InteriorImageLinks:    vehicle.InteriorImageLinks,
		IsWellEquipped:        vehicle.IsWellEquipped,
	}

	return &entities.Vehicle{
		BasicRecord:    vehicleBasicRecord,
		DetailedRecord: vehicleDetailedRecord,
	}
}

// MapFromVehicleBasicRecordEntity maps the entity to the vehicle model
func MapFromVehicleBasicRecordEntity(records *[]entities.VehicleBasicRecord) *models.Vehicles {
	if len(*records) == 0 {
		return nil
	}

	vehicles := models.Vehicles{}
	for _, record := range *records {
		vehicle := models.Vehicle{
			StockID:            record.StockID,
			Title:              record.Title,
			Price:              record.Price,
			Miles:              record.Miles,
			CoverImageLink:     record.CoverImageLink,
			InventoryAddedDate: record.InventoryAddedDate,
			DaysInInventory:    calculateDaysInInventory(record.InventoryAddedDate),
		}
		vehicles = append(vehicles, vehicle)
	}

	return &vehicles
}

// MapFromVehicleEntity maps the entity to the vehicle detailed model
func MapFromVehicleEntity(record *entities.Vehicle) *models.VehicleDetailedRecord {
	if record == nil ||
		(record.BasicRecord == nil && record.DetailedRecord == nil) {
		return nil
	}

	detail := models.VehicleDetailedRecord{}
	if record.BasicRecord != nil {
		detail.Vehicle.StockID = record.BasicRecord.StockID
		detail.Vehicle.Title = record.BasicRecord.Title
		detail.Vehicle.Price = record.BasicRecord.Price
		detail.Vehicle.Miles = record.BasicRecord.Miles
		detail.Vehicle.CoverImageLink = record.BasicRecord.CoverImageLink
		detail.Vehicle.InventoryAddedDate = record.BasicRecord.InventoryAddedDate
		detail.Vehicle.DaysInInventory = calculateDaysInInventory(record.BasicRecord.InventoryAddedDate)
	}

	if record.DetailedRecord != nil {
		detail.Vin = record.DetailedRecord.Vin
		detail.Model = record.DetailedRecord.Model
		detail.Year = record.DetailedRecord.Year
		detail.Brand = record.DetailedRecord.Brand
		detail.Type = record.DetailedRecord.Type
		detail.ExteriorColor = record.DetailedRecord.ExteriorColor
		detail.InteriorColor = record.DetailedRecord.InteriorColor
		detail.MilesPerGallonCity = record.DetailedRecord.MilesPerGallonCity
		detail.MilesPerGallonHighway = record.DetailedRecord.MilesPerGallonHighway
		detail.TransmissionType = record.DetailedRecord.TransmissionType
		detail.Features = record.DetailedRecord.Features
		detail.ExteriorImageLinks = record.DetailedRecord.ExteriorImageLinks
		detail.InteriorImageLinks = record.DetailedRecord.InteriorImageLinks
		detail.IsWellEquipped = record.DetailedRecord.IsWellEquipped
	}

	detail.IsUnderWarranty = isUnderWarranty(detail.Year, detail.Miles)

	return &detail
}

func calculateDaysInInventory(inventoryDate time.Time) int {
	return int((time.Since(inventoryDate).Hours()) / 24)
}

func isUnderWarranty(year, miles int) bool {
	return time.Now().Year()-year < 5 && miles < 50000
}
