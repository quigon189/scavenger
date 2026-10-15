-- +goose Up

CREATE TABLE groups(
    id Serial,
    number Int32,
    name Text NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down

DROP TABLE groups;
