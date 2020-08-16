package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate *validator.Validate

const basePath = "/meetup-manager"

type meetupHandler interface {
	CalculateBeers(w io.Writer, r *http.Request) (*handlerResult, error)
}

type healthHandler struct{}

func New(meetupHandler meetupHandler) http.Handler {
	router := mux.NewRouter()

	hHealth := &healthHandler{}

	router.HandleFunc("/health", hHealth.health).Methods(http.MethodGet)

	router.HandleFunc(fmt.Sprintf("%v/v1/meetups/{id:[0-9]+}/beers", basePath), middleware(meetupHandler.CalculateBeers)).Methods(http.MethodGet)
	//router.HandleFunc(fmt.Sprintf("%v/v1/meetups/{id:[0-9]+}/weather", basePath), middleware(meetupHandler.MeetupWeather)).Methods(http.MethodGet)

	return router
}

func (hh *healthHandler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
