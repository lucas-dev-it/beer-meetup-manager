package weather

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type weatherProvider struct{}

func (w weatherProvider) GetForecastData(country, state, city string, forecastDays uint, client httpClient) (map[string]interface{}, error) {
	data, err := getProviderTestDataJSON(true)
	if err != nil {
		return nil, err
	}

	if country == "chile" {
		data, err = getProviderTestDataJSON(false)
		if err != nil {
			return nil, err
		}
	} else if country == "uruguay" {
		return nil, errors.New("failed request")
	}

	return data, nil
}

func (w weatherProvider) GetAdapter() adapter {
	return weatherBit
}

func getProviderTestDataJSON(proper bool) (map[string]interface{}, error) {
	d := wbTestData
	if !proper {
		d = wbWrongTestData
	}

	var res map[string]interface{}
	if err := json.Unmarshal([]byte(d), &res); err != nil {
		return nil, err
	}

	return res, nil
}

func TestWService_GetForecast(t *testing.T) {
	service, err := NewWeatherService(&weatherProvider{})

	forecast, err := service.GetForecast("argentina", "cordoba", "cordoba", 2)
	if err != nil {
		t.Errorf("unexpected error, got: %v", err)
	}

	expected := &Forecast{
		DateTempMap: map[int64]*DailyForecast{
			1491004800: {
				MaxTemp: 30,
				MinTemp: 26,
			},
			1491091200: {
				MaxTemp: 32,
				MinTemp: 21,
			},
		},
	}

	assert.Equal(t, expected, forecast)
}

func TestWService_GetForecast_MissingTempFields(t *testing.T) {
	service, err := NewWeatherService(&weatherProvider{})

	forecast, err := service.GetForecast("chile", "some", "place", 2)
	if err == nil {
		t.Errorf("expected error, got: %v", forecast)
	}
}

func TestWService_GetForecast_InvalidStatusCode(t *testing.T) {
	service, err := NewWeatherService(&weatherProvider{})

	forecast, err := service.GetForecast("uruguay", "some", "place", 2)
	if err == nil {
		t.Errorf("expected error, got: %v", forecast)
	}
}
