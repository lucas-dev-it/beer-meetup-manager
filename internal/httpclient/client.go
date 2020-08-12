package httpclient

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type RequestData struct {
	Verb        string
	URL         string
	Body        io.Reader
	QueryParams map[string]string
	Headers     map[string]string
}

type Client struct {
	cli *http.Client
}

// New instance for the httpclient
func New(t time.Duration) *Client {
	return &Client{
		cli: &http.Client{
			Timeout: t * time.Second,
		},
	}
}

// PerformRequest validates and triggers the HTTP call as per RequestData parameter
func (c *Client) PerformRequest(rd *RequestData) (*http.Response, error) {
	if err := c.validate(rd); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(rd.Verb, rd.URL, rd.Body)
	if err != nil {
		return nil, err
	}

	// Set headers
	for k, v := range rd.Headers {
		req.Header.Set(k, v)
	}

	// If querystring was provided
	if len(rd.QueryParams) > 0 {
		q := req.URL.Query()
		for k, v := range rd.QueryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	res, err := c.cli.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err := res.Body.Close(); err != nil {
		return nil, err
	}

	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return res, nil
}

// validate checks for the minimum mandatory request composition
func (c *Client) validate(rd *RequestData) error {
	if rd.Verb == "" || rd.URL == "" {
		return errors.New("neither http Verb nor url must not be string zero value")
	}

	if err := isValidUrl(rd.URL); err != nil {
		return fmt.Errorf("invalid provided url %v: %w", rd.URL, err)
	}

	return nil
}

// isValidUrl tests a string to determine if it is a well-structured url or not.
func isValidUrl(URL string) error {
	_, err := url.ParseRequestURI(URL)
	if err != nil {
		return err
	}

	u, err := url.Parse(URL)
	if err != nil {
		return err
	}

	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("there is an error within the URL composition")
	}

	return nil
}
