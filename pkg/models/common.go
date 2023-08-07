package models

// Metadata used to return metadata for a search/list operation
type Metadata struct {
	Offset     int `json:"offset"`
	Limit      int `json:"limit"`
	TotalCount int `json:"total_count"`
}
