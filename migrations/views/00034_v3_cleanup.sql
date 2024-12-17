-- +goose Up
-- +goose StatementBegin
-- additional migration in case of changes between migration and view updates
UPDATE movies
SET genres = ARRAY[genre], version = version + 1
WHERE genres IS NULL AND genre IS NOT NULL;

ALTER TABLE movies
DROP COLUMN genre;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE movies
ADD COLUMN genre text;

UPDATE movies
SET genre = genres[1], version = version + 1
WHERE genres IS NOT NULL AND array_length(genres, 1) > 0;
-- +goose StatementEnd
