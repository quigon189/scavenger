package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DataClient struct {
	client *Client
}

func NewDataClient(baseURL string, timeout time.Duration) *DataClient {
	return &DataClient{
		client: NewClient(baseURL, timeout),
	}
}

type Discipline struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	TeacherID   int    `json:"teacher_id"`
	GroupID     int    `json:"group_id"`
	PeriodID    int    `json:"period_id"`
	Description string `json:"description"`
}

type Lab struct {
	ID           int    `json:"id"`
	DisciplineID int    `json:"discipline_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Deadline     string `json:"deadline"`
}

type LabReport struct {
	ID        int    `json:"id"`
	LabID     int    `json:"lab_id"`
	StudentID int    `json:"student_id"`
	Status    string `json:"status"`
	Grade     int    `json:"grade"`
	Comment   string `json:"comment"`
}

type Student struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	GroupID int    `json:"group_id"`
	Name    string `json:"name"`
}

type Group struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (d *DataClient) GetGroups(ctx context.Context) ([]Group, error) {
	var groups []Group

	err := d.client.doRequest(ctx, http.MethodGet, "/api/groups", nil, nil, &groups)	
	if err != nil {
		return nil, fmt.Errorf("failed to get groups: %w", err)
	}

	return groups, nil
}

func (d *DataClient) GetDisciplines(ctx context.Context, sessionID string) ([]Discipline, error) {
	var disciplines []Discipline

	headers := map[string]string{
		"X-Session-ID": sessionID,
	}

	err := d.client.doRequest(ctx, http.MethodGet, "/api/disciplines", headers, nil, &disciplines)
	if err != nil {
		return nil, fmt.Errorf("failed to get disciplines: %w", err)
	}

	return disciplines, nil
}

func (d *DataClient) GetDiscipline(ctx context.Context, sessionID string, id int) (*Discipline, error) {
	var discipline Discipline

	headers := map[string]string{
		"X-Session-ID": sessionID,
	}

	err := d.client.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/disciplines/%d", id), headers, nil, &discipline)
	if err != nil {
		return nil, fmt.Errorf("failed to get discipline: %w", err)
	}

	return &discipline, nil
}

func (d *DataClient) GetLabsByDiscipline(ctx context.Context, sessionID string, disciplineID int) ([]Lab, error) {
	var labs []Lab

	headers := map[string]string{
		"X-Session-ID": sessionID,
	}

	err := d.client.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/labs/discipline/%d", disciplineID), headers, nil, &labs)
	if err != nil {
		return nil, fmt.Errorf("failed to get labs: %w", err)
	}

	return labs, nil
}

func (d *DataClient) CreateReport(ctx context.Context, sessionID string, labID int, comment string) (*LabReport, error) {
	var report LabReport

	headers := map[string]string{
		"X-Session-ID": sessionID,
	}

	req := map[string]any{
		"lab_id":  labID,
		"comment": comment,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	err = d.client.doRequest(ctx, http.MethodPost, "/api/reports", headers, bytes.NewReader(body), &report)
	if err != nil {
		return nil, fmt.Errorf("failed to create report: %w", err)
	}

	return &report, nil
}

func (d *DataClient) GetReportsByStudent(ctx context.Context, sessionID string, studentID int) ([]LabReport, error) {
	var reports []LabReport

	headers := map[string]string{
		"X-Session-ID": sessionID,
	}

	err := d.client.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/reports/student/%d", studentID), headers, nil, &reports)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports: %w", err)
	}

	return reports, nil
}

// ... добавьте другие методы по мере необходимости
