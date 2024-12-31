package handler

import (
	"net/http"

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

type Handler interface {
	CreateHandler(w http.ResponseWriter, r *http.Request)
	GetHandler(w http.ResponseWriter, r *http.Request)
	UpdateHandler(w http.ResponseWriter, r *http.Request)
	DeleteHandler(w http.ResponseWriter, r *http.Request)
	ListHandler(w http.ResponseWriter, r *http.Request)
}

func NewHandlers(errors *e.Errors, models *data.Models) Handlers {
	return Handlers{
		Views:           views.NewHandlers(errors, &models.Views),
		ExpandDeprecate: expandDeprecate.NewHandlers(errors, &models.ExpandDeprecate),
		Branches:        branches.NewHandlers(errors, &models.Branches),
	}
}
