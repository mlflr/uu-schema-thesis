-- +goose Up
-- +goose StatementBegin
INSERT INTO actors (name, birthdate)
VALUES
  ('Leonardo DiCaprio', '1974-11-11'), -- Inception, The Wolf of Wall Street, etc.
  ('Joseph Gordon-Levitt', '1981-02-17'), -- Inception
  ('Ellen Page', '1987-02-21'), -- Inception
  ('Christian Bale', '1974-01-30'), -- The Dark Knight
  ('Heath Ledger', '1979-04-04'), -- The Dark Knight
  ('Matthew McConaughey', '1969-11-04'), -- Interstellar
  ('Anne Hathaway', '1982-11-12'), -- Interstellar
  ('Marlon Brando', '1924-04-03'), -- The Godfather
  ('Al Pacino', '1940-04-25'), -- The Godfather
  ('Samuel L. Jackson', '1948-12-21'), -- Pulp Fiction
  ('John Travolta', '1954-02-18'), -- Pulp Fiction
  ('Morgan Freeman', '1937-06-01'), -- The Shawshank Redemption
  ('Tim Robbins', '1958-10-16'), -- The Shawshank Redemption
  ('Keanu Reeves', '1964-09-02'), -- The Matrix
  ('Laurence Fishburne', '1961-07-30'), -- The Matrix
  ('Russell Crowe', '1964-04-07'), -- Gladiator
  ('Joaquin Phoenix', '1974-10-28'), -- Gladiator, Joker
  ('Song Kang-ho', '1967-01-17'), -- Parasite
  ('Choi Woo-shik', '1990-03-26'), -- Parasite
  ('Robert De Niro', '1943-08-17'), -- Joker
  ('Zazie Beetz', '1991-06-01'); -- Joker

INSERT INTO movie_actors (movie_id, actor_id, role)
VALUES
    -- Inception
    (1, 1, 'Dom Cobb'),
    (1, 2, 'Arthur'),
    (1, 3, 'Ariadne'),
    -- The Dark Knight
    (2, 4, 'Bruce Wayne / Batman'),
    (2, 5, 'Joker'),
    -- Interstellar
    (3, 6, 'Cooper'),
    (3, 7, 'Brand'),
    -- The Godfather
    (4, 8, 'Vito Corleone'),
    (4, 9, 'Michael Corleone'),
    -- Pulp Fiction
    (5, 10, 'Jules Winnfield'),
    (5, 11, 'Vincent Vega'),
    -- The Shawshank Redemption
    (6, 12, 'Red'),
    (6, 13, 'Andy Dufresne'),
    -- The Matrix
    (7, 14, 'Neo'),
    (7, 15, 'Morpheus'),
    -- Gladiator
    (8, 16, 'Maximus'),
    (8, 17, 'Commodus'),
    -- Parasite
    (9, 18, 'Ki-taek'),
    (9, 19, 'Ki-woo'),
    -- Joker
    (10, 17, 'Arthur Fleck / Joker'),
    (10, 20, 'Murray Franklin'),
    (10, 21, 'Sophie Dumond');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM movie_actors;
DELETE FROM actors;
-- +goose StatementEnd
