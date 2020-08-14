package weather

import (
	"encoding/json"
	"testing"
)

const testData = `
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

const wrongTestData = `
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

func getProviderTestData(proper bool) (map[string]interface{}, error){
	d := testData
	if !proper {
		d = wrongTestData
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(d), &data); err != nil {
		return nil, err
	}

	return data, nil
}

func TestWeatherStackAdapter(t *testing.T) {
	data, err := getProviderTestData(true)
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
	data, err := getProviderTestData(false)
	if err != nil {
		t.Errorf("unexpected error, got %v", err)
	}

	forecast, err := weatherStack(data)
	if err == nil {
		t.Errorf("expected error, got %v", forecast)
	}
}
