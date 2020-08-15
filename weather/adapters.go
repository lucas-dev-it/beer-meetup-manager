package weather

import (
	"time"

	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/mitchellh/mapstructure"
)

type adapter func(response map[string]interface{}) (*Forecast, error)

var weatherStack = func(response map[string]interface{}) (*Forecast, error) {
	type wsForecastDaily struct {
		DateEpoch int64    `mapstructure:"date_epoch"`
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

	resultMap := make(map[int64]*DailyForecast, len(wsf.Forecast))
	for _, day := range wsf.Forecast {
		if day.MaxTemp == nil && day.MinTemp == nil {
			return nil, meetupmanager.ErrResourceMissingData
		}

		df := &DailyForecast{}
		if day.MaxTemp != nil {
			df.MaxTemp = *day.MaxTemp
		}

		if day.MinTemp != nil {
			df.MinTemp = *day.MinTemp
		}

		resultMap[day.DateEpoch] = df
	}

	return &Forecast{DateTempMap: resultMap}, nil
}

var weatherBit = func(response map[string]interface{}) (*Forecast, error) {
	type wbForecastDaily struct {
		Datetime string   `mapstructure:"valid_date"`
		MinTemp  *float64 `mapstructure:"min_temp"`
		MaxTemp  *float64 `mapstructure:"max_temp"`
	}
	type wbForecast struct {
		Forecast []*wbForecastDaily `mapstructure:"data"`
	}

	var wbf wbForecast
	if err := mapstructure.Decode(response, &wbf); err != nil {
		return nil, err
	}

	resultMap := make(map[int64]*DailyForecast, len(wbf.Forecast))
	for _, day := range wbf.Forecast {
		actualDate, err := time.Parse("2006-01-02", day.Datetime)
		if err != nil {
			return nil, err
		}
		ts := actualDate.Unix()

		df := &DailyForecast{}
		if day.MaxTemp != nil {
			df.MaxTemp = *day.MaxTemp
		}

		if day.MinTemp != nil {
			df.MinTemp = *day.MinTemp
		}

		resultMap[ts] = df
	}

	return &Forecast{DateTempMap: resultMap}, nil
}
