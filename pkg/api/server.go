package api

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/santri00/golang-api-example/pkg/apiconfig"
	"github.com/santri00/golang-api-example/pkg/apiconfig/configkeys"
	"github.com/santri00/golang-api-example/pkg/router"
	"github.com/santri00/golang-api-example/pkg/services/vehicle"
	"log"
	"net/http"
)

// Server Defines the API Server construct
type Server struct {
	ctx            context.Context
	cfg            *apiconfig.Config
	redisClient    *redis.Client
	vehicleService vehicle.Service
}

// NewServer constructs the API server instance
func NewServer(ctx context.Context,
	c *apiconfig.Config,
	client *redis.Client,
	vehicleService vehicle.Service) *Server {
	return &Server{
		ctx:            ctx,
		cfg:            c,
		redisClient:    client,
		vehicleService: vehicleService,
	}
}

// Run starts the api server
func (server *Server) Run() {
	log.Printf("Starting the %s Server", server.cfg.Store.Get(configkeys.Name))
	log.Panic(http.ListenAndServe(
		fmt.Sprintf(":%s", server.cfg.Store.Get(configkeys.Port)),
		server.addHandlers()))
}

// adds api handlers/routes
func (server *Server) addHandlers() http.Handler {
	server.configureVehicleHandlerRoutes()
	server.configurePingHandlerRoutes()

	handler := router.NewRouter()

	return handler
}
