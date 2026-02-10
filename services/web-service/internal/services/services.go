package services

import "web-service/pkg/apiclient"

type Services struct {
	Data *DataService
	Auth *AuthService
}

func NewServices(dataClient *apiclient.DataClient, authClient *apiclient.AuthClient) *Services {
	return &Services{
		Data: NewDataService(dataClient),
		Auth: NewAuthService(authClient),
	}
}
