package entity

//Vehicle defines vehicle data
type Vehicle struct {
	StockID       int64  `json:"stock_id"`
	Name          string `json:"name"`
	Vin           string `json:"vin"`
	Model         string `json:"model"`
	Year          string `json:"year"`
	Price         string `json:"price"`
	InventoryDate string `json:"inventory_date"`
}

//VehicleDetail defines vehicle detailed data
type VehicleDetail struct {
	StockID int64
	//details
}
