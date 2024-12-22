package v4

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"thesis.lefler.eu/internal/validator"
)

type MovieActor struct {
	MovieID   int64     `json:"-"`
	ActorID   int64     `json:"actor_id"`
	ActorName string    `json:"actor_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Version   int32     `json:"-"`
}

type MovieActorModel struct {
	DB *sql.DB
}

func ValidateCrew(v *validator.Validator, movieActor *MovieActor) {
	v.Check(movieActor.MovieID > 0, "movie_id", "must be a positive integer")
	v.Check(movieActor.ActorID > 0, "actor_id", "must be a positive integer")

	v.Check(movieActor.Role != "", "role", "must be provided")
	v.Check(len(movieActor.Role) <= 500, "role", "must not be more than 500 bytes long")

}

func (m MovieActorModel) Insert(movieActor *MovieActor) error {
	query := `
		INSERT INTO movie_actors_v1 (movie_id, actor_id, role)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at, version`

	args := []interface{}{movieActor.MovieID, movieActor.ActorID, movieActor.Role}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movieActor.CreatedAt, &movieActor.UpdatedAt, &movieActor.Version)
}

func (m MovieActorModel) Get(movieID, actorID int64) (*MovieActor, error) {
	if movieID < 1 || actorID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT ma.movie_id, ma.actor_id, a.name AS actor_name, ma.role, ma.created_at, ma.updated_at, ma.version
		FROM movie_actors_v1 ma
		LEFT JOIN actors_v1 a ON ma.actor_id = a.id
		WHERE ma.movie_id = $1 AND ma.actor_id = $2`

	var movieActor MovieActor

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, movieID, actorID).Scan(
		&movieActor.MovieID,
		&movieActor.ActorID,
		&movieActor.ActorName,
		&movieActor.Role,
		&movieActor.CreatedAt,
		&movieActor.UpdatedAt,
		&movieActor.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movieActor, nil
}

func (m MovieActorModel) GetForMovie(movieID int64) ([]*MovieActor, error) {
	if movieID < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT ma.movie_id, ma.actor_id, a.name AS actor_name, ma.role, ma.created_at, ma.updated_at, ma.version
		FROM movie_actors_v1 ma
		LEFT JOIN actors_v1 a ON ma.actor_id = a.id
		WHERE ma.movie_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movieActors := []*MovieActor{}

	for rows.Next() {
		var movieActor MovieActor

		err := rows.Scan(
			&movieActor.MovieID,
			&movieActor.ActorID,
			&movieActor.ActorName,
			&movieActor.Role,
			&movieActor.CreatedAt,
			&movieActor.UpdatedAt,
			&movieActor.Version)

		if err != nil {
			return nil, err
		}

		movieActors = append(movieActors, &movieActor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movieActors, nil
}

func (m MovieActorModel) Update(movieActor *MovieActor) error {
	query := `
		UPDATE movie_actors_v1
		SET role = $1, updated_at = CURRENT_TIMESTAMP, version = version + 1
		WHERE movie_id = $2 AND actor_id = $3
		RETURNING updated_at, version`

	args := []interface{}{movieActor.Role, movieActor.MovieID, movieActor.ActorID}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movieActor.UpdatedAt, &movieActor.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m MovieActorModel) Delete(movieID, actorID int64) error {
	if movieID < 1 || actorID < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM movie_actors_v1
		WHERE movie_id = $1 AND actor_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, movieID, actorID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m MovieActorModel) DeleteForMovie(movieID int64) error {
	if movieID < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM movie_actors_v1
		WHERE movie_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, movieID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
