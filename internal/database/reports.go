package database

import (
	"scavenger/internal/models"
	"time"
)

func (d *Database) AddLabReport(report *models.LabReport) error {
	result, err := d.db.Exec(
		CreateLabReportQuery,
		report.StudentID,
		report.DisciplineID,
		report.LabID,
		report.Comment,
		report.TeacherNote,
		report.UploadedAt.Unix(),
		report.UpdatedAt.Unix(),
		report.Status,
		report.Grade,
	)
	if err != nil {
		return err
	}


	id,err := result.LastInsertId()
	if err != nil {
		return err
	}

	report.ID = int(id)


	err = d.AddReportFiles(report.ID, report.Files)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetLabReport(studID, labID int) (*models.LabReport, error) {
	report := models.LabReport{}

	var uploadedAt, updatedAt int64

	err := d.db.QueryRow(GetLabReportQuery, studID, labID).Scan(
		&report.ID,
		&report.StudentID,
		&report.DisciplineID,
		&report.LabID,
		&report.Comment,
		&report.TeacherNote,
		&uploadedAt,
		&updatedAt,
		&report.Status,
		&report.Grade,
	)
	if err != nil {
		return nil, err
	}

	report.UploadedAt = time.Unix(uploadedAt, 0)
	report.UpdatedAt = time.Unix(updatedAt, 0)

	files, err := d.getReportFiles(&report)
	if err != nil {
		return nil, err
	}

	report.Files = append(report.Files, files...)

	return &report, nil
}

func (d *Database) UpdateReport(report *models.LabReport) error {
	_, err := d.db.Exec(UpdateLabReportQuery, report.Comment, report.UpdatedAt.Unix(), report.ID)
	return err
}

func (d *Database) getReportFiles(report *models.LabReport) ([]models.StoredFile, error) {
	files := []models.StoredFile{}
	row, err := d.db.Query(GetReportFilesQuery, report.ID)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		file := models.StoredFile{}
		err = row.Scan(
			&file.ID,
			&file.Path,
			&file.URL,
			&file.Filename,
			&file.Size,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}

	return files, nil
}

func (d *Database) AddReportFiles(repID int, files []models.StoredFile) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	for _, file := range files {
		_, err = tx.Exec(CreateReportFileQuery, repID, file.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
