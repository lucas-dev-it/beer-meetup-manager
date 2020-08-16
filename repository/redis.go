package repository

import (
	"encoding/json"
	"time"

	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/infrastructure/cache"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/weather"
)

type redisRepository struct {
	cache *cache.Redis
}

func NewRedisRepository(redis *cache.Redis) *redisRepository {
	return &redisRepository{cache: redis}
}

func (rr *redisRepository) StoreForecast(key string, forecast *weather.Forecast) error {
	jsonForecast, err := json.Marshal(forecast)
	if err != nil {
		return err
	}

	// TODO move this duration as config
	if err = rr.cache.Cli.Set(key, jsonForecast, time.Hour).Err(); err != nil {
		return err
	}

	return nil
}

func (rr *redisRepository) RetrieveForecast(key string) (*weather.Forecast, error) {
	fBytes, err := rr.cache.Cli.Get(key).Bytes()
	if err != nil {
		return nil, err
	}

	var forecast weather.Forecast
	if err := json.Unmarshal(fBytes, &forecast); err != nil {
		return nil, err
	}

	return &forecast, nil
}
