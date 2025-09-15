package xvalidator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

type XValidator struct {
	Validator *validator.Validate
}

type GlobalErrorHandlerResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var Validate = validator.New()

func (v XValidator) Validate(data interface{}) []ErrorResponse {
	var validationErrors []ErrorResponse

	errs := Validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

// ConvertToMessages converts validation errors to user-friendly messages
func ConvertToMessages(errors []ErrorResponse) map[string]string {
	errorMap := make(map[string]string)

	for _, err := range errors {
		field := strings.ToLower(err.FailedField)
		message := getErrorMessage(err)
		errorMap[field] = message
	}

	return errorMap
}

// getErrorMessage generates user-friendly error messages
func getErrorMessage(err ErrorResponse) string {
	field := err.FailedField
	tag := err.Tag

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "min":
		return fmt.Sprintf("%s is too short", field)
	case "max":
		return fmt.Sprintf("%s is too long", field)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %v", field, err.Value)
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %v", field, err.Value)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

var Validator = &XValidator{
	Validator: Validate,
}
