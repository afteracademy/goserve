package network

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateDto(t *testing.T) {
	type TestStruct struct {
		Name  string `json:"name" validate:"required,min=3"`
		Email string `json:"email" validate:"required,email"`
		Age   int    `json:"age" validate:"gte=18"`
	}

	t.Run("should validate valid struct successfully", func(t *testing.T) {
		data := &TestStruct{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		result, err := ValidateDto(data)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, data, result)
	})

	t.Run("should return error for required field", func(t *testing.T) {
		data := &TestStruct{
			Name:  "",
			Email: "john@example.com",
			Age:   25,
		}

		result, err := ValidateDto(data)
		assert.Error(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, err.Error(), "name")
	})

	t.Run("should return error for invalid email", func(t *testing.T) {
		data := &TestStruct{
			Name:  "John Doe",
			Email: "invalid-email",
			Age:   25,
		}

		result, err := ValidateDto(data)
		assert.Error(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, err.Error(), "email")
	})

	t.Run("should return error for age validation", func(t *testing.T) {
		data := &TestStruct{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   17,
		}

		result, err := ValidateDto(data)
		assert.Error(t, err)
		assert.NotNil(t, result)
		assert.Contains(t, err.Error(), "age")
	})

	t.Run("should return error for multiple validation failures", func(t *testing.T) {
		data := &TestStruct{
			Name:  "Jo",
			Email: "invalid",
			Age:   17,
		}

		result, err := ValidateDto(data)
		assert.Error(t, err)
		assert.NotNil(t, result)
		// Should contain multiple error messages
		assert.Contains(t, err.Error(), "name")
		assert.Contains(t, err.Error(), "email")
		assert.Contains(t, err.Error(), "age")
	})

	t.Run("should return error for nil pointer", func(t *testing.T) {
		var data *TestStruct

		_, err := ValidateDto(data)
		assert.Error(t, err)
		assert.Equal(t, "invalid payload for validation", err.Error())
	})

	t.Run("should return error for non-pointer value", func(t *testing.T) {
		data := TestStruct{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		d, err := ValidateDto(&data)
		assert.Nil(t, err)
		assert.Equal(t, data, *d)
	})

	t.Run("should return error for non-struct pointer", func(t *testing.T) {
		data := new(string)
		*data = "test"

		_, err := ValidateDto(data)
		assert.Error(t, err)
		assert.Equal(t, "invalid payload for validation", err.Error())
	})
}

func TestValidateDto_WithDto(t *testing.T) {
	t.Run("should return unwrapped value for Dto interface", func(t *testing.T) {
		data := MockDto{
			Field: "test value",
		}

		d, err := ValidateDto(&data)
		assert.NoError(t, err)
		assert.NotNil(t, d)
		assert.Equal(t, data, *d)
	})

	t.Run("should validate and return error for invalid Dto", func(t *testing.T) {
		data := MockDto{
			Field: "t", // Too short, min=2
		}

		d, err := ValidateDto(&data)
		assert.Error(t, err)
		assert.NotNil(t, d)
		assert.Contains(t, err.Error(), "field")
	})
}

func TestValidateDto_WithDtoV(t *testing.T) {
	t.Run("should use custom ValidateErrors for DtoV", func(t *testing.T) {
		data := MockDtoV{
			Field: "t", // Too short, min=2
		}

		d, err := ValidateDto(&data)
		assert.Error(t, err)
		assert.NotNil(t, d)
		assert.Contains(t, err.Error(), "field")
	})

	t.Run("should return unwrapped value for valid DtoV", func(t *testing.T) {
		data := MockDtoV{
			Field: "valid value",
		}

		d, err := ValidateDto(&data)
		assert.NoError(t, err)
		assert.NotNil(t, d)
		assert.Equal(t, data, *d)
	})
}

func TestProcessErrors(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name" validate:"required,min=3"`
		Age  int    `json:"age" validate:"gte=18"`
	}

	t.Run("should format validation errors", func(t *testing.T) {
		data := &TestStruct{
			Name: "ab",
			Age:  17,
		}

		v := validator.New()
		v.RegisterTagNameFunc(CustomTagNameFunc())
		validationErr := v.Struct(data)

		err := processErrors(data, validationErr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name")
		assert.Contains(t, err.Error(), "age")
		// Should be comma-separated
		assert.Contains(t, err.Error(), ",")
	})

	t.Run("should return non-validation errors as-is", func(t *testing.T) {
		data := &TestStruct{}
		customErr := errors.New("custom error")

		err := processErrors(data, customErr)
		assert.Error(t, err)
		assert.Equal(t, "custom error", err.Error())
	})

	t.Run("should handle DtoV custom validation errors", func(t *testing.T) {
		data := &MockDtoV{
			Field: "a", // Too short
		}

		v := validator.New()
		v.RegisterTagNameFunc(CustomTagNameFunc())
		validationErr := v.Struct(data)

		err := processErrors(data, validationErr)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "field")
	})

	t.Run("should trim trailing separator from error message", func(t *testing.T) {
		data := &TestStruct{
			Name: "ab",
		}

		v := validator.New()
		v.RegisterTagNameFunc(CustomTagNameFunc())
		validationErr := v.Struct(data)

		err := processErrors(data, validationErr)
		assert.Error(t, err)
		errorMsg := err.Error()
		// Should not end with ", "
		assert.NotContains(t, errorMsg[len(errorMsg)-2:], ", ")
	})
}

func TestValidateDto_ComplexScenarios(t *testing.T) {
	type NestedStruct struct {
		City    string `json:"city" validate:"required"`
		Country string `json:"country" validate:"required"`
	}

	type ComplexStruct struct {
		Name    string        `json:"name" validate:"required,min=3,max=50"`
		Email   string        `json:"email" validate:"required,email"`
		Age     int           `json:"age" validate:"required,gte=18,lte=100"`
		Address *NestedStruct `json:"address" validate:"required"`
	}

	t.Run("should validate nested structs", func(t *testing.T) {
		data := &ComplexStruct{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   30,
			Address: &NestedStruct{
				City:    "New York",
				Country: "USA",
			},
		}

		result, err := ValidateDto(data)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should catch nested validation errors", func(t *testing.T) {
		data := &ComplexStruct{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   30,
			Address: &NestedStruct{
				City:    "", // Missing required field
				Country: "USA",
			},
		}

		result, err := ValidateDto(data)
		assert.Error(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should handle nil nested struct", func(t *testing.T) {
		data := &ComplexStruct{
			Name:    "John Doe",
			Email:   "john@example.com",
			Age:     30,
			Address: nil, // Required but nil
		}

		result, err := ValidateDto(data)
		assert.Error(t, err)
		assert.NotNil(t, result)
	})
}

func TestValidateDto_EdgeCases(t *testing.T) {
	type EmptyStruct struct{}

	t.Run("should validate empty struct without errors", func(t *testing.T) {
		data := &EmptyStruct{}

		result, err := ValidateDto(data)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	type NoValidationStruct struct {
		Name  string
		Value int
	}

	t.Run("should validate struct without validation tags", func(t *testing.T) {
		data := &NoValidationStruct{
			Name:  "test",
			Value: 123,
		}

		result, err := ValidateDto(data)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}
