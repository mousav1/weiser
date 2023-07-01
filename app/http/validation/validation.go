package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validation struct {
	validator *validator.Validate
}

func New() *Validation {
	validate := validator.New()

	// Defining default settings for the validator
	validate.SetTagName("validate")

	return &Validation{validator: validate}
}

func (v *Validation) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func (v *Validation) AddValidation(tag string, fn validator.Func, fields ...string) error {
	for _, field := range fields {
		if err := v.validator.RegisterValidation(fmt.Sprintf("%s_%s", field, tag), fn); err != nil {
			return err
		}
	}
	return nil
}
