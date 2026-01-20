package middleware

import (
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type notFound struct {
}

func NewNotFound() network.RootMiddleware {
	return &notFound{}
}

func (m *notFound) Attach(engine *gin.Engine) {
	engine.NoRoute(m.Handler)
}

func (m *notFound) Handler(ctx *gin.Context) {
	network.SendNotFoundError(ctx, "url not found", nil)
}
