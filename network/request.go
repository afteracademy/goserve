package network

import (
	"github.com/gin-gonic/gin"
)

// ShouldBindJSON in gin internally used go-playground/validator i.e. why we have error with validaiton info
func ReqBody[T any](ctx *gin.Context, payload T) (T, error) {
	if err := ctx.ShouldBindJSON(payload); err != nil {
		e := processErrors(payload, err)
		return payload, e
	}

	return ValidateDto(ctx, payload)
}

func ReqQuery[T any](ctx *gin.Context, payload T) (T, error) {
	if err := ctx.ShouldBindQuery(payload); err != nil {
		e := processErrors(payload, err)
		return payload, e
	}

	return ValidateDto(ctx, payload)
}

func ReqParams[T any](ctx *gin.Context, payload T) (T, error) {
	if err := ctx.ShouldBindUri(payload); err != nil {
		e := processErrors(payload, err)
		return payload, e
	}

	return ValidateDto(ctx, payload)
}

func ReqHeaders[T any](ctx *gin.Context, payload T) (T, error) {
	if err := ctx.ShouldBindHeader(payload); err != nil {
		e := processErrors(payload, err)
		return payload, e
	}

	return ValidateDto(ctx, payload)
}
