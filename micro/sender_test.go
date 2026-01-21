package micro

import (
	"errors"
	"testing"

	"github.com/afteracademy/goserve/v2/network"
	"github.com/nats-io/nats.go/micro"
	"github.com/stretchr/testify/mock"
)

// MockNatsRequest is a mock implementation of NatsRequest interface
type MockNatsRequest struct {
	mock.Mock
}

func (m *MockNatsRequest) Respond(data []byte, opts ...micro.RespondOpt) error {
	args := m.Called(data, opts)
	return args.Error(0)
}

func (m *MockNatsRequest) RespondJSON(v any, opts ...micro.RespondOpt) error {
	args := m.Called(v, opts)
	return args.Error(0)
}

func (m *MockNatsRequest) Data() []byte {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]byte)
}

func (m *MockNatsRequest) Subject() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockNatsRequest) Headers() micro.Headers {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(micro.Headers)
}

func (m *MockNatsRequest) Error(code, description string, data []byte, opts ...micro.RespondOpt) error {
	args := m.Called(code, description, data, opts)
	return args.Error(0)
}

func (m *MockNatsRequest) Reply() string {
	args := m.Called()
	return args.String(0)
}

func TestRespond_Message(t *testing.T) {
	type TestDTO struct {
		Name string `json:"name" validate:"required"`
		Age  int    `json:"age" validate:"gte=18"`
	}

	t.Run("should respond with data when validation passes", func(t *testing.T) {
		mockNatsRequest := new(MockNatsRequest)
		data := TestDTO{Name: "John", Age: 30}

		mockNatsRequest.On("RespondJSON", mock.Anything, mock.Anything).Return(nil).Once()

		RespondNatsMessage(mockNatsRequest, &data)

		mockNatsRequest.AssertExpectations(t)
	})

	t.Run("should respond with error when validation fails", func(t *testing.T) {
		mockNatsRequest := new(MockNatsRequest)
		data := TestDTO{Name: "", Age: 17} // Both fields fail validation

		mockNatsRequest.On("RespondJSON", mock.Anything, mock.Anything).Return(nil).Once()

		RespondNatsMessage(mockNatsRequest, &data)

		mockNatsRequest.AssertExpectations(t)
	})

	t.Run("should respond with validated data with error for simple types", func(t *testing.T) {
		mockNatsRequest := new(MockNatsRequest)
		data := "simple string"

		mockNatsRequest.On("RespondJSON", mock.Anything, mock.Anything).Return(nil).Once()

		RespondNatsMessage(mockNatsRequest, &data)

		mockNatsRequest.AssertExpectations(t)
	})
}

func TestRespond_Error(t *testing.T) {
	t.Run("should respond with standard error", func(t *testing.T) {
		mockNatsRequest := new(MockNatsRequest)
		err := errors.New("standard error")

		mockNatsRequest.On("RespondJSON", mock.Anything, mock.Anything).Return(nil).Once()

		RespondNatsError(mockNatsRequest, err)

		mockNatsRequest.AssertExpectations(t)
	})

	t.Run("should respond with ApiError formatted as code:message", func(t *testing.T) {
		mockNatsRequest := new(MockNatsRequest)
		apiErr := network.NewNotFoundError("resource not found", nil)

		mockNatsRequest.On("RespondJSON", mock.Anything, mock.Anything).Return(nil).Once()

		RespondNatsError(mockNatsRequest, apiErr)

		mockNatsRequest.AssertExpectations(t)
	})

	t.Run("should respond with BadRequestError", func(t *testing.T) {
		mockNatsRequest := new(MockNatsRequest)
		apiErr := network.NewBadRequestError("invalid input", nil)

		mockNatsRequest.On("RespondJSON", mock.Anything, mock.Anything).Return(nil).Once()

		RespondNatsError(mockNatsRequest, apiErr)

		mockNatsRequest.AssertExpectations(t)
	})

	t.Run("should respond with UnauthorizedError", func(t *testing.T) {
		mockNatsRequest := new(MockNatsRequest)
		apiErr := network.NewUnauthorizedError("unauthorized access", nil)

		mockNatsRequest.On("RespondJSON", mock.Anything, mock.Anything).Return(nil).Once()

		RespondNatsError(mockNatsRequest, apiErr)

		mockNatsRequest.AssertExpectations(t)
	})

	t.Run("should respond with InternalServerError", func(t *testing.T) {
		mockNatsRequest := new(MockNatsRequest)
		apiErr := network.NewInternalServerError("internal error", nil)

		mockNatsRequest.On("RespondJSON", mock.Anything, mock.Anything).Return(nil).Once()

		RespondNatsError(mockNatsRequest, apiErr)

		mockNatsRequest.AssertExpectations(t)
	})
}
