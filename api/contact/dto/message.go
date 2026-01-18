package dto

import (
	"time"

	"github.com/afteracademy/goserve/arch/utility"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `json:"id" binding:"required"`
	Type      string             `json:"type" binding:"required"`
	Msg       string             `json:"msg" binding:"required"`
	CreatedAt time.Time          `json:"createdAt" binding:"required"`
}

func EmptyMessage() *Message {
	return &Message{}
}

func (d *Message) GetValue() *Message {
	return d
}

func (d *Message) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	return utility.FormatValidationErrors(errs), nil
}
