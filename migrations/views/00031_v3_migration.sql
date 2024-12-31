-- +goose Up
-- +goose StatementBegin
-- Wrap the single genre in an array
UPDATE movies
SET genres = ARRAY [genre],
  version = version + 1
WHERE genre IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd