-- +goose Up

CREATE TABLE roles (
    id Serial,
    name TEXT NOT NULL,
    description TEXT,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE roles;
