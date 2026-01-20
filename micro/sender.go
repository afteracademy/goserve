package micro

import (
	"errors"
	"fmt"

	"github.com/afteracademy/goserve/v2/network"
)

type sender struct{}

func NewMessageSender() MessageSender {
	return &sender{}
}

func (m *sender) SendNats(req NatsRequest) SendMessage {
	return &send{
		natsRequest: req,
	}
}

type send struct {
	natsRequest NatsRequest
}

func (s *send) Message(data any) {
	d, err := network.ValidateDto(&data)
	if err != nil {
		s.natsRequest.RespondJSON(NewMessage(d, err))
		return
	}
	s.natsRequest.RespondJSON(NewMessage(d, nil))
}

func (s *send) Error(err error) {
	if apiError, ok := err.(network.ApiError); ok {
		msg := fmt.Sprintf("%d:%s", apiError.GetCode(), apiError.GetMessage())
		s.natsRequest.RespondJSON(NewMessage[any](nil, errors.New(msg)))
		return
	}
	s.natsRequest.RespondJSON(NewMessage[any](nil, err))
}
