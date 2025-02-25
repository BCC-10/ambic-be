package limiter

import (
	"ambic/internal/infra/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"time"
)

type LimiterIf interface {
	Set(max int, duration string) fiber.Handler
}

type Limiter struct {
	redis redis.RedisIf
}

func NewLimiter(redis redis.RedisIf) LimiterIf {
	return &Limiter{
		redis,
	}
}

func (l *Limiter) Set(max int, duration string) fiber.Handler {
	d, _ := time.ParseDuration(duration)

	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: d,
		Storage:    l.redis,
	})
}
