package middleware

import (
	res "ambic/internal/infra/response"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) EnsureNotVerified(ctx *fiber.Ctx) error {
	isVerified := ctx.Locals("isVerified")

	if isVerified == true {
		return res.Forbidden(ctx, res.UserVerified)
	}

	return ctx.Next()
}
