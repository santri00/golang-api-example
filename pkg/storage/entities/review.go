package entities

// Review defines vehicle review entity
type Review struct {
	Title   string    `json:"title"`
	Comment string    `json:"comment"`
	Rating  int8      `json:"rating"`
	Pros    *[]string `json:"pros"`
	Cons    *[]string `json:"cons"`
}

// Reviews defines the list of reviews
type Reviews []Review
