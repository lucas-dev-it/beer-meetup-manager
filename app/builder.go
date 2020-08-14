package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/business/service"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/controller"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/infrastructure/postgres"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/repository"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/weather"
)

var weatherProviderName = internal.GetEnv("WEATHER_PROVIDER", "weather-stack")

type builder struct {
	postgresDB *gorm.DB
	controller http.Handler
}

func (b *builder) injectDependencies() error {
	postgres, err := postgres.NewPostgres()
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	b.postgresDB = postgres

	weatherProvider, err := weather.GetProvider(weatherProviderName)
	if err != nil {
		return err
	}

	weatherService, err := weather.NewWeatherService(weatherProvider)
	if err != nil {
		return err
	}

	meetupRepository := repository.NewMeetupRespository(postgres)

	meetupService := service.NewMeetUpService(meetupRepository, weatherService)

	meetupHandler := controller.NewMeetupHandler(meetupService)

	b.controller = controller.New(meetupHandler)

	return nil
}

func (b *builder) close() {
	b.postgresDB.Close()
}
