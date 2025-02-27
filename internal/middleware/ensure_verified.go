package middleware

import (
	res "ambic/internal/infra/response"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func (m *Middleware) EnsureVerified(ctx *fiber.Ctx) error {
	isVerified := ctx.Locals("isVerified")
	fmt.Println(isVerified)

	if isVerified == false {
		return res.Forbidden(ctx, res.UserNotVerified)
	}

	return ctx.Next()
}
