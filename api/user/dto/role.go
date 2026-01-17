package dto

import (
	"github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Role struct {
	ID   uuid.UUID      `json:"id" binding:"required" validate:"required"`
	Code model.RoleCode `json:"code" binding:"required" validate:"required,rolecode"`
}

func NewRole(role *model.Role) *Role {
	return &Role{
		ID:   role.ID,
		Code: role.Code,
	}
}

func EmptyRole() *Role {
	return &Role{}
}

func (d *Role) GetValue() *Role {
	return d
}

func (d *Role) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utils.FormatValidationErrors(errs), nil
}
