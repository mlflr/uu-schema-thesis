-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION handle_movies_write_v3()
RETURNS TRIGGER AS $$
DECLARE
  person_id people_v1.id%TYPE;
BEGIN
    IF TG_OP = 'INSERT' THEN
      INSERT INTO movies (title, release_year, genres, runtime, language, version)
      VALUES (
        NEW.title,
        NEW.release_year,
        NEW.genres,
        NEW.runtime,
        NEW.language,
        COALESCE(NEW.version, 1)
      )
      RETURNING id, title, release_year, genres, NEW.director as director, runtime, language, created_at, updated_at, version
      INTO NEW;
    ELSIF TG_OP = 'UPDATE' THEN
      UPDATE movies
      SET title = NEW.title,
        release_year = NEW.release_year,
        genres = NEW.genres,
        runtime = NEW.runtime,
        language = NEW.language,
        version = NEW.version
      WHERE id = OLD.id
      RETURNING id, title, release_year, genres, NEW.director as director, runtime, language, created_at, updated_at, version
      INTO NEW;
    END IF;

    -- Insert the director into the people table if not already present
    -- We can only match by name here, dq risk of matching wrong person with same name,
    -- but we don't have enough data to match 100%
    WITH s AS (
        SELECT id FROM people_v1 WHERE name = NEW.director LIMIT 1
    ), i as (
        INSERT INTO people_v1 (name)
        SELECT NEW.director
        WHERE NOT EXISTS (SELECT 1 FROM s)
        RETURNING id
    )
    SELECT id FROM i
    UNION ALL
    SELECT id FROM s
    INTO person_id;

    -- Link the movie and director in the crew table
    INSERT INTO crew (movie_id, person_id, crew_type)
    VALUES (
        NEW.id,
        person_id,
        'Director'
    )
    ON CONFLICT DO NOTHING;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_movies_v3_write
INSTEAD OF INSERT OR UPDATE ON movies_v3
FOR EACH ROW EXECUTE FUNCTION handle_movies_write_v3();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_movies_v3_write ON movies_v3;
DROP FUNCTION IF EXISTS handle_movies_write_v3();
-- +goose StatementEnd
