package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Field1 string
	Field2 int
}

func TestCopyAndSetField(t *testing.T) {
	input := &TestStruct{Field1: "original", Field2: 123}
	newValue := "new value"
	result := CopyAndSetField(input, "Field1", &newValue)

	assert.Equal(t, "new value", result.Field1)
	assert.Equal(t, 123, result.Field2)
}

func TestMapTo(t *testing.T) {
	type From struct {
		Field1 string
		Field2 int
	}

	type To struct {
		Field1 string
		Field2 int
	}

	from := &From{Field1: "value", Field2: 42}
	to, err := MapTo[To](from)

	assert.NoError(t, err)
	assert.Equal(t, from.Field1, to.Field1)
	assert.Equal(t, from.Field2, to.Field2)
}
