-- +goose Up
-- +goose StatementBegin
-- Create the v1 movies table
CREATE TABLE movies (
    id SERIAL PRIMARY KEY,          -- Unique identifier for each movie
    title VARCHAR(255) NOT NULL,          -- Title of the movie
    release_year INT NOT NULL,            -- Year the movie was released
    genre VARCHAR(100) NOT NULL,          -- Genre of the movie (e.g., Action, Drama)
    version INT DEFAULT 1 NOT NULL        -- Optimistic concurrency control
);

-- Create the v1 users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,           -- Unique identifier for each user
    username VARCHAR(100) NOT NULL,       -- Name of the user
    email VARCHAR(255) UNIQUE NOT NULL,   -- Email of the user
    version INT DEFAULT 1 NOT NULL        -- Optimistic concurrency control
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the v1 movies table
DROP TABLE movies;
-- Drop the v1 users table
DROP TABLE users;
-- +goose StatementEnd
