package services

import (
	"scavenger/internal/services/authservice"
	"scavenger/internal/services/dataservice"
	"scavenger/internal/services/sessionservice"
)

type Services struct {
	Auth    *authservice.AuthService
	Session *sessionservice.SessionService
	Data    *dataservice.DataService
}

func NewServices(auth *authservice.AuthService, data *dataservice.DataService, session *sessionservice.SessionService) *Services {
	return &Services{
		Auth: auth,
		Data: data,
		Session: session,
	}
}
