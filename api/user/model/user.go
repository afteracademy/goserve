package model

import (
	"time"

	"github.com/google/uuid"
)

const UserTableName = "users"
const UserRoleRelTableName = "user_roles"

type User struct {
	ID            uuid.UUID
	Email         string
	Name          string
	Password      *string
	ProfilePicURL *string
	Roles         []*Role // not stored in DB directly
	Verified      bool
	Status        bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
