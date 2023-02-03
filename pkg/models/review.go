package models

// VehicleReview defines a vehicle review struct
type VehicleReview struct {
	Title   string    `json:"title"`
	Comment string    `json:"comment"`
	Rating  int8      `json:"rating"`
	Pros    *[]string `json:"pros"`
	Cons    *[]string `json:"cons"`
}

// VehicleReviews defines a list of vehicle reviews
type VehicleReviews []VehicleReview
