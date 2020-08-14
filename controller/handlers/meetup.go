package handlers

import (
	"io"
	"net/http"

	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/interfaces"
)

type MeetupHandler struct {
	weatherService interfaces.WeatherService
}

func NewMeetupHandler(ws interfaces.WeatherService) *MeetupHandler {
	return &MeetupHandler{weatherService: ws}
}

func (h *MeetupHandler) CalculateBeers(w io.Writer, r *http.Request) (interface{}, int, error) {
	/*	vars := mux.Vars(r)
		retID := vars["id"]

		ID, err := strconv.ParseUint(retID, 10, 64)
		if err != nil {
			return nil, 400, fmt.Errorf("provided ID: %v, is not a valid ID value", retID)
		}*/
	panic("implement me")
}
