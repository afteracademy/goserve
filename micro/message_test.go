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
		msg := NewMessage(data, nil)
		assert.NotNil(t, msg.Data)
		assert.Equal(t, data, *msg.Data)
		assert.Nil(t, msg.Error)
	})

	t.Run("with data and error", func(t *testing.T) {
		data := "test data"
		err := errors.New("test error")
		msg := NewMessage(data, err)
		assert.NotNil(t, msg.Data)
		assert.Equal(t, data, *msg.Data)
		assert.NotNil(t, msg.Error)
		assert.Equal(t, "test error", *msg.Error)
	})

	t.Run("with nil data and error", func(t *testing.T) {
		var data *string
		err := errors.New("test error")
		msg := NewMessage(data, err)
		assert.NotNil(t, msg.Data)
		assert.Equal(t, data, *msg.Data)
		assert.NotNil(t, msg.Error)
		assert.Equal(t, "test error", *msg.Error)
	})
}

// TestParseMsg tests the ParseMsg function.
func TestParseMsg(t *testing.T) {
	type TestDTO struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"gte=18"`
	}

	t.Run("valid message", func(t *testing.T) {
		dto := TestDTO{Name: "John Doe", Age: 30}
		msg := NewMessage(dto, nil)
		raw, _ := json.Marshal(msg)

		parsed, err := ParseMsg[TestDTO](raw)
		assert.NoError(t, err)
		assert.NotNil(t, parsed)
		assert.Equal(t, dto.Name, parsed.Name)
		assert.Equal(t, dto.Age, parsed.Age)
	})

	t.Run("message with error", func(t *testing.T) {
		errMsg := "an error occurred"
		msg := Message[TestDTO]{
			Data:  nil,
			Error: &errMsg,
		}
		raw, _ := json.Marshal(msg)

		parsed, err := ParseMsg[TestDTO](raw)
		assert.Error(t, err)
		assert.Equal(t, errMsg, err.Error())
		assert.Nil(t, parsed)
	})

	t.Run("invalid json", func(t *testing.T) {
		raw := []byte(`{"data": "invalid"`)
		parsed, err := ParseMsg[TestDTO](raw)
		assert.Error(t, err)
		assert.Nil(t, parsed)
	})

	t.Run("validation error", func(t *testing.T) {
		// Age is less than 18, which should fail validation.
		dto := TestDTO{Name: "Jane Doe", Age: 17}
		msg := NewMessage(dto, nil)
		raw, _ := json.Marshal(msg)

		parsed, err := ParseMsg[TestDTO](raw)
		assert.Error(t, err)
		assert.NotNil(t, parsed)
		assert.Contains(t, err.Error(), "age must be greater than or equal to 18")
	})

	t.Run("validation error for required field", func(t *testing.T) {
		// Name is empty, which should fail validation.
		dto := TestDTO{Name: "", Age: 25}
		msg := NewMessage(dto, nil)
		raw, _ := json.Marshal(msg)

		parsed, err := ParseMsg[TestDTO](raw)
		assert.Error(t, err)
		assert.NotNil(t, parsed)
		assert.Contains(t, err.Error(), "name is required")
	})
}

// TestAnyMessage ensures AnyMessage is an alias for Message[any].
func TestAnyMessage(t *testing.T) {
	data := map[string]interface{}{"key": "value"}
	msg := NewMessage[any](data, nil)
	var anyMsg *AnyMessage
	anyMsg = msg
	assert.NotNil(t, anyMsg)
	assert.Equal(t, data, *anyMsg.Data)
}
