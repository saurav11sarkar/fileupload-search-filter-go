package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidationStruct(value any) error {
	if err := validate.Struct(value); err != nil {

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			message := make([]string, 0, len(validationErrors))
			for _, ve := range validationErrors {
				message = append(message, fmt.Sprintf("%s failed on %s validation", ve.Field(), ve.Tag()))
			}
			return NewAppError(http.StatusBadRequest, "VALIDATION_ERROR", strings.Join(message, ", "))
		}
		return err
	}
	return nil
}
