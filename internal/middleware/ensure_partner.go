package middleware

import "github.com/gofiber/fiber/v2"

func (m *Middleware) EnsurePartner(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"isPartner": ctx.Locals("isPartner"),
	})
}
