package main

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
