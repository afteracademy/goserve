package dto

import (
	"github.com/afteracademy/goserve/utils"
	"github.com/go-playground/validator/v10"
)

type BlogCreate struct {
	Title       string   `json:"title" validate:"required,min=3,max=500"`
	Description string   `json:"description" validate:"required,min=3,max=2000"`
	DraftText   string   `json:"draftText" validate:"required,max=50000"`
	Slug        string   `json:"slug" validate:"required,min=3,max=200"`
	ImgURL      string   `json:"imgUrl" validate:"required,uri,max=200"`
	Tags        []string `json:"tags" validate:"required,min=1,dive,uppercase"`
}

func EmptyBlogCreate() *BlogCreate {
	return &BlogCreate{}
}

func (d *BlogCreate) GetValue() *BlogCreate {
	return d
}

func (d *BlogCreate) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utils.FormatValidationErrors(errs), nil
}
