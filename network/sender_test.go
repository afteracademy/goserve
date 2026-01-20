package network

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestSend_MixedError_Nil(t *testing.T) {
	gin.SetMode(gin.TestMode)

	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	SendMixedError(ctx, nil)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"code":"%s"`, failue_code))
	assert.Contains(t, resp.Body.String(), `"message":"something went wrong"`)
}

func TestSend_MixedError_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	err := errors.New("test error")
	SendMixedError(ctx, err)
	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"code":"%s"`, failue_code))
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"message":"%s"`, err.Error()))
}

func TestSend_MixedError_UnauthorizedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	err := NewUnauthorizedError("test message", nil)
	SendMixedError(ctx, err)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"code":"%s"`, failue_code))
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"message":"%s"`, "test message"))
}

func TestSend_SuccessMsgResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	SendSuccessMsgResponse(ctx, "test message")

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"code":"%s"`, success_code))
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"message":"%s"`, "test message"))
}

func TestSend_SuccessDataResponse_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	data := &MockPayload{
		Field: "t",
	}

	SendSuccessDataResponse(ctx, "test message", data)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"code":"%s"`, failue_code))
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"message":"%s"`, "field must be at least 2 characters"))
	assert.NotContains(t, resp.Body.String(), fmt.Sprintf(`"data":%s`, `{"field":"test data"}`))
}

func TestSend_SuccessDataResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	data := struct {
		Field string `json:"field"`
	}{
		Field: "test data",
	}

	SendSuccessDataResponse(ctx, "test message", &data)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"code":"%s"`, success_code))
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"message":"%s"`, "test message"))
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"data":%s`, `{"field":"test data"}`))
}
