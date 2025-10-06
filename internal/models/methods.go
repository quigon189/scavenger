package models

import (
	"context"
	"strconv"
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

func(g *Group) IDtoStr() string {
	return strconv.Itoa(g.ID)
}
func(d *Discipline) IDtoStr() string {
	return strconv.Itoa(d.ID)
}

func(l *Lab) FormatDeadline() string {
	return l.Deadline.Format("01.02.2006")
}
