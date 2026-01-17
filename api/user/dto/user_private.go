package dto

import (
	"github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserPrivate struct {
	ID            uuid.UUID  `json:"id" binding:"required" validate:"required"`
	Email         string     `json:"email" binding:"required" validate:"required,email"`
	Name          string     `json:"name" binding:"required" validate:"required"`
	ProfilePicURL *string    `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
	Roles         []*Role `json:"roles" validate:"required,dive,required"`
}

func NewUserPrivate(user *model.User, roles []*model.Role) *UserPrivate {
	roless := make([]*Role, len(user.Roles))
	for i, role := range roles {
		roless[i] = NewRole(role)
	}

	return &UserPrivate{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		ProfilePicURL: user.ProfilePicURL,
		Roles:         roless,
	}
}

func (d *UserPrivate) GetValue() *UserPrivate {
	return d
}

func (d *UserPrivate) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utils.FormatValidationErrors(errs), nil
}
