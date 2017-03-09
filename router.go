package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// Create a new router
func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	establishPublicationsRouter(router)
	//establishStudentsRouter(router)

	return router
}

// Helper function for subrouting
func route(routes []Route, router *mux.Router) *mux.Router {

	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
