package service

import (
	"scavenger/internal/repository"
	"scavenger/internal/storage"
)

type Deps struct {
	Repo          repository.Repository
	Storage       storage.Storage
	SessionSecret []byte
	MaxUploadSize int64
}

type Services struct {
	Auth       *AuthService
	User       *UserService
	Group      *GroupService
	Discipline *DisciplineService
	Material   *MaterialService
	File       *FileService
}

func New(d Deps) *Services {
	return &Services{
		Auth:       NewAuth(d.Repo, d.SessionSecret),
		User:       NewUser(d.Repo),
		Group:      NewGroup(d.Repo),
		Discipline: NewDiscipline(d.Repo),
		Material:   NewMaterial(d.Repo),
		File:       NewFile(d.Repo, d.Storage, d.MaxUploadSize),
	}
}
