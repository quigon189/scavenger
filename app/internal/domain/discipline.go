package domain

import (
	"time"
)

type Discipline struct {
	ID          int64
	Title       string
	Description string
	TeacherID   int64
	GroupID     int64
	Archived    bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
