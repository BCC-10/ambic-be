package limiter

import (
	"ambic/internal/infra/redis"
	"ambic/internal/infra/response"
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
		Next: func(ctx *fiber.Ctx) bool {
			return ctx.IP() == "127.0.0.1"
		},
		LimitReached: func(ctx *fiber.Ctx) error {
			return response.TooManyRequests(ctx, response.LimitExceeded)
		},
		Max:        max,
		Expiration: d,
		Storage:    l.redis,
	})
}
