package dto

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/afteracademy/goserve/api/blog/model"
	userModel "github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PrivateBlog struct {
	ID          primitive.ObjectID `json:"_id" binding:"required" validate:"required"`
	Title       string             `json:"title" validate:"required,min=3,max=500"`
	Description string             `json:"description" validate:"required,min=3,max=2000"`
	Text        *string            `json:"text,omitempty" validate:"omitempty,max=50000"`
	DraftText   string             `json:"draftText" validate:"required"`
	Slug        string             `json:"slug" validate:"required,min=3,max=200"`
	Author      *InfoAuthor        `json:"author,omitempty" validate:"required,omitempty"`
	ImgURL      *string            `json:"imgUrl,omitempty" validate:"omitempty,uri,max=200"`
	Score       *float64           `json:"score,omitempty" validate:"omitempty,min=0,max=1"`
	Tags        *[]string          `json:"tags,omitempty" validate:"omitempty,dive,uppercase"`
	Submitted   bool               `json:"submitted" validate:"required"`
	Drafted     bool               `json:"drafted" validate:"required"`
	Published   bool               `json:"published" validate:"required"`
	PublishedAt *time.Time         `json:"publishedAt,omitempty"`
	CreatedAt   time.Time          `json:"createdAt" validate:"required"`
	UpdatedAt   time.Time          `json:"updatedAt" validate:"required"`
}

func EmptyInfoPrivateBlog() *PrivateBlog {
	return &PrivateBlog{}
}

func NewPrivateBlog(blog *model.Blog, author *userModel.User) (*PrivateBlog, error) {
	b, err := utils.MapTo[PrivateBlog](blog)
	if err != nil {
		return nil, err
	}

	b.Author, err = utils.MapTo[InfoAuthor](author)
	if err != nil {
		return nil, err
	}

	return b, err
}

func (d *PrivateBlog) GetValue() *PrivateBlog {
	return d
}

func (b *PrivateBlog) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param()))
		case "url":
			msgs = append(msgs, fmt.Sprintf("%s must be a valid URL", err.Field()))
		case "uri":
			msgs = append(msgs, fmt.Sprintf("%s must be a valid URI", err.Field()))
		case "uppercase":
			msgs = append(msgs, fmt.Sprintf("%s must be uppercase", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
