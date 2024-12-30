-- +goose Up
-- +goose StatementBegin
CREATE VIEW movies_v2 AS
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
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS movies_v2;
-- +goose StatementEnd
