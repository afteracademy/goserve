package network

import "github.com/gin-gonic/gin"

type baseMiddleware struct {
	ResponseSender
}

func NewBaseMiddleware() BaseMiddleware {
	return &baseMiddleware{
		ResponseSender: NewResponseSender(),
	}
}

func (m *baseMiddleware) Debug() bool {
	return gin.Mode() == gin.DebugMode
}
