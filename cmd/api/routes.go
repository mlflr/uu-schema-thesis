package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"thesis.lefler.eu/internal/handler"
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

func registerRoutes(router *httprouter.Router, prefix string, version string, resource string, handler handler.Handler) {
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("%s/%s/%s", prefix, version, resource), handler.ListHandler)
	router.HandlerFunc(http.MethodPost, fmt.Sprintf("%s/%s/%s", prefix, version, resource), handler.CreateHandler)
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("%s/%s/%s/:id", prefix, version, resource), handler.GetHandler)
	router.HandlerFunc(http.MethodPatch, fmt.Sprintf("%s/%s/%s/:id", prefix, version, resource), handler.UpdateHandler)
	router.HandlerFunc(http.MethodDelete, fmt.Sprintf("%s/%s/%s/:id", prefix, version, resource), handler.DeleteHandler)
}
