package routes

import "github.com/gorilla/mux"

//NewRouter instantiates a new mux router
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	//TODO add authentication, Authorization, Throttling, Logging into pipeline
	for _, route := range Routes {
		router.Methods(route.HTTPVerb, "OPTIONS").
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandleFunc)

	}

	return router
}
