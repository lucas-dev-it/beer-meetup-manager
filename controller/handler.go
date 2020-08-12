package controller

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate *validator.Validate

const basePath = "/api"

type healthHandler struct{}

func New() http.Handler {
	router := mux.NewRouter()

	hHealth := &healthHandler{}

	router.HandleFunc("/health", hHealth.health).Methods(http.MethodGet)

	return router
}

func (hh *healthHandler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
