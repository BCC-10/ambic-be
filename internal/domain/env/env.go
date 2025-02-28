package env

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"time"
)

type Env struct {
	AppPort int    `env:"APP_PORT"`
	AppURL  string `env:"APP_URL"`

	OTPLength      int           `env:"OTP_LENGTH"`
	OTPExpiresTime time.Duration `env:"OTP_EXPIRES_TIME"`

	TokenLength      int           `env:"TOKEN_LENGTH"`
	TokenExpiresTime time.Duration `env:"TOKEN_EXPIRES_TIME"`

	StateExpiresTime time.Duration `env:"STATE_EXPIRES_TIME"`

	DBHost     string `env:"DB_HOST"`
	DBPort     int    `env:"DB_PORT"`
	DBName     string `env:"DB_NAME"`
	DBUsername string `env:"DB_USERNAME"`
	DBPassword string `env:"DB_PASSWORD"`

	JWTSecret  string        `env:"JWT_SECRET"`
	JWTExpires time.Duration `env:"JWT_EXPIRES"`

	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     int    `env:"REDIS_PORT"`
	RedisUsername string `env:"REDIS_USERNAME"`
	RedisPassword string `env:"REDIS_PASSWORD"`

	SMTPHost string `env:"SMTP_HOST"`
	SMTPPort int    `env:"SMTP_PORT"`
	SMTPUser string `env:"SMTP_USER"`
	SMTPPass string `env:"SMTP_PASS"`

	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `env:"GOOGLE_REDIRECT_URL"`
}

func New() (*Env, error) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	_env := new(Env)
	err = env.Parse(_env)
	if err != nil {
		return nil, err
	}

	return _env, nil
}
