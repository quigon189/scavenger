package models

import (
	"context"
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
	case "graded" :
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
