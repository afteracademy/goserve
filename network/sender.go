package network

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendCustomResponse[T any](ctx *gin.Context, rescode ResCode, status int, message string, data *T) {
	sendResponse(ctx, NewCustomResponse(rescode, status, message, data))
}

func SendSuccessMsgResponse(ctx *gin.Context, message string) {
	sendResponse(ctx, NewSuccessMsgResponse(message))
}

func SendSuccessDataResponse[T any](ctx *gin.Context, message string, data *T) {
	sendResponse(ctx, NewSuccessDataResponse(message, data))
}

func SendBadRequestError(ctx *gin.Context, message string, err error) {
	sendError(ctx, NewBadRequestError(message, err))
}

func SendForbiddenError(ctx *gin.Context, message string, err error) {
	sendError(ctx, NewForbiddenError(message, err))
}

func SendUnauthorizedError(ctx *gin.Context, message string, err error) {
	sendError(ctx, NewUnauthorizedError(message, err))
}

func SendNotFoundError(ctx *gin.Context, message string, err error) {
	sendError(ctx, NewNotFoundError(message, err))
}

func SendInternalServerError(ctx *gin.Context, message string, err error) {
	sendError(ctx, NewInternalServerError(message, err))
}

func SendMixedError(ctx *gin.Context, err error) {
	if err == nil {
		SendInternalServerError(ctx, "something went wrong", err)
		return
	}

	var apiError ApiError
	if errors.As(err, &apiError) {
		sendError(ctx, apiError)
		return
	}

	SendInternalServerError(ctx, err.Error(), err)
}

func sendResponse[T any](ctx *gin.Context, response Response[T]) {
	data := response.GetData()
	if data != nil {
		_, err := ValidateDto(data)
		if err != nil {
			res := NewInternalServerErrorResponse(err.Error())
			ctx.JSON(int(res.GetStatus()), res)
			ctx.Abort()
			return
		}
	}

	ctx.JSON(int(response.GetStatus()), response)
	// this is needed since gin calls ctx.Next() inside the resposne handeling
	// ref: https://github.com/gin-gonic/gin/issues/2221
	ctx.Abort()
}

func sendError(ctx *gin.Context, err ApiError) {
	var debug = gin.Mode() != gin.ReleaseMode
	var res Response[any]

	switch err.GetCode() {
	case http.StatusBadRequest:
		res = NewBadRequestResponse(err.GetMessage())
	case http.StatusForbidden:
		res = NewForbiddenResponse(err.GetMessage())
	case http.StatusUnauthorized:
		res = NewUnauthorizedResponse(err.GetMessage())
	case http.StatusNotFound:
		res = NewNotFoundResponse(err.GetMessage())
	case http.StatusInternalServerError:
		if debug {
			res = NewInternalServerErrorResponse(err.Unwrap().Error())
		}
	default:
		if debug {
			res = NewInternalServerErrorResponse(err.Unwrap().Error())
		}
	}

	if res == nil {
		res = NewInternalServerErrorResponse("An unexpected error occurred. Please try again later.")
	}

	sendResponse(ctx, res)
}
