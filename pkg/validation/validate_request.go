package requestuestValidation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationErrorResponse struct {
	Error       bool        `json:"error"`
	FailedField string      `json:"field"`
	Tag         string      `json:"tag"`
	Value       interface{} `json:"value"`
	// EmailErr    string      `json:"email,omitempty"`
	// LenErr      string      `json:"len,omitempty"`
}

func init() {
	validate.RegisterValidation("pincode", validatePincode)
}

var validate = validator.New()

func ValidateRequest(req interface{}) []string {
	errResponse := []string{}
	errs := validate.Struct(req)

	if errs == nil {
		return nil
	}

	for _, err := range errs.(validator.ValidationErrors) {
		e := ValidationErrorResponse{
			Error:       true,
			FailedField: err.Field(),
			Tag:         err.Tag(),
			Value:       err.Value(),
		}

		message := fmt.Sprintf("[%s]: '%v' | Needs to implement '%s'", e.FailedField, e.Value, e.Tag)

		errResponse = append(errResponse, message)
	}
	return errResponse
}

func validatePincode(fl validator.FieldLevel) bool {
	value := fl.Field().Uint()

	return (value >= 110000 && value <= 899999)
}
