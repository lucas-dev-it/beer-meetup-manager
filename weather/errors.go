package weather

import "errors"

var (
	ErrResourceMissingData = errors.New("there is missing information from resource request")
	ErrResourceFailedRequest = errors.New("request has failed")
	)
