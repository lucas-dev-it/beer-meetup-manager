package repository

import (
	"fmt"

	"github.com/jinzhu/gorm"
	meetupmanager "github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business/model"
)

type MeetupRepository struct {
	db *gorm.DB
}

// NewMeetupRespository new meetup respository instance
func NewMeetupRespository(DB *gorm.DB) *MeetupRepository {
	return &MeetupRepository{db: DB}
}

// CountMeetupAttendees counts the attendees number for a particular meetup ID
func (mr *MeetupRepository) CountMeetupAttendees(meetupID uint) int {
	meetup := &model.MeetUp{
		Model: gorm.Model{ID: meetupID},
	}
	return mr.db.Find(meetup).Association("Attendees").Count()
}

// FindMeetupByID gets the meetup by providing its ID
func (mr *MeetupRepository) FindMeetupByID(meetupID uint) (*model.MeetUp, error) {
	meetup := &model.MeetUp{
		Model: gorm.Model{ID: meetupID},
	}
	if err := mr.db.Find(meetup).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, meetupmanager.CustomError{
				Cause: err,
				Type:  meetupmanager.ErrNotFound,
				Message: fmt.Sprintf("meetup with id: %v not found", meetupID),
			}
		}
		return nil, err
	}

	return meetup, nil
}
