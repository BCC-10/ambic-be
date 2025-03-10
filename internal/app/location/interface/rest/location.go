package rest

import (
	"ambic/internal/app/location/usecase"
	"ambic/internal/domain/dto"
	res "ambic/internal/infra/response"
	"ambic/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	Validator       *validator.Validate
	LocationUsecase usecase.LocationUsecaseItf
}

func NewLocationHandler(routerGroup fiber.Router, locationUsecase usecase.LocationUsecaseItf, m middleware.MiddlewareIf, validator *validator.Validate) {
	LocationHandler := LocationHandler{
		LocationUsecase: locationUsecase,
		Validator:       validator,
	}

	routerGroup = routerGroup.Group("/locations", m.Authentication)
	routerGroup.Get("/", LocationHandler.AutocompleteLocation)
}

func (h *LocationHandler) AutocompleteLocation(ctx *fiber.Ctx) error {
	req := new(dto.LocationRequest)
	if err := ctx.QueryParser(req); err != nil {
		return res.BadRequest(ctx)
	}

	if err := h.Validator.Struct(req); err != nil {
		return res.ValidationError(ctx, err)
	}

	data, err := h.LocationUsecase.AutocompleteLocation(*req)
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.GetAutoCompleteSuccess, fiber.Map{
		"locations": data,
	})
}
