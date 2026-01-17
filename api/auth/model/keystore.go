package model

import (
	"time"

	"github.com/google/uuid"
)

const KeystoreTableName = "keystores"

type Keystore struct {
	ID           uuid.UUID
	Client       uuid.UUID
	PrimaryKey   string
	SecondaryKey string
	Status       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewKeystore(clientID uuid.UUID, primaryKey string, secondaryKey string) *Keystore {
	now := time.Now()
	k := Keystore{
		Client:       clientID,
		PrimaryKey:   primaryKey,
		SecondaryKey: secondaryKey,
		Status:       true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	return &k
}

func (keystore *Keystore) GetValue() *Keystore {
	return keystore
}
