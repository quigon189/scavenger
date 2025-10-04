-- +goose Up
CREATE TABLE disciplines (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	group_id INTEGER REFERENCES groups(id) ON DELETE SET NULL
);

-- +goose Down
DROP TABLE disciplines;
