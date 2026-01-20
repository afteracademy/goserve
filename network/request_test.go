package network

import (
	"testing"

	coredto "github.com/afteracademy/goserve/v2/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func TestReqParam_uuid_Dto(t *testing.T) {
	id := uuid.New()

	mockHandler := func(ctx *gin.Context) {
		param, err := ReqParams[coredto.UUID](ctx)
		assert.NoError(t, err)
		assert.Equal(t, param.ID.String(), id.String())
	}

	MockTestHandler(t, "POST", "/mock/id/:id", "/mock/id/"+id.String(), "", mockHandler, nil)
}

func TestReqParam_mongo_id_Dto(t *testing.T) {
	id := primitive.NewObjectID()

	mockHandler := func(ctx *gin.Context) {
		param, err := ReqParams[coredto.MongoId](ctx)
		assert.NoError(t, err)
		assert.Equal(t, param.ID.Hex(), id.Hex())
	}

	MockTestHandler(t, "POST", "/mock/id/:id", "/mock/id/"+id.Hex(), "", mockHandler, nil)
}
