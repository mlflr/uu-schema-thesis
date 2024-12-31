-- +goose Up
-- +goose StatementBegin
CREATE TABLE actors (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- Unique identifier for each actor
    name text NOT NULL,              -- Name of the actor
    birthdate date,                  -- Birthdate of the actor (optional)
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), -- Record creation timestamp
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), -- Timestamp of when the record was last updated
    version integer DEFAULT 1 NOT NULL  -- Optimistic concurrency control
);

CREATE TABLE movie_actors (
    movie_id bigint NOT NULL REFERENCES movies(id) ON DELETE CASCADE, -- Foreign key to movies
    actor_id bigint NOT NULL REFERENCES actors(id) ON DELETE CASCADE, -- Foreign key to actors
    role text NOT NULL,                 -- Role played in the movie
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),    -- Record creation timestamp
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), -- Timestamp of when the record was last updated
    version integer DEFAULT 1 NOT NULL,  -- Optimistic concurrency control
    PRIMARY KEY (movie_id, actor_id)    -- Composite primary key
);

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

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_actors_on_update ON actors;
DROP TRIGGER IF EXISTS trg_movie_actors_on_update ON movie_actors;
DROP TABLE IF EXISTS movie_actors;
DROP TABLE IF EXISTS actors;
-- +goose StatementEnd
