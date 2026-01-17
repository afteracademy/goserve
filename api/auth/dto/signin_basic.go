package dto

import (
	"github.com/afteracademy/goserve/utils"
	"github.com/go-playground/validator/v10"
)

type SignInBasic struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=6,max=100"`
}

func EmptySignInBasic() *SignInBasic {
	return &SignInBasic{}
}

func (d *SignInBasic) GetValue() *SignInBasic {
	return d
}

func (d *SignInBasic) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utils.FormatValidationErrors(errs), nil
}
