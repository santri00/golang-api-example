package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stripathy/demo-vehicle-api/models"
	"github.com/stripathy/demo-vehicle-api/routes"
)

//ConfigureVehicleRoutes adds routes related to vehicle features
func (apiServer *APIServer) ConfigureVehicleRoutes() {
	routes.Routes = append(routes.Routes, routes.Route{
		Name:       "GetAllVehicles",
		HTTPVerb:   "GET",
		Pattern:    "/v1/vehicles",
		HandleFunc: apiServer.GetAllVehicles,
	},

		routes.Route{
			Name:       "GetVehicle",
			HTTPVerb:   "GET",
			Pattern:    "/v1/vehicles/{stock_id}",
			HandleFunc: apiServer.GetVehicle,
		},
	)
}

//GetAllVehicles returns all vehicles
func (apiServer *APIServer) GetAllVehicles(w http.ResponseWriter, r *http.Request) {

	vehicles, err := apiServer.vehicleResource.GetVehicles()
	if err != nil {
		log.Printf("ERROR >> %s", err.Error())
		respondJSON(models.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError),
			Detail:     "Somethimg went wrong",
		}, w)

		return
	}

	if vehicles == nil ||
		len(*vehicles) == 0 {
		log.Printf("WARNING >> %s", "No vehicles were found")
		respondJSON(models.ErrorResponse{
			StatusCode: http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound),
			Detail:     "No vehicles were found",
		}, w)

		return
	}

	log.Printf("Responded: %s\t%s\t%s", r.Method, r.URL.String(), r.RemoteAddr)
	respondJSON(vehicles, w)
}

//GetVehicle returns the vehicle data for the given stock Id
func (apiServer *APIServer) GetVehicle(w http.ResponseWriter, r *http.Request) {

	//Read the path parameter
	stockIdParam := mux.Vars(r)["stock_id"]
	stockId, err := strconv.ParseInt(stockIdParam, 10, 64)
	if err != nil {
		log.Printf("ERROR >> %s", err.Error())
		respondJSON(models.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			StatusText: http.StatusText(http.StatusBadRequest),
			Detail:     "Stock Id path parameter must be numeric",
		}, w)

		return
	}

	vehicle, err := apiServer.vehicleResource.GetVehicle(stockId)
	if err != nil {
		log.Printf("ERROR >> %s", err.Error())
		respondJSON(models.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError),
			Detail:     "Somethimg went wrong",
		}, w)

		return
	}

	if vehicle == nil {
		log.Printf("WARNING >> %s", "No vehicle data was found")
		respondJSON(models.ErrorResponse{
			StatusCode: http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound),
			Detail:     "No vehicle data was found",
		}, w)

		return
	}

	log.Printf("Responded: %s\t%s\t%s", r.Method, r.URL.String(), r.RemoteAddr)
	respondJSON(vehicle, w)
}
