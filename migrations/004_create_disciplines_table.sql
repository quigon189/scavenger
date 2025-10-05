-- +goose Up
CREATE TABLE disciplines (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	group_id INTEGER REFERENCES groups(id) ON DELETE SET NULL,
	UNIQUE (name, group_id)
);

-- +goose Down
DROP TABLE disciplines;
