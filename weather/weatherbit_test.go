package weather

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal/httpclient"
	"github.com/stretchr/testify/assert"
)

type httpCliMock struct{}

func (hcm *httpCliMock) PerformRequest(rd *httpclient.RequestData) (*http.Response, error) {
	query := rd.QueryParams["country"]

	if query == "errorland" {
		return nil, errors.New("http communication error")
	}

	data := wbTestData
	status := 200

	if strings.Contains(query, "uruguay") {
		data = ""
		status = 500
	}

	return &http.Response{
		StatusCode: status,
		Body:       ioutil.NopCloser(strings.NewReader(data)),
	}, nil
}

func Test_weatherBitCli_GetForecastData(t *testing.T) {
	wbCli := &weatherBitCli{}
	data, err := wbCli.GetForecastData("argentina", "cordoba", "cordoba", 10, &httpCliMock{})
	if err != nil {
		t.Error("unexpected error")
	}

	d, ok := data["data"].([]interface{	})
	if !ok {
		t.Error("unexpected response data from client")
	}

	assert.Len(t, d, 2)
}

func Test_weatherBitCli_GetForecastData_HTTPError(t *testing.T) {
	wbCli := &weatherBitCli{}
	data, err := wbCli.GetForecastData("errorland", "cordoba", "cordoba", 10, &httpCliMock{})
	if err == nil && data != nil{
		t.Error("expected error")
	}

	assert.IsType(t, meetupmanager.CustomError{}, err)
}

func Test_weatherBitCli_GetForecastData_NotExpectedHTTPStatus(t *testing.T) {
	wbCli := &weatherBitCli{}
	data, err := wbCli.GetForecastData("uruguay", "cordoba", "cordoba", 10, &httpCliMock{})
	if err == nil && data != nil{
		t.Error("expected error")
	}

	assert.IsType(t, meetupmanager.CustomError{}, err)
}
