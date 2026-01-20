package coredto

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestEmptyPagination(t *testing.T) {
	t.Run("should create empty Pagination", func(t *testing.T) {
		pagination := EmptyPagination()

		assert.NotNil(t, pagination)
		assert.Equal(t, int64(0), pagination.Page)
		assert.Equal(t, int64(0), pagination.Limit)
	})
}

func TestPagination_GetValue(t *testing.T) {
	tests := []struct {
		name  string
		page  int64
		limit int64
	}{
		{
			name:  "should return pagination with valid values",
			page:  1,
			limit: 10,
		},
		{
			name:  "should return pagination with max values",
			page:  1000,
			limit: 1000,
		},
		{
			name:  "should return pagination with min values",
			page:  1,
			limit: 1,
		},
		{
			name:  "should return pagination even with zero values",
			page:  0,
			limit: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pagination := &Pagination{
				Page:  tt.page,
				Limit: tt.limit,
			}
			result := pagination.GetValue()

			assert.NotNil(t, result)
			assert.Equal(t, tt.page, result.Page)
			assert.Equal(t, tt.limit, result.Limit)
			assert.Equal(t, pagination, result)
		})
	}
}

func TestPagination_ValidateErrors(t *testing.T) {
	pagination := &Pagination{}
	validate := validator.New()

	tests := []struct {
		name          string
		setupErrors   func() validator.ValidationErrors
		expectedMsgs  []string
		expectedError bool
	}{
		{
			name: "should handle required validation error for Page",
			setupErrors: func() validator.ValidationErrors {
				pagination.Page = 0
				pagination.Limit = 10
				return validate.Struct(pagination).(validator.ValidationErrors)
			},
			expectedMsgs:  []string{"Page is required"},
			expectedError: false,
		},
		{
			name: "should handle required validation error for Limit",
			setupErrors: func() validator.ValidationErrors {
				pagination.Page = 10
				pagination.Limit = 0
				return validate.Struct(pagination).(validator.ValidationErrors)
			},
			expectedMsgs:  []string{"Limit is required"},
			expectedError: false,
		},
		{
			name: "should handle multiple required validation errors",
			setupErrors: func() validator.ValidationErrors {
				pagination.Page = 0
				pagination.Limit = 0
				return validate.Struct(pagination).(validator.ValidationErrors)
			},
			expectedMsgs:  []string{"Page is required", "Limit is required"},
			expectedError: false,
		},
		{
			name: "should handle both required and min validation errors",
			setupErrors: func() validator.ValidationErrors {
				pagination.Page = 0
				pagination.Limit = 10
				return validate.Struct(pagination).(validator.ValidationErrors)
			},
			expectedMsgs:  []string{"Page is required"},
			expectedError: false,
		},
		{
			name: "should handle max validation error",
			setupErrors: func() validator.ValidationErrors {
				pagination.Page = 1001
				pagination.Limit = 10
				return validate.Struct(pagination).(validator.ValidationErrors)
			},
			expectedMsgs:  []string{"Page must be max1000"},
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
			msgs, err := pagination.ValidateErrors(errs)

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

func TestPagination_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name      string
		page      int64
		limit     int64
		wantError bool
	}{
		{
			name:      "should validate with min values",
			page:      1,
			limit:     1,
			wantError: false,
		},
		{
			name:      "should validate with max values",
			page:      1000,
			limit:     1000,
			wantError: false,
		},
		{
			name:      "should validate with typical values",
			page:      5,
			limit:     20,
			wantError: false,
		},
		{
			name:      "should fail validation for page below min",
			page:      0,
			limit:     10,
			wantError: true,
		},
		{
			name:      "should fail validation for page above max",
			page:      1001,
			limit:     10,
			wantError: true,
		},
		{
			name:      "should fail validation for limit below min",
			page:      10,
			limit:     0,
			wantError: true,
		},
		{
			name:      "should fail validation for limit above max",
			page:      10,
			limit:     1001,
			wantError: true,
		},
		{
			name:      "should fail validation for negative page",
			page:      -1,
			limit:     10,
			wantError: true,
		},
		{
			name:      "should fail validation for negative limit",
			page:      10,
			limit:     -1,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pagination := &Pagination{
				Page:  tt.page,
				Limit: tt.limit,
			}

			err := validate.Struct(pagination)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPagination_Integration(t *testing.T) {
	t.Run("should validate and return valid pagination", func(t *testing.T) {
		pagination := EmptyPagination()
		pagination.Page = 1
		pagination.Limit = 10

		validate := validator.New()
		err := validate.Struct(pagination)
		assert.NoError(t, err)

		result := pagination.GetValue()
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.Page)
		assert.Equal(t, int64(10), result.Limit)
	})

	t.Run("should fail validation for invalid pagination", func(t *testing.T) {
		pagination := EmptyPagination()
		pagination.Page = 0
		pagination.Limit = 0

		validate := validator.New()
		err := validate.Struct(pagination)
		assert.Error(t, err)

		validationErrs := err.(validator.ValidationErrors)
		msgs, validationErr := pagination.ValidateErrors(validationErrs)
		assert.NoError(t, validationErr)
		assert.Greater(t, len(msgs), 0)
	})
}
