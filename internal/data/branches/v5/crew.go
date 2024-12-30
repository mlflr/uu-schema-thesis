package v5

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"thesis.lefler.eu/internal/validator"
)

type Crew struct {
	MovieID    int64     `json:"-"`
	PersonID   int64     `json:"person_id"`
	PersonName string    `json:"person_name"`
	CrewType   string    `json:"crew_type"`
	Role       string    `json:"role,omitempty"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
	Version    int32     `json:"-"`
}

type CrewModel struct {
	DB *sql.DB
}

func ValidateCrew(v *validator.Validator, crew *Crew) {
	v.Check(crew.PersonID > 0, "person_id", "must be a positive integer")

	v.Check(crew.CrewType != "", "crew_type", "must be provided")
	v.Check(len(crew.CrewType) <= 500, "crew_type", "must not be more than 500 bytes long")
	v.Check(validator.PermittedValue(crew.CrewType, "Actor", "Director", "Producer"), "crew_type", "must be either 'Actor', 'Director', or 'Producer'")

	if strings.ToLower(crew.CrewType) == "actor" {
		v.Check(crew.Role != "", "role", "must be provided if crew_type is 'actor'")
		v.Check(len(crew.Role) <= 500, "role", "must not be more than 500 bytes long")
	}
}

func (m CrewModel) Insert(crew *Crew) error {
	query := `
		WITH old AS (
			INSERT INTO movie_actors (movie_id, actor_id, role)
			SELECT $1, (SELECT old_actor_id FROM people WHERE id = $2), $4
			WHERE EXISTS (SELECT 1 WHERE $3 = 'Actor')
		)
		INSERT INTO crew (movie_id, person_id, crew_type, role)
			VALUES ($1, $2, $3, $4)
			RETURNING created_at, updated_at, version`

	args := []interface{}{crew.MovieID, crew.PersonID, crew.CrewType, crew.Role}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&crew.CreatedAt, &crew.UpdatedAt, &crew.Version)
}

func (m CrewModel) GetForMovie(movieID int64) ([]*Crew, error) {
	if movieID < 1 {
		return nil, ErrRecordNotFound
	}
	// postgres union eliminates duplicates, but the potential null ID is a problem, we should take care of forward synchronization of actors to crew using a trigger or versioning
	query := `
		SELECT c.movie_id, c.person_id, p.name AS person_name, c.crew_type, c.role, c.created_at, c.updated_at, c.version
			FROM crew c
			LEFT JOIN people p ON c.person_id = p.id
			WHERE c.movie_id = $1
		UNION
		SELECT ma.movie_id, COALESCE(p.id, NULL) AS person_id, COALESCE(p.name, a.name) AS person_name, 'Actor' AS crew_type, ma.role, ma.created_at, ma.updated_at, ma.version
			FROM movie_actors ma
			LEFT JOIN people p ON ma.actor_id = p.old_actor_id
			LEFT JOIN actors a ON ma.actor_id = a.id
			WHERE ma.movie_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	crews := []*Crew{}

	for rows.Next() {
		var crew Crew

		err := rows.Scan(
			&crew.MovieID,
			&crew.PersonID,
			&crew.PersonName,
			&crew.CrewType,
			&crew.Role,
			&crew.CreatedAt,
			&crew.UpdatedAt,
			&crew.Version)

		if err != nil {
			return nil, err
		}

		crews = append(crews, &crew)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return crews, nil
}

func (m CrewModel) DeleteForMovie(movieID int64) error {
	if movieID < 1 {
		return ErrRecordNotFound
	}

	query := `
		WITH deleted AS (
			DELETE FROM movie_actors
			WHERE movie_id = $1
		)
		DELETE FROM crew
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
