package services

import "web-service/pkg/apiclient"

type Services struct {
	Data *DataService
}

func NewServices(dataClient *apiclient.DataClient) *Services {
	return &Services{
		Data: NewDataService(dataClient),
	}
}
