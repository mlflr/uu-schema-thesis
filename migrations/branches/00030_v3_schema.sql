-- +goose Up
-- +goose StatementBegin
CREATE TABLE movies_branch_v3 (
    id bigint PRIMARY KEY, -- Unique identifier for each movie
    genres text[]; -- Array of genres (text type)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS movies_branch_v3;
-- +goose StatementEnd