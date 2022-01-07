package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
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
	// arrange
	issuer := "https://cognito-idp.$AWS_REGION.amazonaws.com/$POOL_ID"
	sampleJwks := `
	{
		"keys": [{
			"kid": "1234example=",
			"alg": "RS256",
			"kty": "RSA",
			"e": "AQAB",
			"n": "1234567890",
			"use": "sig"
		}, {
			"kid": "5678example=",
			"alg": "RS256",
			"kty": "RSA",
			"e": "AQAB",
			"n": "987654321",
			"use": "sig"
		}]
	}`

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, sampleJwks)
	}))

	os.Setenv(AUTH_TOKEN_ISS, issuer)
	os.Setenv(AUTH_JWKS_LOCATION, svr.URL)

	// act
	config := newAuthenticatorConfig()

	// assert
	assert.Equal(t, []byte(sampleJwks), config.KeySetJSON)
	assert.Equal(t, issuer, config.Issuer)
}
