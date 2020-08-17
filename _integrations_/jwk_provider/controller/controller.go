package controller

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	encodedJWK = os.Getenv("JWK_KEY_SET")
	apiKey     = os.Getenv("INTERNAL_API_KEY")
)

type healthHandler struct{}
type jwkHandler struct{}

func New() http.Handler {
	router := mux.NewRouter()

	hHealth := &healthHandler{}
	jHandler := jwkHandler{}

	router.HandleFunc("/health", hHealth.health).Methods(http.MethodGet)

	router.HandleFunc("/.well-known/jwk.json", jwkMiddleware(jHandler.wellKnown))

	return router
}

func (hh *healthHandler) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func jwkMiddleware(h func(io.Writer, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := h(w, r)
		if err != nil {
			data = err.Error()
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if status != http.StatusNoContent {
			if err := json.NewEncoder(w).Encode(data); err != nil {
				log.Printf("could not encode response to output: %v", err)
			}
		}
	}
}

func (jh *jwkHandler) wellKnown(w io.Writer, r *http.Request) (interface{}, int, error) {
	jwkSet, err := base64.StdEncoding.DecodeString(encodedJWK)
	if err != nil {
		return nil, 500, err
	}

	var set map[string]interface{}
	if err := json.Unmarshal(jwkSet, &set); err != nil {
		return nil, 500, err
	}

	return set, 200, nil
}
