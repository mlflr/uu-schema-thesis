package main

import (
	"github.com/julienschmidt/httprouter"
)

func (app *application) routesExpandDeprecate(router *httprouter.Router) {
	// v1/movies routes
	registerRoutes(router, "expand_deprecate", "v1", "movies", &app.handlers.ExpandDeprecate.V1.Movies)

	// v2/movies routes
	registerRoutes(router, "expand_deprecate", "v2", "movies", &app.handlers.ExpandDeprecate.V2.Movies)

	// v3/movies routes
	registerRoutes(router, "expand_deprecate", "v3", "movies", &app.handlers.ExpandDeprecate.V3.Movies)

	// v4/movies and actors routes
	registerRoutes(router, "expand_deprecate", "v4", "movies", &app.handlers.ExpandDeprecate.V4.Movies)
	registerRoutes(router, "expand_deprecate", "v4", "actors", &app.handlers.ExpandDeprecate.V4.Actors)

	// v5/movies and people routes
	registerRoutes(router, "expand_deprecate", "v5", "movies", &app.handlers.ExpandDeprecate.V5.Movies)
	registerRoutes(router, "expand_deprecate", "v5", "people", &app.handlers.ExpandDeprecate.V5.People)
}
