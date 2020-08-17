package repository

import (
	"github.com/jinzhu/gorm"
	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business/model"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRespository(DB *gorm.DB) *userRepository {
	return &userRepository{db: DB}
}

func (ur *userRepository) FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := ur.db.Where("username = ?", username).Preload("Scopes").Find(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, meetupmanager.ErrNotFound
		}
	}

	return &user, nil
}
