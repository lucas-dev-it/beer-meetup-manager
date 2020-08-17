package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

const basePath = "/meetup-manager"

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

	router.HandleFunc(fmt.Sprintf("%v/v1/meetups/{id:[0-9]+}/beers", basePath), middleware(meetupHandler.CalculateBeers, true)).Methods(http.MethodGet)
	router.HandleFunc(fmt.Sprintf("%v/v1/meetups/{id:[0-9]+}/weather", basePath), middleware(meetupHandler.MeetupWeather, true)).Methods(http.MethodGet)

	router.HandleFunc("/auth/token-issue", middleware(userHandler.TokenIssue, false)).Methods(http.MethodPost)

	return router
}

func (hh *healthHandler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
