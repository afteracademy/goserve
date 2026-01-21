package micro

import (
	"errors"
	"fmt"

	"github.com/afteracademy/goserve/v2/network"
)

func RespondNatsMessage[T any](req NatsRequest, data *T) {
	d, err := network.ValidateDto(data)
	if err != nil {
		req.RespondJSON(NewMessage(d, err))
		return
	}
	req.RespondJSON(NewMessage(d, nil))
}

func RespondNatsError(req NatsRequest, err error) {
	if apiError, ok := err.(network.ApiError); ok {
		msg := fmt.Sprintf("%d:%s", apiError.GetCode(), apiError.GetMessage())
		req.RespondJSON(NewMessage[any](nil, errors.New(msg)))
		return
	}
	req.RespondJSON(NewMessage[any](nil, err))
}
