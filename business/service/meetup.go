package service

import "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/interfaces"

type MeetUpService struct {
	weatherFetcher interfaces.WeatherService
}

func NewMeetUpService(wf interfaces.WeatherService) *MeetUpService {
	return &MeetUpService{weatherFetcher: wf}
}
