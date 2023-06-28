package http

import (
	"mime/multipart"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("name")
		if name == "" {
			name = fld.Tag.Get("json")
		}
		return name
	})
	return &Validator{validator: v}
}

func (v *Validator) Validate(data interface{}) error {
	return v.validator.Struct(data)
}

func (v *Validator) ValidateFiles(files map[string][]*multipart.FileHeader) error {
	for _, fileHeaders := range files {
		for _, fileHeader := range fileHeaders {
			err := v.validator.Var(fileHeader.Filename, "required")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *Validator) SetTagName(tagName string) {
	v.validator.SetTagName(tagName)
}
