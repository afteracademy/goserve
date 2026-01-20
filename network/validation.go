package network

import (
	"errors"
	"strings"

	"github.com/afteracademy/goserve/v2/utility"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateDto[T any](ctx *gin.Context, payload T) (T, error) {
	v := validator.New()
	v.RegisterTagNameFunc(CustomTagNameFunc())
	if err := v.Struct(payload); err != nil {
		e := processErrors(payload, err)
		return payload, e
	}

	if dto, ok := any(payload).(Dto[T]); ok {
		return dto.GetValue(), nil
	}

	return payload, nil
}

func processErrors[T any](payload T, err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var msgs []string
		if d, ok := any(payload).(DtoV[T]); ok {
			vmsgs, e := d.ValidateErrors(validationErrors)
			if e != nil {
				return e
			}
			msgs = vmsgs
		} else {
			msgs = utility.FormatValidationErrors(err)
		}

		var msg strings.Builder
		br := ", "
		for _, m := range msgs {
			msg.WriteString(m + br)
		}
		// Remove the trailing separator
		errorMsg := msg.String()
		if len(errorMsg) > 0 {
			errorMsg = errorMsg[:len(errorMsg)-len(br)]
		}
		return errors.New(errorMsg)
	}
	return err
}
