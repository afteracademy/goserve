package coredto

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestEmptyMongoId(t *testing.T) {
	t.Run("should create empty MongoId", func(t *testing.T) {
		mongoId := EmptyMongoId()

		assert.NotNil(t, mongoId)
		assert.Equal(t, "", mongoId.Id)
		assert.Equal(t, primitive.NilObjectID, mongoId.ID)
	})
}

func TestMongoId_GetValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
		checkID  bool
	}{
		{
			name:     "should parse valid 24-character hex string",
			input:    "507f1f77bcf86cd799439011",
			wantErr:  false,
			checkID:  true,
		},
		{
			name:     "should handle invalid mongo id gracefully",
			input:    "invalid",
			wantErr:  true,
			checkID:  false,
		},
		{
			name:     "should handle empty string",
			input:    "",
			wantErr:  true,
			checkID:  false,
		},
		{
			name:     "should handle wrong length string",
			input:    "507f1f77bcf86cd799",
			wantErr:  true,
			checkID:  false,
		},
		{
			name:     "should handle another valid mongo id",
			input:    "65a1b2c3d4e5f6a7b8c9d0e1",
			wantErr:  false,
			checkID:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mongoId := &MongoId{Id: tt.input}
			result := mongoId.GetValue()

			assert.NotNil(t, result)
			assert.Equal(t, tt.input, result.Id)

			if tt.checkID {
				assert.NotEqual(t, primitive.NilObjectID, result.ID)
				assert.Equal(t, tt.input, result.ID.Hex())
			} else {
				assert.Equal(t, primitive.NilObjectID, result.ID)
			}
		})
	}
}

func TestMongoId_ValidateErrors(t *testing.T) {
	mongoId := &MongoId{}
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
				mongoId.Id = ""
				return validate.Struct(mongoId).(validator.ValidationErrors)
			},
			expectedMsgs:  []string{"Id is required"},
			expectedError: false,
		},
		{
			name: "should handle len validation error",
			setupErrors: func() validator.ValidationErrors {
				mongoId.Id = "short"
				return validate.Struct(mongoId).(validator.ValidationErrors)
			},
			expectedMsgs:  []string{"Id must be of length 24"},
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
			msgs, err := mongoId.ValidateErrors(errs)

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

func TestMongoId_Integration(t *testing.T) {
	t.Run("should validate and parse valid mongo id", func(t *testing.T) {
		mongoId := EmptyMongoId()
		mongoId.Id = "507f1f77bcf86cd799439011"

		validate := validator.New()
		err := validate.Struct(mongoId)
		assert.NoError(t, err)

		result := mongoId.GetValue()
		assert.NotNil(t, result)
		assert.NotEqual(t, primitive.NilObjectID, result.ID)
		assert.Equal(t, "507f1f77bcf86cd799439011", result.ID.Hex())
	})

	t.Run("should fail validation for invalid mongo id", func(t *testing.T) {
		mongoId := EmptyMongoId()
		mongoId.Id = "invalid"

		validate := validator.New()
		err := validate.Struct(mongoId)
		assert.Error(t, err)

		validationErrs := err.(validator.ValidationErrors)
		msgs, validationErr := mongoId.ValidateErrors(validationErrs)
		assert.NoError(t, validationErr)
		assert.Greater(t, len(msgs), 0)
	})
}
