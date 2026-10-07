-- +goose Up
-- +goose StatementBegin
CREATE TABLE groups (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	number INTEGER NOT NULL,
	start_year INTEGER NOT NULL CHECK ( start_year BETWEEN 2000 AND 3000 ),
	end_year INTEGER NOT NULL CHECK ( end_year BETWEEN start_year AND 3000 ),
	specialty TEXT NOT NULL,
	short_specialty TEXT NOT NULL,
	created_at DATETIME NOT NULL,
	updated_at DATETIME,

	UNIQUE (number, start_year)
);

CREATE TABLE group_students (
	student_id INTEGER PRIMARY KEY,
	group_id INTEGER NOT NULL,

	FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS group_students;
-- +goose StatementEnd
