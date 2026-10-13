-- +goose Up
-- +goose StatementBegin
CREATE TABLE files (
	id            TEXT PRIMARY KEY,
	storage_key   TEXT NOT NULL,
	original_name TEXT NOT NULL,
	content_type  TEXT NOT NULL,
	size          INTEGER NOT NULL,
	owner_id      INTEGER NOT NULL,
	created_at    DATEIME NOT NULL,

	FOREIGN KEY(owned_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS files;
