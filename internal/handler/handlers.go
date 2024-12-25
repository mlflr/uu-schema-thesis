package handler

import (
	data "thesis.lefler.eu/internal/data"
	e "thesis.lefler.eu/internal/error"
	views "thesis.lefler.eu/internal/handler/views"
)

type Handlers struct {
	Views views.Handlers
}

func NewHandlers(errors *e.Errors, models *data.Models) Handlers {
	return Handlers{
		Views: views.NewHandlers(errors, &models.Views),
	}
}
