CREATE SCHEMA IF NOT EXISTS data;

CREATE TABLE data.periods (
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	half_year INT NOT NULL CHECK (half_year in (1,2)),
	start_date DATE NOT NULL,
	end_date DATE NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE data.groups (
	id SERIAL PRIMARY KEY,
	name VARCHAR(50) NOT NULL UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE data.teachers (
	id INT PRIMARY KEY,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE
);

CREATE TABLE data.students (
	id INT PRIMARY KEY,
	group_id INT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (id) REFERENCES auth.users(id) ON DELETE CASCADE,
	FOREIGN KEY (group_id) REFERENCES data.groups(id) ON DELETE RESTRICT
);

CREATE TABLE data.disciplines (
	id SERIAL PRIMARY KEY,
	name VARCHAR(200) NOT NULL,
	teacher_id INT NOT NULL,
	group_id INT NOT NULL,
	period_id INT NOT NULL,
	description TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (teacher_id) REFERENCES data.teachers(id) ON DELETE CASCADE,
	FOREIGN KEY (group_id) REFERENCES data.groups(id) ON DELETE CASCADE,
	FOREIGN KEY (period_id) REFERENCES data.periods(id) ON DELETE RESTRICT
);

CREATE TABLE data.files (
	id SERIAL PRIMARY KEY,
	uuid VARCHAR(36) NOT NULL UNIQUE,
	filename VARCHAR(255) NOT NULL,
	size BIGINT NOT NULL,
	content_type VARCHAR(100),
	bucket VARCHAR(100) NOT NULL,
	path VARCHAR(500) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE data.labs (
	id SERIAL PRIMARY KEY,
	discipline_id INT NOT NULL,
	md_file_id INT NOT NULL,
	name VARCHAR(200) NOT NULL,
	description TEXT,
	deadline TIMESTAMP NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (discipline_id) REFERENCES data.disciplines(id) ON DELETE CASCADE,
	FOREIGN KEY (md_file_id) REFERENCES data.files(id) ON DELETE RESTRICT
);

CREATE TABLE data.lab_reports (
	id SERIAL PRIMARY KEY,
	lab_id INT NOT NULL,
	student_id INT NOT NULL,
	status VARCHAR(20) DEFAULT 'submitted',
	grade INT,
	comment TEXT,
	teacher_note TEXT,
	graded_at TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

	FOREIGN KEY (lab_id) REFERENCES data.labs(id) ON DELETE CASCADE,
	FOREIGN KEY (student_id) REFERENCES data.students(id) ON DELETE CASCADE,
	UNIQUE(lab_id, student_id)
);

CREATE TABLE data.lab_files (
	id SERIAL PRIMARY KEY,
	lab_id INT NOT NULL,
	file_id INT NOT NULL,

	FOREIGN KEY (lab_id) REFERENCES data.lab_reports(id) ON DELETE CASCADE,
	FOREIGN KEY (file_id) REFERENCES data.files(id) ON DELETE CASCADE
);

CREATE TABLE data.report_files (
	id SERIAL PRIMARY KEY,
	report_id INT NOT NULL,
	file_id INT NOT NULL,

	FOREIGN KEY (report_id) REFERENCES data.lab_reports(id) ON DELETE CASCADE,
	FOREIGN KEY (file_id) REFERENCES data.files(id) ON DELETE CASCADE
);

CREATE INDEX idx_students_id ON data.students(id);
CREATE INDEX idx_students_group_id ON data.students(group_id);
CREATE INDEX idx_disciplines_group_id ON data.disciplines(group_id);
CREATE INDEX idx_disciplines_teacher_id ON data.disciplines(teacher_id);
CREATE INDEX idx_labs_discipline_id ON data.labs(discipline_id);
CREATE INDEX idx_lab_reports_lab_id ON data.lab_reports(lab_id);
CREATE INDEX idx_lab_reports_student_id ON data.lab_reports(student_id);
CREATE INDEX idx_files_uuid ON data.files(uuid);
