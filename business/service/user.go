package service

import (
	"strconv"
	"time"

	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business/model"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	expirationTime = internal.GetEnv("TOKEN_EXPIRATION_TIME", "1")
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

func (us *userService) TokenIssue(ti *business.TokenIssue) (*business.ClaimSet, error) {
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

	accessToken := map[string]interface{}{
		"iss": "meetup-manager",
		"nbf": now.Unix(),
		"exp": expiresAt,
		"jti": uuid,
	}

	if len(user.Scopes) > 0 {
		sNames := make([]model.ScopeName, len(user.Scopes))
		for i, s := range user.Scopes {
			sNames[i] = s.Name
		}
		accessToken["scopes"] = sNames
	}

	return &business.ClaimSet{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}
