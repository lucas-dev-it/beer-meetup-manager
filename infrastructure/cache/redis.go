package cache

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/lucas-dev-it/62252aee-9d11-4149-a0ea-de587cbcd233/internal"
)

var (
	redisHost = internal.GetEnv("REDIS_HOST", "redis")
	redisPort = internal.GetEnv("REDIS_PORT", "6379")
)

type Redis struct {
	Cli *redis.Client
}

func NewRedis() *Redis {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})
	return &Redis{Cli: client}
}
