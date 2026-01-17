package model

import (
	"time"

	"github.com/google/uuid"
)

const ApiKeyTableName = "api_keys"

type Permission string

const (
	GeneralPermission Permission = "GENERAL"
)

type ApiKey struct {
	ID          uuid.UUID
	Key         string
	Version     int
	Permissions []Permission
	Comments    []string
	Status      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
