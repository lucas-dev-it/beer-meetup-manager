package interfaces

import (
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/weather"
)

type WeatherService interface {
	GetForecast(country, state, city string) (*weather.Forecast, error)
}
