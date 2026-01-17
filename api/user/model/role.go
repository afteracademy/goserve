package model

import (
	"time"

	"github.com/google/uuid"
)

const RolesTableName = "roles"

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
