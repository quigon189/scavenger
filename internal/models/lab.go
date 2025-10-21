package models

import "time"

type Lab struct {
	ID           string
	Name         string
	MDFileID     int
	Deadline     time.Time
	Description  string
	DisciplineID int

	MDFile      StoredFile
	StoredFiles []StoredFile
	Reports     []LabReport
}

func (l *Lab) FormatDeadline() string {
	return l.Deadline.Format("02.01.2006")
}

func (l *Lab) GetStatus() string {
	now := time.Now()
	if now.After(l.Deadline) {
		return "Скрок вышел"
	} else if now.Add(7 * 24 * time.Hour).After(l.Deadline) {
		return "Скрок истекает"
	} else {
		return "Активно"
	}
}

func (l *Lab) GetStatusBadge() string {
	now := time.Now()
	if now.After(l.Deadline) {
		return "secondary"
	} else if now.Add(7 * 24 * time.Hour).After(l.Deadline) {
		return "warning"
	} else {
		return "success"
	}
}
