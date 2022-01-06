package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequiredEnv(t *testing.T) {
	// arrange
	varName := "ENV_VAR_NAME"
	varValue := "ENV_VAR_VALUE"

	os.Setenv(varName, varValue)

	// act
	value := getRequiredEnv(varName)

	// assert
	assert.Equal(t, varValue, value)
}

func TestGetRequiredEnvWithMissingEnvironment(t *testing.T) {
	// act & assert
	assert.Panics(t, func() {
		getRequiredEnv("MISSING_ENV_VAR_NAME")
	})
}

func TestNewAuthenticator(t *testing.T) {
}
