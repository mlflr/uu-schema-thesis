package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) routesBranches(router *httprouter.Router) {
	// v1/movies routes
	registerRoutes(router, "branches", "v1", "movies", &app.handlers.Branches.V1.Movies)

	// v2/movies routes
	registerRoutes(router, "branches", "v2", "movies", &app.handlers.Branches.V2.Movies)

	// v3/movies routes
	registerRoutes(router, "branches", "v3", "movies", &app.handlers.Branches.V3.Movies)

	// v4/movies and actors routes
	registerRoutes(router, "branches", "v4", "movies", &app.handlers.Branches.V4.Movies)
	registerRoutes(router, "branches", "v4", "actors", &app.handlers.Branches.V4.Actors)

	// v5/movies and people routes
	registerRoutes(router, "branches", "v5", "movies", &app.handlers.Branches.V5.Movies)
	registerRoutes(router, "branches", "v5", "people", &app.handlers.Branches.V5.People)
}
