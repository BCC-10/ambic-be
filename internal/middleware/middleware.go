package middleware

import (
	"ambic/internal/infra/jwt"
	"github.com/gofiber/fiber/v2"
)

type MiddlewareIf interface {
	Authentication(ctx *fiber.Ctx) error
	EnsurePartner(ctx *fiber.Ctx) error
	EnsureNotPartner(ctx *fiber.Ctx) error
	EnsureVerifiedPartner(ctx *fiber.Ctx) error
}

type Middleware struct {
	jwt jwt.JWTIf
}

func NewMiddleware(jwt jwt.JWTIf) MiddlewareIf {
	return &Middleware{
		jwt: jwt,
	}
}
