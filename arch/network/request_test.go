package network

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestReqBody(t *testing.T) {
	body := `{"field": "test"}`

	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqBody(ctx, &MockDto{})
		assert.NoError(t, err)
		assert.Equal(t, dto.Field, "test")
	}

	MockTestHandler(t, "POST", "/mock", "/mock", body, mockHandler, nil)
}

func TestReqBody_Error(t *testing.T) {
	body := `{"wrong": "test"}`

	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqBody(ctx, &MockDto{})
		assert.Nil(t, dto)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "field is required")
	}

	MockTestHandler(t, "POST", "/mock", "/mock", body, mockHandler, nil)
}

func TestReqQuery(t *testing.T) {
	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqQuery(ctx, &MockDto{})
		assert.NoError(t, err)
		assert.Equal(t, dto.Field, "test")
	}

	MockTestHandler(t, "GET", "/mock", "/mock?field=test", "", mockHandler, nil)
}

func TestReqQuery_Error(t *testing.T) {
	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqQuery(ctx, &MockDto{})
		assert.Nil(t, dto)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "field is required")
	}

	MockTestHandler(t, "GET", "/mock", "/mock?wrong=test", "", mockHandler, nil)
}
