package api

import (
	"encoding/json"
	"net/http"
)

func respondJSON(o interface{}, w http.ResponseWriter) {
	jsonResponse, err := json.Marshal(o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

// errorResponse represents an error response returned by non-success http status response
type errorResponse struct {
	StatusCode int    `json:"status_code"`
	StatusText string `json:"status_text"`
	Detail     string `json:"detail"`
}
