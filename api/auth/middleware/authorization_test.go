package middleware

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	userModel "github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAuthorizationProvider_NoRole(t *testing.T) {
	mockAuthProvider := new(network.MockAuthenticationProvider)
	mockAuthProvider.On("Middleware").Return(gin.HandlerFunc(func(ctx *gin.Context) {
		ctx.Next()
	}))

	rr := network.MockTestAuthorizationProvider(t, "",
		mockAuthProvider,
		NewAuthorizationProvider(),
		network.MockSuccessMsgHandler("success"),
	)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"permission denied: role missing"`)
}

func TestAuthorizationProvider_WrongRole(t *testing.T) {
	role := &userModel.Role{ID: primitive.NewObjectID(), Code: "CORRECT_ROLE"}
	user := &userModel.User{ID: primitive.NewObjectID(), RoleDocs: []*userModel.Role{role}}

	mockAuthProvider := new(network.MockAuthenticationProvider)
	mockAuthProvider.On("Middleware").Return(gin.HandlerFunc(func(ctx *gin.Context) {
		payload := common.NewContextPayload()
		payload.SetUser(ctx, user)
		ctx.Next()
	}))

	rr := network.MockTestAuthorizationProvider(t, "WRONG_ROLE",
		mockAuthProvider,
		NewAuthorizationProvider(),
		network.MockSuccessMsgHandler("success"),
	)

	assert.Equal(t, http.StatusForbidden, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"permission denied: does not have suffient role"`)
}

func TestAuthorizationProvider_Success(t *testing.T) {

	role := &userModel.Role{ID: primitive.NewObjectID(), Code: "CORRECT_ROLE"}
	user := &userModel.User{ID: primitive.NewObjectID(), RoleDocs: []*userModel.Role{role}}

	mockAuthProvider := new(network.MockAuthenticationProvider)
	mockAuthProvider.On("Middleware").Return(gin.HandlerFunc(func(ctx *gin.Context) {
		payload := common.NewContextPayload()
		payload.SetUser(ctx, user)
		ctx.Next()
	}))

	rr := network.MockTestAuthorizationProvider(t, "CORRECT_ROLE",
		mockAuthProvider,
		NewAuthorizationProvider(),
		network.MockSuccessMsgHandler("success"),
	)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), `"message":"success"`)
}
