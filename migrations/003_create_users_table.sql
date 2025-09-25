-- +goose Up
CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	name TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	role_id INTEGER REFERENCES roles(id),
	group_id INTEGER REFERENCES groups(id)
);

-- +goose Down
DROP TABLE users;
