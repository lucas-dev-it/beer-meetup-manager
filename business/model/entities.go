package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string   `gorm:"index;not null" validate:"email"`
	Password string   `gorm:"not null"`
	Scopes   []*Scope `gorm:"many2many:retailer_scopes"`
}

type Scope struct {
	gorm.Model
	Name        string `gorm:"index;unique;not null"`
	Description string
}

type MeetUp struct {
	gorm.Model
	Name        string     `gorm:"not null"`
	Description string     `gorm:"not null"`
	DateOn      *time.Time `gorm:"index;not null"`
	Country     string     `gorm:"not null"`
	State       string     `gorm:"not null"`
	City        string     `gorm:"not null"`
}
