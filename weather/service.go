package weather

import (
	"net/http"

	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal/httpclient"
	"github.com/sirupsen/logrus"
)

const ttl = 30

type HttpClient interface {
	PerformRequest(rd *httpclient.RequestData) (*http.Response, error)
}

type Forecast struct {
	DateTempMap map[int64]*DailyForecast
}

type DailyForecast struct {
	MinTemp float64
	MaxTemp float64
}

type WService struct {
	weatherProvider WeatherProvider
}

func NewWeatherService(weatherProvider WeatherProvider) (*WService, error) {
	return &WService{weatherProvider: weatherProvider}, nil
}

func (ws *WService) GetForecast(country, state, city string, forecastDays uint) (*Forecast, error) {
	logrus.Info("fetching forecast from weather provider")
	client := httpclient.New(ttl)
	forecastData, err := ws.weatherProvider.GetForecastData(country, state, city, forecastDays, client)
	if err != nil {
		return nil, err
	}

	forecast, err := ws.weatherProvider.GetAdapter()(forecastData)
	if err != nil {
		return nil, err
	}

	return forecast, nil
}
