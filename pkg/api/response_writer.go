package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondJSON(httpStatus int, o interface{}, w http.ResponseWriter, r *http.Request) {
	jsonResponse, err := json.Marshal(o)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		//Temp. move to logging pipeline
		log.Printf("Responded: %s\t%s\t%s\t%d", r.Method, r.URL.String(), r.RemoteAddr, http.StatusInternalServerError)
		return
	}

	//Temp. move to logging pipeline
	log.Printf("Responded: %s\t%s\t%s\t%d", r.Method, r.URL.String(), r.RemoteAddr, httpStatus)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(jsonResponse)
}

// errorResponse represents an error response returned by non-success http status response
type errorResponse struct {
	StatusCode       int      `json:"status_code"`
	StatusText       string   `json:"status_text"`
	Detail           string   `json:"detail"`
	ValidationErrors []string `json:"validation_errors"`
}
