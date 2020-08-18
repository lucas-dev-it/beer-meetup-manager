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

const wbURL = "https://api.weatherbit.io"

var wbAccessKey = internal.GetEnv("WEATHERBIT_API_KEY", "myTestAPIKey")

type weatherBitCli struct {
	Adapter adapter
}

func NewWeatherBitResource(adapter adapter) *weatherBitCli {
	return &weatherBitCli{
		Adapter: adapter,
	}
}

func (cli *weatherBitCli) GetAdapter() adapter {
	return cli.Adapter
}

func (cli *weatherBitCli) GetForecastData(country, state, city string, forecastDays uint, httpCli HttpClient) (map[string]interface{}, error) {
	URL := fmt.Sprintf("%v/v2.0/forecast/daily", wbURL)
	request := &httpclient.RequestData{
		Verb: http.MethodGet,
		URL:  URL,
		QueryParams: map[string]string{
			"key":     wbAccessKey,
			"country": country,
			"city":    city,
			"days":    fmt.Sprintf("%v", forecastDays),
		},
	}

	logrus.Infof("sending request to %v", URL)
	response, err := httpCli.PerformRequest(request)
	if err != nil {
		return nil, meetupmanager.CustomError{
			Cause:   meetupmanager.ErrDependencyNotAvailable,
			Type:    meetupmanager.ErrDependencyNotAvailable,
			Message: "weather provider is not available",
		}
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
