package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) routesViews(router *httprouter.Router) {
	// v1/movies routes
	registerRoutes(router, "views", "v1", "movies", &app.handlers.Views.V1.Movies)

	// v2/movies routes
	registerRoutes(router, "views", "v2", "movies", &app.handlers.Views.V2.Movies)

	// v3/movies routes
	registerRoutes(router, "views", "v3", "movies", &app.handlers.Views.V3.Movies)

	// v4/movies and actors routes
	registerRoutes(router, "views", "v4", "movies", &app.handlers.Views.V4.Movies)
	registerRoutes(router, "views", "v4", "actors", &app.handlers.Views.V4.Actors)

	// v5/movies and people routes
	registerRoutes(router, "views", "v5", "movies", &app.handlers.Views.V5.Movies)
	registerRoutes(router, "views", "v5", "people", &app.handlers.Views.V5.People)
}
