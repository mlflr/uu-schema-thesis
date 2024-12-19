package v2

import (
	data "thesis.lefler.eu/internal/data/v2"
	e "thesis.lefler.eu/internal/error"
)

type Handlers struct {
	Movies MovieHandler
}

func NewHandlers(errors *e.Errors, models *data.Models) Handlers {
	return Handlers{
		Movies: MovieHandler{
			errors: errors,
			models: models,
		},
	}
}
