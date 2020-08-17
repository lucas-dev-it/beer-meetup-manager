package controller

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business"
)

type userService interface {
	TokenIssue(ti *business.TokenIssue) (*business.ClaimSet, error)
}

type userHandler struct {
	userService userService
}

func NewUserHandler(userService userService) *userHandler {
	return &userHandler{userService: userService}
}

func (uh *userHandler) TokenIssue(w io.Writer, r *http.Request) (*handlerResult, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var tokenIssue business.TokenIssue
	if err := json.Unmarshal(body, &tokenIssue); err != nil {
		return nil, err
	}

	if string(tokenIssue.Password) == "" || tokenIssue.Username == "" {
		return nil, meetupmanager.CustomError{
			Cause:   meetupmanager.ErrBadRequest,
			Type:    meetupmanager.ErrBadRequest,
			Message: "missing parameters",
		}
	}

	claimSet, err := uh.userService.TokenIssue(&tokenIssue)
	if err != nil {
		return nil, err
	}

	return &handlerResult{body: claimSet, status: http.StatusOK}, nil
}
