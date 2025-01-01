-- +goose Up
-- +goose StatementBegin
CREATE TABLE people (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,           -- Unique identifier for each person
    old_actor_id bigint,                -- Old actor ID from actors table
    name text NOT NULL,                 -- Name of the person
    birthdate date,                     -- Birthdate of the person (optional)
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), -- Record creation timestamp
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), -- Timestamp of when the record was last updated
    version integer DEFAULT 1 NOT NULL  -- Optimistic concurrency control
);

CREATE TABLE crew (
    movie_id bigint NOT NULL REFERENCES movies(id) ON DELETE CASCADE, -- Foreign key to movies
    person_id bigint NOT NULL REFERENCES people(id) ON DELETE CASCADE, -- Foreign key to people
    crew_type text NOT NULL DEFAULT 'Actor',              -- Role type (e.g., Actor, Director, Producer), default to Actor for compatibility reasons
    role text,                            -- Specific role for actors (e.g., "Maximus Decimus Meridius", "Joker")
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),    -- Record creation timestamp
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), -- Timestamp of when the record was last updated
    version integer DEFAULT 1 NOT NULL,  -- Optimistic concurrency control
    PRIMARY KEY (movie_id, person_id, crew_type) -- Composite key to ensure unique entries
);

-- Create a trigger to automatically run the on_update function on row updates
CREATE TRIGGER trg_people_on_update
BEFORE UPDATE ON people
FOR EACH ROW
EXECUTE FUNCTION on_update();

CREATE TRIGGER trg_crew_on_update
BEFORE UPDATE ON crew
FOR EACH ROW
EXECUTE FUNCTION on_update();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_people_on_update ON people;
DROP TRIGGER IF EXISTS trg_crew_on_update ON crew;
DROP TABLE IF EXISTS crew;
DROP TABLE IF EXISTS people;
-- +goose StatementEnd
