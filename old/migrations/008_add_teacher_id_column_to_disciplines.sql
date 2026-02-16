-- +goose Up
ALTER TABLE disciplines
ADD COLUMN teacher_id INTEGER REFERENCES users(id);

-- +goose Down
ALTER TABLE disciplines
DROP COLUMN teacher_id;
