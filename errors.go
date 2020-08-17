package meetupmanager

import (
	"errors"
	"fmt"
)

var (
	ErrResourceMissingData       = errors.New("there is missing information from resource request")
	ErrResourceInvalidStatusCode = errors.New("resource provider responded with an unexpected status code")

	ErrNoWeatherInformationAsYet = errors.New("there is no forecast available for the meetup's date")

	ErrBadRequest             = errors.New("provided parameters are wrong") // 400
	ErrNotFound               = errors.New("the record has not been found") // 404
	ErrDependencyNotAvailable = errors.New("dependency not available")      // 424
)

type CustomError struct {
	Cause   error
	Type    error
	Message string
}

func (c CustomError) Error() string {
	return fmt.Sprintf("errorType: %v with message: %v", c.Type, c.Message)
}

func (c CustomError) GetType() error {
	return c.Type
}
