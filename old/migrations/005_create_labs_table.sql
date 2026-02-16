-- +goose Up
CREATE TABLE stored_files (
	id       INTEGER PRIMARY KEY AUTOINCREMENT,
	path     TEXT NOT NULL,
	url      TEXT NOT NULL,
	filename TEXT NOT NULL,
	size     INTEGER NOT NULL
);

CREATE TABLE labs (
	id            INTEGER PRIMARY KEY AUTOINCREMENT,
	name          TEXT NOT NULL,
	description   TEXT,
	md_id         INTEGER,
	deadline      INTEGER DEFAULT (strftime('%s','now','+7 days')),
	discipline_id INTEGER NOT NULL,
	FOREIGN KEY(discipline_id) REFERENCES disciplines(id),
	FOREIGN KEY(md_id) REFERENCES stored_files(id) ON DELETE SET NULL
);

CREATE TABLE labs_files (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	file_id INTEGER,
	lab_id INTEGER,
	FOREIGN KEY(file_id) REFERENCES stored_files(id) ON DELETE SET NULL,
	FOREIGN KEY(lab_id) REFERENCES labs(id) ON DELETE CASCADE
);
-- +goose Down
DROP TABLE labs;
