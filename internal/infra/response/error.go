package response

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func (e *Err) Error() string {
	return e.Message
}

func ErrNotFound(message ...string) *Err {
	var msg string
	if len(message) == 1 {
		msg = message[0]
	} else {
		msg = fiber.ErrNotFound.Message
	}

	return &Err{Code: fiber.ErrNotFound.Code, Message: msg}
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

func ErrForbidden(message ...string) *Err {
	var msg string
	if len(message) == 1 {
		msg = message[0]
	} else {
		msg = fiber.ErrForbidden.Message
	}

	return &Err{Code: fiber.ErrForbidden.Code, Message: msg}
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

func ValidationError(ctx *fiber.Ctx, err error) error {
	_errors := make(map[string]string)
	old := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		_errors[field] = strings.Trim(fmt.Sprintf("%s: %s %s", field, err.Tag(), err.Param()), " ")
		old[field] = err.Value().(string)
	}

	return ctx.Status(fiber.ErrBadRequest.Code).JSON(Res{
		StatusCode: fiber.ErrBadRequest.Code,
		Message:    fiber.ErrBadRequest.Message,
		Payload: map[string]interface{}{
			"errors": _errors,
			"old":    old,
		},
	})
}

func InternalSeverError(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.ErrInternalServerError.Code).JSON(Res{
		StatusCode: fiber.ErrInternalServerError.Code,
		Message:    fiber.ErrInternalServerError.Message,
	})
}

func Unauthorized(ctx *fiber.Ctx, message ...string) error {
	var msg string
	if len(message) == 1 {
		msg = message[0]
	} else {
		msg = fiber.ErrUnauthorized.Message
	}

	return ctx.Status(fiber.ErrUnauthorized.Code).JSON(Res{
		StatusCode: fiber.ErrUnauthorized.Code,
		Message:    msg,
	})
}

func Forbidden(ctx *fiber.Ctx, message ...string) error {
	var msg string
	if len(message) == 1 {
		msg = message[0]
	} else {
		msg = fiber.ErrForbidden.Message
	}

	return ctx.Status(fiber.ErrForbidden.Code).JSON(Res{
		StatusCode: fiber.ErrForbidden.Code,
		Message:    msg,
	})
}

func TooManyRequests(ctx *fiber.Ctx, message ...string) error {
	var msg string
	if len(message) == 1 {
		msg = message[0]
	} else {
		msg = fiber.ErrTooManyRequests.Message
	}

	return ctx.Status(fiber.ErrTooManyRequests.Code).JSON(Res{
		StatusCode: fiber.ErrTooManyRequests.Code,
		Message:    msg,
	})
}

func Error(ctx *fiber.Ctx, err *Err) error {
	var customErr *Err
	if errors.As(err, &customErr) {
		return ctx.Status(customErr.Code).JSON(Res{
			StatusCode: customErr.Code,
			Message:    customErr.Message,
			Payload:    customErr.Payload,
		})
	}

	return InternalSeverError(ctx)
}
