package v5

import (
	data "thesis.lefler.eu/internal/data/branches/v5"
	e "thesis.lefler.eu/internal/error"
)

type Handlers struct {
	Movies MovieHandler
	People PersonHandler
}

func NewHandlers(errors *e.Errors, models *data.Models) Handlers {
	return Handlers{
		Movies: MovieHandler{
			errors: errors,
			models: models,
		},
		People: PersonHandler{
			errors: errors,
			models: models,
		},
	}
}
