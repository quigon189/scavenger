package services

import (
	"scavenger/core/services/authservice"
	"scavenger/core/services/dataservice"
	"scavenger/core/services/sessionservice"
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
