package network

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ApiError interface {
	GetCode() int
	GetMessage() string
	Error() string
	Unwrap() error
}

type Response interface {
	GetResCode() ResCode
	GetStatus() int
	GetMessage() string
	GetData() any
}

type SendResponse interface {
	SuccessMsgResponse(message string)
	SuccessDataResponse(message string, data any)
	BadRequestError(message string, err error)
	ForbiddenError(message string, err error)
	UnauthorizedError(message string, err error)
	NotFoundError(message string, err error)
	InternalServerError(message string, err error)
	MixedError(err error)
}

type ResponseSender interface {
	Debug() bool
	Send(ctx *gin.Context) SendResponse
}

type BaseController interface {
	ResponseSender
	Path() string
	Authentication() gin.HandlerFunc
	Authorization(role string) gin.HandlerFunc
}

type Controller interface {
	BaseController
	MountRoutes(group *gin.RouterGroup)
}

type BaseService interface {
	Context() context.Context
}

type Dto[T any] interface {
	GetValue() *T
	ValidateErrors(errs validator.ValidationErrors) ([]string, error)
}

type BaseMiddleware interface {
	ResponseSender
	Debug() bool
}

type RootMiddleware interface {
	BaseMiddleware
	Attach(engine *gin.Engine)
	Handler(ctx *gin.Context)
}

type BaseMiddlewareProvider interface {
	BaseMiddleware
}

type Param0MiddlewareProvider interface {
	BaseMiddlewareProvider
	Middleware() gin.HandlerFunc
}

type Param1MiddlewareProvider[T any] interface {
	BaseMiddlewareProvider
	Middleware(param1 T) gin.HandlerFunc
}

type Param2MiddlewareProvider[T any, V any] interface {
	BaseMiddlewareProvider
	Middleware(param1 T, param2 V) gin.HandlerFunc
}

type Param3MiddlewareProvider[T any, V any, W any] interface {
	BaseMiddlewareProvider
	Middleware(param1 T, param2 V, param3 W) gin.HandlerFunc
}

type ParamNMiddlewareProvider[T any] interface {
	BaseMiddlewareProvider
	Middleware(params ...T) gin.HandlerFunc
}

type AuthenticationProvider Param0MiddlewareProvider
type AuthorizationProvider ParamNMiddlewareProvider[string]

type BaseRouter interface {
	GetEngine() *gin.Engine
	RegisterValidationParsers(tagNameFunc validator.TagNameFunc)
	LoadRootMiddlewares(middlewares []RootMiddleware)
	Start(ip string, port uint16)
}

type Router interface {
	BaseRouter
	LoadControllers(controllers []Controller)
}

type BaseModule[T any] interface {
	GetInstance() *T
	RootMiddlewares() []RootMiddleware
	AuthenticationProvider() AuthenticationProvider
	AuthorizationProvider() AuthorizationProvider
}

type Module[T any] interface {
	BaseModule[T]
	Controllers() []Controller
}
