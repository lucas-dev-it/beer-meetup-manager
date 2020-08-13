package service

type WeatherFetcher interface {
	FetchWeather(country, state, city string) (float64, error)
}

type MeetUpService struct {
	weatherFetcher WeatherFetcher
}

func NewMeetUpService(wf WeatherFetcher) *MeetUpService {
	return &MeetUpService{weatherFetcher: wf}
}


