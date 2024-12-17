-- +goose Up
-- +goose StatementBegin
CREATE VIEW people_v1 AS
SELECT
    id,
    name,
    birthdate,
    created_at,
    updated_at,
    version
FROM people;

CREATE VIEW crew_v1 AS
SELECT
    movie_id,
    person_id,
    crew_type,
    role,
    created_at,
    updated_at,
    version
FROM crew;

CREATE OR REPLACE VIEW actors_v1 AS
SELECT
    id,
    name,
    birthdate,
    created_at,
    updated_at,
    version
FROM people_v1;

CREATE OR REPLACE VIEW movie_actors_v1 AS
SELECT
    movie_id,
    person_id as actor_id,
    role,
    created_at,
    updated_at,
    version
FROM crew_v1
WHERE crew_type = 'Actor';

CREATE OR REPLACE VIEW movies_v4 AS
SELECT
  id,
  title,
  release_year,
  genres,
  runtime,
  language,
  created_at,
  updated_at,
  version
FROM movies;

CREATE OR REPLACE VIEW movies_v3 AS
SELECT
  id,
  title,
  release_year,
  genres,
  COALESCE((
    SELECT name FROM people_v1 p
    JOIN crew_v1 c ON p.id = c.person_id
    WHERE c.movie_id = m.id AND c.crew_type = 'Director'
    LIMIT 1
  ), NULL) AS director,
  runtime,
  language,
  created_at,
  updated_at,
  version
FROM movies_v4 m;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE VIEW movies_v3 AS
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

DROP VIEW IF EXISTS movies_v4;

CREATE OR REPLACE VIEW actors_v1 AS
SELECT
    id,
    name,
    birthdate,
    created_at,
    updated_at,
    version
FROM actors;

CREATE OR REPLACE VIEW movie_actors_v1 AS
SELECT
    movie_id,
    actor_id,
    role,
    created_at,
    updated_at,
    version
FROM movie_actors;

DROP VIEW IF EXISTS crew_v1;
DROP VIEW IF EXISTS people_v1;

-- +goose StatementEnd
