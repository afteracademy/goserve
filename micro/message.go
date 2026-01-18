package micro

import (
	"encoding/json"
	"errors"
)

type Message[T any] struct {
	Data  T       `json:"data,omitempty"`
	Error *string `json:"error,omitempty"`
}

type AnyMessage = Message[any]

func NewAnyMessage(data any, err error) *AnyMessage {
	var e *string
	if err != nil {
		er := err.Error()
		e = &er
	}

	return &AnyMessage{
		Data:  data,
		Error: e,
	}
}

func NewMessage[T any](data T, err error) *Message[T] {
	var e *string
	if err != nil {
		er := err.Error()
		e = &er
	}

	return &Message[T]{
		Data:  data,
		Error: e,
	}
}

func ParseMsg[T any](data []byte) (*T, error) {
	var msg Message[*T]
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	if msg.Error != nil {
		err = errors.New(*msg.Error)
	}

	return msg.Data, err
}
