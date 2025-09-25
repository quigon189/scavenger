-- +goose Up
CREATE TABLE labs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	file_path TEXT NOT NULL,
	deadline INTEGER DEFAULT (strftime('%s','now','+7 days')),
	discipline_id INTEGER REFERENCES disciplines(id)
);

-- +goose Down
DROP TABLE labs;
