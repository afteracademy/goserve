package middleware

import (
	"net/http"
	"testing"

	"github.com/afteracademy/goserve/arch/network"
	"github.com/stretchr/testify/assert"
)

func TestNotFoundMiddleware(t *testing.T) {
	rr := network.MockTestRootMiddlewareWithUrl(t, "/test", "/wrong", NewNotFound(), network.MockSuccessMsgHandler("success"), nil)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"url not found"`)
}
