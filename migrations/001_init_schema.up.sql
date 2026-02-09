CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.users (
	id SERIAL PRIMARY KEY,
	username VARCHAR(50) UNIQUE NOT NULL,
	email VARCHAR(100) UNIQUE,
	name VARCHAR(100) NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	role VARCHAR(20) NOT NULL CHECK (role IN ('student', 'teacher', 'admin')),
	status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'blocked')),
	group_id INT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- admin admin123
-- stud stud123
-- teach teach123
INSERT INTO auth.users (username, password_hash, name, email, role, status, group_id) VALUES
    ('admin', '$2a$10$UeSM2lg6ALPQnjc/d2R3/Ou4xSZanVeBsIaxjkgYwMwDOvoGqD1bq', 'Admin', 'admin@localhost', 'admin', 'active', NULL),
    ('stud', '$2a$10$VZC0.fidHnzEqCqZeTDzSeYUfy4Rll9gZE/saEB9Mazkysfk/OTKe', 'Student', 'stud@localhost', 'student', 'active', 1),
    ('teach', '$2a$10$9lo5QPzCXB/sIej/VLclHumVos6pifEkUmvKN13pyXtOAGogv3FT.', 'Teacher', 'teach@localhost', 'teacher', 'active', NULL);


CREATE INDEX idx_users_username ON auth.users(username);
CREATE INDEX idx_users_email ON auth.users(email);
CREATE INDEX idx_users_role ON auth.users(role);
