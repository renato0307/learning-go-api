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

type AuthenticatorConfig struct {
	KeySetJSON []byte
	Issuer     string
}

const (
	TokenUseKey = "token_use"
	ClientIdKey = "client_id"
)

func Authenticator(ac *AuthenticatorConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authentication")

		if authHeader == "" {
			log.Debug().Msg("JWT not found")
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				apierror.New("Not authorized"))
			return
		}

		token, err := validateToken(ac, authHeader)
		if err != nil {
			log.Debug().Err(err).Msg("JWT not valid")
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				apierror.New("Not authorized"))
			return
		}

		ci, _ := token.Get(ClientIdKey)
		c.Set(ClientIdKey, ci)
	}
}

func validateToken(ac *AuthenticatorConfig, tokenString string) (jwt.Token, error) {
	keySet, err := jwk.Parse(ac.KeySetJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse keyset: %s", err)
	}

	// Step 1: Confirm the structure of the JWT
	// Step 2: Validate the JWT signature
	token, err := jwt.Parse(
		[]byte(tokenString),
		jwt.WithKeySet(keySet),
	)
	if err != nil {
		log.Debug().Err(err).Msg("error parsing the token")
		return nil, fmt.Errorf("invalid token: %s", err)
	}

	// Step 3: Verify the claims
	clientId, _ := token.Get(ClientIdKey)
	err = jwt.Validate(token,
		jwt.WithClaimValue(TokenUseKey, "access"),
		jwt.WithClaimValue(jwt.IssuerKey, ac.Issuer),
		jwt.WithRequiredClaim(ClientIdKey),
		jwt.WithRequiredClaim(jwt.SubjectKey),
		jwt.WithClaimValue(jwt.SubjectKey, clientId),
	)
	if err != nil {
		log.Debug().Err(err).Msg("error validating the token")
		return nil, fmt.Errorf("invalid token: %s", err)
	}

	return token, nil
}
