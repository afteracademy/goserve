package dto

import (
	"github.com/afteracademy/goserve/arch/utility"
	"github.com/go-playground/validator/v10"
)

type MessageCreate struct {
	Type string `json:"type" binding:"required,min=2,max=50"`
	Msg  string `json:"msg" binding:"required,min=0,max=2000"`
}

func EmptyMessageCreate() *MessageCreate {
	return &MessageCreate{}
}

func (d *MessageCreate) GetValue() *MessageCreate {
	return d
}

func (d *MessageCreate) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
