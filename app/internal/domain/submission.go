package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrInvalidSubmissionID         = errors.New("invalid submission id")
	ErrInvalidSubmissionStudentID  = errors.New("invalid submission student id")
	ErrInvalidSubmissionMaterialID = errors.New("invalid submission student id")
	ErrInvalidGradeValue           = errors.New("invalid grade value")
)

type Grade struct {
	Value     int
	Comment   string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type Submission interface {
	GetID() string
	GetStudentID() int64
	GetMaterialID() string
	GetGrade() *Grade
	SetGrade(value int, comment string) error
	GetCreatedAt() time.Time
	GetUpdatedAt() *time.Time
}

type baseSubmission struct {
	id         string
	studentID  int64
	materialID string
	grade      *Grade
	createdAt  time.Time
	updatedAt  *time.Time
}

func (s *baseSubmission) GetID() string           { return s.id }
func (s *baseSubmission) GetStudentID() int64     { return s.studentID }
func (s *baseSubmission) GetMaterialID() string   { return s.materialID }
func (s *baseSubmission) GetCreatedAt() time.Time { return s.createdAt }
func (s *baseSubmission) GetUpdatedAt() *time.Time {
	updatedAt := *s.updatedAt
	return &updatedAt
}
func (s *baseSubmission) GetGrade() *Grade {
	if s.grade == nil {
		return nil
	}
	copied := *s.grade
	return &copied
}
func (s *baseSubmission) SetGrade(value int, comment string) error {
	if value < 2 || value > 5 {
		return ErrInvalidGradeValue
	}
	if s.grade == nil {
		s.grade = &Grade{
			Value:     value,
			Comment:   comment,
			CreatedAt: time.Now(),
		}
		return nil
	}
	s.grade.Value = value
	if comment != "" {
		s.grade.Comment = comment
	}
	updatedAt := time.Now()
	s.grade.UpdatedAt = &updatedAt
	return nil
}

type Report struct {
	baseSubmission
	text            string
	attachedFileIDs map[string]struct{}
}

func NewReport(id string, studentID int64, materialID string, text string, attachedFileIDs []string) (*Report, error) {
	if id == "" {
		return nil, ErrInvalidSubmissionID
	}
	if studentID == 0 {
		return nil, ErrInvalidSubmissionStudentID
	}
	if materialID == "" {
		return nil, ErrInvalidSubmissionMaterialID
	}
	report := Report{
		baseSubmission: baseSubmission{
			id:        id,
			studentID: studentID,
			createdAt: time.Now(),
		},
		text: text,
	}
	for _, fID := range attachedFileIDs {
		report.attachedFileIDs[fID] = struct{}{}
	}
	return &report, nil
}

func (r *Report) GetText() string { return r.text }
func (r *Report) GetAttachedFileIDs() []string {
	var attachedFileIDs []string
	for fID := range r.attachedFileIDs {
		attachedFileIDs = append(attachedFileIDs, fID)
	}
	return attachedFileIDs
}
func (r *Report) UpdateText(newText string) {
	r.text = newText
	updatedAt := time.Now()
	r.updatedAt = &updatedAt
}
func (r *Report) AttachFileIDs(fileIDs ...string) {
	for _, fID := range fileIDs {
		r.attachedFileIDs[fID] = struct{}{}
	}
	updatedAt := time.Now()
	r.updatedAt = &updatedAt
}

type SubmissionRepo interface {
	Save(context.Context, Submission) error
	ByID(context.Context, string) (Submission, error)
	ByStudentID(context.Context, int64) ([]Submission, error)
	ByStudentIDs(context.Context, []int64) ([]Submission, error)
	ByMaterialID(context.Context, string) ([]Submission, error)
	ByMaterialIDs(context.Context, []string) ([]Submission, error)
}
