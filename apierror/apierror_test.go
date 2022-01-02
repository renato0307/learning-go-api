package apierror

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	errorMessage := "this is a fake error message"
	err := New(errorMessage)

	assert.Equal(t, errorMessage, err.Message)
}

func TestAssertIsValid(t *testing.T) {
	// arrange
	err := New("this is a fake error message")
	data, _ := json.Marshal(err)
	tt := testing.T{}

	// act
	AssertIsValid(&tt, data)

	// assert
	assert.False(t, tt.Failed())
}

func TestAssertIsValidWithInvalidJson(t *testing.T) {
	// arrange
	data := []byte("this is an invalid json")
	tt := testing.T{}

	// act
	AssertIsValid(&tt, data)

	// assert
	assert.True(t, tt.Failed())
}
