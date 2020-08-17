package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
)

type responseWrapper struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type errorDetails struct {
	Message string `json:"message"`
	// TODO errorTypes
}

type errorWrapper struct {
	Success bool         `json:"success"`
	Error   errorDetails `json:"error"`
}

type handlerResult struct {
	status int32
	body   interface{}
}

func middleware(h func(io.Writer, *http.Request) (*handlerResult, error), wrapper bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response interface{}

		var status int32
		hr, err := h(w, r)
		if err != nil {
			response = errorWrapper{
				Success: false,
				Error:   errorDetails{Message: err.Error()},
			}
			status = handleErrors(err)
		} else {
			if !wrapper {
				response = hr.body
			} else {
				response = responseWrapper{Data: hr.body, Success: true}
			}
			status = hr.status
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(int(status))

		if status != http.StatusNoContent {
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Printf("could not encode response to output: %v", err)
			}
		}
	}
}

func handleErrors(err error) int32 {
	switch t := err.(type) {
	case meetupmanager.CustomError:
		switch t.Type {
		case meetupmanager.ErrDependencyNotAvailable:
			return http.StatusFailedDependency
		case meetupmanager.ErrBadRequest:
			return http.StatusBadRequest
		case meetupmanager.ErrNotFound:
			return http.StatusNotFound
		case meetupmanager.ErrNoWeatherInformationAsYet:
			return http.StatusNotAcceptable
		default:
			return http.StatusInternalServerError
		}
	default:
		return http.StatusInternalServerError
	}
}
