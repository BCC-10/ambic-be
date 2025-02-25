package middleware

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (m *Middleware) Authentication(ctx *fiber.Ctx) error {
	authToken := ctx.GetReqHeaders()["Authorization"]

	if len(authToken) < 1 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "token is not provided",
		})
	}

	bearerToken := authToken[0]
	token := strings.Split(bearerToken, " ")

	userId, isAdmin, err := m.jwt.ValidateToken(token[1])
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "invalid token",
		})
	}

	ctx.Locals("userId", userId)
	ctx.Locals("isAdmin", isAdmin)

	return ctx.Next()
}
