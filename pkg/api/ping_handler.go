package api

import (
	"github.com/santri00/golang-api-example/pkg/router"
	"log"
	"net/http"
)

// adds routes related to ping handler
func (server *Server) configurePingHandlerRoutes() {
	router.Routes = append(router.Routes, router.Route{
		Name:       "Ping",
		HTTPVerb:   "GET",
		Pattern:    "/inventory/v1/ping",
		HandleFunc: server.ping,
	})
}

// Handles the ping response
func (server *Server) ping(w http.ResponseWriter, r *http.Request) {
	//Redis client
	client := server.redisClient
	//ping
	_, e := client.Ping(server.ctx).Result()
	if e != nil {
		log.Printf("ERROR >> %s", e.Error())
		respondJSON(http.StatusInternalServerError, errorResponse{
			StatusCode: http.StatusInternalServerError,
			StatusText: http.StatusText(http.StatusInternalServerError),
			Detail:     "Failed to ping Redis server"}, w, r)
		return
	}

	respondJSON(http.StatusOK, struct {
		IP string `json:"ip"`
	}{r.RemoteAddr}, w, r)
}
