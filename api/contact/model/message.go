package model

import (
	"time"

	"github.com/google/uuid"
)

const MessageTableName = "messages"

type Message struct {
	ID        uuid.UUID
	Type      string
	Msg       string
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
