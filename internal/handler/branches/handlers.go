package handler

import (
	data "thesis.lefler.eu/internal/data/branches"
	e "thesis.lefler.eu/internal/error"
	v1 "thesis.lefler.eu/internal/handler/branches/v1"
	v2 "thesis.lefler.eu/internal/handler/branches/v2"
	v3 "thesis.lefler.eu/internal/handler/branches/v3"
	v4 "thesis.lefler.eu/internal/handler/branches/v4"
	v5 "thesis.lefler.eu/internal/handler/branches/v5"
)

type Handlers struct {
	V1 v1.Handlers
	V2 v2.Handlers
	V3 v3.Handlers
	V4 v4.Handlers
	V5 v5.Handlers
}

func NewHandlers(errors *e.Errors, models *data.Models) Handlers {
	return Handlers{
		V1: v1.NewHandlers(errors, &models.V1),
		V2: v2.NewHandlers(errors, &models.V2),
		V3: v3.NewHandlers(errors, &models.V3),
		V4: v4.NewHandlers(errors, &models.V4),
		V5: v5.NewHandlers(errors, &models.V5),
	}
}