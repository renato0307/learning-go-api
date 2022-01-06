package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/renato0307/learning-go-api/internal/apierror"
	"github.com/rs/zerolog/log"
)

func Authenticator(jwks []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authentication")

		if authHeader == "" {
			log.Debug().Msg("JWT not found")
			c.JSON(http.StatusUnauthorized, apierror.New("Not authorized"))
			return
		}

		token, err := validateToken(jwks, authHeader)
		if err != nil {
			log.Debug().Err(err).Msg("JWT not valid")
			c.JSON(http.StatusUnauthorized, apierror.New("Not authorized"))
			return
		}

		_ = token
	}
}

func validateToken(jwks []byte, tokenString string) (*jwt.Token, error) {
	keySet, err := jwk.Parse(jwks)
	if err != nil {
		return nil, fmt.Errorf("failed to parse keyset: %s", err)
	}

	token, err := jwt.Parse([]byte(tokenString), jwt.WithKeySet(keySet))
	if err != nil {
		return nil, fmt.Errorf("failed to parse payload: %s", err)
	}

	return &token, nil
}
