-- +goose Up
-- +goose StatementBegin
ALTER TABLE movies
ADD COLUMN genres text[]; -- Array of genres (text type)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE movies
DROP COLUMN IF EXISTS genres;
-- +goose StatementEnd
