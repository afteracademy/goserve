package dto

import (
	"time"

	"github.com/afteracademy/goserve/api/blog/model"
	"github.com/afteracademy/goserve/api/user/dto"
	"github.com/afteracademy/goserve/utils"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogPublic struct {
	ID          primitive.ObjectID `json:"id" binding:"required" validate:"required"`
	Title       string             `json:"title" validate:"required,min=3,max=500"`
	Description string             `json:"description" validate:"required,min=3,max=2000"`
	Text        string             `json:"text" validate:"required,max=50000"`
	Slug        string             `json:"slug" validate:"required,min=3,max=200"`
	Author      *dto.UserPublic    `json:"author,omitempty" validate:"required,omitempty"`
	ImgURL      *string            `json:"imgUrl,omitempty" validate:"omitempty,uri,max=200"`
	Score       *float64           `json:"score,omitempty" validate:"omitempty,min=0,max=1"`
	Tags        *[]string          `json:"tags,omitempty" validate:"omitempty,dive,uppercase"`
	PublishedAt *time.Time         `json:"publishedAt,omitempty"`
}

func EmptyBlogPublic() *BlogPublic {
	return &BlogPublic{}
}

func NewBlogPublic(blog *model.Blog, author *dto.UserPublic) (*BlogPublic, error) {
	b, err := utils.MapTo[BlogPublic](blog)
	if err != nil {
		return nil, err
	}

	b.Author = author

	return b, err
}

func (d *BlogPublic) GetValue() *BlogPublic {
	return d
}

func (b *BlogPublic) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utils.FormatValidationErrors(errs), nil
}
