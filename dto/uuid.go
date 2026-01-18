package coredto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func EmptyUUID() *UUID {
	return &UUID{}
}

type UUID struct {
	Id string    `uri:"id" binding:"required" validate:"required,uuid"`
	ID uuid.UUID `uri:"-" validate:"-"`
}

func (d *UUID) GetValue() *UUID {
	id, err := uuid.Parse(d.Id)
	if err == nil {
		d.ID = id
	}
	return d
}

func (d *UUID) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "uuid":
			msgs = append(msgs, fmt.Sprintf("%s must be a valid UUID", err.Field()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
