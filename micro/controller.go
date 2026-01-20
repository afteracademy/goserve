package micro

import (
	"github.com/afteracademy/goserve/v2/network"
)

type controller struct {
	network.Controller
}

func NewController(
	basePath string,
	authProvider network.AuthenticationProvider,
	authorizeProvider network.AuthorizationProvider,
) Controller {
	return &controller{
		Controller: network.NewController(basePath, authProvider, authorizeProvider),
	}
}

func (c *controller) MountNats(group NatsGroup) {
	// To be implemented by specific controllers
}
