package blogs

import (
	"github.com/gin-gonic/gin"
	"github.com/afteracademy/goserve/api/blogs/dto"
	coredto "github.com/afteracademy/goserve/arch/dto"
	"github.com/afteracademy/goserve/arch/network"
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
		BaseController: network.NewBaseController("/blogs", authMFunc, authorizeMFunc),
		service:        service,
	}
}

func (c *controller) MountRoutes(group *gin.RouterGroup) {
	group.GET("/latest", c.getLatestBlogsHandler)
	group.GET("/tag/:tag", c.getTaggedBlogsHandler)
	group.GET("/similar/id/:id", c.getSimilarBlogsHandler)

}

func (c *controller) getLatestBlogsHandler(ctx *gin.Context) {
	pagination, err := network.ReqQuery(ctx, coredto.EmptyPagination())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blogs, err := c.service.GetPaginatedLatestBlogs(pagination)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blogs)
}

func (c *controller) getTaggedBlogsHandler(ctx *gin.Context) {
	tag, err := network.ReqParams(ctx, dto.EmptyTag())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	pagination, err := network.ReqQuery(ctx, coredto.EmptyPagination())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blogs, err := c.service.GetPaginatedTaggedBlogs(tag.Tag, pagination)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blogs)
}

func (c *controller) getSimilarBlogsHandler(ctx *gin.Context) {
	mongoId, err := network.ReqParams(ctx, coredto.EmptyMongoId())
	if err != nil {
		c.Send(ctx).BadRequestError(err.Error(), err)
		return
	}

	blogs, err := c.service.GetSimilarBlogsDtoCache(mongoId.ID)
	if err == nil {
		c.Send(ctx).SuccessDataResponse("success", blogs)
		return
	}

	blogs, err = c.service.GetSimilarBlogs(mongoId.ID)
	if err != nil {
		c.Send(ctx).MixedError(err)
		return
	}

	c.Send(ctx).SuccessDataResponse("success", blogs)
	c.service.SetSimilarBlogsDtoCache(mongoId.ID, blogs)
}
