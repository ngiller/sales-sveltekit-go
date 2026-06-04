package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ErrorResponseValidation struct {
	FailedField string `json:"field"`
	Tag         string `json:"tag"`
	Value       string `json:"value"`
}

func ValidateStruct(data interface{}) []ErrorResponseValidation {
	var errors []ErrorResponseValidation
	err := validate.Struct(data)
	if err != nil {
		// Handle both ValidationErrors and InvalidValidationError
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				var element ErrorResponseValidation
				element.FailedField = fieldErr.StructNamespace()
				element.Tag = fieldErr.Tag()
				element.Value = fmt.Sprintf("%v", fieldErr.Value())
				errors = append(errors, element)
			}
		}
	}
	return errors
}
