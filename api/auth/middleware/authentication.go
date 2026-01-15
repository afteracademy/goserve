package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/afteracademy/goserve/api/auth"
	"github.com/afteracademy/goserve/api/user"
	"github.com/afteracademy/goserve/arch/mongo"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/common"
	"github.com/afteracademy/goserve/utils"
)

type authenticationProvider struct {
	network.ResponseSender
	common.ContextPayload
	authService auth.Service
	userService user.Service
}

func NewAuthenticationProvider(authService auth.Service, userService user.Service) network.AuthenticationProvider {
	return &authenticationProvider{
		ResponseSender: network.NewResponseSender(),
		ContextPayload: common.NewContextPayload(),
		authService:    authService,
		userService:    userService,
	}
}

func (m *authenticationProvider) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(network.AuthorizationHeader)
		if len(authHeader) == 0 {
			m.Send(ctx).UnauthorizedError("permission denied: missing Authorization", nil)
			return
		}

		token := utils.ExtractBearerToken(authHeader)
		if token == "" {
			m.Send(ctx).UnauthorizedError("permission denied: invalid Authorization", nil)
			return
		}

		claims, err := m.authService.VerifyToken(token)
		if err != nil {
			m.Send(ctx).UnauthorizedError(err.Error(), err)
			return
		}

		valid := m.authService.ValidateClaims(claims)
		if !valid {
			m.Send(ctx).UnauthorizedError("permission denied: invalid claims", nil)
			return
		}

		userId, err := mongo.NewObjectID(claims.Subject)
		if err != nil {
			m.Send(ctx).UnauthorizedError("permission denied: invalid claims subject", nil)
			return
		}

		user, err := m.userService.FindUserById(userId)
		if err != nil {
			m.Send(ctx).UnauthorizedError("permission denied: claims subject does not exists", err)
			return
		}

		keystore, err := m.authService.FindKeystore(user, claims.ID)
		if err != nil || keystore == nil {
			m.Send(ctx).UnauthorizedError("permission denied: invalid access token", err)
			return
		}

		m.SetUser(ctx, user)
		m.SetKeystore(ctx, keystore)

		ctx.Next()
	}
}
