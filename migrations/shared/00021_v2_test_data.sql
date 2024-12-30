-- +goose Up
-- +goose StatementBegin
-- Update existing records with test data for the new columns
UPDATE movies
SET 
    director = CASE id
        WHEN 1 THEN 'Christopher Nolan'
        WHEN 2 THEN 'Christopher Nolan'
        WHEN 3 THEN 'Christopher Nolan'
        WHEN 4 THEN 'Francis Ford Coppola'
        WHEN 5 THEN 'Quentin Tarantino'
        WHEN 6 THEN 'Frank Darabont'
        WHEN 7 THEN 'The Wachowskis'
        WHEN 8 THEN 'Ridley Scott'
        WHEN 9 THEN 'Bong Joon-ho'
        WHEN 10 THEN 'Todd Phillips'
    END,
    runtime = CASE id
        WHEN 1 THEN 148
        WHEN 2 THEN 152
        WHEN 3 THEN 169
        WHEN 4 THEN 175
        WHEN 5 THEN 154
        WHEN 6 THEN 142
        WHEN 7 THEN 136
        WHEN 8 THEN 155
        WHEN 9 THEN 132
        WHEN 10 THEN 122
    END,
    language = CASE id
        WHEN 1 THEN 'English'
        WHEN 2 THEN 'English'
        WHEN 3 THEN 'English'
        WHEN 4 THEN 'English'
        WHEN 5 THEN 'English'
        WHEN 6 THEN 'English'
        WHEN 7 THEN 'English'
        WHEN 8 THEN 'English'
        WHEN 9 THEN 'Korean'
        WHEN 10 THEN 'English'
    END,
    version = version + 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE movies
SET 
    director = NULL,
    runtime = NULL,
    language = NULL,
    version = version + 1
WHERE id BETWEEN 1 AND 10;
-- +goose StatementEnd
