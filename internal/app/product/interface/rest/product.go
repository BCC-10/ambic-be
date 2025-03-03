package rest

import (
	"ambic/internal/app/product/usecase"
	res "ambic/internal/infra/response"
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

	routerGroup = routerGroup.Group("/products")
	routerGroup.Post("/create", m.Authentication, m.EnsurePartner, ProductHandler.CreateProduct)
}

func (h ProductHandler) CreateProduct(ctx *fiber.Ctx) error {
	form, err := ctx.MultipartForm()
	if err != nil {
		return res.BadRequest(ctx)
	}

	return res.SuccessResponse(ctx, "success", form)
}
