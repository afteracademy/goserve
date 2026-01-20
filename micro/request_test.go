package micro

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestNewRequestBuilder(t *testing.T) {
	s := RunNatsServerOnPort(t, -1)
	defer s.Shutdown()

	nc, err := nats.Connect(s.ClientURL())
	assert.NoError(t, err)
	defer nc.Close()

	mockNatsClient := &natsClient{
		Conn:    nc,
		Timeout: 1 * time.Second,
	}

	subject := "test.subject"

	t.Run("should create a request builder", func(t *testing.T) {
		builder := NewRequestBuilder[any](mockNatsClient, subject)
		assert.NotNil(t, builder)
		assert.Equal(t, mockNatsClient, builder.NatsClient())
	})

	t.Run("should create builder with correct subject", func(t *testing.T) {
		subject := "user.service.get"
		builder := NewRequestBuilder[any](mockNatsClient, subject)
		assert.NotNil(t, builder)
	})
}

func TestRequestBuilder_Request(t *testing.T) {
	s := RunNatsServerOnPort(t, -1)
	defer s.Shutdown()

	nc, err := nats.Connect(s.ClientURL())
	assert.NoError(t, err)
	defer nc.Close()

	mockNatsClient := &natsClient{
		Conn:    nc,
		Timeout: 1 * time.Second,
	}

	subject := "test.request"

	t.Run("should create a request object", func(t *testing.T) {
		builder := NewRequestBuilder[any](mockNatsClient, subject)
		req := builder.Request("test data")
		assert.NotNil(t, req)
	})

	t.Run("should create request with struct data", func(t *testing.T) {
		type UserRequest struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}

		builder := NewRequestBuilder[any](mockNatsClient, subject)
		data := UserRequest{ID: 1, Name: "test"}
		req := builder.Request(data)
		assert.NotNil(t, req)
	})
}

func TestRequest_Nats(t *testing.T) {
	s := RunNatsServerOnPort(t, -1)
	defer s.Shutdown()

	nc, err := nats.Connect(s.ClientURL())
	assert.NoError(t, err)
	defer nc.Close()

	mockNatsClient := &natsClient{
		Conn:    nc,
		Timeout: 2 * time.Second,
	}

	type TestResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	t.Run("should send request and receive successful response", func(t *testing.T) {
		subject := "test.success"

		// Mock subscriber responding with success
		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			var receivedMsg Message[any]
			json.Unmarshal(m.Data, &receivedMsg)

			respData := TestResponse{Message: "success", Status: "ok"}
			resp := NewMessage(&respData, nil)
			rawResp, _ := json.Marshal(resp)
			m.Respond(rawResp)
		})
		assert.NoError(t, err)

		builder := NewRequestBuilder[TestResponse](mockNatsClient, subject)
		req := builder.Request(map[string]string{"action": "test"})

		resp, err := req.Nats()
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "success", resp.Message)
		assert.Equal(t, "ok", resp.Status)
	})

	t.Run("should handle timeout when no subscriber", func(t *testing.T) {
		subject := "test.timeout"

		mockClientTimeout := &natsClient{
			Conn:    nc,
			Timeout: 1 * time.Nanosecond,
		}

		builder := NewRequestBuilder[TestResponse](mockClientTimeout, subject)
		req := builder.Request("request data")

		resp, err := req.Nats()
		assert.Error(t, err)
		assert.ErrorIs(t, err, nats.ErrTimeout)
		assert.Nil(t, resp)
	})

	t.Run("should handle error in response message", func(t *testing.T) {
		subject := "test.error"

		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			errMsg := "service error occurred"
			var data *string = nil
			resp := NewMessage(data, errors.New(errMsg))
			rawResp, _ := json.Marshal(resp)
			m.Respond(rawResp)
		})
		assert.NoError(t, err)

		builder := NewRequestBuilder[TestResponse](mockNatsClient, subject)
		req := builder.Request("request data")

		resp, err := req.Nats()
		assert.Error(t, err)
		assert.Equal(t, "service error occurred", err.Error())
		assert.Nil(t, resp)
	})

	t.Run("should handle invalid JSON response", func(t *testing.T) {
		subject := "test.invalid.json"

		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			m.Respond([]byte(`{invalid json}`))
		})
		assert.NoError(t, err)

		builder := NewRequestBuilder[TestResponse](mockNatsClient, subject)
		req := builder.Request("request data")

		resp, err := req.Nats()
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("should handle response with data and error", func(t *testing.T) {
		subject := "test.data.error"

		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			data := TestResponse{Message: "partial", Status: "error"}
			resp := NewMessage(&data, errors.New("partial error"))
			rawResp, _ := json.Marshal(resp)
			m.Respond(rawResp)
		})
		assert.NoError(t, err)

		builder := NewRequestBuilder[TestResponse](mockNatsClient, subject)
		req := builder.Request("request data")

		resp, err := req.Nats()
		assert.Error(t, err)
		assert.Equal(t, "partial error", err.Error())
		assert.NotNil(t, resp)
		assert.Equal(t, "partial", resp.Message)
	})

	t.Run("should handle different data types", func(t *testing.T) {
		subject := "test.string"

		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			data := "simple string response"
			resp := NewMessage(&data, nil)
			rawResp, _ := json.Marshal(resp)
			m.Respond(rawResp)
		})
		assert.NoError(t, err)

		builder := NewRequestBuilder[string](mockNatsClient, subject)
		req := builder.Request("request")

		resp, err := req.Nats()
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "simple string response", *resp)
	})
}
