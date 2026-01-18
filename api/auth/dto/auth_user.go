package dto

import (
	"github.com/afteracademy/goserve/api/user/dto"
	"github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/arch/utility"
	"github.com/go-playground/validator/v10"
)

type UserAuth struct {
	User   *dto.UserPrivate `json:"user" validate:"required"`
	Tokens *Tokens          `json:"tokens" validate:"required"`
}

func NewUserAuth(user *model.User, tokens *Tokens) *UserAuth {
	return &UserAuth{
		User:   dto.NewUserPrivate(user),
		Tokens: tokens,
	}
}

func (d *UserAuth) GetValue() *UserAuth {
	return d
}

func (d *UserAuth) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
