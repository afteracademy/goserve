package coredto

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEmptyUUID(t *testing.T) {
	t.Run("should create empty UUID", func(t *testing.T) {
		uuidDto := EmptyUUID()

		assert.NotNil(t, uuidDto)
		assert.Equal(t, "", uuidDto.Id)
		assert.Equal(t, uuid.Nil, uuidDto.ID)
	})
}

func TestUUID_GetValue(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		checkID bool
	}{
		{
			name:    "should parse valid UUID v4",
			input:   "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
			checkID: true,
		},
		{
			name:    "should parse valid UUID v1",
			input:   "c232ab00-9414-11ec-b909-0242ac120002",
			wantErr: false,
			checkID: true,
		},
		{
			name:    "should handle invalid uuid gracefully",
			input:   "invalid-uuid",
			wantErr: true,
			checkID: false,
		},
		{
			name:    "should handle empty string",
			input:   "",
			wantErr: true,
			checkID: false,
		},
		{
			name:    "should handle malformed uuid",
			input:   "550e8400-e29b-41d4-a716",
			wantErr: true,
			checkID: false,
		},
		{
			name:    "should handle uuid without hyphens",
			input:   "550e8400e29b41d4a716446655440000",
			wantErr: false,
			checkID: true,
		},
		{
			name:    "should handle uppercase UUID",
			input:   "550E8400-E29B-41D4-A716-446655440000",
			wantErr: false,
			checkID: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuidDto := &UUID{Id: tt.input}
			result := uuidDto.GetValue()

			assert.NotNil(t, result)
			assert.Equal(t, tt.input, result.Id)

			if tt.checkID {
				assert.NotEqual(t, uuid.Nil, result.ID)
			} else {
				assert.Equal(t, uuid.Nil, result.ID)
			}
		})
	}
}

func TestUUID_ValidateErrors(t *testing.T) {
	uuidDto := &UUID{}
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
				uuidDto.Id = ""
				return validate.Struct(uuidDto).(validator.ValidationErrors)
			},
			expectedMsgs:  []string{"Id is required"},
			expectedError: false,
		},
		{
			name: "should handle uuid validation error",
			setupErrors: func() validator.ValidationErrors {
				uuidDto.Id = "invalid-uuid"
				return validate.Struct(uuidDto).(validator.ValidationErrors)
			},
			expectedMsgs:  []string{"Id must be a valid UUID"},
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
			msgs, err := uuidDto.ValidateErrors(errs)

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

func TestUUID_Validation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name      string
		id        string
		wantError bool
	}{
		{
			name:      "should validate standard UUID v4",
			id:        "550e8400-e29b-41d4-a716-446655440000",
			wantError: false,
		},
		{
			name:      "should validate UUID v1",
			id:        "c232ab00-9414-11ec-b909-0242ac120002",
			wantError: false,
		},

		{
			name:      "should fail validation for empty string",
			id:        "",
			wantError: true,
		},
		{
			name:      "should fail validation for invalid uuid",
			id:        "invalid-uuid",
			wantError: true,
		},
		{
			name:      "should fail validation for malformed uuid",
			id:        "550e8400-e29b-41d4-a716",
			wantError: true,
		},
		{
			name:      "should fail validation for random string",
			id:        "not-a-uuid-at-all",
			wantError: true,
		},
		{
			name:      "should fail validation for partial uuid",
			id:        "550e8400-e29b",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uuidDto := &UUID{Id: tt.id}

			err := validate.Struct(uuidDto)
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUUID_Integration(t *testing.T) {
	t.Run("should validate and parse valid uuid", func(t *testing.T) {
		uuidDto := EmptyUUID()
		uuidDto.Id = "550e8400-e29b-41d4-a716-446655440000"

		validate := validator.New()
		err := validate.Struct(uuidDto)
		assert.NoError(t, err)

		result := uuidDto.GetValue()
		assert.NotNil(t, result)
		assert.NotEqual(t, uuid.Nil, result.ID)
		assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", result.ID.String())
	})

	t.Run("should fail validation for invalid uuid", func(t *testing.T) {
		uuidDto := EmptyUUID()
		uuidDto.Id = "invalid"

		validate := validator.New()
		err := validate.Struct(uuidDto)
		assert.Error(t, err)

		validationErrs := err.(validator.ValidationErrors)
		msgs, validationErr := uuidDto.ValidateErrors(validationErrs)
		assert.NoError(t, validationErr)
		assert.Greater(t, len(msgs), 0)
	})

	t.Run("should parse uuid but fail validation", func(t *testing.T) {
		uuidDto := EmptyUUID()
		uuidDto.Id = "invalid-uuid"

		// Parse will fail silently in GetValue
		result := uuidDto.GetValue()
		assert.Equal(t, uuid.Nil, result.ID)

		// But validation should catch it
		validate := validator.New()
		err := validate.Struct(uuidDto)
		assert.Error(t, err)
	})
}

func TestUUID_NewUUID(t *testing.T) {
	t.Run("should create new UUID and validate", func(t *testing.T) {
		newUUID := uuid.New()
		uuidDto := &UUID{
			Id: newUUID.String(),
		}

		validate := validator.New()
		err := validate.Struct(uuidDto)
		assert.NoError(t, err)

		result := uuidDto.GetValue()
		assert.NotNil(t, result)
		assert.Equal(t, newUUID, result.ID)
	})
}
