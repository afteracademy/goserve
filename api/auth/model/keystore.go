package model

import (
	"time"

	"github.com/google/uuid"
)

const KeystoreTableName = "keystores"

type Keystore struct {
	ID           uuid.UUID
	UserID     uuid.UUID
	PrimaryKey   string
	SecondaryKey string
	Status       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
