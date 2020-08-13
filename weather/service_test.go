package weather

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal/httpclient"
)

type httpCliMock struct{}

func (hcm *httpCliMock) PerformRequest(rd *httpclient.RequestData) (*http.Response, error) {
	query := rd.QueryParams["query"]

	data := getProviderTestDataJSON(true)
	status := 200

	if strings.Contains(query, "chile") {
		data = getProviderTestDataJSON(false)
	} else if strings.Contains(query, "uruguay") {
		data = []byte{}
		status = 500
	}

	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(strings.NewReader(string(data))),
	}, nil
}

func getProviderTestDataJSON(proper bool) []byte {
	d := testData
	if !proper {
		d = wrongTestData
	}

	return []byte(d)
}

func TestWService_GetForecast(t *testing.T) {
	service, err := NewWeatherService(&httpCliMock{})

	forecast, err := service.GetForecast("argentina", "cordoba", "cordoba")
	if err != nil {
		t.Errorf("unexpected error, got: %v", err)
	}

	expected := &Forecast{
		dateTempMap: map[uint]*DailyForecast{
			1567814400: {
				maxTemp: 25,
				minTemp: 17,
			},
			1567911600: {
				maxTemp: 20,
				minTemp: 30,
			},
		},
	}

	assert.IsEqual(forecast, expected)
}

func TestWService_GetForecast_MissingTempFields(t *testing.T) {
	service, err := NewWeatherService(&httpCliMock{})

	forecast, err := service.GetForecast("chile", "some", "place")
	if err == nil {
		t.Errorf("expected error, got: %v", forecast)
	}
}

func TestWService_GetForecast_InvalidStatusCode(t *testing.T) {
	service, err := NewWeatherService(&httpCliMock{})

	forecast, err := service.GetForecast("uruguay", "some", "place")
	if err == nil {
		t.Errorf("expected error, got: %v", forecast)
	}
}
