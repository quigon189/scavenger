package services

import (
	"scavenger/core/services/authservice"
	"scavenger/core/services/dataservice"
)

type Services struct {
	Auth    *authservice.AuthService
	Data    *dataservice.DataService
}

func NewServices(auth *authservice.AuthService, data *dataservice.DataService) *Services {
	return &Services{
		Auth: auth,
		Data: data,
	}
}
