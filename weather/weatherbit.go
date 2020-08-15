package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal/httpclient"
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
	request := &httpclient.RequestData{
		Verb: http.MethodGet,
		URL:  fmt.Sprintf("%v/v2.0/forecast/daily", wbURL),
		QueryParams: map[string]string{
			"key":     wbAccessKey,
			"country": country,
			"city":    city,
			"days":    fmt.Sprintf("%v", forecastDays),
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
