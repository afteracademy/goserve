package dto

import (
	"github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Author struct {
	ID            uuid.UUID `json:"id" binding:"required" validate:"required"`
	Name          string    `json:"name" binding:"required" validate:"required"`
	ProfilePicURL *string   `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
}

func NewAuthor(user *model.User) *Author {
	return &Author{
		ID:            user.ID,
		Name:          user.Name,
		ProfilePicURL: user.ProfilePicURL,
	}
}

func (d *Author) GetValue() *Author {
	return d
}

func (d *Author) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utils.FormatValidationErrors(errs), nil
}
