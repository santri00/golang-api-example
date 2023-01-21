package api

import (
	"github.com/santri00/golang-api-example/pkg/router"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// adds routes related to vehicle resource handler
func (apiServer *Server) configureVehicleHandlerRoutes() {
	router.Routes = append(router.Routes, router.Route{
		Name:       "ListVehicles",
		HTTPVerb:   "GET",
		Pattern:    "/inventory/v1/vehicles",
		HandleFunc: apiServer.listVehicles,
	},

		router.Route{
			Name:       "GetVehicleByStockId",
			HTTPVerb:   "GET",
			Pattern:    "/inventory/v1/vehicles/{stock_id}",
			HandleFunc: apiServer.getVehicleByStockId,
		},
	)
}

// Returns all vehicles
func (apiServer *Server) listVehicles(w http.ResponseWriter, r *http.Request) {

	vehicles, err := apiServer.vehicleService.ListVehicles()
	if err != nil {
		log.Printf("ERROR >> %s", err.Error())
		respondJSON(errorResponse{
			StatusCode: http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError),
			Detail:     "Something went wrong",
		}, w)

		return
	}

	if vehicles == nil ||
		len(*vehicles) == 0 {
		log.Printf("WARNING >> %s", "No vehicles were found")
		respondJSON(errorResponse{
			StatusCode: http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound),
			Detail:     "No vehicles were found",
		}, w)

		return
	}

	log.Printf("Responded: %s\t%s\t%s", r.Method, r.URL.String(), r.RemoteAddr)
	respondJSON(vehicles, w)
}

// Returns the vehicle data for the given stock ID
func (apiServer *Server) getVehicleByStockId(w http.ResponseWriter, r *http.Request) {

	//Read the path parameter
	stockIdParam := mux.Vars(r)["stock_id"]
	stockId, err := strconv.Atoi(stockIdParam)
	if err != nil {
		log.Printf("ERROR >> %s", err.Error())
		respondJSON(errorResponse{
			StatusCode: http.StatusBadRequest,
			StatusText: http.StatusText(http.StatusBadRequest),
			Detail:     "Stock Id path parameter must be numeric",
		}, w)

		return
	}

	vehicle, err := apiServer.vehicleService.GetVehicleByStockId(stockId)
	if err != nil {
		log.Printf("ERROR >> %s", err.Error())
		respondJSON(errorResponse{
			StatusCode: http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError),
			Detail:     "Something went wrong",
		}, w)

		return
	}

	if vehicle == nil {
		log.Printf("WARNING >> %s", "No vehicle data was found")
		respondJSON(errorResponse{
			StatusCode: http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound),
			Detail:     "No vehicle data was found",
		}, w)

		return
	}

	log.Printf("Responded: %s\t%s\t%s", r.Method, r.URL.String(), r.RemoteAddr)
	respondJSON(vehicle, w)
}
