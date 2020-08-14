package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type responseWrapper struct {
	Success bool
	Data    interface{}
}

type errorDetails struct {
	message string
	// TODO errorTypes
}

type errorWrapper struct {
	Success bool
	Error   errorDetails
}

func middleware(h func(io.Writer, *http.Request) (interface{}, int, error), wrapResponse bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// execute actual retHandler
		data, status, err := h(w, r)
		if err != nil {
			data = err.Error()
		}

		// wrap the response data
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if status != http.StatusNoContent {
			var response interface{}
			if wrapResponse {
				response = responseWrapper{Data: data, Success: err == nil}
			} else {
				response = errorWrapper{
					Success: false,
					Error:   errorDetails{message: err.Error()},
				}
			}

			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Printf("could not encode response to output: %v", err)
			}
		}
	}
}
