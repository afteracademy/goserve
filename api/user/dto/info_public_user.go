package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/afteracademy/goserve/api/user/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InfoPublicUser struct {
	ID            primitive.ObjectID `json:"_id" binding:"required" validate:"required"`
	Name          string             `json:"name" binding:"required" validate:"required"`
	ProfilePicURL *string            `json:"profilePicUrl,omitempty" validate:"omitempty,url"`
}

func NewInfoPublicUser(user *model.User) *InfoPublicUser {
	roles := make([]*InfoRole, len(user.Roles))
	for i, role := range user.RoleDocs {
		roles[i] = NewInfoRole(role)
	}

	return &InfoPublicUser{
		ID:            user.ID,
		Name:          user.Name,
		ProfilePicURL: user.ProfilePicURL,
	}
}

func (d *InfoPublicUser) GetValue() *InfoPublicUser {
	return d
}

func (d *InfoPublicUser) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
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
