package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.errors.NotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.errors.MethodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	app.routesViews(router)           // views routes
	app.routesExpandDeprecate(router) // expand_deprecate routes
	app.routesBranches(router)        // branches routes

	return app.recoverPanic(router)
}
