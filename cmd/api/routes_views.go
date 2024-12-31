package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) routesViews(router *httprouter.Router) {
	// v1/movies routes
	registerMovieRoutes(router, "views", "v1", &app.handlers.Views.V1.Movies)

	// v2/movies routes
	registerMovieRoutes(router, "views", "v2", &app.handlers.Views.V2.Movies)

	// v3/movies routes
	registerMovieRoutes(router, "views", "v3", &app.handlers.Views.V3.Movies)

	// v4/movies and actors routes
	registerMovieRoutes(router, "views", "v4", &app.handlers.Views.V4.Movies)
	registerActorRoutes(router, "views", "v4", &app.handlers.Views.V4.Actors)

	// v5/movies and people routes
	registerMovieRoutes(router, "views", "v5", &app.handlers.Views.V5.Movies)
	registerPersonRoutes(router, "views", "v5", &app.handlers.Views.V5.People)
}
