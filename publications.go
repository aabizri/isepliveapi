package main

import (
	"github.com/gorilla/mux"
)

func establishPublicationsRouter(router *mux.Router) {
	var routes = Routes{
		Route{
			"Publish",
			"POST",
			"/",
			postPublications,
		},
		Route{
			"Available Groups",
			"GET",
			"/groups",
			getGroups,
		},
	}

	publicationsRouter := router.PathPrefix("/publications").Subrouter()
	route(routes, publicationsRouter)
}
