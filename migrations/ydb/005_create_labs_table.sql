-- +goose Up
CREATE TABLE stored_files (
    id       Serial,
    path     Text NOT NULL,
    url      Text NOT NULL,
    filename Text NOT NULL,
    size     Int32 NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE labs (
    id            Serial,
    name          Text NOT NULL,
    description   Text,
    md_id         Int32,
    deadline      Int32,
    discipline_id Int32 NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE labs_files (
    id     Serial,
    file_id Int32,
    lab_id  Int32,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE labs;
