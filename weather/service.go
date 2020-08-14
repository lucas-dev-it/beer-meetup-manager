package weather

import (
	"net/http"
	"strconv"

	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal/httpclient"
)

var forecastDays = internal.GetEnv("FORECAST_DAYS", "5")

type HttpClient interface {
	PerformRequest(rd *httpclient.RequestData) (*http.Response, error)
}

type Forecast struct {
	DateTempMap map[uint]*DailyForecast
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

func (ws *WService) GetForecast(country, state, city string) (*Forecast, error) {
	fd, err := strconv.ParseUint(forecastDays, 10, 64)
	if err != nil {
		return nil, err
	}

	client := httpclient.New(30) // TODO move this to env var
	response, err := ws.weatherProvider.GetForecastData(country, state, city, uint(fd), client)
	if err != nil {
		return nil, err
	}

	forecast, err := ws.weatherProvider.GetAdapter()(response)
	if err != nil {
		return nil, err
	}

	return forecast, nil
}
