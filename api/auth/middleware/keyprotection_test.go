package middleware

import (
	"errors"
	"net/http"
	"testing"

	"github.com/afteracademy/goserve/api/auth"
	"github.com/afteracademy/goserve/api/auth/model"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/common"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestKeyProtectionMiddleware_NoApiKey(t *testing.T) {
	mockAuthService := new(auth.MockService)

	rr := network.MockTestRootMiddleware(
		t,
		NewKeyProtection(mockAuthService),
		network.MockSuccessMsgHandler("success"),
		nil,
	)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"permission denied: missing x-api-key header"`)
}

func TestKeyProtectionMiddleware_WrongApiKey(t *testing.T) {
	mockAuthService := new(auth.MockService)
	key := "wrong"
	mockAuthService.On("FetchApiKey", key).Return(nil, errors.New(""))

	rr := network.MockTestRootMiddleware(
		t,
		NewKeyProtection(mockAuthService),
		network.MockSuccessMsgHandler("success"),
		map[string]string{network.ApiKeyHeader: key},
	)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"permission denied: invalid x-api-key"`)
}

func TestKeyProtectionMiddleware_CorrectApiKey(t *testing.T) {
	mockAuthService := new(auth.MockService)
	key := "correct"
	mockAuthService.On("FetchApiKey", key).Return(&model.ApiKey{Key: key}, nil)

	mockHandler := func(ctx *gin.Context) {
		assert.Equal(t, common.NewContextPayload().MustGetApiKey(ctx).Key, key)
		network.NewResponseSender().Send(ctx).SuccessMsgResponse("success")
	}

	rr := network.MockTestRootMiddleware(
		t,
		NewKeyProtection(mockAuthService),
		mockHandler,
		map[string]string{network.ApiKeyHeader: key},
	)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"success"`)
}
