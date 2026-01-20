package coredto

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestEmptySlug(t *testing.T) {
	t.Run("should create empty Slug", func(t *testing.T) {
		slug := EmptySlug()

		assert.NotNil(t, slug)
		assert.Equal(t, "", slug.Slug)
	})
}

func TestSlug_GetValue(t *testing.T) {
	tests := []struct {
		name string
		slug string
	}{
		{
			name: "should return slug with valid value",
			slug: "test-slug",
		},
		{
			name: "should return slug with minimum length",
			slug: "abc",
		},
		{
			name: "should return slug with long value",
			slug: "this-is-a-very-long-slug-with-many-words-and-hyphens",
		},
		{
			name: "should return slug even with empty value",
			slug: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := &Slug{Slug: tt.slug}
			result := slug.GetValue()

			assert.NotNil(t, result)
			assert.Equal(t, tt.slug, result.Slug)
			assert.Equal(t, slug, result)
		})
	}
}

func TestSlug_ValidateErrors(t *testing.T) {
	slug := &Slug{}
	validate := validator.New()

	tests := []struct {
		name          string
		setupErrors   func() validator.ValidationErrors
		expectedMsgs  []string
		expectedError bool
	}{
		{
			name: "should handle required validation error",
			setupErrors: func() validator.ValidationErrors {
				slug.Slug = ""
				return validate.Struct(slug).(validator.ValidationErrors)
			},
			expectedMsgs:  []string{"Slug is required"},
			expectedError: false,
		},
		{
			name: "should handle min validation error",
			setupErrors: func() validator.ValidationErrors {
				slug.Slug = "ab"
				return validate.Struct(slug).(validator.ValidationErrors)
			},
			expectedMsgs:  []string{"Slug must be at least 3 characters"},
			expectedError: false,
		},
		{
			name: "should handle max validation error",
			setupErrors: func() validator.ValidationErrors {
				slug.Slug = string(make([]byte, 201))
				return validate.Struct(slug).(validator.ValidationErrors)
			},
			expectedMsgs:  []string{"Slug must be at most 200 characters"},
			expectedError: false,
		},
		{
			name: "should return empty messages for nil errors",
			setupErrors: func() validator.ValidationErrors {
				return nil
			},
			expectedMsgs:  []string{},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.setupErrors()
			msgs, err := slug.ValidateErrors(errs)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, len(tt.expectedMsgs), len(msgs))
			for i, expectedMsg := range tt.expectedMsgs {
				if i < len(msgs) {
					assert.Equal(t, expectedMsg, msgs[i])
				}
			}
		})
	}
}

func TestSlug_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name      string
		slug      string
		wantError bool
	}{
		{
			name:      "should validate slug with minimum length",
			slug:      "abc",
			wantError: false,
		},
		{
			name:      "should validate slug with typical value",
			slug:      "test-slug-example",
			wantError: false,
		},
		{
			name:      "should validate slug with maximum length",
			slug:      string(make([]byte, 200)),
			wantError: false,
		},
		{
			name:      "should fail validation for empty slug",
			slug:      "",
			wantError: true,
		},
		{
			name:      "should fail validation for slug below min length",
			slug:      "ab",
			wantError: true,
		},
		{
			name:      "should fail validation for slug above max length",
			slug:      string(make([]byte, 201)),
			wantError: true,
		},
		{
			name:      "should validate slug with numbers",
			slug:      "slug-123",
			wantError: false,
		},
		{
			name:      "should validate slug with special characters",
			slug:      "slug-with-underscores_and_hyphens-123",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug := &Slug{Slug: tt.slug}

			err := validate.Struct(slug)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSlug_Integration(t *testing.T) {
	t.Run("should validate and return valid slug", func(t *testing.T) {
		slug := EmptySlug()
		slug.Slug = "test-slug"

		validate := validator.New()
		err := validate.Struct(slug)
		assert.NoError(t, err)

		result := slug.GetValue()
		assert.NotNil(t, result)
		assert.Equal(t, "test-slug", result.Slug)
	})

	t.Run("should fail validation for invalid slug", func(t *testing.T) {
		slug := EmptySlug()
		slug.Slug = "ab"

		validate := validator.New()
		err := validate.Struct(slug)
		assert.Error(t, err)

		validationErrs := err.(validator.ValidationErrors)
		msgs, validationErr := slug.ValidateErrors(validationErrs)
		assert.NoError(t, validationErr)
		assert.Greater(t, len(msgs), 0)
	})
}
