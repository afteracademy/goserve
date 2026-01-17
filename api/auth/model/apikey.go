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

func NewApiKey(key string, version int, permissions []Permission, comments []string) *ApiKey {
	now := time.Now()
	return &ApiKey{
		Key:         key,
		Version:     version,
		Permissions: permissions,
		Comments:    comments,
		Status:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (apikey *ApiKey) GetValue() *ApiKey {
	return apikey
}
