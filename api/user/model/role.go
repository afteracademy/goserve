package model

import (
	"time"

	"github.com/google/uuid"
)

type RoleCode string

const (
	RoleCodeLearner RoleCode = "LEARNER"
	RoleCodeAdmin   RoleCode = "ADMIN"
	RoleCodeAuthor  RoleCode = "AUTHOR"
	RoleCodeEditor  RoleCode = "EDITOR"
)

type Role struct {
	ID        uuid.UUID
	Code      RoleCode
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

const RolesTableName = "roles"

func NewRole(code RoleCode) *Role {
	now := time.Now()
	r := Role{
		Code:      code,
		Status:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return &r
}
