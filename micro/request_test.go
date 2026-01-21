package micro

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func TestRequestNats(t *testing.T) {
	s := RunNatsServerOnPort(t, -1)
	defer s.Shutdown()

	nc, err := nats.Connect(s.ClientURL())
	assert.NoError(t, err)
	defer nc.Close()

	mockNatsClient := &natsClient{
		Conn:    nc,
		Timeout: 2 * time.Second,
	}

	type TestRequest struct {
		Action string `json:"action" validate:"required"`
	}

	type TestResponse struct {
		Message string `json:"message" validate:"required"`
		Status  string `json:"status" validate:"required"`
	}

	t.Run("should send request and receive successful response", func(t *testing.T) {
		subject := "test.success"

		// Mock subscriber responding with success
		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			var receivedMsg message[TestRequest]
			json.Unmarshal(m.Data, &receivedMsg)

			respData := TestResponse{Message: "success", Status: "ok"}
			resp := NewMessage(&respData, nil)
			rawResp, _ := json.Marshal(resp)
			m.Respond(rawResp)
		})
		assert.NoError(t, err)

		reqData := TestRequest{Action: "test"}
		resp, err := RequestNats[TestRequest, TestResponse](mockNatsClient, subject, &reqData)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "success", resp.Message)
		assert.Equal(t, "ok", resp.Status)
	})

	t.Run("should handle timeout when no subscriber", func(t *testing.T) {
		subject := "test.timeout"

		mockClientTimeout := &natsClient{
			Conn:    nc,
			Timeout: 0 * time.Second,
		}

		reqData := TestRequest{Action: "timeout"}
		resp, err := RequestNats[TestRequest, TestResponse](mockClientTimeout, subject, &reqData)
		assert.Error(t, err)
		assert.ErrorIs(t, err, nats.ErrTimeout)
		assert.Nil(t, resp)
	})

	t.Run("should handle error in response message", func(t *testing.T) {
		subject := "test.error"

		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			errMsg := "service error occurred"
			resp := NewMessage[TestResponse](nil, errors.New(errMsg))
			rawResp, _ := json.Marshal(resp)
			m.Respond(rawResp)
		})
		assert.NoError(t, err)

		reqData := TestRequest{Action: "error"}
		resp, err := RequestNats[TestRequest, TestResponse](mockNatsClient, subject, &reqData)
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

		reqData := TestRequest{Action: "invalid"}
		resp, err := RequestNats[TestRequest, TestResponse](mockNatsClient, subject, &reqData)
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

		reqData := TestRequest{Action: "partial"}
		resp, err := RequestNats[TestRequest, TestResponse](mockNatsClient, subject, &reqData)
		assert.Error(t, err)
		assert.Equal(t, "partial error", err.Error())
		assert.NotNil(t, resp)
		assert.Equal(t, "partial", resp.Message)
	})

	t.Run("should not handle string types", func(t *testing.T) {
		subject := "test.string"

		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			data := "simple string response"
			resp := NewMessage(&data, nil)
			rawResp, _ := json.Marshal(resp)
			m.Respond(rawResp)
		})
		assert.NoError(t, err)

		reqData := "request"
		resp, err := RequestNats[string, string](mockNatsClient, subject, &reqData)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, "only struct payloads are valid for validation", err.Error())
	})

	t.Run("should handle validation error on request", func(t *testing.T) {
		subject := "test.validation"

		// Request with missing required field
		reqData := TestRequest{Action: ""}
		resp, err := RequestNats[TestRequest, TestResponse](mockNatsClient, subject, &reqData)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "action is required")
	})

	t.Run("should handle validation error on response", func(t *testing.T) {
		subject := "test.response.validation"

		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			// Response with missing required field
			data := TestResponse{Message: "incomplete", Status: ""}
			resp := NewMessage(&data, nil)
			rawResp, _ := json.Marshal(resp)
			m.Respond(rawResp)
		})
		assert.NoError(t, err)

		reqData := TestRequest{Action: "validate"}
		resp, err := RequestNats[TestRequest, TestResponse](mockNatsClient, subject, &reqData)
		assert.Error(t, err)
		assert.NotNil(t, resp)
		assert.Contains(t, err.Error(), "status is required")
	})

	t.Run("should handle nil request data", func(t *testing.T) {
		subject := "test.nil.request"

		var reqData *TestRequest = nil
		resp, err := RequestNats[TestRequest, TestResponse](mockNatsClient, subject, reqData)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestRequestNatsRaw(t *testing.T) {
	s := RunNatsServerOnPort(t, -1)
	defer s.Shutdown()

	nc, err := nats.Connect(s.ClientURL())
	assert.NoError(t, err)
	defer nc.Close()

	mockNatsClient := &natsClient{
		Conn:    nc,
		Timeout: 2 * time.Second,
	}

	type TestRequest struct {
		Action string `json:"action"`
	}

	type TestResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	t.Run("should send raw request and receive response", func(t *testing.T) {
		subject := "test.raw.success"

		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			var receivedReq TestRequest
			json.Unmarshal(m.Data, &receivedReq)

			respData := TestResponse{Message: "raw success", Status: "ok"}
			rawResp, _ := json.Marshal(respData)
			m.Respond(rawResp)
		})
		assert.NoError(t, err)

		reqData := TestRequest{Action: "test"}
		resp, msg, err := RequestNatsRaw[TestRequest, TestResponse](mockNatsClient, subject, &reqData)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, msg)
		assert.Equal(t, "raw success", resp.Message)
		assert.Equal(t, "ok", resp.Status)
	})

	t.Run("should handle timeout", func(t *testing.T) {
		subject := "test.raw.timeout"

		mockClientTimeout := &natsClient{
			Conn:    nc,
			Timeout: 0 * time.Second,
		}

		reqData := TestRequest{Action: "timeout"}
		resp, msg, err := RequestNatsRaw[TestRequest, TestResponse](mockClientTimeout, subject, &reqData)
		assert.Error(t, err)
		assert.ErrorIs(t, err, nats.ErrTimeout)
		assert.Nil(t, resp)
		assert.Nil(t, msg)
	})

	t.Run("should handle invalid JSON response", func(t *testing.T) {
		subject := "test.raw.invalid.json"

		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			m.Respond([]byte(`{invalid json}`))
		})
		assert.NoError(t, err)

		reqData := TestRequest{Action: "invalid"}
		resp, msg, err := RequestNatsRaw[TestRequest, TestResponse](mockNatsClient, subject, &reqData)
		assert.Error(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, msg)
	})

	t.Run("should handle marshal error on request", func(t *testing.T) {
		subject := "test.raw.marshal.error"

		// Using a channel which cannot be marshaled to JSON
		type InvalidRequest struct {
			Ch chan int `json:"ch"`
		}

		reqData := InvalidRequest{Ch: make(chan int)}
		resp, msg, err := RequestNatsRaw[InvalidRequest, TestResponse](mockNatsClient, subject, &reqData)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Nil(t, msg)
	})

	t.Run("should return raw message for further inspection", func(t *testing.T) {
		subject := "test.raw.inspect"

		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			respData := TestResponse{Message: "inspect", Status: "ok"}
			rawResp, _ := json.Marshal(respData)
			m.Respond(rawResp)
		})
		assert.NoError(t, err)

		reqData := TestRequest{Action: "inspect"}
		resp, msg, err := RequestNatsRaw[TestRequest, TestResponse](mockNatsClient, subject, &reqData)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, msg)
		assert.NotEmpty(t, msg.Data)
	})

	t.Run("should handle different data types", func(t *testing.T) {
		subject := "test.raw.types"

		_, err := nc.Subscribe(subject, func(m *nats.Msg) {
			resp := 42
			rawResp, _ := json.Marshal(resp)
			m.Respond(rawResp)
		})
		assert.NoError(t, err)

		reqData := 10
		resp, msg, err := RequestNatsRaw[int, int](mockNatsClient, subject, &reqData)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, msg)
		assert.Equal(t, 42, *resp)
	})
}
