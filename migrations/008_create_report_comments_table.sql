-- +goose Up
CREATE TABLE report_comments (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	message TEXT,
	report_id INTEGER REFERENCES lab_reports(id)
);

-- +goose Down
DROP TABLE report_comments;
