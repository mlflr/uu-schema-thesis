package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) routesBranches(router *httprouter.Router) {
	// v1/movies routes
	registerMovieRoutes(router, "branches", "v1", &app.handlers.Branches.V1.Movies)

	// v2/movies routes
	registerMovieRoutes(router, "branches", "v2", &app.handlers.Branches.V2.Movies)

	// v3/movies routes
	registerMovieRoutes(router, "branches", "v3", &app.handlers.Branches.V3.Movies)

	// v4/movies and actors routes
	registerMovieRoutes(router, "branches", "v4", &app.handlers.Branches.V4.Movies)
	registerActorRoutes(router, "branches", "v4", &app.handlers.Branches.V4.Actors)

	// v5/movies and people routes
	registerMovieRoutes(router, "branches", "v5", &app.handlers.Branches.V5.Movies)
	registerPersonRoutes(router, "branches", "v5", &app.handlers.Branches.V5.People)
}
