package editor

import (
	"github.com/gin-gonic/gin"
	userModel "github.com/afteracademy/goserve/api/user/model"
	coredto "github.com/afteracademy/goserve/arch/dto"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/common"
)

type controller struct {
	network.BaseController
	common.ContextPayload
	service Service
}

func NewController(
	authMFunc network.AuthenticationProvider,
	authorizeMFunc network.AuthorizationProvider,
	service Service,
) network.Controller {
	return &controller{
		BaseController: network.NewBaseController("/blog/editor", authMFunc, authorizeMFunc),
		ContextPayload: common.NewContextPayload(),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.Use(c.Authentication(), c.Authorization(string(userModel.RoleCodeEditor)))
	group.GET("/id/:id", c.getBlogHandler)
	group.PUT("/publish/id/:id", c.publishBlogHandler)
	group.PUT("/unpublish/id/:id", c.unpublishBlogHandler)
	group.GET("/submitted", c.getSubmittedBlogsHandler)
	group.GET("/published", c.getPublishedBlogsHandler)
}

func (c *controller) getBlogHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blog, err := c.service.GetBlogById(mongoId.ID)
	if err != nil {
		c.Send(ctx).NotFoundError(mongoId.Id+" not found", err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blog)
}

func (c *controller) publishBlogHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	err = c.service.BlogPublication(mongoId.ID, user, true)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessMsgResponse("blog published successfully")
}

func (c *controller) unpublishBlogHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	user := c.MustGetUser(ctx)

	err = c.service.BlogPublication(mongoId.ID, user, false)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessMsgResponse("blog unpublished successfully")
}

func (c *controller) getSubmittedBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery(ctx, coredto.EmptyPagination())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blog, err := c.service.GetPaginatedSubmitted(pagination)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blog)
}

func (c *controller) getPublishedBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery(ctx, coredto.EmptyPagination())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blogs, err := c.service.GetPaginatedPublished(pagination)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blogs)
}
