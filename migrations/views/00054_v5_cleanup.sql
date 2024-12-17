-- +goose Up
-- +goose StatementBegin
-- Remove director column from movies table
ALTER TABLE movies
DROP COLUMN IF EXISTS director;

-- Drop actors and movie_actors tables
DROP TRIGGER IF EXISTS trg_actors_on_update ON actors;
DROP TRIGGER IF EXISTS trg_movie_actors_on_update ON movie_actors;
DROP TABLE IF EXISTS movie_actors;
DROP TABLE IF EXISTS actors;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE movies
ADD COLUMN director text;

-- fill in the director column with the first director from the crew table
UPDATE movies
SET director = COALESCE((
    SELECT name
    FROM people
    JOIN crew ON people.id = crew.person_id
    WHERE crew.movie_id = movies.id
    AND crew.crew_type = 'Director'
    LIMIT 1
), NULL), version = version + 1;

-- Recreate actors and movie_actors tables
CREATE TABLE actors (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- Unique identifier for each actor
    name text NOT NULL,              -- Name of the actor
    birthdate date,                  -- Birthdate of the actor (optional)
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), -- Record creation timestamp
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), -- Timestamp of when the record was last updated
    version integer DEFAULT 1 NOT NULL  -- Optimistic concurrency control
);

-- Copy actors from the people table to the actors table
INSERT INTO actors (id, name, birthdate, created_at, updated_at, version)
OVERRIDING SYSTEM VALUE
SELECT id, name, birthdate, created_at, updated_at, version
FROM people;

CREATE TABLE movie_actors (
    movie_id bigint NOT NULL REFERENCES movies(id) ON DELETE CASCADE, -- Foreign key to movies
    actor_id bigint NOT NULL REFERENCES actors(id) ON DELETE CASCADE, -- Foreign key to actors
    role text NOT NULL,                 -- Role played in the movie
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),    -- Record creation timestamp
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), -- Timestamp of when the record was last updated
    version integer DEFAULT 1 NOT NULL,  -- Optimistic concurrency control
    PRIMARY KEY (movie_id, actor_id)    -- Composite primary key
);

-- Copy actors from the cast table to the movie_actors table
INSERT INTO movie_actors (movie_id, actor_id, role, created_at, updated_at, version)
SELECT movie_id, person_id, role, created_at, updated_at, version
FROM crew
WHERE crew_type = 'Actor';

-- Create a trigger to automatically run the on_update function on row updates
CREATE TRIGGER trg_actors_on_update
BEFORE UPDATE ON actors
FOR EACH ROW
EXECUTE FUNCTION on_update();

CREATE TRIGGER trg_movie_actors_on_update
BEFORE UPDATE ON movie_actors
FOR EACH ROW
EXECUTE FUNCTION on_update();
-- +goose StatementEnd
