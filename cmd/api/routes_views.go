package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routesViews(router *httprouter.Router) {
	// v1/movies routes
	v1Handler := &app.handlers.Views.V1.Movies
	router.HandlerFunc(http.MethodGet, "/views/v1/movies", v1Handler.ListMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/views/v1/movies", v1Handler.CreateMovieHandler)
	router.HandlerFunc(http.MethodGet, "/views/v1/movies/:id", v1Handler.GetMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/views/v1/movies/:id", v1Handler.UpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/views/v1/movies/:id", v1Handler.DeleteMovieHandler)

	// v2/movies routes
	v2Handler := &app.handlers.Views.V2.Movies
	router.HandlerFunc(http.MethodGet, "/views/v2/movies", v2Handler.ListMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/views/v2/movies", v2Handler.CreateMovieHandler)
	router.HandlerFunc(http.MethodGet, "/views/v2/movies/:id", v2Handler.GetMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/views/v2/movies/:id", v2Handler.UpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/views/v2/movies/:id", v2Handler.DeleteMovieHandler)

	// v3/movies routes
	v3Handler := &app.handlers.Views.V3.Movies
	router.HandlerFunc(http.MethodGet, "/views/v3/movies", v3Handler.ListMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/views/v3/movies", v3Handler.CreateMovieHandler)
	router.HandlerFunc(http.MethodGet, "/views/v3/movies/:id", v3Handler.GetMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/views/v3/movies/:id", v3Handler.UpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/views/v3/movies/:id", v3Handler.DeleteMovieHandler)

	// v4/movies routes
	v4Handler := &app.handlers.Views.V4.Movies
	v4ActorHandler := &app.handlers.Views.V4.Actors
	router.HandlerFunc(http.MethodGet, "/views/v4/movies", v4Handler.ListMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/views/v4/movies", v4Handler.CreateMovieHandler)
	router.HandlerFunc(http.MethodGet, "/views/v4/movies/:id", v4Handler.GetMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/views/v4/movies/:id", v4Handler.UpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/views/v4/movies/:id", v4Handler.DeleteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/views/v4/actors", v4ActorHandler.ListActorsHandler)
	router.HandlerFunc(http.MethodPost, "/views/v4/actors", v4ActorHandler.CreateActorHandler)
	router.HandlerFunc(http.MethodGet, "/views/v4/actors/:id", v4ActorHandler.GetActorHandler)
	router.HandlerFunc(http.MethodPatch, "/views/v4/actors/:id", v4ActorHandler.UpdateActorHandler)
	router.HandlerFunc(http.MethodDelete, "/views/v4/actors/:id", v4ActorHandler.DeleteActorHandler)

	// v5/movies routes
	v5Handler := &app.handlers.Views.V5.Movies
	v5PersonHandler := &app.handlers.Views.V5.People
	router.HandlerFunc(http.MethodGet, "/views/v5/movies", v5Handler.ListMoviesHandler)
	router.HandlerFunc(http.MethodPost, "/views/v5/movies", v5Handler.CreateMovieHandler)
	router.HandlerFunc(http.MethodGet, "/views/v5/movies/:id", v5Handler.GetMovieHandler)
	router.HandlerFunc(http.MethodPatch, "/views/v5/movies/:id", v5Handler.UpdateMovieHandler)
	router.HandlerFunc(http.MethodDelete, "/views/v5/movies/:id", v5Handler.DeleteMovieHandler)
	router.HandlerFunc(http.MethodGet, "/views/v5/people", v5PersonHandler.ListPeopleHandler)
	router.HandlerFunc(http.MethodPost, "/views/v5/people", v5PersonHandler.CreatePersonHandler)
	router.HandlerFunc(http.MethodGet, "/views/v5/people/:id", v5PersonHandler.GetPersonHandler)
	router.HandlerFunc(http.MethodPatch, "/views/v5/people/:id", v5PersonHandler.UpdatePersonHandler)
	router.HandlerFunc(http.MethodDelete, "/views/v5/people/:id", v5PersonHandler.DeletePersonHandler)

}