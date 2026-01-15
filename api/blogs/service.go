package blogs

import (
	"time"

	"github.com/afteracademy/goserve/api/blog/model"
	"github.com/afteracademy/goserve/api/blogs/dto"
	coredto "github.com/afteracademy/goserve/arch/dto"
	"github.com/afteracademy/goserve/arch/mongo"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/arch/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	SetSimilarBlogsDtoCache(blogId primitive.ObjectID, blogs []*dto.ItemBlog) error
	GetSimilarBlogsDtoCache(blogId primitive.ObjectID) ([]*dto.ItemBlog, error)
	GetPaginatedLatestBlogs(p *coredto.Pagination) ([]*dto.ItemBlog, error)
	GetPaginatedTaggedBlogs(tag string, p *coredto.Pagination) ([]*dto.ItemBlog, error)
	GetSimilarBlogs(blogId primitive.ObjectID) ([]*dto.ItemBlog, error)
	getPublicPaginated(filter bson.M, p *coredto.Pagination) ([]*dto.ItemBlog, error)
	getPaginated(filter bson.M, p *coredto.Pagination, opts *options.FindOptions) ([]*dto.ItemBlog, error)
}

type service struct {
	network.BaseService
	blogQueryBuilder mongo.QueryBuilder[model.Blog]
	itemBlogCache    redis.Cache[dto.ItemBlog]
}

func NewService(db mongo.Database, store redis.Store) Service {
	return &service{
		BaseService:      network.NewBaseService(),
		blogQueryBuilder: mongo.NewQueryBuilder[model.Blog](db, model.CollectionName),
		itemBlogCache:    redis.NewCache[dto.ItemBlog](store),
	}
}

func (s *service) SetSimilarBlogsDtoCache(blogId primitive.ObjectID, blogs []*dto.ItemBlog) error {
	key := "similar_blogs_" + blogId.Hex()
	return s.itemBlogCache.SetJSONList(key, blogs, 6*time.Hour)
}

func (s *service) GetSimilarBlogsDtoCache(blogId primitive.ObjectID) ([]*dto.ItemBlog, error) {
	key := "similar_blogs_" + blogId.Hex()
	return s.itemBlogCache.GetJSONList(key)
}

func (s *service) GetPaginatedLatestBlogs(p *coredto.Pagination) ([]*dto.ItemBlog, error) {
	filter := bson.M{"status": true, "published": true}
	return s.getPublicPaginated(filter, p)
}

func (s *service) GetPaginatedTaggedBlogs(tag string, p *coredto.Pagination) ([]*dto.ItemBlog, error) {
	filter := bson.M{"status": true, "published": true, "tags": tag}
	return s.getPublicPaginated(filter, p)
}

func (s *service) GetSimilarBlogs(blogId primitive.ObjectID) ([]*dto.ItemBlog, error) {
	filter := bson.M{"_id": blogId, "published": true, "status": true}
	blog, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, network.NewNotFoundError("blog not found", err)
	}

	filter = bson.M{
		"$text":     bson.M{"$search": blog.Title, "$caseSensitive": false},
		"status":    true,
		"published": true,
		"_id":       bson.M{"$ne": blog.ID},
	}

	opts := options.Find()
	opts.SetProjection(bson.M{"similarity": bson.M{"$meta": "textScore"}})
	opts.SetSort(bson.D{
		{Key: "similarity", Value: bson.M{"$meta": "textScore"}},
		{Key: "updatedAt", Value: -1},
		{Key: "score", Value: -1},
	})

	pagination := &coredto.Pagination{
		Page:  1,
		Limit: 6,
	}

	return s.getPaginated(filter, pagination, opts)
}

func (s *service) getPublicPaginated(filter bson.M, p *coredto.Pagination) ([]*dto.ItemBlog, error) {
	projection := bson.D{{Key: "draftText", Value: 0}}
	opts := options.Find().SetProjection(projection)
	opts.SetSort(bson.D{{Key: "updatedAt", Value: -1}, {Key: "score", Value: -1}})
	return s.getPaginated(filter, p, opts)
}

func (s *service) getPaginated(filter bson.M, p *coredto.Pagination, opts *options.FindOptions) ([]*dto.ItemBlog, error) {
	blogs, err := s.blogQueryBuilder.SingleQuery().FindPaginated(filter, p.Page, p.Limit, opts)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.ItemBlog, len(blogs))

	for i, b := range blogs {
		d, err := dto.NewItemBlog(b)
		if err != nil {
			return nil, err
		}
		dtos[i] = d
	}

	return dtos, nil
}
