package network

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type sender struct{}

func NewResponseSender() ResponseSender {
	return &sender{}
}

func (m *sender) Debug() bool {
	return gin.Mode() != gin.ReleaseMode
}

func (m *sender) Send(ctx *gin.Context) SendResponse {
	return &send{
		debug:   m.Debug(),
		context: ctx,
	}
}

type send struct {
	debug   bool
	context *gin.Context
}

func (s *send) SuccessMsgResponse(message string) {
	s.sendResponse(NewSuccessMsgResponse(message))
}

func (s *send) SuccessDataResponse(message string, data any) {
	s.sendResponse(NewSuccessDataResponse(message, data))
}

func (s *send) BadRequestError(message string, err error) {
	s.sendError(NewBadRequestError(message, err))
}

func (s *send) ForbiddenError(message string, err error) {
	s.sendError(NewForbiddenError(message, err))
}

func (s *send) UnauthorizedError(message string, err error) {
	s.sendError(NewUnauthorizedError(message, err))
}

func (s *send) NotFoundError(message string, err error) {
	s.sendError(NewNotFoundError(message, err))
}

func (s *send) InternalServerError(message string, err error) {
	s.sendError(NewInternalServerError(message, err))
}

func (s *send) MixedError(err error) {
	if err == nil {
		s.InternalServerError("something went wrong", err)
		return
	}

	var apiError ApiError
	if errors.As(err, &apiError) {
		s.sendError(apiError)
		return
	}

	s.InternalServerError(err.Error(), err)
}

func (s *send) sendResponse(response Response) {
	s.context.JSON(int(response.GetStatus()), response)
	// this is needed since gin calls ctx.Next() inside the resposne handeling
	// ref: https://github.com/gin-gonic/gin/issues/2221
	s.context.Abort() 
}

func (s *send) sendError(err ApiError) {
	var res Response

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
		if s.debug {
			res = NewInternalServerErrorResponse(err.Unwrap().Error())
		}
	default:
		if s.debug {
			res = NewInternalServerErrorResponse(err.Unwrap().Error())
		}
	}

	if res == nil {
		res = NewInternalServerErrorResponse("An unexpected error occurred. Please try again later.")
	}

	s.sendResponse(res)
}
