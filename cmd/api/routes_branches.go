package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routesBranches(router *httprouter.Router) {
	// v1/movies routes
	v1Handler := &app.handlers.Branches.V1.Movies
	router.HandlerFunc(http.MethodGet, "/branches/v1/movies", v1Handler.ListMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/branches/v1/movies", v1Handler.CreateMovieHandler)
	router.HandlerFunc(http.MethodGet, "/branches/v1/movies/:id", v1Handler.GetMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/branches/v1/movies/:id", v1Handler.UpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/branches/v1/movies/:id", v1Handler.DeleteMovieHandler)

	// v2/movies routes
	v2Handler := &app.handlers.Branches.V2.Movies
	router.HandlerFunc(http.MethodGet, "/branches/v2/movies", v2Handler.ListMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/branches/v2/movies", v2Handler.CreateMovieHandler)
	router.HandlerFunc(http.MethodGet, "/branches/v2/movies/:id", v2Handler.GetMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/branches/v2/movies/:id", v2Handler.UpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/branches/v2/movies/:id", v2Handler.DeleteMovieHandler)

	// v3/movies routes
	v3Handler := &app.handlers.Branches.V3.Movies
	router.HandlerFunc(http.MethodGet, "/branches/v3/movies", v3Handler.ListMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/branches/v3/movies", v3Handler.CreateMovieHandler)
	router.HandlerFunc(http.MethodGet, "/branches/v3/movies/:id", v3Handler.GetMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/branches/v3/movies/:id", v3Handler.UpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/branches/v3/movies/:id", v3Handler.DeleteMovieHandler)

	// v4/movies routes
	v4Handler := &app.handlers.Branches.V4.Movies
	v4ActorHandler := &app.handlers.Branches.V4.Actors
	router.HandlerFunc(http.MethodGet, "/branches/v4/movies", v4Handler.ListMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/branches/v4/movies", v4Handler.CreateMovieHandler)
	router.HandlerFunc(http.MethodGet, "/branches/v4/movies/:id", v4Handler.GetMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/branches/v4/movies/:id", v4Handler.UpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/branches/v4/movies/:id", v4Handler.DeleteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/branches/v4/actors", v4ActorHandler.ListActorsHandler)
	router.HandlerFunc(http.MethodPost, "/branches/v4/actors", v4ActorHandler.CreateActorHandler)
	router.HandlerFunc(http.MethodGet, "/branches/v4/actors/:id", v4ActorHandler.GetActorHandler)
	router.HandlerFunc(http.MethodPatch, "/branches/v4/actors/:id", v4ActorHandler.UpdateActorHandler)
	router.HandlerFunc(http.MethodDelete, "/branches/v4/actors/:id", v4ActorHandler.DeleteActorHandler)

	// v5/movies routes
	v5Handler := &app.handlers.Branches.V5.Movies
	v5PersonHandler := &app.handlers.Branches.V5.People
	router.HandlerFunc(http.MethodGet, "/branches/v5/movies", v5Handler.ListMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/branches/v5/movies", v5Handler.CreateMovieHandler)
	router.HandlerFunc(http.MethodGet, "/branches/v5/movies/:id", v5Handler.GetMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/branches/v5/movies/:id", v5Handler.UpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/branches/v5/movies/:id", v5Handler.DeleteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/branches/v5/people", v5PersonHandler.ListPeopleHandler)
	router.HandlerFunc(http.MethodPost, "/branches/v5/people", v5PersonHandler.CreatePersonHandler)
	router.HandlerFunc(http.MethodGet, "/branches/v5/people/:id", v5PersonHandler.GetPersonHandler)
	router.HandlerFunc(http.MethodPatch, "/branches/v5/people/:id", v5PersonHandler.UpdatePersonHandler)
	router.HandlerFunc(http.MethodDelete, "/branches/v5/people/:id", v5PersonHandler.DeletePersonHandler)
}