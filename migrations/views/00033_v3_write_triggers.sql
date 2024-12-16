-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION handle_movies_write()
RETURNS TRIGGER AS $$
DECLARE
    new_rec movies%ROWTYPE; -- Variable to store new record
BEGIN
    new_rec.genres := ARRAY[NEW.genre];
    IF TG_TABLE_NAME = 'movies_v2' THEN
      new_rec.director := NEW.director;
      new_rec.runtime := NEW.runtime;
      new_rec.language := NEW.language;
    END IF;
    
    IF TG_OP = 'INSERT' THEN
      INSERT INTO movies (title, release_year, genres, director, runtime, language, version)
      VALUES (
        NEW.title,
        NEW.release_year,
        new_rec.genres,
        new_rec.director,
        new_rec.runtime,
        new_rec.language,
        COALESCE(NEW.version, 1)
      )
      RETURNING * INTO new_rec;
    ELSIF TG_OP = 'UPDATE' THEN
      UPDATE movies
      SET title = NEW.title,
        release_year = NEW.release_year,
        -- merge singular v1/v2 genre to existing array to prevent losing data and remove duplicates
        genres = ARRAY(SELECT DISTINCT g FROM unnest(new_rec.genres || genres) as g),
        -- coalesce on v2 columns to prevent removal when updated from v1
        director = COALESCE(new_rec.director, director),
        runtime = COALESCE(new_rec.runtime, runtime),
        language = COALESCE(new_rec.language, language),
        version = NEW.version
      WHERE id = OLD.id
      RETURNING * INTO new_rec;
    ELSIF TG_OP = 'DELETE' THEN
      DELETE FROM movies WHERE id = OLD.id;
      RETURN OLD;
    END IF;

    -- these columns could have changed in the process, make sure we return the up to date version
    NEW.id := new_rec.id;
    NEW.genre := new_rec.genres[1];
    NEW.created_at := new_rec.created_at;
    NEW.updated_at := new_rec.updated_at;
    NEW.version := new_rec.version;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_movies_v1_write
INSTEAD OF INSERT OR UPDATE ON movies_v1
FOR EACH ROW EXECUTE FUNCTION handle_movies_write();

CREATE TRIGGER trg_movies_v2_write
INSTEAD OF INSERT OR UPDATE ON movies_v2
FOR EACH ROW EXECUTE FUNCTION handle_movies_write();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_movies_v1_write ON movies_v1;
DROP TRIGGER IF EXISTS trg_movies_v2_write ON movies_v2;
DROP FUNCTION IF EXISTS handle_movies_view_write();
-- +goose StatementEnd
