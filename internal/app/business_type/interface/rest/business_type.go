package rest

import (
	"ambic/internal/app/business_type/usecase"
	res "ambic/internal/infra/response"
	"ambic/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type BusinessTypeHandler struct {
	BusinessTypeUseCase usecase.BusinessTypeUsecaseItf
}

func NewBusinessTypeHandler(routerGroup fiber.Router, businessTypeUseCase usecase.BusinessTypeUsecaseItf, middleware middleware.MiddlewareIf) {
	BusinessTypeHandler := BusinessTypeHandler{
		BusinessTypeUseCase: businessTypeUseCase,
	}

	routerGroup = routerGroup.Group("/business-types")
	routerGroup.Get("/", middleware.Authentication, BusinessTypeHandler.GetBusinessTypes)
}

func (h BusinessTypeHandler) GetBusinessTypes(ctx *fiber.Ctx) error {
	resp, err := h.BusinessTypeUseCase.Get()
	if err != nil {
		return res.Error(ctx, err)
	}

	return res.SuccessResponse(ctx, res.GetBusinessTypeSuccess, fiber.Map{
		"business_types": resp,
	})
}
