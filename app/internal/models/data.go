package models

import "time"

type Group struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`

	Students    []Student    `json:"students"`
	Disciplines []Discipline `json:"disciplines"`
}

type Teacher struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `json:"user"`

	Disciplines []Discipline `json:"disciplines"`
}


type Student struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User  User  `json:"user"`
	Group Group `json:"group"`
}

type Discipline struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	TeacherID   int       `json:"teacher_id"`
	GroupID     int       `json:"group_id"`
	PeriodID    int       `json:"period_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Teacher Teacher `json:"teacher"`
	Group   Group   `json:"group"`
	Period  Period  `json:"period"`

	Labs []Lab `json:"labs"`
}

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

type Period struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	HalfYear  int       `json:"half_year"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
}

type File struct {
	ID          int       `json:"id"`
	UUID        string    `json:"uuid"`
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Bucket      string    `json:"bucket"`
	Path        string    `json:"path"`
	CreatedAt   time.Time `json:"created_at"`

	SignedURL string `json:"signed_url,omitempty"`
}

type LabFile struct {
	ID     int `json:"id"`
	LabID  int `json:"lab_id"`
	FileID int `json:"file_id"`

	Lab  Lab  `json:"lab"`
	File File `json:"file"`
}

type ReportFile struct {
	ID       int `json:"id"`
	ReportID int `json:"report_id"`
	FileID   int `json:"file_id"`

	Report LabReport `json:"report"`
	File   File      `json:"file"`
}
