package dto

import (
	"github.com/afteracademy/goserve/arch/utility"
	"github.com/go-playground/validator/v10"
)

type Tokens struct {
	AccessToken  string `json:"accessToken" binding:"required" validate:"required"`
	RefreshToken string `json:"refreshToken" binding:"required" validate:"required"`
}

func NewTokens(access string, refresh string) *Tokens {
	return &Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}
}

func (d *Tokens) GetValue() *Tokens {
	return d
}

func (d *Tokens) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
