package v4

import (
	data "thesis.lefler.eu/internal/data/v4"
	e "thesis.lefler.eu/internal/error"
)

type Handlers struct {
	Movies MovieHandler
	Actors ActorHandler
}

func NewHandlers(errors *e.Errors, models *data.Models) Handlers {
	return Handlers{
		Movies: MovieHandler{
			errors: errors,
			models: models,
		},
		Actors: ActorHandler{
			errors: errors,
			models: models,
		},
	}
}
