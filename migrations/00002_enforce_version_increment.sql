-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION enforce_version_increment()
RETURNS TRIGGER AS $$
BEGIN
    -- Check if the version column is incremented
    IF NEW.version <= OLD.version THEN
        RAISE EXCEPTION 'Updates must increment the version column.';
    END IF;

    -- Allow the update if the version is incremented
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_enforce_version_increment
BEFORE UPDATE ON movies
FOR EACH ROW EXECUTE FUNCTION enforce_version_increment();

CREATE TRIGGER trg_enforce_version_increment
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION enforce_version_increment();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER trg_enforce_version_increment ON movies;
DROP TRIGGER trg_enforce_version_increment ON users;
DROP FUNCTION enforce_version_increment();
-- +goose StatementEnd
