package dto

import (
	"github.com/afteracademy/goserve/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type BlogUpdate struct {
	ID          uuid.UUID `json:"id" binding:"required" validate:"required"`
	Title       *string   `json:"title" validate:"omitempty,min=3,max=500"`
	Description *string   `json:"description" validate:"omitempty,min=3,max=2000"`
	DraftText   *string   `json:"draftText" validate:"omitempty,max=50000"`
	Slug        *string   `json:"slug" validate:"omitempty,min=3,max=200"`
	ImgURL      *string   `json:"imgUrl" validate:"omitempty,uri,max=200"`
	Tags        *[]string `json:"tags" validate:"omitempty,min=1,dive,uppercase"`
}

func EmptyBlogUpdate() *BlogUpdate {
	return &BlogUpdate{}
}

func (d *BlogUpdate) GetValue() *BlogUpdate {
	return d
}

func (b *BlogUpdate) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utils.FormatValidationErrors(errs), nil
}
