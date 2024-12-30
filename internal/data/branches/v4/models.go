package v4

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies      MovieModel
	Actors      ActorModel
	MovieActors MovieActorModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      MovieModel{DB: db},
		Actors:      ActorModel{DB: db},
		MovieActors: MovieActorModel{DB: db},
	}
}
