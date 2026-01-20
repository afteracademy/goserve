package micro

import (
	"testing"

	"github.com/afteracademy/goserve/v2/network"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseController(t *testing.T) {
	t.Run("should create a new base controller with valid parameters", func(t *testing.T) {
		basePath := "/api/v1"
		authProvider := &network.MockAuthenticationProvider{}
		authorizeProvider := &network.MockAuthorizationProvider{}

		bc := NewBaseController(basePath, authProvider, authorizeProvider)

		assert.NotNil(t, bc)
		assert.Equal(t, basePath, bc.Path())
		assert.NotNil(t, bc.(MessageSender))
		assert.NotNil(t, bc.(network.BaseController))
	})

	t.Run("should create base controller with empty path", func(t *testing.T) {
		basePath := ""
		authProvider := &network.MockAuthenticationProvider{}
		authorizeProvider := &network.MockAuthorizationProvider{}

		bc := NewBaseController(basePath, authProvider, authorizeProvider)

		assert.NotNil(t, bc)
		assert.Equal(t, basePath, bc.Path())
	})

	t.Run("should create base controller with nil providers", func(t *testing.T) {
		basePath := "/api"

		bc := NewBaseController(basePath, nil, nil)

		assert.NotNil(t, bc)
		assert.Equal(t, basePath, bc.Path())
	})
}
