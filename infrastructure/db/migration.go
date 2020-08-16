package db

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business/model"
)

type dbMigration func(db *gorm.DB) error

var prepareTestMigration = func(db *gorm.DB) error {

	user := &model.User{
		Model: gorm.Model{ID: 1},
	}
	if err := db.Find(&user).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}

	if user.Username == "username_0" {
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

		// TODO check existence first before calling this
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
		validMeetup := &model.MeetUp{
			Name:        "test validMeetup",
			Description: "test validMeetup",
			StartDate:   &start,
			EndDate:     &end,
			Country:     "Argentina",
			State:       "Córdoba",
			City:        "Córdoba",
			Attendees:   attendees,
		}
		tx.Create(validMeetup)

		start = time.Now().AddDate(0, 0, 9)
		end = start.Add(time.Hour * time.Duration(2))
		otherLocation := &model.MeetUp{
			Name:        "test other location validMeetup",
			Description: "test other location validMeetup",
			StartDate:   &start,
			EndDate:     &end,
			Country:     "Brasil",
			State:       "Sao Pablo",
			City:        "Sao Pablo",
			Attendees:   attendees,
		}
		tx.Create(otherLocation)


		start = time.Now().AddDate(1, 0, 0)
		end = start.Add(time.Hour * time.Duration(2))
		futureMeetup := &model.MeetUp{
			Name:        "test future",
			Description: "test future",
			StartDate:   &start,
			EndDate:     &end,
			Country:     "argentina",
			State:       "cordoba",
			City:        "cordoba",
			Attendees:   attendees,
		}
		tx.Create(futureMeetup)

		return nil
	})
}
