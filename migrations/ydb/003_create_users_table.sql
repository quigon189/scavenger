-- +goose Up
CREATE TABLE users (
    id Serial,
    username TEXT NOT NULL,
    name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role_id INTEGER NOT NULL,
    PRIMARY KEY (id),
);

CREATE TABLE students (
    user_id Int32,
    group_id Int32 NOT NULL,
    PRIMARY KEY (user_id)
);

-- +goose Down
DROP TABLE users;
DROP TABLE students;
