package network

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestReqBody_Payload(t *testing.T) {
	body := `{"field": "test"}`

	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqBody[MockPayload](ctx)
		assert.NoError(t, err)
		assert.Equal(t, dto.Field, "test")
	}

	MockTestHandler(t, "POST", "/mock", "/mock", body, mockHandler, nil)
}

func TestReqBody_Dto(t *testing.T) {
	body := `{"field": "test"}`

	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqBody[MockDto](ctx)
		assert.NoError(t, err)
		assert.Equal(t, dto.Field, "test")
	}

	MockTestHandler(t, "POST", "/mock", "/mock", body, mockHandler, nil)
}

func TestReqBody_DtoV(t *testing.T) {
	body := `{"field": "test"}`

	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqBody[MockDtoV](ctx)
		assert.NoError(t, err)
		assert.Equal(t, dto.Field, "test")
	}

	MockTestHandler(t, "POST", "/mock", "/mock", body, mockHandler, nil)
}

func TestReqBody_Validation_Payload_Err(t *testing.T) {
	body := `{"field": "t"}`

	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqBody[MockPayload](ctx)
		assert.Equal(t, err.Error(), "field must be at least 2 characters")
		assert.Equal(t, dto.Field, "t")
	}

	MockTestHandler(t, "POST", "/mock", "/mock", body, mockHandler, nil)
}

func TestReqBody_Validation_Dto_Err(t *testing.T) {
	body := `{"field": "t"}`

	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqBody[MockDto](ctx)
		assert.Equal(t, err.Error(), "field must be at least 2 characters")
		assert.Equal(t, dto.Field, "t")
	}

	MockTestHandler(t, "POST", "/mock", "/mock", body, mockHandler, nil)
}

func TestReqBody_Validation_DtoV_Err(t *testing.T) {
	body := `{"field": "t"}`

	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqBody[MockDtoV](ctx)
		assert.Equal(t, err.Error(), "field must be at least 2 characters")
		assert.Equal(t, dto.Field, "t")
	}

	MockTestHandler(t, "POST", "/mock", "/mock", body, mockHandler, nil)
}

func TestReqBody_Error(t *testing.T) {
	body := `{"wrong": "test"}`

	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqBody[MockDto](ctx)
		assert.NotNil(t, dto)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "field is required")
	}

	MockTestHandler(t, "POST", "/mock", "/mock", body, mockHandler, nil)
}

func TestReqQuery(t *testing.T) {
	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqQuery[MockDto](ctx)
		assert.NoError(t, err)
		assert.Equal(t, dto.Field, "test")
	}

	MockTestHandler(t, "GET", "/mock", "/mock?field=test", "", mockHandler, nil)
}

func TestReqQuery_Error(t *testing.T) {
	mockHandler := func(ctx *gin.Context) {
		dto, err := ReqQuery[MockDto](ctx)
		assert.NotNil(t, dto)
		assert.Error(t, err)
		assert.Equal(t, err.Error(), "field is required")
	}

	MockTestHandler(t, "GET", "/mock", "/mock?wrong=test", "", mockHandler, nil)
}
