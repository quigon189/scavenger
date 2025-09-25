-- +goose Up
CREATE TABLE lab_reports (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	student_id INTEGER REFERENCES user(id),
	lab_id INTEGER REFERENCES labs(id),
	status TEXT DEFAULT 'submitted',
	grade INTEGER
);

-- +goose Down
DROP TABLE lab_reports;
