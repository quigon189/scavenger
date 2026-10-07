-- +goose Up
-- +goose StatementBegin
CREATE TABLE materials (
	id TEXT PRIMARY KEY,
	discipline_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	display_order INTEGER NOT NULL,
	visible BOOLEAN NOT NULL,
	type TEXT NOT NULL,
	payload TEXT NOT NULL,
	created_at DATETIME NOT NULL,
	updated_at DATETIME,

	FOREIGN KEY (discipline_id) REFERENCES disciplines(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS materials;
-- +goose StatementEnd
