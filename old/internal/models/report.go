package models

import (
	"strconv"
	"time"
)

type LabReport struct {
	ID           int
	StudentID    int
	DisciplineID int
	LabID        int
	Comment      string
	TeacherNote  string
	UploadedAt   time.Time
	UpdatedAt    time.Time
	Status       string
	Grade        int

	Student    User
	Discipline Discipline
	Lab        Lab
	Files      []StoredFile
}

func (r *LabReport) GetStatusText() string {
	switch r.Status {
	case "graded":
		return "Проверено"
	case "draft":
		return "Черновик"
	case "submitted":
		return "Ожидает проверки"
	default:
		return r.Status
	}
}

func (r *LabReport) GetStatusBadge() string {
	switch r.Status {
	case "draft":
		return "secondary"
	case "submitted":
		return "warning"
	case "graded":
		return "success"
	default:
		return "secondary"
	}
}

func (r *LabReport) IDtoStr() string {
	return strconv.Itoa(r.ID)
}
