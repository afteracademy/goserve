package mongo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestIsValidObjectID(t *testing.T) {
	validID := primitive.NewObjectID().Hex()
	invalidID := "invalid"

	assert.True(t, IsValidObjectID(validID))
	assert.False(t, IsValidObjectID(invalidID))
}
