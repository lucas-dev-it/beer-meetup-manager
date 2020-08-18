package service

import (
	"errors"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business/model"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal/token"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

type uRepo struct{}

func (u uRepo) FindUserByUsername(username string) (*model.User, error) {
	var scopes []*model.Scope
	if username == "admin@mail.com" {
		scopes = []*model.Scope{{
			Model:       gorm.Model{ID: 1},
			Name:        "ADMIN",
			Description: "ADMIN",
		}, {
			Model:       gorm.Model{ID: 2},
			Name:        "USER",
			Description: "USER",
		}}
	} else if username == "user@mail.com" {
		scopes = []*model.Scope{{
			Model:       gorm.Model{ID: 2},
			Name:        "USER",
			Description: "USER",
		}}
	} else {
		return nil, errors.New("user not found")
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("123qwe"), 12)

	return &model.User{
		Model:    gorm.Model{ID: 100},
		Username: username,
		Password: string(password),
		Scopes:   scopes,
	}, nil
}

func getUserService() *userService {
	return &userService{userRepository: uRepo{}}
}

func Test_userService_TokenIssue(t *testing.T) {
	s := getUserService()

	ti, err := s.TokenIssue(&business.TokenIssue{
		Username: "admin@mail.com",
		Password: "123qwe",
	})
	if err != nil {
		t.Fatal("unexpected error")
	}
	assert.NotNil(t, ti)

	tkn, err := token.ParseTokenString(ti.AccessToken, "testSigningString")
	if err != nil {
		t.Fatal("unexpected error")
	}

	assert.Nil(t, tkn.Claims.Valid())

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("unexpected error, missing claims")
	}

	ss, ok := claims["scopes"]
	if !ok {
		t.Fatal("unexpected error, missing scopes")
	}

	scopes := ss.([]interface{})

	assert.Len(t, scopes, 2)
}

func Test_userService_TokenIssue_WrongPassword(t *testing.T) {
	s := getUserService()

	ti, err := s.TokenIssue(&business.TokenIssue{
		Username: "admin@mail.com",
		Password: "123qwe----",
	})
	if err == nil && ti != nil {
		t.Fatal("expected error")
	}
	assert.Nil(t, ti)
	assert.IsType(t, meetupmanager.CustomError{}, err)
}
