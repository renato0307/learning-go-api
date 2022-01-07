package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-api/internal/middleware"
	"github.com/renato0307/learning-go-api/pkg/finance"
	"github.com/renato0307/learning-go-api/pkg/programming"
	financelib "github.com/renato0307/learning-go-lib/finance"
	programminglib "github.com/renato0307/learning-go-lib/programming"
	"github.com/rs/zerolog/log"
)

const (
	CURRCONV_API_KEY   = "CURRCONV_API_KEY"
	AUTH_TOKEN_ISS     = "AUTH_TOKEN_ISS"
	AUTH_JWKS_LOCATION = "AUTH_JWKS_LOCATION"
)

func main() {
	// Initialize Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.DefaultStructuredLogger())
	r.Use(middleware.Authenticator(newAuthenticatorConfig()))
	r.Use(gin.Recovery())

	// Default route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, welcome to the learning-go-api",
		})
	})

	// Utility functions routes
	base := r.Group("/v1")

	p := programminglib.ProgrammingFunctions{}
	programming.SetRouterGroup(&p, base)

	useDefaultUrl := ""
	apiKey := getRequiredEnv(CURRCONV_API_KEY)
	f := financelib.NewFinanceFunctions(useDefaultUrl, apiKey)
	finance.SetRouterGroup(&f, base)

	// Start serving request
	r.Run()
}

func getRequiredEnv(key string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		panic(fmt.Sprintf("error: %s environment variable was not defined", key))
	}

	return value
}

// newAuthenticatorConfig gathers all authentication related information to set
// up the Authenticator middleware configuration.
//
// It requires the definition two environment variables:
//
// AUTH_JWKS_LOCATION: https://cognito-idp.$AWS_REGION.amazonaws.com/$POOL_ID/.well-known/jwks.json
//
// AUTH_TOKEN_ISS: https://cognito-idp.$AWS_REGION.amazonaws.com/$POOL_ID
func newAuthenticatorConfig() *middleware.AuthenticatorConfig {

	// Gets the JSON Web Key Set download URL
	jwksLocation := getRequiredEnv(AUTH_JWKS_LOCATION)
	r, err := http.Get(jwksLocation)
	if err != nil {
		msg := "cannot get the JWKS content for authentication"
		log.Error().Err(err).Msg(msg)
		panic(msg)
	}
	defer r.Body.Close()

	// Downloads the JSON Web Key Set
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "cannot read the JWKS content for authentication"
		log.Error().Err(err).Msg(msg)
		panic(msg)
	}

	// Creates the AuthenticatorConfig structure
	config := middleware.AuthenticatorConfig{
		KeySetJSON: body,
		Issuer:     getRequiredEnv(AUTH_TOKEN_ISS),
	}

	log.Debug().
		Str("auth_token_iss", config.Issuer).
		Str("auth_jwks_location", jwksLocation).
		Msg("authenticator config loaded")

	return &config
}
