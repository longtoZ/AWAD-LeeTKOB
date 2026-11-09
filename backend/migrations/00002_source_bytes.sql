-- +goose Up
ALTER TABLE submissions ALTER COLUMN source TYPE bytea USING convert_to(source, 'UTF8');

-- +goose Down
ALTER TABLE submissions ALTER COLUMN source TYPE text USING convert_from(source, 'UTF8');
