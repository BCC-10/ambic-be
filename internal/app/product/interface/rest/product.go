package rest

import (
	"ambic/internal/app/product/usecase"
	"ambic/internal/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	ProductUsecase usecase.ProductUsecaseItf
	validator      *validator.Validate
}

func NewProductHandler(routerGroup fiber.Router, productUsecase usecase.ProductUsecaseItf, validator *validator.Validate, m middleware.MiddlewareIf) {
	ProductHandler := ProductHandler{
		ProductUsecase: productUsecase,
		validator:      validator,
	}

	routerGroup.Group("/products")
}
