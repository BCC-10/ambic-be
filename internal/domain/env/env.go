package env

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"time"
	"os"
	"log"
)

type Env struct {
<<<<<<< HEAD
        AppEnv		       string `env:"APP_ENV"`
	AppPort                int    `env:"APP_PORT"`
	AppURL                 string `env:"APP_URL"`
	AppLogoPath            string `env:"APP_LOGO_PATH"`
	AppLogoURL             string
	MaxUploadSize          int64   `env:"MAX_UPLOAD_SIZE"`
	DefaultPaginationLimit int     `env:"DEFAULT_PAGINATION_LIMIT"`
	DefaultPaginationPage  int     `env:"DEFAULT_PAGINATION_PAGE"`
	DefaultUserLatitude    float64 `env:"DEFAULT_USER_LATITUDE"`
	DefaultUserLongitude   float64 `env:"DEFAULT_USER_LONGITUDE"`
	DefaultUserRadius      float64 `env:"DEFAULT_USER_RADIUS"`

	DefaultProfilePhotoPath        string `env:"DEFAULT_PROFILE_PHOTO_PATH"`
	DefaultProfilePhotoURL         string
	DefaultPartnerProfilePhotoPath string `env:"DEFAULT_PARTNER_PROFILE_PHOTO_PATH"`
	DefaultPartnerProfilePhotoURL  string

	PartnerVerificationToken string `env:"PARTNER_VERIFICATION_TOKEN"`

	TokenLength                     int           `env:"TOKEN_LENGTH"`
	TokenExpiresTime                time.Duration `env:"TOKEN_EXPIRES_TIME"`
	PartnerVerificationTokenExpires time.Duration `env:"PARTNER_VERIFICATION_TOKEN_EXPIRES"`

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

	MidtransServerKey          string `env:"MIDTRANS_SERVER_KEY"`
	MidtransMaxPaymentDuration int64  `env:"MIDTRANS_MAX_PAYMENT_DURATION"`
	MidtransEnvironment        string `env:"MIDTRANS_ENVIRONMENT"`

	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN"`
	TelegramChatID   int64  `env:"TELEGRAM_CHAT_ID"`
}

func New() (*Env, error) {
	if os.Getenv("APP_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, this is fine in Docker.")
		}
        }

	_env := new(Env)
	if err := env.Parse(_env); err != nil {
		return nil, err
	}

	_env.DefaultProfilePhotoURL = fmt.Sprintf("%s/storage/v1/object/public/%s/%s", _env.SupabaseURL, _env.SupabaseBucket, _env.DefaultProfilePhotoPath)

	_env.DefaultPartnerProfilePhotoURL = fmt.Sprintf("%s/storage/v1/object/public/%s/%s", _env.SupabaseURL, _env.SupabaseBucket, _env.DefaultPartnerProfilePhotoPath)

	_env.AppLogoURL = fmt.Sprintf("%s/storage/v1/object/public/%s/%s", _env.SupabaseURL, _env.SupabaseBucket, _env.AppLogoPath)

	return _env, nil
}
