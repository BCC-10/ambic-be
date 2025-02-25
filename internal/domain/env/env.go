package env

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Env struct {
	AppPort int `env:"APP_PORT"`

	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBUsername string `env:"DB_USERNAME"`
	DBPassword string `env:"DB_PASSWORD"`

	JWTSecret  string `env:"JWT_SECRET"`
	JWTExpires int    `env:"JWT_EXPIRES"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     int    `env:"REDIS_PORT"`
	RedisUsername string `env:"REDIS_USERNAME"`
	RedisPassword string `env:"REDIS_PASSWORD"`
}

func New() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env")
	}

	_env := new(Env)
	err = env.Parse(_env)
	if err != nil {
		return nil, err
	}

	return _env, nil
}
