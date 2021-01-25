package models

//Vehicle models a vehicle
type Vehicle struct {
	StockID         int64  `json:"stock_id"`
	Name            string `json:"name"`
	Vin             string `json:"vin"`
	Model           string `json:"model"`
	Year            string `json:"year"`
	Price           string `json:"price"`
	DaysInInventory int    `json:"days_in_inventory"`
}

//Vehicles a list of vehicle
type Vehicles []Vehicle
