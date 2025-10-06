-- +goose Up
CREATE TABLE labs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	description TEXT,
	md_path TEXT NOT NULL,
	deadline INTEGER DEFAULT (strftime('%s','now','+7 days')),
	discipline_id INTEGER REFERENCES disciplines(id)
);

CREATE TABLE pdf_files (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	lab_id INTEGER NOT NULL,
	file_path TEXT NOT NULL,
	FOREIGN KEY(lab_id) REFERENCES labs(id)
);
-- +goose Down
DROP TABLE labs;
