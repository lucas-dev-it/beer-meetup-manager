package weather

import (
	"fmt"
)

type WeatherProvider interface {
	GetForecastData(country, state, city string, forecastDays uint, client HttpClient) (map[string]interface{}, error)
	GetAdapter() adapter
}

func GetProvider(providerName string) (WeatherProvider, error) {
	switch providerName {
	case "weather-stack":
		return NewWeatherStackResource(weatherStack), nil
	case "weather-bit":
		return NewWeatherBitResource(weatherBit), nil
	default:
		return nil, fmt.Errorf("there is no such %v defined resource to fetch the weather forecast", providerName)
	}
}
