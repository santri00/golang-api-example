package api

import (
	"github.com/santri00/golang-api-example/pkg/router"
	"log"
	"net/http"
)

// adds routes related to ping handler
func (apiServer *Server) configurePingHandlerRoutes() {
	router.Routes = append(router.Routes, router.Route{
		Name:       "Ping",
		HTTPVerb:   "GET",
		Pattern:    "/inventory/v1/ping",
		HandleFunc: apiServer.ping,
	})
}

// Handles the ping response
func (apiServer *Server) ping(w http.ResponseWriter, r *http.Request) {
	log.Printf("Responded: %s\t%s\t%s", r.Method, r.URL.String(), r.RemoteAddr)
	respondJSON(struct {
		IP string `json:"ip"`
	}{r.RemoteAddr}, w)
}
