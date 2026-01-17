package dto

import (
	"github.com/afteracademy/goserve/api/blog/model"
	"github.com/afteracademy/goserve/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BlogInfo struct {
	ID          uuid.UUID `json:"id" binding:"required" validate:"required"`
	Title       string    `json:"title" validate:"required,min=3,max=500"`
	Description string    `json:"description" validate:"required,min=3,max=2000"`
	Slug        string    `json:"slug" validate:"required,min=3,max=200"`
	ImgURL      *string   `json:"imgUrl,omitempty" validate:"omitempty,uri,max=200"`
	Score       float64   `json:"score," validate:"required,min=0,max=1"`
	Tags        []string  `json:"tags" validate:"required,dive,uppercase"`
}

func NewBlogInfo(blog *model.Blog) (*BlogInfo, error) {
	return utils.MapTo[BlogInfo](blog)
}

func EmptyBlogInfo() *BlogInfo {
	return &BlogInfo{}
}

func (d *BlogInfo) GetValue() *BlogInfo {
	return d
}

func (d *BlogInfo) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utils.FormatValidationErrors(errs), nil
}
