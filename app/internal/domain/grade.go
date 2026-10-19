package domain

import "time"

type TargetType string

const (
	PracticalGrade TargetType = "practical_grade"
)

type Grade struct {
	ID         string
	TargetID   string
	TargetType string
	Value      int
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}
