package middleware

import (
	"errors"
	"net/http"
	"testing"

	"github.com/afteracademy/goserve/network"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestErrorCatcherMiddleware(t *testing.T) {
	mockHandler := func(ctx *gin.Context) {
		panic(errors.New("panic test"))
	}

	rr := network.MockTestRootMiddleware(t, NewErrorCatcher(), mockHandler, nil)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"panic test"`)
}

func TestErrorCatcherMiddleware_NonError(t *testing.T) {
	mockHandler := func(ctx *gin.Context) {
		panic(1)
	}

	rr := network.MockTestRootMiddleware(t, NewErrorCatcher(), mockHandler, nil)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"something went wrong"`)
}
