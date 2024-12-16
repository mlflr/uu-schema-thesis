-- +goose Up
-- +goose StatementBegin
CREATE TABLE movies (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY, -- Unique identifier for each movie
    title text NOT NULL,                -- Title of the movie
    release_year integer NOT NULL,      -- Year the movie was released
    genre text,                         -- Genre of the movie (e.g., Action, Drama)
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), -- Timestamp of when the record was created
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(), -- Timestamp of when the record was last updated
    version integer DEFAULT 1 NOT NULL  -- Optimistic concurrency control
);

ALTER TABLE movies ADD CONSTRAINT movies_year_check CHECK (release_year BETWEEN 1888 AND date_part('year', now()));

-- On update trigger for data quality checks
CREATE OR REPLACE FUNCTION on_update()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the version column is incremented
    IF NEW.version <= OLD.version THEN
        RAISE EXCEPTION 'Updates must increment the version column.';
    END IF;

    -- Check if the created_at column is modified
    IF NEW.created_at <> OLD.created_at THEN
        RAISE WARNING 'Changing the created_at column is not allowed.';
        NEW.created_at = OLD.created_at;
    END IF;

    -- Update the updated_at column
    NEW.updated_at = NOW();
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger to automatically run the on_update function on row updates
CREATE TRIGGER trg_on_update
BEFORE UPDATE ON movies
FOR EACH ROW
EXECUTE FUNCTION on_update();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_on_update ON movies;
DROP FUNCTION IF EXISTS on_update();
DROP TABLE IF EXISTS movies;
-- +goose StatementEnd
