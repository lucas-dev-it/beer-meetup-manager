package meetupmanager

import (
	"errors"
	"fmt"
)

var (
	ErrResourceMissingData = errors.New("there is missing information from resource request")

	ErrNoWeatherInformationAsYet = errors.New("there is no forecast available for the meetup's date")

	ErrUnauthorizedAccess     = errors.New("unauthorized access")
	ErrForbiddenAccess        = errors.New("access not allowed")
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
	return fmt.Sprintf("%v. Details: %v", c.Type, c.Message)
}

func (c CustomError) GetType() error {
	return c.Type
}
