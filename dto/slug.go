package coredto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func EmptySlug() *Slug {
	return &Slug{}
}

type Slug struct {
	Slug string `uri:"slug" validate:"required,min=3,max=200"`
}

func (d *Slug) GetValue() *Slug {
	return d
}

func (b *Slug) ValidateErrors(errs validator.ValidationErrors) ([]string, error) {
	var msgs []string
	for _, err := range errs {
		switch err.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", err.Field()))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param()))
		case "max":
			msgs = append(msgs, fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", err.Field()))
		}
	}
	return msgs, nil
}
