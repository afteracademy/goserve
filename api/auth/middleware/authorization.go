package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/common"
)

type authorizationProvider struct {
	network.ResponseSender
	common.ContextPayload
}

func NewAuthorizationProvider() network.AuthorizationProvider {
	return &authorizationProvider{
		ResponseSender: network.NewResponseSender(),
		ContextPayload: common.NewContextPayload(),
	}
}

func (m *authorizationProvider) Middleware(roleNames ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(roleNames) == 0 {
			m.Send(ctx).ForbiddenError("permission denied: role missing", nil)
			return
		}

		user := m.MustGetUser(ctx)

		hasRole := false
		for _, code := range roleNames {
			for _, role := range user.RoleDocs {
				if role.Code == model.RoleCode(code) {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			m.Send(ctx).ForbiddenError("permission denied: does not have suffient role", nil)
			return
		}

		ctx.Next()
	}
}
