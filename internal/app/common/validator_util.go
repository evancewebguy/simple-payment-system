package common

import "github.com/go-playground/validator/v10"

func ValidateModel(r interface{}) error {
	validate := validator.New()
	return validate.Struct(r)
}
