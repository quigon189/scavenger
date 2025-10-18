-- +goose Up
ALTER TABLE users ADD COLUMN theme TEXT DEFAULT 'light';

-- +goose Down
ALTER TABLE users DROP COLUMN theme;
