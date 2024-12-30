package data

import (
	"database/sql"

	v1 "thesis.lefler.eu/internal/data/expand_deprecate/v1"
	v2 "thesis.lefler.eu/internal/data/expand_deprecate/v2"
	v3 "thesis.lefler.eu/internal/data/expand_deprecate/v3"
	v4 "thesis.lefler.eu/internal/data/expand_deprecate/v4"
	v5 "thesis.lefler.eu/internal/data/expand_deprecate/v5"
)

type Models struct {
	V1 v1.Models
	V2 v2.Models
	V3 v3.Models
	V4 v4.Models
	V5 v5.Models
}

func NewModels(db *sql.DB) Models {
	return Models{
		V1: v1.NewModels(db),
		V2: v2.NewModels(db),
		V3: v3.NewModels(db),
		V4: v4.NewModels(db),
		V5: v5.NewModels(db),
	}
}
