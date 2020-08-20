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
	TokenIssue(ti *business.TokenIssue) (*business.AccessToken, error)
}

type userHandler struct {
	userService userService
}

// NewUserHandler gets a new instances for this handler
func NewUserHandler(userService userService) *userHandler {
	return &userHandler{userService: userService}
}

// TokenIssue parses request data to issue a new JWT token
func (uh *userHandler) TokenIssue(w io.Writer, r *http.Request) (*handlerResult, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var tokenIssue business.TokenIssue
	if err := json.Unmarshal(body, &tokenIssue); err != nil {
		return nil, err
	}

	if tokenIssue.Password == "" || tokenIssue.Username == "" {
		return nil, meetupmanager.CustomError{
			Cause:   meetupmanager.ErrBadRequest,
			Type:    meetupmanager.ErrBadRequest,
			Message: "missing parameters",
		}
	}

	token, err := uh.userService.TokenIssue(&tokenIssue)
	if err != nil {
		return nil, err
	}

	return &handlerResult{body: token, status: http.StatusOK}, nil
}
