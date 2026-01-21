package micro

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewMessage tests the NewMessage function.
func TestNewMessage(t *testing.T) {
	t.Run("with data and no error", func(t *testing.T) {
		data := "test data"
		msg := NewMessage(&data, nil)
		assert.NotNil(t, msg.Data)
		assert.Equal(t, data, *msg.Data)
		assert.Nil(t, msg.Error)
	})

	t.Run("with data and error", func(t *testing.T) {
		data := "test data"
		err := errors.New("test error")
		msg := NewMessage(&data, err)
		assert.NotNil(t, msg.Data)
		assert.Equal(t, data, *msg.Data)
		assert.NotNil(t, msg.Error)
		assert.Equal(t, "test error", *msg.Error)
	})

	t.Run("with nil data and error", func(t *testing.T) {
		var data *string = nil
		err := errors.New("test error")
		msg := NewMessage(data, err)
		assert.Nil(t, msg.Data)
		assert.NotNil(t, msg.Error)
		assert.Equal(t, "test error", *msg.Error)
	})

	t.Run("with nil error", func(t *testing.T) {
		data := 42
		msg := NewMessage(&data, nil)
		assert.NotNil(t, msg.Data)
		assert.Equal(t, 42, *msg.Data)
		assert.Nil(t, msg.Error)
	})
}

// TestNewMsgToJson tests the NewMsgToJson function.
func TestNewMsgToJson(t *testing.T) {
	type TestDTO struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"gte=18"`
	}

	t.Run("valid dto", func(t *testing.T) {
		dto := TestDTO{Name: "John Doe", Age: 30}
		jsonBytes, err := NewMsgToJson(&dto)
		assert.NoError(t, err)
		assert.NotNil(t, jsonBytes)

		var msg message[TestDTO]
		err = json.Unmarshal(jsonBytes, &msg)
		assert.NoError(t, err)
		assert.NotNil(t, msg.Data)
		assert.Equal(t, "John Doe", msg.Data.Name)
		assert.Equal(t, 30, msg.Data.Age)
		assert.Nil(t, msg.Error)
	})

	t.Run("validation error - missing required field", func(t *testing.T) {
		dto := TestDTO{Name: "", Age: 25}
		jsonBytes, err := NewMsgToJson(&dto)
		assert.Error(t, err)
		assert.NotNil(t, jsonBytes)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("validation error - age constraint", func(t *testing.T) {
		dto := TestDTO{Name: "Jane Doe", Age: 17}
		jsonBytes, err := NewMsgToJson(&dto)
		assert.Error(t, err)
		assert.NotNil(t, jsonBytes)
		assert.Contains(t, err.Error(), "age must be greater than or equal to 18")
	})

	t.Run("with nil dto", func(t *testing.T) {
		var dto *TestDTO = nil
		jsonBytes, err := NewMsgToJson(dto)
		assert.Error(t, err)
		assert.NotNil(t, jsonBytes)
	})

	t.Run("with simple type", func(t *testing.T) {
		data := "hello world"
		jsonBytes, err := NewMsgToJson(&data)
		assert.Error(t, err)
		assert.NotNil(t, jsonBytes)

		var msg message[string]
		err = json.Unmarshal(jsonBytes, &msg)
		assert.NoError(t, err)
		assert.NotNil(t, msg.Data)
		assert.Equal(t, "hello world", *msg.Data)
	})
}

// TestNewJsonToMsg tests the NewJsonToMsg function.
func TestNewJsonToMsg(t *testing.T) {
	type TestDTO struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"gte=18"`
	}

	t.Run("valid json message", func(t *testing.T) {
		dto := TestDTO{Name: "John Doe", Age: 30}
		msg := NewMessage(&dto, nil)
		jsonBytes, _ := json.Marshal(msg)

		parsed, err := NewJsonToMsg[TestDTO](jsonBytes)
		assert.NoError(t, err)
		assert.NotNil(t, parsed)
		assert.Equal(t, "John Doe", parsed.Name)
		assert.Equal(t, 30, parsed.Age)
	})

	t.Run("message with error", func(t *testing.T) {
		errMsg := "an error occurred"
		msg := message[TestDTO]{
			Data:  nil,
			Error: &errMsg,
		}
		jsonBytes, _ := json.Marshal(msg)

		parsed, err := NewJsonToMsg[TestDTO](jsonBytes)
		assert.Error(t, err)
		assert.Equal(t, errMsg, err.Error())
		assert.Nil(t, parsed)
	})

	t.Run("invalid json", func(t *testing.T) {
		jsonBytes := []byte(`{"data": "invalid"`)
		parsed, err := NewJsonToMsg[TestDTO](jsonBytes)
		assert.Error(t, err)
		assert.Nil(t, parsed)
	})

	t.Run("validation error - age constraint", func(t *testing.T) {
		dto := TestDTO{Name: "Jane Doe", Age: 17}
		msg := NewMessage(&dto, nil)
		jsonBytes, _ := json.Marshal(msg)

		parsed, err := NewJsonToMsg[TestDTO](jsonBytes)
		assert.Error(t, err)
		assert.NotNil(t, parsed)
		assert.Contains(t, err.Error(), "age must be greater than or equal to 18")
	})

	t.Run("validation error - missing required field", func(t *testing.T) {
		dto := TestDTO{Name: "", Age: 25}
		msg := NewMessage(&dto, nil)
		jsonBytes, _ := json.Marshal(msg)

		parsed, err := NewJsonToMsg[TestDTO](jsonBytes)
		assert.Error(t, err)
		assert.NotNil(t, parsed)
		assert.Contains(t, err.Error(), "name is required")
	})

	t.Run("empty json object", func(t *testing.T) {
		jsonBytes := []byte(`{}`)
		parsed, err := NewJsonToMsg[TestDTO](jsonBytes)
		assert.Error(t, err)
		assert.Nil(t, parsed)
	})

	t.Run("with simple type", func(t *testing.T) {
		data := "hello world"
		msg := NewMessage(&data, nil)
		jsonBytes, _ := json.Marshal(msg)

		parsed, err := NewJsonToMsg[string](jsonBytes)
		assert.Error(t, err)
		assert.NotNil(t, parsed)
		assert.Equal(t, "hello world", *parsed)
	})
}
