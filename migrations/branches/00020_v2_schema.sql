-- +goose Up
-- +goose StatementBegin

CREATE TABLE movies_branch_v2 (
    id bigint PRIMARY KEY NOT NULL REFERENCES movies(id) ON DELETE CASCADE, -- Unique identifier for each movie
    director text,         -- Name of the director
    runtime integer,       -- Runtime of the movie in minutes
    language text         -- Language of the movie
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS movies_branch_v2;
-- +goose StatementEnd
