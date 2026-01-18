package network

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBadRequestError(t *testing.T) {
	message := "Bad request"
	err := errors.New("underlying error")
	apiErr := NewBadRequestError(message, err)

	assert.Equal(t, http.StatusBadRequest, apiErr.GetCode())
	assert.Equal(t, message, apiErr.GetMessage())
	assert.EqualError(t, apiErr, fmt.Sprintf("%d - %s: %v", http.StatusBadRequest, message, err))
	assert.ErrorIs(t, apiErr, err)
}

func TestNewForbiddenError(t *testing.T) {
	message := "Forbidden"
	err := errors.New("access denied")
	apiErr := NewForbiddenError(message, err)

	assert.Equal(t, http.StatusForbidden, apiErr.GetCode())
	assert.Equal(t, message, apiErr.GetMessage())
	assert.EqualError(t, apiErr, fmt.Sprintf("%d - %s: %v", http.StatusForbidden, message, err))
	assert.ErrorIs(t, apiErr, err)
}

func TestNewUnauthorizedError(t *testing.T) {
	message := "Unauthorized"
	err := errors.New("authentication failed")
	apiErr := NewUnauthorizedError(message, err)

	assert.Equal(t, http.StatusUnauthorized, apiErr.GetCode())
	assert.Equal(t, message, apiErr.GetMessage())
	assert.EqualError(t, apiErr, fmt.Sprintf("%d - %s: %v", http.StatusUnauthorized, message, err))
	assert.ErrorIs(t, apiErr, err)
}

func TestNewNotFoundError(t *testing.T) {
	message := "Not found"
	err := errors.New("resource not found")
	apiErr := NewNotFoundError(message, err)

	assert.Equal(t, http.StatusNotFound, apiErr.GetCode())
	assert.Equal(t, message, apiErr.GetMessage())
	assert.EqualError(t, apiErr, fmt.Sprintf("%d - %s: %v", http.StatusNotFound, message, err))
	assert.ErrorIs(t, apiErr, err)
}

func TestNewInternalServerError(t *testing.T) {
	message := "Internal server error"
	err := errors.New("server crashed")
	apiErr := NewInternalServerError(message, err)

	assert.Equal(t, http.StatusInternalServerError, apiErr.GetCode())
	assert.Equal(t, message, apiErr.GetMessage())
	assert.EqualError(t, apiErr, fmt.Sprintf("%d - %s: %v", http.StatusInternalServerError, message, err))
	assert.ErrorIs(t, apiErr, err)
}