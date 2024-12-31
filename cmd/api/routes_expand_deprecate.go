package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) routesExpandDeprecate(router *httprouter.Router) {
	// v1/movies routes
	registerMovieRoutes(router, "expand_deprecate", "v1", &app.handlers.ExpandDeprecate.V1.Movies)

	// v2/movies routes
	registerMovieRoutes(router, "expand_deprecate", "v2", &app.handlers.ExpandDeprecate.V2.Movies)

	// v3/movies routes
	registerMovieRoutes(router, "expand_deprecate", "v3", &app.handlers.ExpandDeprecate.V3.Movies)

	// v4/movies and actors routes
	registerMovieRoutes(router, "expand_deprecate", "v4", &app.handlers.ExpandDeprecate.V4.Movies)
	registerActorRoutes(router, "expand_deprecate", "v4", &app.handlers.ExpandDeprecate.V4.Actors)

	// v5/movies and people routes
	registerMovieRoutes(router, "expand_deprecate", "v5", &app.handlers.ExpandDeprecate.V5.Movies)
	registerPersonRoutes(router, "expand_deprecate", "v5", &app.handlers.ExpandDeprecate.V5.People)
}
