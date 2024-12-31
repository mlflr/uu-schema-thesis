-- +goose Up
-- +goose StatementBegin
INSERT INTO movies (title, release_year, genre)
VALUES 
    ('Inception', 2010, 'Sci-Fi'),
    ('The Dark Knight', 2008, 'Action'),
    ('Interstellar', 2014, 'Sci-Fi'),
    ('The Godfather', 1972, 'Crime'),
    ('Pulp Fiction', 1994, 'Drama'),
    ('The Shawshank Redemption', 1994, 'Drama'),
    ('The Matrix', 1999, 'Sci-Fi'),
    ('Gladiator', 2000, 'Action'),
    ('Parasite', 2019, 'Thriller'),
    ('Joker', 2019, 'Drama');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM movies;
-- +goose StatementEnd
