package blog

import (
	coredto "github.com/afteracademy/goserve/arch/dto"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/gin-gonic/gin"
)

type controller struct {
	network.BaseController
	service Service
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		BaseController: network.NewBaseController("/blog", authMFunc, authorizeMFunc),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/id/:id", c.getBlogByIdHandler)
	group.GET("/slug/:slug", c.getBlogBySlugHandler)
}

func (c *controller) getBlogByIdHandler(ctx *gin.Context) {
	uuidParam, err := network.ReqParams(ctx, coredto.EmptyUUID())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blog, err := c.service.GetBlogDtoCacheById(uuidParam.ID)
	if err == nil {
		c.Send(ctx).SuccessDataResponse("success", blog)
		return
	}

	blog, err = c.service.GetPublisedBlogById(uuidParam.ID)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blog)
	c.service.SetBlogDtoCacheById(blog)
}

func (c *controller) getBlogBySlugHandler(ctx *gin.Context) {
	slug, err := network.ReqParams(ctx, coredto.EmptySlug())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blog, err := c.service.GetBlogDtoCacheBySlug(slug.Slug)
	if err == nil {
		c.Send(ctx).SuccessDataResponse("success", blog)
		return
	}

	blog, err = c.service.GetPublishedBlogBySlug(slug.Slug)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blog)
	c.service.SetBlogDtoCacheBySlug(blog)
}
