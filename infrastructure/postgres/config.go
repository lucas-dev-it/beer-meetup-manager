package postgres

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business/model"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
)

var (
	host        = internal.GetEnv("DB_HOST", "localhost")
	user        = internal.GetEnv("DB_USER", "local")
	name        = internal.GetEnv("DB_NAME", "beer_meetup")
	port        = internal.GetEnv("DB_PORT", "5440")
	password    = internal.GetEnv("DB_PASSWORD", "password")
	environment = internal.GetEnv("ENVIRONMENT", "dev")
)

// NewPostgres return new connection to db and execute migration.
func NewPostgres() (*gorm.DB, error) {
	connectStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		user,
		name,
		password,
	)

	db, err := gorm.Open("postgres", connectStr)
	if err != nil {
		return nil, err
	}

	db.DB().SetConnMaxLifetime(5 * time.Minute)
	db.DB().SetMaxOpenConns(20)
	db.DB().SetMaxIdleConns(20)

	db.AutoMigrate(model.User{}, model.Scope{}, model.MeetUp{})

	if environment == "dev" {
		if err := prepareTestMigration(db); err != nil {
			fmt.Printf("error when running test migration\n")
		}

	}

	return db, nil
}
