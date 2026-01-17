package model

import (
	"time"

	"github.com/google/uuid"
)

const UserTableName = "users"

type User struct {
	ID            uuid.UUID
	Name          string
	Email         string
	Password      *string
	ProfilePicURL *string
	Roles         []*Role
	Verified      bool
	Status        bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewUser(email string, pwdHash string, name string, profilePicUrl *string, roles []*Role) *User {
	now := time.Now()
	u := User{
		Email:         email,
		Password:      &pwdHash,
		Name:          name,
		ProfilePicURL: profilePicUrl,
		Roles:         roles,
		Verified:      false,
		Status:        true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	return &u
}
