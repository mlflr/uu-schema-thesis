-- +goose Up
-- +goose StatementBegin
-- Wrap the single genre in an array
INSERT INTO movies_branch_v3 (id, genres)
SELECT id, ARRAY [genre]
FROM movies
WHERE genre IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM movies_branch_v3;
-- +goose StatementEnd