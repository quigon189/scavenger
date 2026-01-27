package models

import "time"

type Lab struct {
	ID           int       `json:"id"`
	DisciplineID int       `json:"discipline_id"`
	MDFileID     int       `json:"md_file_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Deadline     time.Time `json:"deadline"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Discipline Discipline `json:"discipline"`
	MDFile     File       `json:"md_file"`

	Files []File `json:"files"`
}

type LabReport struct {
	ID          int        `json:"id"`
	LabID       int        `json:"lab_id"`
	StudentID   int        `json:"student_id"`
	Status      string     `json:"status"`
	Grade       int        `json:"grade"`
	Comment     string     `json:"comment"`
	TeacherNote string     `json:"teacher_note"`
	GradedAt    *time.Time `json:"graded_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Lab     Lab     `json:"lab"`
	Student Student `json:"student"`

	Files []File `json:"files"`
}
