package weather

import (
	"github.com/mitchellh/mapstructure"
)

type adapter func(response map[string]interface{}) (*Forecast, error)

var weatherStack = func(response map[string]interface{}) (*Forecast, error) {
	type wsForecastDaily struct {
		DateEpoch uint     `mapstructure:"date_epoch"`
		MinTemp   *float64 `mapstructure:"mintemp"`
		MaxTemp   *float64 `mapstructure:"maxtemp"`
	}
	type wsForecast struct {
		Forecast map[string]*wsForecastDaily `mapstructure:"forecast"`
	}

	var wsf wsForecast
	if err := mapstructure.Decode(response, &wsf); err != nil {
		return nil, err
	}

	resultMap := make(map[uint]*DailyForecast, len(wsf.Forecast))
	for _, day := range wsf.Forecast {
		if day.MaxTemp == nil && day.MinTemp == nil {
			return nil, ErrResourceMissingData
		}

		df := &DailyForecast{}
		if day.MaxTemp != nil {
			df.maxTemp = *day.MaxTemp
		}

		if day.MinTemp != nil {
			df.minTemp = *day.MinTemp
		}

		resultMap[day.DateEpoch] = df
	}

	return &Forecast{dateTempMap: resultMap}, nil
}
