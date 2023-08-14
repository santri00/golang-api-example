package models

import "time"

// AddVehicleResponse response for adding the vehicle record
type AddVehicleResponse struct {
	StockId string `json:"stock_id"`
}

// ListVehicleResponse response for returning list of vehicles
type ListVehicleResponse struct {
	Metadata Metadata  `json:"_metadata"`
	Records  *Vehicles `json:"records"`
}

// Vehicle defines a vehicle basic record
type Vehicle struct {
	StockID            string    `json:"stock_id"`
	Title              string    `json:"title"`
	Price              string    `json:"price"`
	Miles              int       `json:"miles"`
	CoverImageLink     string    `json:"cover_image_link"`
	InventoryAddedDate time.Time `json:"inventory_added_date"`
	DaysInInventory    int       `json:"days_in_inventory"`
}

// Vehicles defines the list of vehicles
type Vehicles []Vehicle

// VehicleDetailedRecord defines a vehicle detailed record
type VehicleDetailedRecord struct {
	Vehicle                         //inherits the basic vehicle data
	Vin                   string    `json:"vin"`
	Model                 string    `json:"model"`
	Year                  int       `json:"year"`
	Brand                 string    `json:"brand"`
	Type                  string    `json:"type"`
	ExteriorColor         string    `json:"exterior_color"`
	InteriorColor         string    `json:"interior_color"`
	MilesPerGallonCity    int8      `json:"miles_per_gallon_city"`
	MilesPerGallonHighway int8      `json:"miles_per_gallon_highway"`
	TransmissionType      string    `json:"transmission_type"`
	Features              *[]string `json:"features,omitempty"`
	//Reviews               *VehicleReviews `json:"reviews,omitempty"`
	ExteriorImageLinks *[]string `json:"exterior_image_links,omitempty"`
	InteriorImageLinks *[]string `json:"interior_image_links,omitempty"`
	IsWellEquipped     bool      `json:"is_well_equipped"`
	IsUnderWarranty    bool      `json:"is_under_warranty"`
}
