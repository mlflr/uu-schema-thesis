-- +goose Up
-- +goose StatementBegin
CREATE VIEW movies_v3 AS
SELECT
    id,
    title,
    release_year,
    genres,
    director,
    runtime,
    language,
    created_at,
    updated_at,
    version
FROM movies;

-- Use v3 as the base view for v1 and v2
CREATE OR REPLACE VIEW movies_v1 AS
SELECT
    id,
    title,
    release_year,
    COALESCE(genres[1], '') AS genre,
    created_at,
    updated_at,
    version
FROM movies_v3;

CREATE OR REPLACE VIEW movies_v2 AS
SELECT
    id,
    title,
    release_year,
    COALESCE(genres[1], '') AS genre,
    director,
    runtime,
    language,
    created_at,
    updated_at,
    version
FROM movies_v3;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE VIEW movies_v1 AS
SELECT
    id,
    title,
    release_year,
    genre,
    created_at,
    updated_at,
    version
FROM movies;
CREATE OR REPLACE VIEW movies_v2 AS
SELECT
    id,
    title,
    release_year,
    genre,
    director,
    runtime,
    language,
    created_at,
    updated_at,
    version
FROM movies;
DROP VIEW IF EXISTS movies_v3;
-- +goose StatementEnd
