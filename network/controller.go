package network

import (
	"github.com/gin-gonic/gin"
)

type controller struct {
	basePath          string
	authProvider      AuthenticationProvider
	authorizeProvider AuthorizationProvider
}

func NewController(
	basePath string,
	authProvider AuthenticationProvider,
	authorizeProvider AuthorizationProvider,
) Controller {
	return &controller{
		basePath:          basePath,
		authProvider:      authProvider,
		authorizeProvider: authorizeProvider,
	}
}

func (c *controller) Path() string {
	return c.basePath
}

func (c *controller) Authentication() gin.HandlerFunc {
	return c.authProvider.Middleware()
}

func (c *controller) Authorization(role string) gin.HandlerFunc {
	return c.authorizeProvider.Middleware(role)
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	// To be implemented by specific controllers
}
