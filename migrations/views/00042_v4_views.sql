-- +goose Up
-- +goose StatementBegin
CREATE VIEW actors_v1 AS
SELECT
    id,
    name,
    birthdate,
    created_at,
    updated_at,
    version
FROM actors;

CREATE VIEW movie_actors_v1 AS
SELECT
    movie_id,
    actor_id,
    role,
    created_at,
    updated_at,
    version
FROM movie_actors;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS movie_actors_v1;
DROP VIEW IF EXISTS actors_v1;
-- +goose StatementEnd
