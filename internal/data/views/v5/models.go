package v5

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies MovieModel
	People PersonModel
	Crew   CrewModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
		People: PersonModel{DB: db},
		Crew:   CrewModel{DB: db},
	}
}
