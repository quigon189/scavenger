-- +goose Up
CREATE TABLE groups (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	number INTEGER,
	name TEXT NOT NULL
);

-- +goose Down
DROP TABLE groups;
