package controller

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business"
	"github.com/stretchr/testify/assert"
)

type uService struct{}

func (u *uService) TokenIssue(ti *business.TokenIssue) (*business.AccessToken, error) {
	if ti.Password == "wronglyTyped" {
		return nil, errors.New("heads up bro")
	}

	if ti.Password == "" && ti.Username == "" {
		return nil, errors.New("heads up bro")
	}

	return &business.AccessToken{
		AccessToken: "token",
		ExpiresAt:   1,
	}, nil
}

func Test_userHandler_TokenIssue(t *testing.T) {
	jn := `{"username":"lucas","password":"pass"}`

	r := &http.Request{
		Body: ioutil.NopCloser(strings.NewReader(jn)),
	}

	handler := userHandler{userService: &uService{}}

	result, err := handler.TokenIssue(nil, r)
	if err != nil {
		t.Error("unexpected error")
	}

	assert.Equal(t, 200, int(result.status))

	token, ok := result.body.(*business.AccessToken)
	if !ok {
		t.Error("unexpected error from type conversion")
	}

	assert.Equal(t, "token", token.AccessToken, )
	assert.Equal(t, 1, int(token.ExpiresAt))
}

func Test_userHandler_TokenIssue_WrongPass(t *testing.T) {
	jn := `{"username":"lucas","password":"wronglyTyped"}`

	r := &http.Request{
		Body: ioutil.NopCloser(strings.NewReader(jn)),
	}

	handler := userHandler{userService: &uService{}}

	result, err := handler.TokenIssue(nil, r)
	if err == nil && result != nil {
		t.Error("expected error")
	}
}

func Test_userHandler_TokenIssue_EmptyFields(t *testing.T) {
	jn := `{"username":"","password":""}`

	r := &http.Request{
		Body: ioutil.NopCloser(strings.NewReader(jn)),
	}

	handler := userHandler{userService: &uService{}}

	result, err := handler.TokenIssue(nil, r)
	if err == nil && result != nil {
		t.Error("expected error")
	}

	assert.IsType(t, meetupmanager.CustomError{}, err)
}
