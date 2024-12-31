-- +goose Up
-- +goose StatementBegin
ALTER TABLE movies
ADD COLUMN director text,   -- Name of the director
ADD COLUMN runtime integer, -- Runtime of the movie in minutes
ADD COLUMN language text;   -- Language of the movie
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE movies
DROP COLUMN IF EXISTS director,
DROP COLUMN IF EXISTS runtime,
DROP COLUMN IF EXISTS language;
-- +goose StatementEnd
