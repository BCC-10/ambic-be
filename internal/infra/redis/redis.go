package redis

import (
	"ambic/internal/domain/env"
	"github.com/gofiber/storage/redis/v3"
)

func New(env *env.Env) *redis.Storage {
	store := redis.New(redis.Config{
		Host:     env.RedisHost,
		Port:     env.RedisPort,
		Username: env.RedisUsername,
		Password: env.RedisPassword,
	})

	return store
}
