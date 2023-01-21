package router

import "net/http"

// Route defines a API route
type Route struct {
	Name       string
	HTTPVerb   string
	Pattern    string
	Scope      string
	HandleFunc http.HandlerFunc
}

// Routes holds all the routes defined
var Routes []Route
