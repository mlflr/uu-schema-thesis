package data

import (
	"database/sql"

	views "thesis.lefler.eu/internal/data/views"
)

type Models struct {
	Views views.Models
}

func NewModels(db *sql.DB) Models {
	return Models{
		Views: views.NewModels(db),
	}
}
