package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func validationErrorToText(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s characters", e.Field(), e.Param())
	case "email":
		return "Invalid email format"
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
	}

	return fmt.Sprintf("%s is not valid", e.Field())
}

func parseError(err error) map[string]string {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		out := make(map[string]string, len(ve))

		for _, fe := range ve {
			out[strings.ToLower(fe.Field())] = validationErrorToText(fe)
		}

		return out
	} else {
		return map[string]string{"error": err.Error()}
	}
}
