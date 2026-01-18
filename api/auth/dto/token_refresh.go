package dto

import (
	"github.com/afteracademy/goserve/arch/utility"
	"github.com/go-playground/validator/v10"
)

type TokenRefresh struct {
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

func EmptyTokenRefresh() *TokenRefresh {
	return &TokenRefresh{}
}

func (d *TokenRefresh) GetValue() *TokenRefresh {
	return d
}

func (d *TokenRefresh) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
