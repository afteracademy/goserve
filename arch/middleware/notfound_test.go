package middleware

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/afteracademy/goserve/arch/network"
)

func TestNotFoundMiddleware(t *testing.T) {
	rr := network.MockTestRootMiddlewareWithUrl(t, "/test", "/wrong", NewNotFound(), network.MockSuccessMsgHandler("success"))

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"url not found"`)
}
