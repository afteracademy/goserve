package micro

import (
	"encoding/json"
	"errors"

	"github.com/afteracademy/goserve/v2/network"
)

type message[T any] struct {
	Data  *T      `json:"data,omitempty"`
	Error *string `json:"error,omitempty"`
}

func NewMessage[T any](data *T, err error) *message[T] {
	var e *string
	if err != nil {
		er := err.Error()
		e = &er
	}

	return &message[T]{
		Data:  data,
		Error: e,
	}
}

func MsgToJson[T any](data *T) ([]byte, error) {
	msg := NewMessage(data, nil)

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return jsonMsg, err
	}

	_, err = network.ValidateDto(data)
	if err != nil {
		return jsonMsg, err
	}

	return jsonMsg, nil
}

func JsonToMsg[T any](data []byte) (*T, error) {
	var msg message[T]
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}

	if msg.Error != nil {
		err = errors.New(*msg.Error)
		return msg.Data, err
	}

	_, err = network.ValidateDto(msg.Data)
	if err != nil {
		return msg.Data, err
	}

	return msg.Data, err
}
