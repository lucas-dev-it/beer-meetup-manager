package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
)

var signingString = internal.GetEnv("INTERNAL_API_KEY", "testSigningString")

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
				sendResponse(w, http.StatusUnauthorized, response)
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

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingString), nil
	})
	if err != nil {
		return err
	}

	if err := validJWT(token, requiredScopes); err != nil {
		return err
	}

	return nil
}

func validJWT(token *jwt.Token, requiredScopes map[string]interface{}) error {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		valid := claims.VerifyExpiresAt(time.Now().Unix(), true)
		if !valid {
			return errors.New("invalid token - expired")
		}

		scopes, ok := claims["scopes"].([]interface{})
		if !ok {
			return errors.New("invalid token - missing scopes")
		}

		for _, s := range scopes {
			sn := s.(string)
			if _, ok := requiredScopes[sn]; !ok {
				return errors.New("invalid token - not authorized scopes")
			}
		}

		return nil
	}

	return errors.New("invalid token")
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
		default:
			return http.StatusInternalServerError
		}
	default:
		return http.StatusInternalServerError
	}
}
