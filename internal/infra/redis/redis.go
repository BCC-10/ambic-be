package redis

import (
	"ambic/internal/domain/env"
	"github.com/gofiber/storage/redis/v3"
	"time"
)

type RedisIf interface {
	Get(key string) ([]byte, error)
	Set(key string, val []byte, exp time.Duration) error
	Delete(key string) error
	Reset() error
	Close() error
}

type Redis struct {
	store *redis.Storage
}

func NewRedis(env *env.Env) RedisIf {
	store := redis.New(redis.Config{
		Host:     env.RedisHost,
		Port:     env.RedisPort,
		Username: env.RedisUsername,
		Password: env.RedisPassword,
	})

	return &Redis{
		store,
	}
}

func (r *Redis) Get(key string) ([]byte, error) {
	value, err := r.store.Get(key)
	if err != nil {
		return make([]byte, 0), err
	}

	return value, nil
}

func (r *Redis) Set(key string, val []byte, exp time.Duration) error {
	return r.store.Set(key, val, exp)
}

func (r *Redis) Delete(key string) error {
	return r.store.Delete(key)
}

func (r *Redis) Reset() error {
	return r.store.Reset()
}

func (r *Redis) Close() error {
	return r.store.Close()
}
