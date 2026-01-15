package model

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/arch/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongod "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const CollectionName = "blogs"

type Blog struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title" validate:"required,max=500"`
	Description string             `bson:"description" validate:"required,max=2000"`
	Text        *string            `bson:"text,omitempty"`
	DraftText   string             `bson:"draftText" validate:"required"`
	Tags        []string           `bson:"tags" validate:"required"`
	Author      primitive.ObjectID `bson:"author" validate:"required"`
	ImgURL      *string            `bson:"imgUrl,omitempty"`
	Slug        string             `bson:"slug" validate:"required,min=3,max=200"`
	Score       float64            `bson:"score" validate:"min=0,max=1"`
	Submitted   bool               `bson:"submitted"`
	Drafted     bool               `bson:"drafted"`
	Published   bool               `bson:"published"`
	Status      bool               `bson:"status"`
	PublishedAt *time.Time         `bson:"publishedAt,omitempty"`
	CreatedBy   primitive.ObjectID `bson:"createdBy" validate:"required"`
	UpdatedBy   primitive.ObjectID `bson:"updatedBy" validate:"required"`
	CreatedAt   time.Time          `bson:"createdAt" validate:"required"`
	UpdatedAt   time.Time          `bson:"updatedAt" validate:"required"`
}

func NewBlog(slug, title, description, draftText string, tags []string, author *model.User) (*Blog, error) {
	now := time.Now()
	b := Blog{
		Title:       title,
		Description: description,
		DraftText:   draftText,
		Tags:        tags,
		Author:      author.ID,
		Slug:        slug,
		Score:       0.01,
		Submitted:   false,
		Drafted:     true,
		Published:   false,
		Status:      true,
		CreatedBy:   author.ID,
		UpdatedBy:   author.ID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := b.Validate(); err != nil {
		return nil, err
	}
	return &b, nil
}

func (blog *Blog) GetValue() *Blog {
	return blog
}

func (blog *Blog) Validate() error {
	validate := validator.New()
	return validate.Struct(blog)
}

func (*Blog) EnsureIndexes(db mongo.Database) {
	indexes := []mongod.IndexModel{
		{
			Keys: bson.D{
				{Key: "slug", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "title", Value: "text"}, {Key: "description", Value: "text"}},
			Options: options.Index().SetWeights(bson.M{
				"title":       3,
				"description": 1,
			}),
		},
		{Keys: bson.D{{Key: "_id", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "published", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "_id", Value: 1}, {Key: "published", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "slug", Value: 1}, {Key: "published", Value: 1}, {Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "tags", Value: 1}, {Key: "published", Value: 1}, {Key: "status", Value: 1}}},
	}

	mongo.NewQueryBuilder[Blog](db, CollectionName).Query(context.Background()).CreateIndexes(indexes)
}
