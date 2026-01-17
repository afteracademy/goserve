package dto

import (
	"github.com/afteracademy/goserve/utils"
	"github.com/go-playground/validator/v10"
)

type SignUpBasic struct {
	Email         string  `json:"email" binding:"required" validate:"required,email"`
	Password      string  `json:"password" binding:"required" validate:"required,min=6,max=100"`
	Name          string  `json:"name" binding:"required" validate:"required,min=2,max=200"`
	ProfilePicUrl *string `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
}

func EmptySignUpBasic() *SignUpBasic {
	return &SignUpBasic{}
}

func (d *SignUpBasic) GetValue() *SignUpBasic {
	return d
}

func (d *SignUpBasic) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utils.FormatValidationErrors(errs), nil
}
