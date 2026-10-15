-- +goose Up
CREATE TABLE lab_reports (
    id            Serial,
    student_id    Int32 NOT NULL,
    discipline_id Int32 NOT NULL,
    lab_id        Int32 NOT NULL,
    comment       Text,
    teacher_note  Text,
    uploaded_at   Int32,
    updated_at    Int32,
    status        Text DEFAULT "submitted",
    grade         Int32,
    PRIMARY KEY (id),
);

CREATE TABLE report_files (
    id        Serial,
    report_id Int32 NOT NULL,
    file_id   Int32 NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
DROP TABLE lab_reports;
