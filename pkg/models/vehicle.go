package models

// Vehicle defines a vehicle basic record
type Vehicle struct {
	StockID           int    `json:"stock_id"`
	Title             string `json:"title"`
	Price             string `json:"price"`
	Miles             string `json:"miles"`
	CoverImageLink    string `json:"cover_image_link"`
	VehicleRecordLink string `json:"vehicle_record_link"`
}

// Vehicles defines the list of vehicles
type Vehicles []Vehicle

// VehicleDetailedRecord defines a vehicle detailed record
type VehicleDetailedRecord struct {
	Vehicle                               //inherits the basic vehicle data
	Vin                   string          `json:"vin"`
	Model                 string          `json:"model"`
	Year                  string          `json:"year"`
	Brand                 string          `json:"brand"`
	Type                  string          `json:"type"`
	ExteriorColor         string          `json:"exterior_color"`
	InteriorColor         string          `json:"interior_color"`
	MilesPerGallonCity    int8            `json:"miles_per_gallon_city"`
	MilesPerGallonHighway int8            `json:"miles_per_gallon_highway"`
	TransmissionType      string          `json:"transmission_type"`
	Features              *[]string       `json:"features,omitempty"`
	Reviews               *VehicleReviews `json:"reviews,omitempty"`
	ExteriorImageLinks    *[]string       `json:"exterior_image_links,omitempty"`
	InteriorImageLinks    *[]string       `json:"interior_image_links,omitempty"`
	DaysInInventory       int             `json:"days_in_inventory"`
	IsWellEquipped        bool            `json:"is_well_equipped"`
	IsUnderWarranty       bool            `json:"is_under_warranty"`
}
