package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type movieHandler interface {
	ListMoviesHandler(w http.ResponseWriter, r *http.Request)
	CreateMovieHandler(w http.ResponseWriter, r *http.Request)
	GetMovieHandler(w http.ResponseWriter, r *http.Request)
	UpdateMovieHandler(w http.ResponseWriter, r *http.Request)
	DeleteMovieHandler(w http.ResponseWriter, r *http.Request)
}

type actorHandler interface {
	ListActorsHandler(w http.ResponseWriter, r *http.Request)
	CreateActorHandler(w http.ResponseWriter, r *http.Request)
	GetActorHandler(w http.ResponseWriter, r *http.Request)
	UpdateActorHandler(w http.ResponseWriter, r *http.Request)
	DeleteActorHandler(w http.ResponseWriter, r *http.Request)
}

type personHandler interface {
	ListPeopleHandler(w http.ResponseWriter, r *http.Request)
	CreatePersonHandler(w http.ResponseWriter, r *http.Request)
	GetPersonHandler(w http.ResponseWriter, r *http.Request)
	UpdatePersonHandler(w http.ResponseWriter, r *http.Request)
	DeletePersonHandler(w http.ResponseWriter, r *http.Request)
}

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

func registerMovieRoutes(router *httprouter.Router, prefix string, version string, handler movieHandler) {
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("%s/%s/movies", prefix, version), handler.ListMoviesHandler)
	router.HandlerFunc(http.MethodPost, fmt.Sprintf("%s/%s/movies", prefix, version), handler.CreateMovieHandler)
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("%s/%s/movies/:id", prefix, version), handler.GetMovieHandler)
	router.HandlerFunc(http.MethodPatch, fmt.Sprintf("%s/%s/movies/:id", prefix, version), handler.UpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, fmt.Sprintf("%s/%s/movies/:id", prefix, version), handler.DeleteMovieHandler)
}

func registerActorRoutes(router *httprouter.Router, prefix string, version string, handler actorHandler) {
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("%s/%s/actors", prefix, version), handler.ListActorsHandler)
	router.HandlerFunc(http.MethodPost, fmt.Sprintf("%s/%s/actors", prefix, version), handler.CreateActorHandler)
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("%s/%s/actors/:id", prefix, version), handler.GetActorHandler)
	router.HandlerFunc(http.MethodPatch, fmt.Sprintf("%s/%s/actors/:id", prefix, version), handler.UpdateActorHandler)
	router.HandlerFunc(http.MethodDelete, fmt.Sprintf("%s/%s/actors/:id", prefix, version), handler.DeleteActorHandler)
}

func registerPersonRoutes(router *httprouter.Router, prefix string, version string, handler personHandler) {
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("%s/%s/people", prefix, version), handler.ListPeopleHandler)
	router.HandlerFunc(http.MethodPost, fmt.Sprintf("%s/%s/people", prefix, version), handler.CreatePersonHandler)
	router.HandlerFunc(http.MethodGet, fmt.Sprintf("%s/%s/people/:id", prefix, version), handler.GetPersonHandler)
	router.HandlerFunc(http.MethodPatch, fmt.Sprintf("%s/%s/people/:id", prefix, version), handler.UpdatePersonHandler)
	router.HandlerFunc(http.MethodDelete, fmt.Sprintf("%s/%s/people/:id", prefix, version), handler.DeletePersonHandler)
}
