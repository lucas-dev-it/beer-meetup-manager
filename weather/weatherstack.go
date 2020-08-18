package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal/httpclient"
	"github.com/sirupsen/logrus"
)

const wsURL = "http://api.weatherstack.com"

var wsAccessKey = internal.GetEnv("WEATHERSTACK_API_KEY", "myTestAPIKey")

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
	URL := fmt.Sprintf("%v/forecast", wsURL)
	request := &httpclient.RequestData{
		Verb: http.MethodGet,
		URL:  URL,
		QueryParams: map[string]string{
			"access_key":    wsAccessKey,
			"query":         fmt.Sprintf("%v,%v,%v", country, state, city),
			"forecast_days": fmt.Sprintf("%v", forecastDays),
			"units":         "m",
		},
	}

	logrus.Infof("sending request to %v", URL)
	response, err := httpCli.PerformRequest(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, meetupmanager.CustomError{
			Cause:   meetupmanager.ErrDependencyNotAvailable,
			Type:    meetupmanager.ErrDependencyNotAvailable,
			Message: fmt.Sprintf("weather provider responded with an invalid status code: %d", response.StatusCode),
		}
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
