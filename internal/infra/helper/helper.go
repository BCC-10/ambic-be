package helper

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"mime/multipart"
	"reflect"
	"strconv"
)

func ParseForm(ctx *fiber.Ctx, target interface{}) error {
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

		// Jika formKey ada di form.Value (teks)
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

		// Jika formKey ada di form.File (file upload)
		if files, ok := form.File[formKey]; ok && len(files) > 0 {
			if field.Type() == reflect.TypeOf((*multipart.FileHeader)(nil)) {
				// Set file pertama ke field Photo (*multipart.FileHeader)
				field.Set(reflect.ValueOf(files[0]))
			} else if field.Type() == reflect.TypeOf([]*multipart.FileHeader{}) {
				// Set semua file ke field Docs ([]*multipart.FileHeader)
				field.Set(reflect.ValueOf(files))
			}
		}
	}

	return nil
}
