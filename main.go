package main

import (
	"log"
	"net/http"

	"github.com/stripathy/demo-vehicle-api/dataaccess"
	"github.com/stripathy/demo-vehicle-api/routes"
	"github.com/stripathy/demo-vehicle-api/vehicles"
)

//APIServer models the server object
type APIServer struct {
	vehicleResource vehicles.IVehicleResource
}

func main() {
	server := initialize()
	server.run()
}

func (apiServer *APIServer) run() {
	log.Printf("Starting the API Server")
	http.ListenAndServe(":3001", apiServer.handler())
}

func (apiServer *APIServer) handler() http.Handler {
	apiServer.ConfigureVehicleRoutes()

	r := routes.NewRouter()
	r.HandleFunc("/ping", pingHandler)

	return r
}

func initialize() APIServer {
	apiServer := APIServer{
		vehicleResource: vehicles.New(dataaccess.New()),
	}

	return apiServer
}

func pingHandler(w http.ResponseWriter, r *http.Request) {

	log.Printf("Responded: %s\t%s\t%s", r.Method, r.URL.String(), r.RemoteAddr)
	respondJSON(struct {
		IP string `json:"ip"`
	}{r.RemoteAddr}, w)
}
