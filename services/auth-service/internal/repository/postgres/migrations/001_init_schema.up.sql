CREATE SCHEMA IF NOT EXISTS auth;

CREATE TABLE auth.users (
	id SERIAL PRIMARY KEY,
	username VARCHAR(50) UNIQUE NOT NULL,
	email VARCHAR(100) UNIQUE,
	name VARCHAR(100) NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	role VARCHAR(20) NOT NULL CHECK (role IN ('student', 'teacher', 'admin')),
	theme VARCHAR(10) DEFAULT 'light',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE auth.sessions (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
	token VARCHAR(512) NOT NULL,
	expires_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(user_id)
);

CREATE INDEX idx_users_username ON auth.users(username);
CREATE INDEX idx_users_email ON auth.users(email);
CREATE INDEX idx_users_role ON auth.users(role);
CREATE INDEX idx_sessions_token ON auth.sessions(token);
CREATE INDEX idx_sessions_expires_at ON auth.sessions(expires_at);
