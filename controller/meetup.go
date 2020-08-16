package controller

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
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

func (mh *MeetupHandler) CalculateBeers(w io.Writer, r *http.Request) (*handlerResult, error) {
	vars := mux.Vars(r)
	retID := vars["id"]

	ID, err := strconv.ParseUint(retID, 10, 64)
	if err != nil {
		return nil, meetupmanager.CustomError{
			Cause:   err,
			Type:    meetupmanager.ErrBadRequest,
			Message: fmt.Sprintf("provided ID: %v, is not a valid ID value", retID),
		}
	}

	meetupBeersData, err := mh.meetupService.CalculateBeerPacksForMeetup(uint(ID))
	if err != nil {
		return nil, err
	} else if meetupBeersData == nil {
		return &handlerResult{
			status: 204,
		}, nil
	}

	return &handlerResult{
		status: 200,
		body:   meetupBeersData,
	}, nil
}
