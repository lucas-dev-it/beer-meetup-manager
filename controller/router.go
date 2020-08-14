package controller

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/controller/handlers"
)

var validate *validator.Validate

const basePath = "/meetup-manager"

type healthHandler struct{}

func New() http.Handler {
	router := mux.NewRouter()

	hHealth := &healthHandler{}
	hMeetup := &handlers.MeetupHandler{}

	router.HandleFunc("/health", hHealth.health).Methods(http.MethodGet)

	router.HandleFunc("/v1/meetups/{id:[0-9]+}/beers", middleware(hMeetup.CalculateBeers, true)).Methods(http.MethodGet)

	return router
}

func (hh *healthHandler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
