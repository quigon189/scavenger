-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	full_name TEXT NOT NULL,
	role TEXT NOT NULL CHECK(role in ('student', 'teacher')),
	is_active BOOLEAN NOT NULL DEFAULT 1 CHECK(is_active IN (0, 1)),
	created_at DATETIME NOT NULL,
	updated_at DATETIME
);

CREATE TABLE sessions (
	id TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL,
	expires_at DATETIME NOT NULL,
	created_at DATETIME NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS sessions;
-- +goose StatementEnd

