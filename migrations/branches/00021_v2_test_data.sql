-- +goose Up
-- +goose StatementBegin
-- Update existing records with test data for the new columns
INSERT INTO movies_branch_v2 (id, director, runtime, language)
VALUES 
    (1, 'Christopher Nolan', 148, 'English'),
    (2, 'Christopher Nolan', 152, 'English'),
    (3, 'Christopher Nolan', 169, 'English'),
    (4, 'Francis Ford Coppola', 175, 'English'),
    (5, 'Quentin Tarantino', 154, 'English'),
    (6, 'Frank Darabont', 142, 'English'),
    (7, 'The Wachowskis', 136, 'English'),
    (8, 'Ridley Scott', 155, 'English'),
    (9, 'Bong Joon-ho', 132, 'Korean'),
    (10, 'Todd Phillips', 122, 'English');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM movies_branch_v2;
-- +goose StatementEnd
