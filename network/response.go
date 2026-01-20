package network

import (
	"net/http"
)

type ResCode string

const (
	success_code ResCode = "10000"
	failue_code  ResCode = "10001"
)

type response[T any] struct {
	ResCode ResCode `json:"code" binding:"required"`
	Status  int     `json:"status" binding:"required"`
	Message string  `json:"message" binding:"required"`
	Data    *T      `json:"data,omitempty" binding:"required,omitempty"`
}

func (r *response[T]) GetResCode() ResCode {
	return r.ResCode
}

func (r *response[T]) GetStatus() int {
	return r.Status
}

func (r *response[T]) GetMessage() string {
	return r.Message
}

func (r *response[T]) GetData() *T {
	return r.Data
}

func NewCustomResponse[T any](rescode ResCode, status int, message string, data *T) Response[T] {
	return &response[T]{
		ResCode: rescode,
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func NewSuccessDataResponse[T any](message string, data *T) Response[T] {
	return &response[T]{
		ResCode: success_code,
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}
}

func NewSuccessMsgResponse(message string) Response[any] {
	return &response[any]{
		ResCode: success_code,
		Status:  http.StatusOK,
		Message: message,
		Data:    nil,
	}
}

func NewBadRequestResponse(message string) Response[any] {
	return &response[any]{
		ResCode: failue_code,
		Status:  http.StatusBadRequest,
		Message: message,
		Data:    nil,
	}
}

func NewForbiddenResponse(message string) Response[any] {
	return &response[any]{
		ResCode: failue_code,
		Status:  http.StatusForbidden,
		Message: message,
		Data:    nil,
	}
}

func NewUnauthorizedResponse(message string) Response[any] {
	return &response[any]{
		ResCode: failue_code,
		Status:  http.StatusUnauthorized,
		Message: message,
		Data:    nil,
	}
}

func NewNotFoundResponse(message string) Response[any] {
	return &response[any]{
		ResCode: failue_code,
		Status:  http.StatusNotFound,
		Message: message,
		Data:    nil,
	}
}

func NewInternalServerErrorResponse(message string) Response[any] {
	return &response[any]{
		ResCode: failue_code,
		Status:  http.StatusInternalServerError,
		Message: message,
		Data:    nil,
	}
}
