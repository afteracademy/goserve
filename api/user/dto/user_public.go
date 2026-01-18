package dto

import (
	"github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/arch/utility"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserPublic struct {
	ID            uuid.UUID `json:"id" binding:"required" validate:"required"`
	Name          string    `json:"name" binding:"required" validate:"required"`
	ProfilePicURL *string   `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
}

func NewUserPublic(user *model.User) *UserPublic {
	return &UserPublic{
		ID:            user.ID,
		Name:          user.Name,
		ProfilePicURL: user.ProfilePicURL,
	}
}

func (d *UserPublic) GetValue() *UserPublic {
	return d
}

func (d *UserPublic) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
