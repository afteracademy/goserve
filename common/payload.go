package common

import (
	"errors"

	"github.com/gin-gonic/gin"
	authModel "github.com/afteracademy/goserve/api/auth/model"
	userModel "github.com/afteracademy/goserve/api/user/model"
)

const (
	payloadApiKey   string = "apikey"
	payloadUser     string = "user"
	payloadKeystore string = "keystore"
)

type ContextPayload interface {
	SetApiKey(ctx *gin.Context, value *authModel.ApiKey)
	MustGetApiKey(ctx *gin.Context) *authModel.ApiKey
	SetUser(ctx *gin.Context, value *userModel.User)
	MustGetUser(ctx *gin.Context) *userModel.User
	SetKeystore(ctx *gin.Context, value *authModel.Keystore)
	MustGetKeystore(ctx *gin.Context) *authModel.Keystore
}

type payload struct{}

func NewContextPayload() ContextPayload {
	return &payload{}
}

func (u *payload) SetApiKey(ctx *gin.Context, value *authModel.ApiKey) {
	ctx.Set(payloadApiKey, value)
}

func (u *payload) MustGetApiKey(ctx *gin.Context) *authModel.ApiKey {
	value, ok := ctx.MustGet(payloadApiKey).(*authModel.ApiKey)
	if !ok {
		panic(errors.New(payloadApiKey + " missing in context"))
	}
	return value
}

func (u *payload) SetUser(ctx *gin.Context, value *userModel.User) {
	ctx.Set(payloadUser, value)
}

func (u *payload) MustGetUser(ctx *gin.Context) *userModel.User {
	value, ok := ctx.MustGet(payloadUser).(*userModel.User)
	if !ok {
		panic(errors.New(payloadUser + " missing for context"))
	}
	return value
}

func (u *payload) SetKeystore(ctx *gin.Context, value *authModel.Keystore) {
	ctx.Set(payloadKeystore, value)
}

func (u *payload) MustGetKeystore(ctx *gin.Context) *authModel.Keystore {
	value, ok := ctx.MustGet(payloadKeystore).(*authModel.Keystore)
	if !ok {
		panic(errors.New(payloadKeystore + " missing for context"))
	}
	return value
}
