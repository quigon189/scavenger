package sqlite

import (
	"context"
	"encoding/json"
	"errors"
	"scavenger/internal/domain"
	"time"
)

type submissionRepo struct{ q queryer }

type payloadReport struct {
	Text    string   `json:"text"`
	FileIDs []string `json:"file_ids"`
}

type grade struct {
	Value     int       `json:"value"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func (r submissionRepo) Save(ctx context.Context, s domain.Submission) error {
	var payload any
	switch v := s.(type) {
	case *domain.Report:
		payload = payloadReport{
			Text:    v.GetText(),
			FileIDs: v.GetAttachedFileIDs(),
		}
	default:
		return errors.New("unknown submission type")
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	gradeJSON, err := json.Marshal(s.GetGrade())
	if err != nil {
		return err
	}

	_, err = r.q.ExecContext(ctx,
		`INSERT INTO submissions(id, student_id, material_id, grade, payload, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		s.GetID(),
		s.GetStudentID(),
		s.GetMaterialID(),
		gradeJSON,
		payloadJSON,
		s.GetCreatedAt(),
	)
	return nil
}

func (r *submissionRepo) List(ctx context.Context) ([]domain.Submission, error) {

}

func scanSubmission(s rowScanner) (domain.Submission, error) {
	var (
		id, materialID, text, payloadRaw, gradeRaw string
		studentID                                  int64
		createdAt                                  time.Time
		updatedAt                                  *time.Time
	)
}
