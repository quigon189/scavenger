package service

import "scavenger/internal/repository"

type Deps struct {
	Repo          repository.Repository
	SessionSecret []byte
}

type Services struct {
	Auth *AuthService
	User *UserService
}

func New(d Deps) *Services {
	return &Services{
		Auth: NewAuth(d.Repo, d.SessionSecret),
		User: NewUser(d.Repo),
	}
}
