package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business/model"
)

type dbMigration func(db *gorm.DB) error

var prepareTestMigration = func(db *gorm.DB) error {

	var first model.User
	if err := db.Find(&first).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}

	if first.Username == "username_0" {
		return nil
	}

	return db.Transaction(func(tx *gorm.DB) error {
		var attendees []*model.User
		userScope := &model.Scope{
			Name:        model.UserScope,
			Description: "user scope description",
		}
		adminScope := &model.Scope{
			Name:        model.AdminScope,
			Description: "admin scope description",
		}

		tx.Create(userScope)
		tx.Create(adminScope)

		for i := 0; i < 10; i++ {
			var scopes []*model.Scope
			if i < 3 {
				scopes = append(scopes, userScope, adminScope)
			} else {
				scopes = append(scopes, userScope)
			}

			user := &model.User{
				Username: fmt.Sprintf("username_%v", i),
				Password: fmt.Sprintf("password_%v", i),
				Scopes:   scopes,
			}

			if err := tx.Create(user).Error; err != nil {
				return err
			}
			attendees = append(attendees, user)
		}

		start := time.Now()
		end := start.Add(time.Hour * time.Duration(2))
		meetup := &model.MeetUp{
			Name:        "test meetup",
			Description: "test meetup",
			StartDate:   &start,
			EndDate:     &end,
			Country:     "argentina",
			State:       "cordoba",
			City:        "cordoba",
			Attendees:   attendees,
		}

		tx.Create(meetup)
		return nil
	})
}
