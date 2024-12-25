package v3

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"thesis.lefler.eu/internal/validator"
)

// V2 fields still as pointers to allow for null values from pre-v2 records
// genres does not need to be a pointer as it is a slice and can be nil
type Movie struct {
	ID        int64     `json:"id"`                 // Unique identifier
	CreatedAt time.Time `json:"-"`                  // Timestamp of when the movie was added to the database
	UpdatedAt time.Time `json:"-"`                  // Timestamp of when the movie record was last updated
	Title     string    `json:"title"`              // Title of the movie
	Year      int32     `json:"year,omitempty"`     // Release year
	Genres    []string  `json:"genres,omitempty"`   // Slice of genres
	Director  *string   `json:"director,omitempty"` // Director name
	Runtime   *int32    `json:"runtime,omitempty"`  // Runtime in minutes
	Language  *string   `json:"language,omitempty"` // Language
	Version   int32     `json:"version"`            // Version number, starts at 1 and increments each time the movie is updated
}

type MovieModel struct {
	DB *sql.DB
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")

	v.Check(*movie.Director != "", "director", "must be provided")
	v.Check(len(*movie.Director) <= 100, "director", "must not be more than 100 bytes long")

	v.Check(*movie.Runtime != 0, "runtime", "must be provided")
	v.Check(*movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(*movie.Language != "", "language", "must be provided")
	v.Check(len(*movie.Language) <= 50, "language", "must not be more than 50 bytes long")
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `
        INSERT INTO movies_v3 (title, release_year, genres, director, runtime, language) 
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at, version`

	args := []any{movie.Title, movie.Year, pq.Array(movie.Genres), movie.Director, movie.Runtime, movie.Language}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.UpdatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
        SELECT id, created_at, updated_at, title, release_year, genres, director, runtime, language, version
        FROM movies_v3
        WHERE id = $1`

	var movie Movie

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.UpdatedAt,
		&movie.Title,
		&movie.Year,
		pq.Array(&movie.Genres),
		&movie.Director,
		&movie.Runtime,
		&movie.Language,
		&movie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m MovieModel) GetAll() ([]*Movie, error) {
	query := `
        SELECT id, created_at, updated_at, title, release_year, genres, director, runtime, language, version
        FROM movies_v3
        ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	movies := []*Movie{}

	for rows.Next() {
		var movie Movie

		err := rows.Scan(
			&movie.ID,
			&movie.CreatedAt,
			&movie.UpdatedAt,
			&movie.Title,
			&movie.Year,
			pq.Array(&movie.Genres),
			&movie.Director,
			&movie.Runtime,
			&movie.Language,
			&movie.Version,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (m MovieModel) Update(movie *Movie) error {
	query := `
        UPDATE movies_v3
        SET title = $1, release_year = $2, genres = $3, 
					director = $4, runtime = $5, language = $6, 
					version = version + 1
        WHERE id = $7 AND version = $8
        RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		pq.Array(movie.Genres),
		movie.Director,
		movie.Runtime,
		movie.Language,
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)
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

func (m MovieModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
        DELETE FROM movies_v3
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
