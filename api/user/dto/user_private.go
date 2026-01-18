package dto

import (
	"github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/arch/utility"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UserPrivate struct {
	ID            uuid.UUID   `json:"id" binding:"required" validate:"required"`
	Email         string      `json:"email" binding:"required" validate:"required,email"`
	Name          string      `json:"name" binding:"required" validate:"required"`
	ProfilePicURL *string     `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
	Roles         []*RoleInfo `json:"roles" validate:"required,dive,required"`
}

func NewUserPrivate(user *model.User) *UserPrivate {
	var roles []*RoleInfo
	for _, role := range user.Roles {
		roles = append(roles, NewRoleInfo(role))
	}

	return &UserPrivate{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		ProfilePicURL: user.ProfilePicURL,
		Roles:         roles,
	}
}

func (d *UserPrivate) GetValue() *UserPrivate {
	return d
}

func (d *UserPrivate) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
