-- +goose Up
ALTER TABLE disciplines
ADD COLUMN teacher_id Int32;

-- +goose Down
ALTER TABLE disciplines
DROP COLUMN teacher_id;
