package service

import (
	"errors"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business/model"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/weather"
	"github.com/stretchr/testify/assert"
)

type mRepo struct{}

func (m mRepo) CountMeetupAttendees(meetupID uint) int {
	return int(meetupID)
}

func (m mRepo) FindMeetupByID(meetupID uint) (*model.MeetUp, error) {
	if meetupID < 100 {
		now := time.Now()
		after := now.Add(1)
		return &model.MeetUp{
			Model:       gorm.Model{ID: meetupID},
			Name:        "test",
			Description: "test",
			StartDate:   &now,
			EndDate:     &after,
			Country:     "argentina",
			State:       "cordoba",
			City:        "cordoba",
		}, nil
	} else if meetupID < 200 {
		now := time.Now()
		after := now.Add(1)
		return &model.MeetUp{
			Model:       gorm.Model{ID: meetupID},
			Name:        "test",
			Description: "test",
			StartDate:   &now,
			EndDate:     &after,
			Country:     "argentinaCached",
			State:       "cordoba",
			City:        "cordoba",
		}, nil
	} else if meetupID < 300 {
		now := time.Now().AddDate(1, 0, 0)
		after := now.Add(1)
		return &model.MeetUp{
			Model:       gorm.Model{ID: meetupID},
			Name:        "test future",
			Description: "test future",
			StartDate:   &now,
			EndDate:     &after,
			Country:     "argentina",
			State:       "cordoba",
			City:        "cordoba",
		}, nil
	}

	return nil, errors.New("intended error")
}

type cacheRepo struct{}

func (c cacheRepo) StoreForecast(key string, forecast *weather.Forecast) error {
	if key == "wrongKey" {
		return errors.New("error from cache layer")
	}

	return nil
}

func (c cacheRepo) RetrieveForecast(key string) (*weather.Forecast, error) {
	if key == "argentinacached-cordoba" {
		return getValidForecast()
	}

	return nil, errors.New("no present in cache")
}

type wService struct{}

func (w wService) GetForecast(country, state, city string, forecastDays uint) (*weather.Forecast, error) {
	if country == "argentina" {
		return getValidForecast()
	}

	return nil, errors.New("weather not available from provider")
}

func getValidForecast() (*weather.Forecast, error) {
	now := time.Now()
	tomorrow := now.AddDate(0, 0, 1)
	nowUnix := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	tomorrowUnix := time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, tomorrow.Location()).Unix()

	return &weather.Forecast{DateTempMap: map[int64]*weather.DailyForecast{
		nowUnix: {
			MinTemp: 25,
			MaxTemp: 35,
		},
		tomorrowUnix: {
			MinTemp: 26,
			MaxTemp: 36,
		},
	}}, nil
}

func getService() *MeetUpService {
	return NewMeetUpService(mRepo{}, cacheRepo{}, wService{})
}

func TestMeetUpService_CalculateBeerPacksForMeetup(t *testing.T) {
	s := getService()

	data, err := s.CalculateBeerPacksForMeetup(99)
	if err != nil {
		t.Error("unexpected error")
	}

	assert.Equal(t, 49.5, *data.BeerPacks)
	assert.Equal(t, float64(25), data.MinTemperature)
	assert.Equal(t, float64(35), data.MaxTemperature)
}

func TestMeetUpService_CalculateBeerPacksForMeetup_FromCache(t *testing.T) {
	s := getService()

	data, err := s.CalculateBeerPacksForMeetup(199)
	if err != nil {
		t.Fatalf("unexpected error, got %+v", err)
	}

	assert.Equal(t, 99.5, *data.BeerPacks)
	assert.Equal(t, float64(25), data.MinTemperature)
	assert.Equal(t, float64(35), data.MaxTemperature)
}

func TestMeetUpService_CalculateBeerPacksForMeetup_NotExistingMeeting(t *testing.T) {
	s := getService()

	data, err := s.CalculateBeerPacksForMeetup(9999)
	if err == nil && data != nil {
		t.Fatalf("expected error, got %+v", data)
	}
}

func TestMeetUpService_CalculateBeerPacksForMeetup_FutureMeeting(t *testing.T) {
	s := getService()

	data, err := s.CalculateBeerPacksForMeetup(299)
	if err == nil && data != nil {
		t.Fatalf("expected error, got %+v", data)
	}

	assert.IsType(t, meetupmanager.CustomError{}, err)
}

func TestMeetUpService_GetMeetupWeather(t *testing.T) {
	s := getService()

	data, err := s.GetMeetupWeather(99)
	if err != nil {
		t.Error("unexpected error")
	}

	assert.Equal(t, float64(25), data.MinTemperature)
	assert.Equal(t, float64(35), data.MaxTemperature)
}

func TestMeetUpService_GetMeetupWeather_FromCache(t *testing.T) {
	s := getService()

	data, err := s.GetMeetupWeather(199)
	if err != nil {
		t.Fatalf("unexpected error, got %+v", err)
	}

	assert.Equal(t, float64(25), data.MinTemperature)
	assert.Equal(t, float64(35), data.MaxTemperature)
}

func TestMeetUpService_GetMeetupWeather_NotExistingMeeting(t *testing.T) {
	s := getService()

	data, err := s.GetMeetupWeather(9999)
	if err == nil && data != nil {
		t.Fatalf("expected error, got %+v", data)
	}
}

func TestMeetUpService_GetMeetupWeather_FutureMeeting(t *testing.T) {
	s := getService()

	data, err := s.GetMeetupWeather(299)
	if err == nil && data != nil {
		t.Fatalf("expected error, got %+v", data)
	}

	assert.IsType(t, meetupmanager.CustomError{}, err)
}
