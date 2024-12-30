-- +goose Up
-- +goose StatementBegin
-- Insert directors into the people table
INSERT INTO people (name)
SELECT DISTINCT director
FROM movies_branch_v2
WHERE director IS NOT NULL;

-- Link directors to movies in the crew table
INSERT INTO crew (movie_id, person_id, crew_type)
SELECT m.id, p.id, 'Director'
FROM movies_branch_v2 m
JOIN people p ON m.director = p.name;

-- Create a temporary mapping table to link old actor IDs to new person IDs
CREATE TEMP TABLE actor_to_person_map (
  old_actor_id bigint,    -- Old actor ID from actors table
  new_person_id bigint    -- New person ID from people table
);

-- temporarily alter the people table to add a column to store the old actor ID
ALTER TABLE people ADD COLUMN old_actor_id bigint;

-- Insert actors into the people table
WITH inserted AS (
  INSERT INTO people (old_actor_id, name, birthdate, created_at, updated_at, version)
  SELECT id, name, birthdate, created_at, updated_at, version
  FROM actors
  RETURNING old_actor_id, id AS new_person_id
)
INSERT INTO actor_to_person_map (old_actor_id, new_person_id)
SELECT old_actor_id, new_person_id FROM inserted;

-- Remove the old actor ID column from the people table
ALTER TABLE people DROP COLUMN old_actor_id;

-- Use the mapping table to correctly insert data into the crew table
INSERT INTO crew (movie_id, person_id, crew_type, role, created_at)
SELECT ma.movie_id, map.new_person_id, 'Actor', ma.role, ma.created_at
FROM movie_actors ma
JOIN actor_to_person_map map
ON ma.actor_id = map.old_actor_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM crew;
DELETE FROM people;
-- +goose StatementEnd
