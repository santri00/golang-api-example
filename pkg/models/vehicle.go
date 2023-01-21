package models

// Vehicle represents a vehicle record
type Vehicle struct {
	StockID         int    `json:"stock_id"`
	Name            string `json:"name"`
	Vin             string `json:"vin"`
	Model           string `json:"model"`
	Year            string `json:"year"`
	Price           string `json:"price"`
	Brand           string `json:"brand"`
	Type            string `json:"type"`
	DaysInInventory int    `json:"days_in_inventory"`
}

// Vehicles represents the list of vehicles
type Vehicles []Vehicle
