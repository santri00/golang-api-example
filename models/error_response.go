package models

//ErrorResponse models api error response
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
	Detail     string `json:"detail"`
}
