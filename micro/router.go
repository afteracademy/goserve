package micro

import (
	"fmt"
	"strings"

	"github.com/afteracademy/goserve/v2/network"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type router struct {
	netRouter  network.Router
	natsClient NatsClient
}

func NewRouter(mode string, natsClient NatsClient) Router {
	return &router{
		netRouter:  network.NewRouter(mode),
		natsClient: natsClient,
	}
}

func (r *router) GetEngine() *gin.Engine {
	return r.netRouter.GetEngine()
}

func (r *router) NatsClient() NatsClient {
	return r.natsClient
}

func (r *router) LoadRootMiddlewares(middlewares []network.RootMiddleware) {
	r.netRouter.LoadRootMiddlewares(middlewares)
}

func (r *router) LoadControllers(controllers []Controller) {
	nc := make([]network.Controller, len(controllers))
	for i, c := range controllers {
		nc[i] = c.(network.Controller)
	}
	r.netRouter.LoadControllers(nc)

	natsClient := r.natsClient.GetInstance()

	for _, c := range controllers {
		baseSub := natsClient.Service.Info().Name
		endpoint := strings.ReplaceAll(c.Path(), "/", ".")
		if len(endpoint) > 1 {
			baseSub = fmt.Sprintf(`%s%s`, baseSub, endpoint)
		}

		ng := natsClient.Service.AddGroup(baseSub)
		c.MountNats(ng)
	}
}

func (r *router) Start(ip string, port uint16) {
	r.netRouter.Start(ip, port)
}

func (r *router) RegisterValidationParsers(tagNameFunc validator.TagNameFunc) {
	r.netRouter.RegisterValidationParsers(tagNameFunc)
}
