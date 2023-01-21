package api

import (
	"context"
	"github.com/santri00/golang-api-example/pkg/router"
	"github.com/santri00/golang-api-example/pkg/services/vehicle"
	"log"
	"net/http"
)

// Server Defines the API Server construct
type Server struct {
	ctx            context.Context
	vehicleService vehicle.Service
}

// NewServer constructs the API server instance
func NewServer(ctx context.Context,
	vehicleService vehicle.Service) *Server {
	return &Server{
		ctx:            ctx,
		vehicleService: vehicleService,
	}
}

// Run starts the api server
func (apiServer *Server) Run() {
	log.Printf("Starting the API Server")
	http.ListenAndServe(":3001", apiServer.addHandlers())
}

// adds api handlers/routes
func (apiServer *Server) addHandlers() http.Handler {
	apiServer.configureVehicleHandlerRoutes()
	apiServer.configurePingHandlerRoutes()

	handler := router.NewRouter()

	return handler
}
