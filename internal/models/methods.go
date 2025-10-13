package models

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

func (cfg *Config) GetGroupDisciplines(groupName string) []Discipline {
	var groupDisciplnes []Discipline
	return groupDisciplnes
}

func (cfg *Config) GetStudentReports(username string) []LabReport {
	var studReports []LabReport
	return studReports
}

func (cfg *Config) GetDiscepline(id string) *Discipline {
	disc := &Discipline{}
	return disc
}

func (cfg *Config) GetLab(id string) *Lab {
	lab := &Lab{}
	return lab
}

func GetUserFromContext(ctx context.Context) *User {
	if user, ok := ctx.Value("user").(User); ok {
		return &user
	}
	return &User{}

}

func (g *Group) StudentCount() int {
	return len(g.Students)
}

func (g *Group) IDtoStr() string {
	return strconv.Itoa(g.ID)
}
func (d *Discipline) IDtoStr() string {
	return strconv.Itoa(d.ID)
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

func (f *ReportFilterParams) Parse(r *http.Request) {
	f.Page = 1
	f.PageSize = 20
	f.SortBy = "uploaded_at"
	f.SortOrder = "desc"

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			f.Page = page
		}
	}

	if disciplineIDStr := r.URL.Query().Get("discipline_id"); disciplineIDStr != "" {
		if id, err := strconv.Atoi(disciplineIDStr); err == nil {
			f.DisciplineID = id
		}
	}

	if labIDStr := r.URL.Query().Get("lab_id"); labIDStr != "" {
		if id, err := strconv.Atoi(labIDStr); err == nil {
			f.LabID = id
		}
	}

	if status := r.URL.Query().Get("status"); status != "" {
		f.Status = status
	}

	if gradeStr := r.URL.Query().Get("grade"); gradeStr != "" {
		if grade, err := strconv.Atoi(gradeStr); err == nil {
			f.Grade = grade
		}
	}

	if studentSearch := r.URL.Query().Get("student_search"); studentSearch != "" {
		f.StudentSearch = studentSearch
	}

	if period := r.URL.Query().Get("period"); period != "" {
		f.Period = period
	}

	if sortBy := r.URL.Query().Get("sort_by"); sortBy != "" {
		f.SortBy = sortBy
	}

	if sortOrder := r.URL.Query().Get("sort_order"); sortOrder != "" {
		f.SortOrder = sortOrder
	}
}
