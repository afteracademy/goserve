package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSuccessDataResponse(t *testing.T) {
	message := "Success with data"
	data := map[string]interface{}{
		"key": "value",
	}
	resp := NewSuccessDataResponse(message, data)

	assert.Equal(t, success_code, resp.GetResCode())
	assert.Equal(t, "Success with data", resp.GetMessage())
	assert.Equal(t, 200, resp.GetStatus())
	assert.Equal(t, data, resp.GetData())
}

func TestNewSuccessMsgResponse(t *testing.T) {
	message := "Success message"
	resp := NewSuccessMsgResponse(message)

	assert.Equal(t, success_code, resp.GetResCode())
	assert.Equal(t, "Success message", resp.GetMessage())
	assert.Equal(t, 200, resp.GetStatus())
	assert.Nil(t, resp.GetData())
}

func TestNewBadRequestResponse(t *testing.T) {
	message := "Bad request"
	resp := NewBadRequestResponse(message)

	assert.Equal(t, failue_code, resp.GetResCode())
	assert.Equal(t, "Bad request", resp.GetMessage())
	assert.Equal(t, 400, resp.GetStatus())
	assert.Nil(t, resp.GetData())
}

func TestNewForbiddenResponse(t *testing.T) {
	message := "Forbidden"
	resp := NewForbiddenResponse(message)

	assert.Equal(t, failue_code, resp.GetResCode())
	assert.Equal(t, "Forbidden", resp.GetMessage())
	assert.Equal(t, 403, resp.GetStatus())
	assert.Nil(t, resp.GetData())
}

func TestNewUnauthorizedResponse(t *testing.T) {
	message := "Unauthorized"
	resp := NewUnauthorizedResponse(message)

	assert.Equal(t, failue_code, resp.GetResCode())
	assert.Equal(t, "Unauthorized", resp.GetMessage())
	assert.Equal(t, 401, resp.GetStatus())
	assert.Nil(t, resp.GetData())
}

func TestNewNotFoundResponse(t *testing.T) {
	message := "Not found"
	resp := NewNotFoundResponse(message)

	assert.Equal(t, failue_code, resp.GetResCode())
	assert.Equal(t, "Not found", resp.GetMessage())
	assert.Equal(t, 404, resp.GetStatus())
	assert.Nil(t, resp.GetData())
}

func TestNewInternalServerErrorResponse(t *testing.T) {
	message := "Internal server error"
	resp := NewInternalServerErrorResponse(message)

	assert.Equal(t, failue_code, resp.GetResCode())
	assert.Equal(t, "Internal server error", resp.GetMessage())
	assert.Equal(t, 500, resp.GetStatus())
	assert.Nil(t, resp.GetData())
}
