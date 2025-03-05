package env

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"time"
)

type Env struct {
	AppPort       int    `env:"APP_PORT"`
	AppURL        string `env:"APP_URL"`
	MaxUploadSize int64  `env:"MAX_UPLOAD_SIZE"`

	DefaultProfilePhotoPath string `env:"DEFAULT_PROFILE_PHOTO_PATH"`
	DefaultProfilePhotoURL  string

	PartnerVerificationToken string `env:"PARTNER_VERIFICATION_TOKEN"`

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
	GoogleMapsApiKey   string `env:"GOOGLE_MAPS_API_KEY"`

	SupabaseBucket string `env:"SUPABASE_BUCKET"`
	SupabaseURL    string `env:"SUPABASE_URL"`
	SupabaseSecret string `env:"SUPABASE_SECRET"`

	MidtransServerKey string `env:"MIDTRANS_SERVER_KEY"`
}

func New() (*Env, error) {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	_env := new(Env)
	if err := env.Parse(_env); err != nil {
		return nil, err
	}

	_env.DefaultProfilePhotoURL = fmt.Sprintf("%s/%s/%s", _env.SupabaseURL, _env.SupabaseBucket, _env.DefaultProfilePhotoPath)

	return _env, nil
}
