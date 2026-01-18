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

	sender := NewResponseSender()
	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	sender.Send(ctx).MixedError(nil)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"code":"%s"`, failue_code))
	assert.Contains(t, resp.Body.String(), `"message":"something went wrong"`)
}

func TestSend_MixedError_Err(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sender := NewResponseSender()
	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	err := errors.New("test error")
	sender.Send(ctx).MixedError(err)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"code":"%s"`, failue_code))
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"message":"%s"`, err.Error()))
}

func TestSend_MixedError_UnauthorizedError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sender := NewResponseSender()
	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	err := NewUnauthorizedError("test message", nil)
	sender.Send(ctx).MixedError(err)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"code":"%s"`, failue_code))
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"message":"%s"`, "test message"))
}

func TestSend_SuccessMsgResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sender := NewResponseSender()
	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	sender.Send(ctx).SuccessMsgResponse("test message")

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"code":"%s"`, success_code))
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"message":"%s"`, "test message"))
}

func TestSend_SuccessDataResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	sender := NewResponseSender()
	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	data := struct {
		Field string `json:"field"`
	}{
		Field: "test data",
	}

	sender.Send(ctx).SuccessDataResponse("test message", data)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"code":"%s"`, success_code))
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"message":"%s"`, "test message"))
	assert.Contains(t, resp.Body.String(), fmt.Sprintf(`"data":%s`, `{"field":"test data"}`))
}
