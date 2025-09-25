-- +goose Up
CREATE TABLE report_files (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	file_name TEXT NOT NULL,
	report_id INTEGER REFERENCES lab_reports(id),
	created_at TEXT DEFAULT (datetime('now'))
);

-- +goose Down
DROP TABLE report_files;
