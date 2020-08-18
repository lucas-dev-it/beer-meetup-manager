package service

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business/model"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	expirationTime = internal.GetEnv("TOKEN_EXPIRATION_TIME", "1")
	signingString  = internal.GetEnv("INTERNAL_API_KEY", "testSigningString")
)

type userRepository interface {
	FindUserByUsername(username string) (*model.User, error)
}

type userService struct {
	userRepository userRepository
}

func NewUserService(userRepository userRepository) *userService {
	return &userService{userRepository: userRepository}
}

func (us *userService) TokenIssue(ti *business.TokenIssue) (*business.AccessToken, error) {
	user, err := us.userRepository.FindUserByUsername(ti.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(ti.Password)); err != nil {
		return nil, meetupmanager.CustomError{
			Cause:   err,
			Type:    meetupmanager.ErrNotFound,
			Message: "username/password combination is invalid",
		}
	}

	now := time.Now()
	uuid := uuid.NewV4()

	et, err := strconv.Atoi(expirationTime)
	if err != nil {
		return nil, err
	}

	expiresAt := now.Add(time.Hour * time.Duration(et)).Unix()
	cm := jwt.MapClaims{
		"iss": "auth-service",
		"nbf": now.Unix(),
		"exp": expiresAt,
		"jti": uuid,
	}

	if len(user.Scopes) > 0 {
		sNames := make([]model.ScopeName, len(user.Scopes))
		for i, s := range user.Scopes {
			sNames[i] = s.Name
		}
		cm["scopes"] = sNames
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, cm)
	tokenString, err := at.SignedString([]byte(signingString))
	if err != nil {
		return nil, err
	}

	return &business.AccessToken{
		AccessToken: tokenString,
		ExpiresAt:   expiresAt,
	}, nil
}
