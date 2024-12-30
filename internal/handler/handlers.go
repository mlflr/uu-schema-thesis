package handler

import (
	data "thesis.lefler.eu/internal/data"
	e "thesis.lefler.eu/internal/error"
	branches "thesis.lefler.eu/internal/handler/branches"
	expandDeprecate "thesis.lefler.eu/internal/handler/expand_deprecate"
	views "thesis.lefler.eu/internal/handler/views"
)

type Handlers struct {
	Views           views.Handlers
	ExpandDeprecate expandDeprecate.Handlers
	Branches        branches.Handlers
}

func NewHandlers(errors *e.Errors, models *data.Models) Handlers {
	return Handlers{
		Views:           views.NewHandlers(errors, &models.Views),
		ExpandDeprecate: expandDeprecate.NewHandlers(errors, &models.ExpandDeprecate),
		Branches:        branches.NewHandlers(errors, &models.Branches),
	}
}
