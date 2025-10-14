package database

import (
	"sort"
	"strings"
	"time"

	"scavenger/internal/models"
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

	id, err := result.LastInsertId()
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

func (d *Database) GetLabReportByID(id int) (*models.LabReport, error) {
	report := models.LabReport{}

	var uploadedAt, updatedAt int64

	err := d.db.QueryRow(GetLabReportByIDQuery, id).Scan(
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

func (d *Database) GetAllReports() ([]models.LabReport, error) {
	var reports []models.LabReport

	row, err := d.db.Query(GetAllLabReportsQuery)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var report models.LabReport
		var uploadedAt, updatedAt int64
		err = row.Scan(
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

		report.UpdatedAt = time.Unix(updatedAt, 0)
		report.UploadedAt = time.Unix(uploadedAt, 0)

		student, err := d.GetStudentByID(report.StudentID)
		if err != nil {
			return nil, err
		}
		disc, err := d.GetDisciplineByID(report.DisciplineID)
		if err != nil {
			return nil, err
		}
		lab, err := d.GetLabByID(report.LabID)
		if err != nil {
			return nil, err
		}

		report.Student = *student
		report.Discipline = *disc
		report.Lab = *lab

		reports = append(reports, report)
	}

	return reports, nil
}

func (d *Database) GetFilteredReports(filter models.ReportFilterParams) ([]models.LabReport, error) {
	reports, err := d.GetAllReports()
	if err != nil {
		return nil, err
	}

	if filter.IsEmpty() {
		return reports, nil
	}

	var filtered []models.LabReport

	for _, report := range reports {
		if filter.DisciplineID > 0 && filter.DisciplineID != report.DisciplineID {
			continue
		}
		if filter.LabID > 0 && filter.LabID != report.LabID {
			continue
		}
		if filter.Status != "" && filter.Status != report.Status {
			continue
		}
		if filter.Grade > 0 && filter.Grade != report.Grade {
			continue
		}
		if filter.StudentSearch != "" {
			searshLower := strings.ToLower(filter.StudentSearch)
			studentNameLower := strings.ToLower(report.Student.Name)
			studentUsernameLower := strings.ToLower(report.Student.Username)

			if !strings.Contains(studentNameLower, searshLower) &&
				!strings.Contains(studentUsernameLower, searshLower) {
				continue
			}
		}
		if !matchesPeriod(report.UpdatedAt, filter.Period) {
			continue
		}


		filtered = append(filtered, report)
	}

	filtered = sortReports(filtered, filter.SortBy, filter.SortOrder)

	if filter.Page > 0 && filter.PageSize > 0 {
        start := (filter.Page - 1) * filter.PageSize
        end := start + filter.PageSize
        
        if start >= len(filtered) {
            return []models.LabReport{}, nil
        }
        
        if end > len(filtered) {
            end = len(filtered)
        }
        
        return filtered[start:end], nil
    }

	return filtered, nil
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

func matchesPeriod(uploadedAt time.Time, period string) bool {
	now := time.Now()

	switch period {
	case "today":
		return uploadedAt.Year() == now.Year() &&
			uploadedAt.Month() == now.Month() &&
			uploadedAt.Day() == now.Day()

	case "week":
		weekAgo := now.AddDate(0, 0, -7)
		return uploadedAt.After(weekAgo)

	case "month":
		monthAgo := now.AddDate(0,-1,0)
		return uploadedAt.After(monthAgo)

	default:
		 return true
	}
}

func sortReports(reports []models.LabReport, sortBy, sortOrder string) []models.LabReport {
	sort.Slice(reports, func(i, j int) bool {
		switch sortBy {
		case "student":
			if sortOrder == "asc" {
				return reports[i].Student.Name < reports[j].Student.Name
			}
			return reports[i].Student.Name > reports[j].Student.Name
		case "discipline":
			if sortOrder == "asc" {
				return reports[i].Discipline.Name < reports[j].Discipline.Name
			}
			return reports[i].Discipline.Name > reports[j].Discipline.Name
		case "uploaded_at":
			if sortOrder == "asc" {
				return reports[i].UploadedAt.Before(reports[j].UploadedAt)
			}
			return reports[i].UploadedAt.After(reports[j].UploadedAt)
		case "deadline":
			if sortOrder == "asc" {
				return reports[i].Lab.Deadline.Before(reports[j].Lab.Deadline)
			}
			return reports[i].Lab.Deadline.After(reports[j].Lab.Deadline)
		case "grade":
			if sortOrder == "asc" {
				return reports[i].Grade < reports[j].Grade
			}
			return reports[i].Grade > reports[j].Grade
		default:
			return reports[i].UpdatedAt.After(reports[j].UpdatedAt)
		}
	})
	
	return reports
}
