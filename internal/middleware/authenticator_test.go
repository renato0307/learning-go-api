package middleware

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/renato0307/learning-go-api/internal/apierror"
	"github.com/renato0307/learning-go-api/internal/apitesting"
)

func TestAuthenticatorNoAuthHeader(t *testing.T) {

	// arrange - init gin to use the structured logger middleware
	r := gin.New()
	r.Use(Authenticator())
	r.Use(gin.Recovery())

	// arrange - set the routes
	r.GET("/example", func(c *gin.Context) {})

	// act
	w := apitesting.PerformRequest(r, "GET", "/example?a=100")

	// assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	apierror.AssertIsValid(t, w.Body.Bytes())
}

func TestAuthenticatorInvalidJwt(t *testing.T) {

	// arrange - init gin to use the structured logger middleware
	r := gin.New()
	r.Use(Authenticator())
	r.Use(gin.Recovery())

	// arrange - set the routes
	r.GET("/example", func(c *gin.Context) {})

	// arrange - headers
	header := http.Header{}
	header.Add("Authentication", "InvalidJWT")

	// act
	w := apitesting.PerformRequestWithHeader(r, "GET", "/example?a=100", header)

	// assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	apierror.AssertIsValid(t, w.Body.Bytes())
}
