package handlers

import (
	"io"
	"net/http"
)

type MeetupHandler struct {
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
