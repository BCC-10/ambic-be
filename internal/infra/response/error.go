package response

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func (e *Err) Error() string {
	return e.Message
}

func ErrNotFound(str ...string) *Err {
	var message string
	if len(str) == 1 {
		message = str[0] + " " + fiber.ErrNotFound.Message
	} else if len(str) == 2 {
		message = str[0] + " " + str[1]
	} else {
		message += fiber.ErrNotFound.Message
	}
	return &Err{Code: fiber.StatusNotFound, Message: message}
}

func ErrBadRequest(message ...string) *Err {
	var msg string
	if len(message) == 1 {
		msg = message[0]
	} else {
		msg = fiber.ErrBadRequest.Message
	}

	return &Err{Code: fiber.ErrBadRequest.Code, Message: msg}
}

func ErrInternalServer() *Err {
	return &Err{Code: fiber.ErrInternalServerError.Code, Message: fiber.ErrInternalServerError.Message}
}

func ErrUnauthorized(message ...string) *Err {
	var msg string
	if len(message) == 1 {
		msg = message[0]
	} else {
		msg = fiber.ErrUnauthorized.Message
	}

	return &Err{Code: fiber.ErrUnauthorized.Code, Message: msg}
}

func BadRequest(ctx *fiber.Ctx, message ...string) error {
	var msg string
	if len(message) == 1 {
		msg = message[0]
	} else {
		msg = fiber.ErrBadRequest.Message
	}

	return ctx.Status(fiber.ErrBadRequest.Code).JSON(Res{
		StatusCode: fiber.ErrBadRequest.Code,
		Message:    msg,
	})
}

func InternalSeverError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.ErrInternalServerError.Code).JSON(Res{
		StatusCode: fiber.ErrInternalServerError.Code,
		Message:    fiber.ErrInternalServerError.Message,
	})
}

func Error(ctx *fiber.Ctx, err *Err) error {
	var customErr *Err
	if errors.As(err, &customErr) {
		return ctx.Status(customErr.Code).JSON(Res{
			StatusCode: customErr.Code,
			Message:    customErr.Message,
		})
	}

	return InternalSeverError(ctx)
}
