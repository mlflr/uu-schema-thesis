-- +goose Up
-- +goose StatementBegin
CREATE VIEW movies_v1 AS
SELECT
    id,
    title,
    release_year,
    genre,
    created_at,
    updated_at,
    version
FROM movies;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS movies_v1;
-- +goose StatementEnd
