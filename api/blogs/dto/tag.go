package dto

import (
	"github.com/afteracademy/goserve/utils"
	"github.com/go-playground/validator/v10"
)

func EmptyTag() *Tag {
	return &Tag{}
}

type Tag struct {
	Tag string `uri:"tag" validate:"required,uppercase"`
}

func (d *Tag) GetValue() *Tag {
	return d
}

func (b *Tag) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utils.FormatValidationErrors(errs), nil
}
