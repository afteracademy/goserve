package micro

import (
	"encoding/json"
	"errors"
	"time"
)

type RequestBuilder[T any] interface {
	NatsClient() NatsClient
	Request(data any) Request[T]
}

type requestBuilder[T any] struct {
	natsClient NatsClient
	subject    string
	timeout    time.Duration
}

func NewRequestBuilder[T any](natsClient NatsClient, subject string) RequestBuilder[T] {
	return &requestBuilder[T]{
		natsClient: natsClient,
		subject:    subject,
		timeout:    natsClient.GetInstance().Timeout,
	}
}

func (c *requestBuilder[T]) NatsClient() NatsClient {
	return c.natsClient
}


func (c *requestBuilder[T]) Request(data any) Request[T] {
	return newRequest(c, data)
}

type Request[T any] interface {
	Nats() (*T, error)
}

type request[T any] struct {
	builder *requestBuilder[T]
	data       any
}

func newRequest[T any](builder *requestBuilder[T], data any) Request[T] {
	return &request[T]{
		builder:    builder,
		data: data,
	}
}

func (r *request[T]) Nats() (*T, error) {
	sendMsg := NewMessage(r.data, nil)
	sendPayload, err := json.Marshal(sendMsg)
	if err != nil {
		return nil, err
	}

	msg, err := r.builder.natsClient.GetInstance().Conn.Request(r.builder.subject, sendPayload, r.builder.timeout)
	if err != nil {
		return nil, err
	}

	var receiveMsg Message[*T]
	err = json.Unmarshal(msg.Data, &receiveMsg)
	if err != nil {
		return nil, err
	}

	if receiveMsg.Error != nil {
		err = errors.New(*receiveMsg.Error)
	}

	return receiveMsg.Data, err
}
