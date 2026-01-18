package network

import (
	"github.com/gin-gonic/gin"
)

type baseController struct {
	ResponseSender
	basePath          string
	authProvider      AuthenticationProvider
	authorizeProvider AuthorizationProvider
}

func NewBaseController(basePath string, authProvider AuthenticationProvider, authorizeProvider AuthorizationProvider) BaseController {
	return &baseController{
		ResponseSender:    NewResponseSender(),
		basePath:          basePath,
		authProvider:      authProvider,
		authorizeProvider: authorizeProvider,
	}
}

func (c *baseController) Path() string {
	return c.basePath
}

func (c *baseController) Authentication() gin.HandlerFunc {
	return c.authProvider.Middleware()
}

func (c *baseController) Authorization(role string) gin.HandlerFunc {
	return c.authorizeProvider.Middleware(role)
}
