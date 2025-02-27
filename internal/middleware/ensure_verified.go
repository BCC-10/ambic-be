package middleware

import (
	res "ambic/internal/infra/response"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) EnsureVerified(ctx *fiber.Ctx) error {
	isVerified := ctx.Locals("isVerified")

	if isVerified == false {
		return res.Forbidden(ctx)
	}

	return ctx.Next()
}
