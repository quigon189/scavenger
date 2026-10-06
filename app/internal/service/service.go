package service

import "scavenger/internal/repository"

type Deps struct {
	Repo          repository.Repository
	SessionSecret []byte
}

type Services struct {
	Auth       *AuthService
	User       *UserService
	Group      *GroupService
	Discipline *DisciplineService
}

func New(d Deps) *Services {
	return &Services{
		Auth:       NewAuth(d.Repo, d.SessionSecret),
		User:       NewUser(d.Repo),
		Group:      NewGroup(d.Repo),
		Discipline: NewDiscipline(d.Repo),
	}
}
