package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/controller"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/infrastructure/postgres"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/weather"
)

var weatherProviderName = internal.GetEnv("WEATHER_PROVIDER", "weather-stack")

type builder struct {
	postgresDB *gorm.DB
	controller http.Handler
}

func (b *builder) injectDependencies() error {
	postgresConn, err := postgres.NewPostgres()
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	b.postgresDB = postgresConn

	weatherProvider, err := weather.GetProvider(weatherProviderName)
	if err != nil {
		return err
	}

	weatherService, err := weather.NewWeatherService(weatherProvider)
	if err != nil {
		return err
	}

	b.controller = controller.New(weatherService)

	return nil
}

func (b *builder) close() {
	b.postgresDB.Close()
}
