package utility

import (
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		length uint64
	}{
		{16},
		{32},
		{64},
	}

	for _, tt := range tests {
		t.Run("Length"+strconv.FormatUint(tt.length, 10), func(t *testing.T) {
			result, err := GenerateRandomString(tt.length)

			// Use testify assertions
			assert.NoError(t, err)

			// Verify the length of the generated string
			expectedLength := int(tt.length) * 2 // because hex encoding doubles the size
			assert.Equal(t, expectedLength, len(result))

			// Verify the string is a valid hex string
			_, err = hex.DecodeString(result)
			assert.NoError(t, err, "expected valid hex string, got invalid hex string")
		})
	}

}
