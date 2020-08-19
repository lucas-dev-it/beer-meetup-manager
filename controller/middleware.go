package controller

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"

	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
	jwtToken "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal/token"
	"github.com/sirupsen/logrus"
)

var signingString = internal.GetEnv("TOKEN_SIGNING_KEY", "testSigningString")

type responseWrapper struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type errorDetails struct {
	Message string `json:"message"`
}

type errorWrapper struct {
	Success bool         `json:"success"`
	Error   errorDetails `json:"error"`
}

type handlerResult struct {
	status int32
	body   interface{}
}

func middleware(h func(io.Writer, *http.Request) (*handlerResult, error), requiredScopes map[string]interface{}, wrapper bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response interface{}
		var status int32

		if requiredScopes != nil {
			if err := authorize(r, requiredScopes); err != nil {
				response = errorWrapper{
					Success: false,
					Error:   errorDetails{Message: err.Error()},
				}
				logrus.Errorf("an error occurred during request, got: %v", err)
				sendResponse(w, http.StatusForbidden, response)
				return
			}
		}

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

		sendResponse(w, status, response)
	}
}

func authorize(r *http.Request, requiredScopes map[string]interface{}) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return errors.New("missing JWT within Authentication header")
	}

	if !strings.Contains(authHeader, "Bearer") {
		return errors.New("invalid token type")
	}

	tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))

	token, err := jwtToken.ParseTokenString(tokenString, signingString)
	if err != nil {
		return err
	}

	if err := jwtToken.ValidJWT(token, requiredScopes); err != nil {
		return err
	}

	return nil
}

func sendResponse(w http.ResponseWriter, status int32, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(status))
	if status != http.StatusNoContent {
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("could not encode response to output: %v", err)
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
		case meetupmanager.ErrForbiddenAccess:
			return http.StatusForbidden
		default:
			return http.StatusInternalServerError
		}
	default:
		return http.StatusInternalServerError
	}
}
