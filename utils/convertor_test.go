package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertUint16(t *testing.T) {
	tests := []struct {
		input    string
		expected uint16
	}{
		{"65535", 65535},
		{"0", 0},
		{"12345", 12345},
		{"invalid", 0},
	}

	for _, tt := range tests {
		result := ConvertUint16(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}

func TestConvertUint8(t *testing.T) {
	tests := []struct {
		input    string
		expected uint8
	}{
		{"255", 255},
		{"0", 0},
		{"123", 123},
		{"invalid", 0},
	}

	for _, tt := range tests {
		result := ConvertUint8(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}

type TestStruct struct {
	Field1 string
	Field2 int
}

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Bearer token123", "token123"},
		{"Bearer ", ""},
		{"Invalid token123", ""},
		{"BearerBearer token123", ""},
	}

	for _, tt := range tests {
		result := ExtractBearerToken(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}

func TestFormatEndpoint(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"endpoint /path?query", "endpoint-pathquery"},
		{"no changes", "nochanges"},
		{"spaces only", "spacesonly"},
		{"slashes/only", "slashes-only"},
	}

	for _, tt := range tests {
		result := FormatEndpoint(tt.input)
		assert.Equal(t, tt.expected, result)
	}
}
