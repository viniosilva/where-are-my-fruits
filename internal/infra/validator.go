package infra

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func NewValidator() *validator.Validate {
	validate := validator.New()

	validate.RegisterValidation("dgte", ValidateDecimalGreaterThanOrEqual)

	return validate
}

func ValidateDecimalGreaterThanOrEqual(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(decimal.Decimal)
	if !ok {
		return false
	}

	baseValue, err := decimal.NewFromString(fl.Param())
	if err != nil {
		return false
	}

	return value.GreaterThanOrEqual(baseValue)
}
