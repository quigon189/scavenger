-- +goose Up
-- +goose StatementBegin
CREATE TABLE disciplines (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	description TEXT,
	teacher_id INTEGER,
	group_id INTEGER,
	archived BOOLEAN NOT NULL,
	created_at DATETIME NOT NULL,
	updated_at DATETIME,

	FOREIGN KEY (teacher_id) REFERENCES users(id) ON DELETE SET NULL,
	FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS disciplinesl
-- +goose StatementEnd
