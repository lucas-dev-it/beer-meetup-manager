package weather

import (
	"fmt"
)

type WeatherProvider interface {
	GetForecastData(country, state, city string, forecastDays uint) (map[string]interface{}, error)
	GetAdapter() adapter
}

func GetProvider(resourceName string, client httpClient) (WeatherProvider, error) {
	switch resourceName {
	case "weather-stack":
		return NewWeatherStackResource(client, weatherStack), nil
	default:
		return nil, fmt.Errorf("there is no such %v defined resource to fetch the weather forecast", resourceName)
	}
}
