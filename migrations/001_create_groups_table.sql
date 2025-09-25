-- +goose Up
CREATE TABLE groups (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	number INTEGER NOT NULL,
	name TEXT NOT NULL
);

-- +goose Down
DROP TABLE groups;
