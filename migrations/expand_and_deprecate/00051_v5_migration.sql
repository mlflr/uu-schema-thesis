-- +goose Up
-- +goose StatementBegin
-- Insert directors into the people table
INSERT INTO people (name)
SELECT DISTINCT director
FROM movies
WHERE director IS NOT NULL;

-- Link directors to movies in the crew table
INSERT INTO crew (movie_id, person_id, crew_type)
SELECT m.id, p.id, 'Director'
FROM movies m
JOIN people p ON m.director = p.name;


ALTER TABLE people ADD COLUMN old_actor_id bigint;

-- Insert actors into the people table
INSERT INTO people (old_actor_id, name, birthdate, created_at, updated_at, version)
SELECT id, name, birthdate, created_at, updated_at, version
FROM actors
RETURNING old_actor_id, id AS new_person_id;

-- Use the mapping table to correctly insert data into the crew table
INSERT INTO crew (movie_id, person_id, crew_type, role, created_at)
SELECT ma.movie_id, p.id, 'Actor', ma.role, ma.created_at
FROM movie_actors ma
JOIN people p
ON ma.actor_id = p.old_actor_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM crew;
DELETE FROM people;
-- +goose StatementEnd