package middleware

import (
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
)

type errorCatcher struct {
}

func NewErrorCatcher() network.RootMiddleware {
	return &errorCatcher{}
}

func (m *errorCatcher) Attach(engine *gin.Engine) {
	engine.Use(m.Handler)
}

func (m *errorCatcher) Handler(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				network.SendInternalServerError(ctx, err.Error(), err)
			} else {
				network.SendInternalServerError(ctx, "something went wrong", nil)
			}
			ctx.Abort()
		}
	}()
	ctx.Next()
}
