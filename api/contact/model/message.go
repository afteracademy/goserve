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

func NewMessage(msgType string, msgTxt string) *Message {
	time := time.Now()
	m := Message{
		Type:      msgType,
		Msg:       msgTxt,
		Status:    true,
		CreatedAt: time,
		UpdatedAt: time,
	}
	return &m
}

func (message *Message) GetValue() *Message {
	return message
}
