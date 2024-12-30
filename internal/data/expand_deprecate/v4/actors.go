package v4

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"cloud.google.com/go/civil"

	"thesis.lefler.eu/internal/validator"
)

type Actor struct {
	ID        int64       `json:"id"`                  // Unique identifier
	CreatedAt time.Time   `json:"-"`                   // Timestamp of when the movie was added to the database
	UpdatedAt time.Time   `json:"-"`                   // Timestamp of when the movie record was last updated
	Name      string      `json:"name"`                // Actor name
	Birthdate *civil.Date `json:"birthdate,omitempty"` // Actor birthdate
	Version   int32       `json:"version"`             // Version number, starts at 1 and increments each time the movie is updated
}

type ActorModel struct {
	DB *sql.DB
}

func ValidateActor(v *validator.Validator, actor *Actor) {
	v.Check(actor.Name != "", "name", "must be provided")
	v.Check(len(actor.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(actor.Birthdate.IsValid(), "birthdate", "is not valid")
	v.Check(actor.Birthdate.Before(civil.DateOf(time.Now())), "birthdate", "must not be in the future")
}

func (m ActorModel) Insert(actor *Actor) error {
	query := `
				INSERT INTO actors (name, birthdate) 
				VALUES ($1, $2)
				RETURNING id, created_at, updated_at, version`

	args := []any{actor.Name, actor.Birthdate.In(time.UTC)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&actor.ID, &actor.CreatedAt, &actor.UpdatedAt, &actor.Version)
}

func (m ActorModel) Get(id int64) (*Actor, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, updated_at, name, birthdate, version
		FROM actors
		WHERE id = $1`

	var actor Actor

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var birthdate *time.Time

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&actor.ID,
		&actor.CreatedAt,
		&actor.UpdatedAt,
		&actor.Name,
		&birthdate,
		&actor.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if birthdate != nil {
		civilDate := civil.DateOf(*birthdate)
		actor.Birthdate = &civilDate
	}

	return &actor, nil
}

func (m ActorModel) GetAll() ([]*Actor, error) {
	query := `
		SELECT id, created_at, updated_at, name, birthdate, version
		FROM actors
		ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	actors := []*Actor{}

	for rows.Next() {
		var actor Actor
		var birthdate *time.Time

		err := rows.Scan(
			&actor.ID,
			&actor.CreatedAt,
			&actor.UpdatedAt,
			&actor.Name,
			&birthdate,
			&actor.Version)

		if err != nil {
			return nil, err
		}

		if birthdate != nil {
			civilDate := civil.DateOf(*birthdate)
			actor.Birthdate = &civilDate
		}
		actors = append(actors, &actor)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}

func (m ActorModel) Update(actor *Actor) error {
	query := `
		UPDATE actors
		SET name = $1, birthdate = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version`

	args := []any{
		actor.Name,
		actor.Birthdate.In(time.UTC),
		actor.ID,
		actor.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&actor.Version)
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

func (m ActorModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM actors
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
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
