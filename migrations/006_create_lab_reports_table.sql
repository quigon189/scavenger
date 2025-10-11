-- +goose Up
CREATE TABLE lab_reports (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	student_id INTEGER REFERENCES user(id),
	discipline_id INTEGER REFERENCES disciplines(id),
	lab_id INTEGER REFERENCES labs(id),
	comment TEXT,
	teacher_note TEXT,
	uploaded_at INTEGER,
	updated_at INTEGER,
	status TEXT DEFAULT 'submitted',
	grade INTEGER,

	UNIQUE (student_id, lab_id)
);

-- +goose Down
DROP TABLE lab_reports;
