package dto

import (
	"time"

	"github.com/afteracademy/goserve/api/blog/model"
	"github.com/afteracademy/goserve/api/user/dto"
	userModel "github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BlogPrivate struct {
	ID          uuid.UUID       `json:"id" binding:"required" validate:"required"`
	Title       string          `json:"title" validate:"required,min=3,max=500"`
	Description string          `json:"description" validate:"required,min=3,max=2000"`
	Text        *string         `json:"text,omitempty" validate:"omitempty,max=50000"`
	DraftText   string          `json:"draftText" validate:"required"`
	Slug        string          `json:"slug" validate:"required,min=3,max=200"`
	Author      *dto.UserPublic `json:"author,omitempty" validate:"required,omitempty"`
	ImgURL      *string         `json:"imgUrl,omitempty" validate:"omitempty,uri,max=200"`
	Score       *float64        `json:"score,omitempty" validate:"omitempty,min=0,max=1"`
	Tags        *[]string       `json:"tags,omitempty" validate:"omitempty,dive,uppercase"`
	Submitted   bool            `json:"submitted" validate:"required"`
	Drafted     bool            `json:"drafted" validate:"required"`
	Published   bool            `json:"published" validate:"required"`
	PublishedAt *time.Time      `json:"publishedAt,omitempty"`
	CreatedAt   time.Time       `json:"createdAt" validate:"required"`
	UpdatedAt   time.Time       `json:"updatedAt" validate:"required"`
}

func EmptyBlogPrivate() *BlogPrivate {
	return &BlogPrivate{}
}

func NewBlogPrivate(blog *model.Blog, author *userModel.User) (*BlogPrivate, error) {
	b, err := utils.MapTo[BlogPrivate](blog)
	if err != nil {
		return nil, err
	}

	b.Author, err = utils.MapTo[dto.UserPublic](author)
	if err != nil {
		return nil, err
	}

	return b, err
}

func (d *BlogPrivate) GetValue() *BlogPrivate {
	return d
}

func (b *BlogPrivate) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utils.FormatValidationErrors(errs), nil
}
