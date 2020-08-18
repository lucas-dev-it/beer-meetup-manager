package weather

import (
	"encoding/json"
	"testing"
)

const wsTestData = `
{
   "current":{
      "observation_time":"03:38 PM",
      "temperature":18,
      "weather_code":113
   },
   "forecast":{
      "2019-09-07":{
         "date":"2019-09-07",
         "date_epoch":1567814400,
         "mintemp":17,
         "maxtemp":25,
         "avgtemp":21
      },
      "2019-09-08":{
         "date":"2019-09-08",
         "date_epoch":1567911600,
         "mintemp":20,
         "maxtemp":30,
         "avgtemp":25
      }
   }
}
`

const wsWrongTestData = `
{
   "current":{
      "observation_time":"03:38 PM",
      "temperature":18,
      "weather_code":113
   },
   "forecast":{
      "2019-09-07":{
         "date":"2019-09-07",
         "date_epoch":1567814400,
         "avgtemp":21
      },
      "2019-09-08":{
         "date":"2019-09-08",
         "date_epoch":1567911600,
         "mintemp":20,
         "maxtemp":30,
         "avgtemp":25
      }
   }
}
`

const wbTestData = `
{
    "data":[
        {
            "valid_date":"2017-04-01",
            "max_temp":30,
            "min_temp":26
        },
        {
            "valid_date":"2017-04-02",
            "max_temp":32,
            "min_temp":21
        }
    ],
    "city_name":"Raleigh",
    "lon":"-78.63861",
    "timezone":"America\/New_York",
    "lat":"35.7721",
    "country_code":"US",
    "state_code":"NC"
}
`

const wbWrongTestData = `
{
    "data":[
        {
            "max_temp":30,
            "min_temp":26
        },
        {
            "valid_date":"2017-04-02",
            "max_temp":32,
            "min_temp":21
        }
    ],
    "city_name":"Raleigh",
    "lon":"-78.63861",
    "timezone":"America\/New_York",
    "lat":"35.7721",
    "country_code":"US",
    "state_code":"NC"
}
`

func getProviderTestData(proper bool, provider string) (map[string]interface{}, error) {
	var d string
	if provider == "weather-stack" {
		d = wsTestData
		if !proper {
			d = wsWrongTestData
		}
	} else {
		d = wbTestData
		if !proper {
			d = wbWrongTestData
		}
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(d), &data); err != nil {
		return nil, err
	}

	return data, nil
}

func TestWeatherStackAdapter(t *testing.T) {
	data, err := getProviderTestData(true, "weather-stack")
	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	forecast, err := weatherStack(data)
	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	if len(forecast.DateTempMap) != 2 {
		t.Errorf("unexpected result, got %v, ", forecast)
	}
}

func TestWeatherStackAdapter_MissingFields(t *testing.T) {
	data, err := getProviderTestData(false, "weather-stack")
	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	forecast, err := weatherStack(data)
	if err == nil {
		t.Errorf("expected error, got %v", forecast)
	}
}

func TestWeatherBitAdapter(t *testing.T) {
	data, err := getProviderTestData(true, "weather-bit")
	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	forecast, err := weatherBit(data)
	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	if len(forecast.DateTempMap) != 2 {
		t.Errorf("unexpected result, got %v, ", forecast)
	}
}

func TestWeatherBitAdapter_MissingFields(t *testing.T) {
	data, err := getProviderTestData(false, "weather-bit")
	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	forecast, err := weatherBit(data)
	if err == nil {
		t.Errorf("expected error, got %v", forecast)
	}
}
