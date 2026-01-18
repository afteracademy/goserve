package micro

import (
	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go/micro"
)

type NatsGroup = micro.Group
type NatsHandlerFunc = micro.HandlerFunc
type NatsRequest = micro.Request

type SendMessage interface {
	Message(data any)
	Error(err error)
}

type MessageSender interface {
	SendNats(req NatsRequest) SendMessage
}

type BaseController interface {
	MessageSender
	network.BaseController
}

type Controller interface {
	BaseController
	MountNats(group NatsGroup)
	MountRoutes(group *gin.RouterGroup)
}

type Router interface {
	network.BaseRouter
	NatsClient() NatsClient
	LoadControllers(controllers []Controller)
}

type Module[T any] interface {
	network.BaseModule[T]
	Controllers() []Controller
}
