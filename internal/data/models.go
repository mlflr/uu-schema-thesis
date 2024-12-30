package data

import (
	"database/sql"

	branches "thesis.lefler.eu/internal/data/branches"
	expandDeprecate "thesis.lefler.eu/internal/data/expand_deprecate"
	views "thesis.lefler.eu/internal/data/views"
)

type Models struct {
	Views           views.Models
	ExpandDeprecate expandDeprecate.Models
	Branches        branches.Models
}

type DbConns struct {
	Views           *sql.DB
	ExpandDeprecate *sql.DB
	Branches        *sql.DB
}

func NewModels(db DbConns) Models {
	return Models{
		Views:           views.NewModels(db.Views),
		ExpandDeprecate: expandDeprecate.NewModels(db.ExpandDeprecate),
		Branches:        branches.NewModels(db.Branches),
	}
}
