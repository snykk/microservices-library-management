package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var mapHelper = map[string]string{
	"required":    "is a required field",
	"email":       "is not a valid email address",
	"min":         "must be at least %s characters long",
	"max":         "must be less than %s characters",
	"uuid4":       "must be a valid UUID v4",
	"len":         "must be exactly %s characters",
	"oneof":       "must be one of %s",
	"contains":    "must contain '%s'",
	"containsany": "must contain at least one symbol of '%s'",
	"securepwd":   "must contain at least 8 characters, including lowercase, uppercase, a number, and a special character",
}

var needParam = []string{"min", "max", "len", "oneof", "contains", "containsany"}

// ValidatePayloads validates a payload using go-playground validator
func ValidatePayloads(payload interface{}) (map[string]string, error) {
	validate := validator.New()

	// Register custom validation for secure passwords
	validate.RegisterValidation("securepwd", ValidatePassword)

	err := validate.Struct(payload)

	if err == nil {
		return nil, nil // No validation errors
	}

	validationErrors := make(map[string]string)
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			field := strings.ToLower(e.Field())
			tag := e.Tag()
			value := e.Value()
			param := e.Param()

			fmt.Println(field, tag, value, param)

			if value != "" {
				value = fmt.Sprintf("'%s' ", value)
			}

			// Construct the message based on the tag and parameters
			if msgTemplate, exists := mapHelper[tag]; exists {
				if contains(needParam, tag) {
					validationErrors[field] = fmt.Sprintf(msgTemplate, param)
				} else {
					validationErrors[field] = fmt.Sprintf("%v%s", value, msgTemplate)
				}
			} else {
				validationErrors[field] = "Invalid value"
			}
		}
	}

	return validationErrors, fmt.Errorf("validation error")
}

// ValidatePassword validates password to be strong
func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Password rules
	hasMinLength := len(password) >= 8
	hasLower := strings.ContainsAny(password, "abcdefghijklmnopqrstuvwxyz")
	hasUpper := strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	hasSpecial := strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:'\",.<>?/`")
	hasNumber := strings.ContainsAny(password, "0123456789")

	return hasMinLength && hasLower && hasUpper && hasSpecial && hasNumber
}

func contains(array []string, item string) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}
