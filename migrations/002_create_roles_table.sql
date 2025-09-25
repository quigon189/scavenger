-- +goose Up
CREATE TABLE roles (
	id INTEGER PRIMARY KEY AUTOINCREMENT
	name TEXT NOT NULL UNIQUE,
	description TEXT
);

-- +goose Down
DROP TABLE roles;
