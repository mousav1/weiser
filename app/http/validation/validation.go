package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validation struct {
	validator *validator.Validate
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

func New() *Validation {
	validate := validator.New()

	// Defining default settings for the validator
	validate.SetTagName("validate")

	return &Validation{validator: validate}
}

func (v *Validation) Validate(i interface{}) ([]ValidationError, error) {
	var errors []ValidationError
	if err := v.validator.Struct(i); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fieldName := err.Field()

			errorMessage := fmt.Sprintf("Field validation for %s failed on the %s", fieldName, err.Tag())
			errors = append(errors, ValidationError{Field: fieldName, Error: errorMessage})
		}
		return errors, fmt.Errorf("validation error")
	}
	return errors, nil
}

type ErrorResponse struct {
	Message string              `json:"message"`
	Errors  map[string][]string `json:"errors"`
}

func (v *Validation) CreateErrorResponse(validationErrors []ValidationError) *ErrorResponse {
	errorMap := make(map[string][]string)
	for _, validationErr := range validationErrors {
		if _, ok := errorMap[validationErr.Field]; ok {
			errorMap[validationErr.Field] = append(errorMap[validationErr.Field], validationErr.Error)
		} else {
			errorMap[validationErr.Field] = []string{validationErr.Error}
		}
	}
	message := fmt.Sprintf("Validation error. (%d errors)", len(validationErrors))
	return &ErrorResponse{
		Message: message,
		Errors:  errorMap,
	}
}

func (v *Validation) AddValidation(tag string, fn validator.Func, fields ...string) error {
	for _, field := range fields {
		if err := v.validator.RegisterValidation(fmt.Sprintf("%s_%s", field, tag), fn); err != nil {
			return err
		}
	}
	return nil
}
