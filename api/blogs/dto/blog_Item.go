package dto

import (
	"time"

	"github.com/afteracademy/goserve/api/blog/model"
	"github.com/afteracademy/goserve/arch/utility"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogItem struct {
	ID          primitive.ObjectID `json:"id" binding:"required" validate:"required"`
	Title       string             `json:"title" validate:"required,min=3,max=500"`
	Description string             `json:"description" validate:"required,min=3,max=2000"`
	Slug        string             `json:"slug" validate:"required,min=3,max=200"`
	ImgURL      *string            `json:"imgUrl,omitempty" validate:"omitempty,uri,max=200"`
	Score       float64            `json:"score," validate:"required,min=0,max=1"`
	Tags        []string           `json:"tags" validate:"required,dive,uppercase"`
	PublishedAt *time.Time         `json:"publishedAt,omitempty"`
}

func NewBlogItem(blog *model.Blog) (*BlogItem, error) {
	return utility.MapTo[BlogItem](blog)
}

func EmptyBlogItem() *BlogItem {
	return &BlogItem{}
}

func (d *BlogItem) GetValue() *BlogItem {
	return d
}

func (b *BlogItem) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
