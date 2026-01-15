package author

import (
	"time"

	"github.com/afteracademy/goserve/api/blog"
	"github.com/afteracademy/goserve/api/blog/dto"
	"github.com/afteracademy/goserve/api/blog/model"
	userModel "github.com/afteracademy/goserve/api/user/model"
	coredto "github.com/afteracademy/goserve/arch/dto"
	"github.com/afteracademy/goserve/arch/mongo"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Service interface {
	CreateBlog(createBlogDto *dto.CreateBlog, author *userModel.User) (*dto.PrivateBlog, error)
	UpdateBlog(updateBlogDto *dto.UpdateBlog, author *userModel.User) (*dto.PrivateBlog, error)
	DeactivateBlog(blogId primitive.ObjectID, author *userModel.User) error
	BlogSubmission(blogId primitive.ObjectID, author *userModel.User, submit bool) error
	GetBlogById(id primitive.ObjectID, author *userModel.User) (*dto.PrivateBlog, error)
	GetPaginatedDrafts(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error)
	GetPaginatedPublished(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error)
	GetPaginatedSubmitted(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error)
	getPaginated(filter bson.M, p *coredto.Pagination, opts *options.FindOptions) ([]*dto.InfoBlog, error)
}

type service struct {
	network.BaseService
	blogQueryBuilder mongo.QueryBuilder[model.Blog]
	blogService      blog.Service
}

func NewService(db mongo.Database, blogService blog.Service) Service {
	return &service{
		BaseService:      network.NewBaseService(),
		blogQueryBuilder: mongo.NewQueryBuilder[model.Blog](db, model.CollectionName),
		blogService:      blogService,
	}
}

func (s *service) CreateBlog(b *dto.CreateBlog, author *userModel.User) (*dto.PrivateBlog, error) {
	b.Slug = utils.FormatEndpoint(b.Slug)

	exists := s.blogService.BlogSlugExists(b.Slug)
	if exists {
		return nil, network.NewBadRequestError("Blog with slug: "+b.Slug+" already exists", nil)
	}

	blog, err := model.NewBlog(b.Slug, b.Title, b.Description, b.DraftText, b.Tags, author)
	if err != nil {
		return nil, err
	}

	created, err := s.blogQueryBuilder.SingleQuery().InsertAndRetrieveOne(blog)
	if err != nil {
		return nil, err
	}

	return dto.NewPrivateBlog(created, author)
}

func (s *service) UpdateBlog(b *dto.UpdateBlog, author *userModel.User) (*dto.PrivateBlog, error) {
	filter := bson.M{"_id": b.ID, "author": author.ID, "status": true}
	blog, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, network.NewNotFoundError("Blog with id: "+b.ID.Hex()+" does not exists", nil)
	}

	updates := bson.M{}

	if b.Slug != nil {
		slug := utils.FormatEndpoint(*b.Slug)
		if slug != blog.Slug {
			exists := s.blogService.BlogSlugExists(slug)
			if exists {
				return nil, network.NewBadRequestError("Blog with slug: "+slug+" already exists", nil)
			}
			updates["slug"] = slug
		}
	}

	if b.Title != nil {
		updates["title"] = *b.Title
	}

	if b.Description != nil {
		updates["description"] = *b.Description
	}

	if b.DraftText != nil {
		updates["draftText"] = *b.DraftText
	}

	if b.Tags != nil {
		updates["tags"] = *b.Tags
	}

	if b.ImgURL != nil {
		updates["imgUrl"] = *b.ImgURL
	}

	updates["updatedBy"] = author.ID
	updates["updatedAt"] = time.Now()

	set := bson.M{"$set": updates}
	_, err = s.blogQueryBuilder.SingleQuery().UpdateOne(filter, set)
	if err != nil {
		return nil, err
	}

	return s.GetBlogById(blog.ID, author)
}

func (s *service) DeactivateBlog(blogId primitive.ObjectID, author *userModel.User) error {
	filter := bson.M{"_id": blogId, "author": author.ID, "status": true}
	update := bson.M{"$set": bson.M{"status": false, "updatedBy": author.ID, "updatedAt": time.Now()}}
	result, err := s.blogQueryBuilder.SingleQuery().UpdateOne(filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return network.NewNotFoundError("blog not found", nil)
	}

	return nil
}

func (s *service) BlogSubmission(blogId primitive.ObjectID, author *userModel.User, submit bool) error {
	filter := bson.M{"_id": blogId, "author": author.ID, "status": true}
	update := bson.M{"$set": bson.M{"submitted": submit, "updatedBy": author.ID, "updatedAt": time.Now()}}
	result, err := s.blogQueryBuilder.SingleQuery().UpdateOne(filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return network.NewNotFoundError("blog not found", nil)
	}

	return nil
}

func (s *service) GetBlogById(id primitive.ObjectID, author *userModel.User) (*dto.PrivateBlog, error) {
	filter := bson.M{"_id": id, "author": author.ID, "status": true}

	blog, err := s.blogQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	return dto.NewPrivateBlog(blog, author)
}

func (s *service) GetPaginatedDrafts(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"author": author.ID, "status": true, "drafted": true}
	return s.getPaginated(filter, p, nil)
}

func (s *service) GetPaginatedPublished(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"author": author.ID, "status": true, "published": true}
	return s.getPaginated(filter, p, nil)
}

func (s *service) GetPaginatedSubmitted(author *userModel.User, p *coredto.Pagination) ([]*dto.InfoBlog, error) {
	filter := bson.M{"author": author.ID, "status": true, "submitted": true}
	return s.getPaginated(filter, p, nil)
}

func (s *service) getPaginated(filter bson.M, p *coredto.Pagination, opts *options.FindOptions) ([]*dto.InfoBlog, error) {
	blogs, err := s.blogQueryBuilder.SingleQuery().FindPaginated(filter, p.Page, p.Limit, opts)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.InfoBlog, len(blogs))

	for i, b := range blogs {
		d, err := dto.NewInfoBlog(b)
		if err != nil {
			return nil, err
		}
		dtos[i] = d
	}

	return dtos, nil
}
