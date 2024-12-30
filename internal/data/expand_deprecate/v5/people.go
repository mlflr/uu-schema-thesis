package v5

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"cloud.google.com/go/civil"

	"thesis.lefler.eu/internal/validator"
)

type Person struct {
	ID        int64       `json:"id"`                  // Unique identifier
	CreatedAt time.Time   `json:"-"`                   // Timestamp of when the movie was added to the database
	UpdatedAt time.Time   `json:"-"`                   // Timestamp of when the movie record was last updated
	Name      string      `json:"name"`                // Person name
	Birthdate *civil.Date `json:"birthdate,omitempty"` // Person birthdate
	Version   int32       `json:"version"`             // Version number, starts at 1 and increments each time the movie is updated
}

type PersonModel struct {
	DB *sql.DB
}

func ValidatePerson(v *validator.Validator, person *Person) {
	v.Check(person.Name != "", "name", "must be provided")
	v.Check(len(person.Name) <= 500, "name", "must not be more than 500 bytes long")

	v.Check(person.Birthdate.IsValid(), "birthdate", "is not valid")
	v.Check(person.Birthdate.Before(civil.DateOf(time.Now())), "birthdate", "must not be in the future")
}

func (m PersonModel) Insert(person *Person) error {
	query := `
				WITH old AS (
					INSERT INTO actors (name, birthdate)
					VALUES ($1, $2)
					RETURNING id, name, birthdate
				)
				INSERT INTO people (name, birthdate, old_actor_id)
				VALUES ($1, $2, (SELECT id FROM old WHERE name = $1 AND birthdate = $2 ))
				RETURNING id, created_at, updated_at, version`

	args := []any{person.Name, person.Birthdate.In(time.UTC)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&person.ID, &person.CreatedAt, &person.UpdatedAt, &person.Version)
}

func (m PersonModel) Get(id int64) (*Person, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, created_at, updated_at, name, birthdate, version
		FROM people
		WHERE id = $1`

	var person Person

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var birthdate *time.Time

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&person.ID,
		&person.CreatedAt,
		&person.UpdatedAt,
		&person.Name,
		&birthdate,
		&person.Version)

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
		person.Birthdate = &civilDate
	}

	return &person, nil
}

func (m PersonModel) GetAll() ([]*Person, error) {
	query := `
		SELECT id, created_at, updated_at, name, birthdate, version
		FROM people
		ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	persons := []*Person{}

	for rows.Next() {
		var person Person
		var birthdate *time.Time

		err := rows.Scan(
			&person.ID,
			&person.CreatedAt,
			&person.UpdatedAt,
			&person.Name,
			&birthdate,
			&person.Version)

		if err != nil {
			return nil, err
		}

		if birthdate != nil {
			civilDate := civil.DateOf(*birthdate)
			person.Birthdate = &civilDate
		}
		persons = append(persons, &person)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return persons, nil
}

func (m PersonModel) Update(person *Person) error {
	query := `
		WITH new AS (
			UPDATE people
			SET name = $1, birthdate = $2, version = version + 1
			WHERE id = $3 AND version = $4
			RETURNING old_actor_id, name, birthdate, version
		),
		old AS (UPDATE actors
		SET name = $1, birthdate = $2, version = version + 1
		WHERE id = (SELECT old_actor_id FROM new WHERE name = $1 AND birthdate = $2))
		SELECT version from new`

	args := []any{
		person.Name,
		person.Birthdate.In(time.UTC),
		person.ID,
		person.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&person.Version)
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

func (m PersonModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM people
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
