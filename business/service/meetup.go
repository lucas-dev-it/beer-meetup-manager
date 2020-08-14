package service

import (
	"strconv"
	"time"

	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business/model"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/weather"
)

var unitsPerPack = internal.GetEnv("PACK_UNITS", "6")

type WeatherService interface {
	GetForecast(country, state, city string) (*weather.Forecast, error)
}

type meetupRepository interface {
	CountMeetupAttendees(meetupID uint) int
	FindMeetupByID(meetupID uint) (*model.MeetUp, error)
}

type MeetUpService struct {
	weatherService   WeatherService
	meetupRepository meetupRepository
}

func NewMeetUpService(mr meetupRepository, ws WeatherService) *MeetUpService {
	return &MeetUpService{
		meetupRepository: mr,
		weatherService:   ws,
	}
}

func (ms *MeetUpService) CalculateBeerPacksForMeetup(meetupID uint) (*business.MeetupBeersData, error) {
	attendeesCount := ms.meetupRepository.CountMeetupAttendees(meetupID)
	meetup, err := ms.meetupRepository.FindMeetupByID(meetupID)
	if err != nil {
		return nil, err
	}

	// TODO check cache first if no results call the weather provider
	forecast, err := ms.weatherService.GetForecast(meetup.Country, meetup.State, meetup.City)
	if err != nil {
		// TODO evaluate errors here
		return nil, err
	}
	// TODO save in cache what is returned as forecast use country-state-city as key (api have to receive standard names for location)

	upp, err := strconv.ParseUint(unitsPerPack, 10, 64)
	if err != nil {
		return nil, err
	}

	msd := meetup.StartDate
	// use 0 hour for the lookup
	unixDate := time.Date(msd.Year(), msd.Month(), msd.Day(), 0, 0, 0, 0, msd.Location()).Unix()
	dailyForecast, ok := forecast.DateTempMap[uint(unixDate)]
	if !ok {
		return nil, meetupmanager.ErrForecastNotAvailable
	}

	packsQuantity, err := BeerPacksQuantity(uint(attendeesCount), uint(upp), dailyForecast.MaxTemp)
	if err != nil {
		return nil, err
	}

	return &business.MeetupBeersData{
		BeerPacks:      packsQuantity,
		MaxTemperature: dailyForecast.MaxTemp,
		MinTemperature: dailyForecast.MinTemp,
		AttendeesCount: attendeesCount,
	}, nil
}
