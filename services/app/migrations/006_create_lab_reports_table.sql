-- +goose Up
CREATE TABLE lab_reports (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	student_id INTEGER NOT NULL,
	discipline_id INTEGER NOT NULL,
	lab_id INTEGER NOT NULL,
	comment TEXT,
	teacher_note TEXT,
	uploaded_at INTEGER,
	updated_at INTEGER,
	status TEXT DEFAULT 'submitted',
	grade INTEGER,

	UNIQUE (student_id, lab_id),
	FOREIGN KEY(student_id) REFERENCES users(id),
	FOREIGN KEY(discipline_id) REFERENCES disciplines(id),
	FOREIGN KEY(lab_id) REFERENCES labs(id)
);

CREATE TABLE report_files (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	report_id INTEGER NOT NULL,
	file_id INTEGER NOT NULL,

	FOREIGN KEY(report_id) REFERENCES lab_reports(id) ON DELETE CASCADE,
	FOREIGN KEY(file_id) REFERENCES stored_files(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE lab_reports;
