package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

const basePath = "/meetup-manager"

var (
	ALL_SCOPES = map[string]interface{}{
		"USER":  struct{}{},
		"ADMIN": struct{}{},
	}
	ADMIN = map[string]interface{}{
		"ADMIN": struct{}{},
	}
)

type mHandler interface {
	CalculateBeers(w io.Writer, r *http.Request) (*handlerResult, error)
	MeetupWeather(w io.Writer, r *http.Request) (*handlerResult, error)
}

type uHandler interface {
	TokenIssue(w io.Writer, r *http.Request) (*handlerResult, error)
}

type healthHandler struct{}

func New(userHandler uHandler, meetupHandler mHandler) http.Handler {
	router := mux.NewRouter()

	hHealth := &healthHandler{}

	router.HandleFunc("/health", hHealth.health).Methods(http.MethodGet)

	router.HandleFunc(fmt.Sprintf("%v/v1/meetups/{id:[0-9]+}/beers", basePath), middleware(meetupHandler.CalculateBeers, ADMIN, true)).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("%v/v1/meetups/{id:[0-9]+}/weather", basePath), middleware(meetupHandler.MeetupWeather, ALL_SCOPES, true)).Methods(http.MethodGet)

	router.HandleFunc("/auth/token-issue", middleware(userHandler.TokenIssue, nil, false)).Methods(http.MethodPost)

	return router
}

func (hh *healthHandler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
