package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal/httpclient"
)

const URL = "http://api.weatherstack.com"

var accessKey = internal.GetEnv("WEATHERSTACK_API_KEY", "myTestAPIKey")

type weatherStackCli struct {
	Adapter adapter
}

func NewWeatherStackResource(adapter adapter) *weatherStackCli {
	return &weatherStackCli{
		Adapter: adapter,
	}
}

func (cli *weatherStackCli) GetAdapter() adapter {
	return cli.Adapter
}

func (cli *weatherStackCli) GetForecastData(country, state, city string, forecastDays uint, httpCli HttpClient) (map[string]interface{}, error) {
	request := &httpclient.RequestData{
		Verb: http.MethodGet,
		URL:  fmt.Sprintf("%v/forecast", URL),
		QueryParams: map[string]string{
			"access_key":    accessKey,
			"query":         fmt.Sprintf("%v,%v,%v", country, state, city),
			"forecast_days": fmt.Sprintf("%v", forecastDays),
			"units":         "m",
		},
	}

	response, err := httpCli.PerformRequest(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("invalid status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var bodyMap map[string]interface{}

	if err := json.Unmarshal(body, &bodyMap); err != nil {
		return nil, err
	}

	return bodyMap, nil
}
