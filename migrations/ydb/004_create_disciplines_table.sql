-- +goose Up
CREATE TABLE disciplines (
    id Serial,
    name Text NOT NULL,
    group_id Int32,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE disciplines;
