package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type ScopeName string

var (
	AdminScope ScopeName = "ADMIN"
	UserScope  ScopeName = "USER"
)

type User struct {
	gorm.Model
	Username string   `gorm:"index;not null" validate:"email"`
	Password string   `gorm:"not null"`
	Scopes   []*Scope `gorm:"many2many:user_scopes"`
}

type Scope struct {
	gorm.Model
	Name        ScopeName `gorm:"index;unique;not null"`
	Description string
}

type MeetUp struct {
	gorm.Model
	Name        string     `gorm:"not null"`
	Description string     `gorm:"not null"`
	StartDate   *time.Time `gorm:"index;not null"`
	EndDate     *time.Time `gorm:"index;not null"`
	Country     string     `gorm:"not null"`
	State       string     `gorm:"not null"`
	City        string     `gorm:"not null"`
	Attendees   []*User    `gorm:"many2many:meetup_users"`
}
