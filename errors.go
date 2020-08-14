package meetupmanager

import "errors"

var (
	ErrResourceMissingData = errors.New("there is missing information from resource request")
	ErrResourceFailedRequest = errors.New("request has failed")
	ErrDBRecordNotFound = errors.New("the record has not been found")
	ErrForecastNotAvailable = errors.New("forecast not available for selected date")
	)
