-- +goose Up
CREATE TABLE stored_files (
    id       Serial,
    path     Text NOT NULL,
    url      Text NOT NULL,
    filename Text NOT NULL,
    size     Int64 NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE labs (
    id            Serial,
    name          Text NOT NULL,
    description   Text,
    md_id         Int32 NOT NULL,
    deadline      Int64 NOT NULL,
    discipline_id Int32 NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE labs_files (
    id     Serial,
    file_id Int32 NOT NULL,
    lab_id  Int32 NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE labs;
