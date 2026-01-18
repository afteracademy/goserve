package micro

import (
	"github.com/afteracademy/goserve/v2/network"
)

type baseController struct {
	MessageSender
	network.BaseController
}

func NewBaseController(basePath string, authProvider network.AuthenticationProvider, authorizeProvider network.AuthorizationProvider) BaseController {
	return &baseController{
		MessageSender:  NewMessageSender(),
		BaseController: network.NewBaseController(basePath, authProvider, authorizeProvider),
	}
}
