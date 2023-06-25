package api

import (
	"github.com/gorilla/mux"
	"github.com/santri00/golang-api-example/pkg/router"
	"log"
	"net/http"
)

// adds routes related to vehicle resource handler
func (server *Server) configureVehicleHandlerRoutes() {
	router.Routes = append(router.Routes,

		router.Route{
			Name:       "AddVehicle",
			HTTPVerb:   "POST",
			Pattern:    "/inventory/v1/vehicles",
			HandleFunc: server.addVehicles,
		},
		router.Route{
			Name:       "ListVehicles",
			HTTPVerb:   "GET",
			Pattern:    "/inventory/v1/vehicles",
			HandleFunc: server.listVehicles,
		},
		router.Route{
			Name:       "GetVehicleByStockId",
			HTTPVerb:   "GET",
			Pattern:    "/inventory/v1/vehicles/:stock_id",
			HandleFunc: server.getVehicleByStockId,
		},
	)
}

// AddVehicles adds requested vehicle records in the storage
func (server *Server) addVehicles(w http.ResponseWriter, r *http.Request) {
	//todo
}

// Returns all vehicles
func (server *Server) listVehicles(w http.ResponseWriter, r *http.Request) {

	vehicles, err := server.vehicleService.ListVehicles(server.ctx)
	if err != nil {
		log.Printf("ERROR >> %s", err.Error())
		respondJSON(http.StatusInternalServerError, errorResponse{
			StatusCode: http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError),
			Detail:     "Something went wrong",
		}, w, r)

		return
	}

	if vehicles == nil ||
		len(*vehicles) == 0 {
		log.Printf("WARNING >> %s", "No vehicles were found")
		respondJSON(http.StatusNotFound, errorResponse{
			StatusCode: http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound),
			Detail:     "No vehicles were found",
		}, w, r)

		return
	}

	respondJSON(http.StatusOK, vehicles, w, r)
}

// Returns the vehicle data for the given stock ID
func (server *Server) getVehicleByStockId(w http.ResponseWriter, r *http.Request) {
	//Read the path parameter
	stockId := mux.Vars(r)["stock_id"]
	if len(stockId) == 0 {
		log.Printf("ERROR >> %s", "request is invalid. missing a valid stock Id in the url")
		respondJSON(http.StatusBadRequest, errorResponse{
			StatusCode: http.StatusBadRequest,
			StatusText: http.StatusText(http.StatusBadRequest),
			Detail:     "Stock Id must be provided",
		}, w, r)

		return
	}

	vehicle, err := server.vehicleService.GetVehicleByStockId(server.ctx, stockId)
	if err != nil {
		log.Printf("ERROR >> %s", err.Error())
		respondJSON(http.StatusInternalServerError, errorResponse{
			StatusCode: http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError),
			Detail:     "Something went wrong",
		}, w, r)

		return
	}

	if vehicle == nil {
		log.Printf("WARNING >> %s", "No vehicle data was found")
		respondJSON(http.StatusNotFound, errorResponse{
			StatusCode: http.StatusNotFound,
			StatusText: http.StatusText(http.StatusNotFound),
			Detail:     "No vehicle data was found",
		}, w, r)

		return
	}

	respondJSON(http.StatusOK, vehicle, w, r)
}
