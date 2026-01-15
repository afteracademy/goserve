package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/afteracademy/goserve/api/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InfoPrivateUser struct {
	ID            primitive.ObjectID `json:"_id" binding:"required" validate:"required"`
	Email         string             `json:"email" binding:"required" validate:"required,email"`
	Name          string             `json:"name" binding:"required" validate:"required"`
	ProfilePicURL *string            `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
	Roles         []*InfoRole        `json:"roles" validate:"required,dive,required"`
}

func NewInfoPrivateUser(user *model.User) *InfoPrivateUser {
	roles := make([]*InfoRole, len(user.Roles))
	for i, role := range user.RoleDocs {
		roles[i] = NewInfoRole(role)
	}

	return &InfoPrivateUser{
		ID:            user.ID,
		Email:         user.Email,
		Name:          user.Name,
		ProfilePicURL: user.ProfilePicURL,
		Roles:         roles,
	}
}

func (d *InfoPrivateUser) GetValue() *InfoPrivateUser {
	return d
}

func (d *InfoPrivateUser) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be min %s", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s must be max%s", err.Field(), err.Param()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
