package dto

import (
	"github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/arch/utility"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type RoleInfo struct {
	ID   uuid.UUID      `json:"id" binding:"required" validate:"required"`
	Code model.RoleCode `json:"code" binding:"required" validate:"required,rolecode"`
}

func NewRoleInfo(role *model.Role) *RoleInfo {
	return &RoleInfo{
		ID:   role.ID,
		Code: role.Code,
	}
}

func (d *RoleInfo) GetValue() *RoleInfo {
	return d
}

func (d *RoleInfo) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
