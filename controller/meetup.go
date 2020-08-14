package controller

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business"
)

type meetupService interface {
	CalculateBeerPacksForMeetup(meetupID uint) (*business.MeetupBeersData, error)
}

type MeetupHandler struct {
	meetupService meetupService
}

func NewMeetupHandler(meetupService meetupService) *MeetupHandler {
	return &MeetupHandler{meetupService: meetupService}
}

func (mh *MeetupHandler) CalculateBeers(w io.Writer, r *http.Request) (interface{}, int, error) {
	vars := mux.Vars(r)
	retID := vars["id"]

	ID, err := strconv.ParseUint(retID, 10, 64)
	if err != nil {
		return nil, 400, fmt.Errorf("provided ID: %v, is not a valid ID value", retID)
	}

	meetupBeersData, err := mh.meetupService.CalculateBeerPacksForMeetup(uint(ID))
	if err != nil {
		// TODO analize error types here
		return nil, 500, err
	}

	return meetupBeersData, 200, nil
}
