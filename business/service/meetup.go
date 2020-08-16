package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business/model"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/weather"
)

var unitsPerPack = internal.GetEnv("PACK_UNITS", "6")
var forecastDays = internal.GetEnv("FORECAST_DAYS", "10")

type weatherService interface {
	GetForecast(country, state, city string, forecastDays uint) (*weather.Forecast, error)
}

type cacheRepository interface {
	StoreForecast(key string, forecast *weather.Forecast) error
	RetrieveForecast(key string) (*weather.Forecast, error)
}

type meetupRepository interface {
	CountMeetupAttendees(meetupID uint) int
	FindMeetupByID(meetupID uint) (*model.MeetUp, error)
}

type MeetUpService struct {
	weatherService   weatherService
	meetupRepository meetupRepository
	cacheRepository  cacheRepository
}

func NewMeetUpService(mr meetupRepository, cr cacheRepository, ws weatherService) *MeetUpService {
	return &MeetUpService{
		meetupRepository: mr,
		cacheRepository:  cr,
		weatherService:   ws,
	}
}

func (ms *MeetUpService) CalculateBeerPacksForMeetup(meetupID uint) (*business.MeetupBeersData, error) {
	attendeesCount := ms.meetupRepository.CountMeetupAttendees(meetupID)
	meetup, err := ms.meetupRepository.FindMeetupByID(meetupID)
	if err != nil {
		return nil, err
	}

	fd, err := strconv.ParseUint(forecastDays, 10, 64)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	xDaysFromNow := now.AddDate(0, 0, int(fd))
	if meetup.StartDate.After(xDaysFromNow) {
		fmt.Printf("meetup after, data: %v", meetup)
		return nil, nil
	}

	forecast, err := ms.getForecast(meetup, uint(fd))

	upp, err := strconv.ParseUint(unitsPerPack, 10, 64)
	if err != nil {
		return nil, err
	}

	msd := meetup.StartDate
	// use 0 hour for the lookup
	unixDate := time.Date(msd.Year(), msd.Month(), msd.Day(), 0, 0, 0, 0, msd.Location()).Unix()
	dailyForecast, ok := forecast.DateTempMap[unixDate]
	if !ok {
		return nil, meetupmanager.ErrDependencyNotAvailable
	}

	packsQuantity, err := BeerPacksQuantity(uint(attendeesCount), uint(upp), dailyForecast.MaxTemp)
	if err != nil {
		return nil, err
	}

	return &business.MeetupBeersData{
		BeerPacks:      packsQuantity,
		MaxTemperature: dailyForecast.MaxTemp,
		MinTemperature: dailyForecast.MinTemp,
		MeetupMetadata: &business.MeetupMetadata{
			Name:           meetup.Name,
			StartDate:      meetup.StartDate,
			EndDate:        meetup.EndDate,
			Country:        meetup.Country,
			State:          meetup.State,
			City:           meetup.City,
			AttendeesCount: attendeesCount,
		},
	}, nil
}

func (ms *MeetUpService) getForecast(meetup *model.MeetUp, fd uint) (*weather.Forecast, error) {
	var forecast *weather.Forecast
	key := fmt.Sprintf("%v-%v", strings.ToLower(meetup.Country), strings.ToLower(meetup.City))

	// first look it up within the cache
	forecast, err := ms.cacheRepository.RetrieveForecast(key)
	if err != nil || forecast == nil {
		// get from provider and then
		forecast, err = ms.weatherService.GetForecast(meetup.Country, meetup.State, meetup.City, fd)
		if err != nil {
			return nil, err
		}

		// store it within the cache
		if err = ms.cacheRepository.StoreForecast(key, forecast); err != nil {
			return nil, err
		}
	}

	return forecast, nil
}
