package helper

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/env"
	res "ambic/internal/infra/response"
	"crypto/rand"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"math/big"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"
)

type HelperIf interface {
	FormParser(ctx *fiber.Ctx, target interface{}) error
	ValidateImage(file *multipart.FileHeader) *res.Err
	GenerateInvoiceNumber() string
	CreatePagination(pagination dto.PaginationRequest) dto.PaginationRequest
	CalculatePagination(request dto.PaginationRequest, totalItems int64) dto.PaginationResponse
}

type Helper struct {
	env *env.Env
}

func New(env *env.Env) HelperIf {
	return &Helper{
		env: env,
	}
}

func (h Helper) FormParser(ctx *fiber.Ctx, target interface{}) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	val = val.Elem()
	typ := val.Type()

	form, err := ctx.MultipartForm()
	if err != nil {
		return err
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)
		formKey := structField.Tag.Get("form")

		if formKey == "" || !field.CanSet() {
			continue
		}

		if values, ok := form.Value[formKey]; ok && len(values) > 0 {
			valueStr := values[0]

			switch field.Kind() {
			case reflect.String:
				field.SetString(valueStr)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				intValue, _ := strconv.ParseInt(valueStr, 10, 64)
				field.SetInt(intValue)
			case reflect.Float32, reflect.Float64:
				floatValue, _ := strconv.ParseFloat(valueStr, 64)
				field.SetFloat(floatValue)
			case reflect.Bool:
				boolValue, _ := strconv.ParseBool(valueStr)
				field.SetBool(boolValue)
			}
		}

		if files, ok := form.File[formKey]; ok && len(files) > 0 {
			if field.Type() == reflect.TypeOf((*multipart.FileHeader)(nil)) {
				field.Set(reflect.ValueOf(files[0]))
			} else if field.Type() == reflect.TypeOf([]*multipart.FileHeader{}) {
				field.Set(reflect.ValueOf(files))
			}
		}
	}

	return nil
}

func (h Helper) ValidateImage(file *multipart.FileHeader) *res.Err {
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return res.ErrUnprocessableEntity(res.PhotoOnly)
	}

	if file.Size > h.env.MaxUploadSize*1024*1024 {
		return res.ErrEntityTooLarge(int(h.env.MaxUploadSize), res.PhotoSizeLimit)
	}

	return nil
}

func (h Helper) GenerateInvoiceNumber() string {
	num, _ := rand.Int(rand.Reader, big.NewInt(999999999999))

	return fmt.Sprintf("%012d", num)
}

func (h Helper) CreatePagination(pagination dto.PaginationRequest) dto.PaginationRequest {
	if pagination.Limit < 1 {
		pagination.Limit = h.env.DefaultPaginationLimit
	}

	if pagination.Page < 1 {
		pagination.Page = h.env.DefaultPaginationPage
	}

	pagination.Offset = (pagination.Page - 1) * pagination.Limit

	return pagination
}

func (h Helper) CalculatePagination(request dto.PaginationRequest, totalItems int64) dto.PaginationResponse {
	totalPages := float64(totalItems) / float64(request.Limit)
	if totalPages > float64(int(totalPages)) {
		totalPages++
	}

	return dto.PaginationResponse{
		Page:      request.Page,
		Limit:     request.Limit,
		TotalData: int(totalItems),
		TotalPage: int(totalPages),
	}
}
