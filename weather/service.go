package weather

import (
	"net/http"
	"strconv"

	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal/httpclient"
)

var (
	providerName = internal.GetEnv("WEATHER_PROVIDER", "weather-stack")
	forecastDays = internal.GetEnv("FORECAST_DAYS", "5")
)

type httpClient interface {
	PerformRequest(rd *httpclient.RequestData) (*http.Response, error)
}

type Forecast struct {
	dateTempMap map[uint]*DailyForecast
}

type DailyForecast struct {
	minTemp float64
	maxTemp float64
}

type WService struct {
	weatherProvider WeatherProvider
}

func NewWeatherService(httpCli httpClient) (*WService, error) {
	wProvider, err := GetProvider(providerName, httpCli)
	if err != nil {
		return nil, err
	}

	return &WService{weatherProvider: wProvider}, nil
}

func (ws *WService) GetForecast(country, state, city string) (*Forecast, error) {
	fd, err := strconv.ParseUint(forecastDays, 10, 64)
	if err != nil {
		return nil, err
	}

	response, err := ws.weatherProvider.GetForecastData(country, state, city, uint(fd))
	if err != nil {
		return nil, err
	}

	forecast, err := ws.weatherProvider.GetAdapter()(response)
	if err != nil {
		return nil, err
	}

	return forecast, nil
}
