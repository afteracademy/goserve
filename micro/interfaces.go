package micro

import (
	"github.com/afteracademy/goserve/v2/network"
	"github.com/nats-io/nats.go/micro"
)

type NatsGroup = micro.Group
type NatsHandlerFunc = micro.HandlerFunc
type NatsRequest = micro.Request

type Controller interface {
	network.Controller
	MountNats(group NatsGroup)
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
