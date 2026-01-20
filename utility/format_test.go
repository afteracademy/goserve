package utility

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var validate = validator.New()

type User struct {
	Email    string `validate:"email"`
	Username string `validate:"min=3,max=10"`
	Age      int    `validate:"gte=18"`
	Role     string `validate:"oneof=admin user"`
}

func TestFormatValidationErrors_NonValidationError(t *testing.T) {
	err := errors.New("some random error")

	msgs := FormatValidationErrors(err)

	assert.Equal(t, []string{"invalid request"}, msgs)
}

func TestFormatValidationErrors_SingleFieldSingleError(t *testing.T) {
	u := User{
		Email: "",
	}

	err := validate.Struct(u)
	assert.Error(t, err)

	msgs := FormatValidationErrors(err)

	assert.ElementsMatch(t, []string{
		"email is not a valid email",
		"username must be at least 3 characters",
		"age must be greater than or equal to 18",
		"role must be one of admin user",
	}, msgs)
}

func TestFormatValidationErrors_EmailCorrectRestError(t *testing.T) {
	u := User{
		Email: "goserve@afteracademy.com",
	}

	err := validate.Struct(u)
	assert.Error(t, err)

	msgs := FormatValidationErrors(err)

	assert.ElementsMatch(t, []string{
		"username must be at least 3 characters",
		"age must be greater than or equal to 18",
		"role must be one of admin user",
	}, msgs)
}

func TestFormatValidationErrors_MultipleErrors(t *testing.T) {
	u := User{
		Email:    "invalid-email",
		Username: "ab",
		Age:      16,
		Role:     "guest",
	}

	err := validate.Struct(u)
	assert.Error(t, err)

	msgs := FormatValidationErrors(err)

	assert.ElementsMatch(t, []string{
		"email is not a valid email",
		"username must be at least 3 characters",
		"age must be greater than or equal to 18",
		"role must be one of admin user",
	}, msgs)
}

func TestFormatValidationErrors_FieldNameLowercased(t *testing.T) {
	type Test struct {
		MyField string `validate:"required"`
	}

	err := validate.Struct(Test{})
	assert.Error(t, err)

	msgs := FormatValidationErrors(err)

	assert.Equal(t, []string{"myfield is required"}, msgs)
}

func TestFormatValidationErrors_MultipleSameTag(t *testing.T) {
	type Test struct {
		A string `validate:"required"`
		B string `validate:"required"`
	}

	err := validate.Struct(Test{})
	assert.Error(t, err)

	msgs := FormatValidationErrors(err)

	assert.ElementsMatch(t, []string{
		"a is required",
		"b is required",
	}, msgs)
}
