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

func (d *Discipline) Avg() string {
	avg := 0
	counter := 0
	for _, lab := range d.Labs {
		for _, report := range lab.Reports {
			if report.Grade != 0 {
				avg = avg + report.Grade
				counter ++
			}
		}
	}

	if avg == 0 {
		return "-"
	}

	return strconv.FormatFloat(float64(avg) / float64(counter), 'f', 2, 64)
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
	f.DisciplineID, _ = strconv.Atoi(r.URL.Query().Get("discipline_id"))	
	f.LabID, _ = strconv.Atoi(r.URL.Query().Get("lab_id"))
	f.Status = r.URL.Query().Get("status")
	f.Grade, _ =strconv.Atoi(r.URL.Query().Get("grade"))
	f.StudentSearch = r.URL.Query().Get("student_search")
	f.Period = r.URL.Query().Get("period")

	f.SortBy = r.URL.Query().Get("sort_by")
	f.SortOrder = r.URL.Query().Get("sort_order")
	if f.SortOrder == "" {
		f.SortOrder = "desk"
	}

	f.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	if f.Page == 0 {
		f.Page = 1
	}

	f.PageSize, _ = strconv.Atoi(r.URL.Query().Get("page_size"))
	if f.PageSize == 0 {
		f.PageSize = 20
	}
}

func (f *ReportFilterParams) IsEmpty() bool {
	return f.DisciplineID == 0 &&
		   f.LabID == 0 &&
		   f.Status == "" &&
		   f.Grade == 0 &&
		   f.StudentSearch == "" &&
		   f.Period == "" &&
		   f.SortBy == ""
}
