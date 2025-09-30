package models

import "context"

func (cfg *Config) GetGroupDisciplines(groupName string) []Discipline {
	var groupDisciplnes []Discipline

	var group Group
	for _, g := range cfg.Groups {
		if g.Name == groupName {
			group = g
			break
		}
	}

	for _, discName := range group.Disciplines {
		for _, disc := range cfg.Disciplines {
			if disc.Name == discName {
				groupDisciplnes = append(groupDisciplnes, disc)
			}
		}
	}

	return groupDisciplnes
}

func (cfg *Config) GetStudentReports(username string) []LabReport {
	var studReports []LabReport

	for _, rep := range cfg.LabReports {
		if rep.Student == username {
			studReports = append(studReports, rep)
		}
	}

	return studReports
}

func (cfg *Config) GetDiscepline(id string) *Discipline {
	disc := &Discipline{}

	for _, d := range cfg.Disciplines {
		if d.ID == id {
			disc = &d
			break
		}
	}

	return disc
}

func (cfg *Config) GetLab(id string) *Lab {
	lab := &Lab{}

	for _, d := range cfg.Disciplines {
		for _, l := range d.Labs {
			if l.ID == id {
				lab = &l
				break
			}
		}
	}

	return lab
}

func GetUserFromContext(ctx context.Context) *User {
	if user, ok := ctx.Value("user").(User); ok {
		return &user
	}
	return &User{}

}
