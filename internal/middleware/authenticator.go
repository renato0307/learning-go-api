package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/renato0307/learning-go-api/internal/apierror"
)

func Authenticator() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authentication")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, apierror.New("Not authorized"))
			return
		}

		validationGetter(&jwt.Token{Raw: authHeader})

	}
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func validationGetter(token *jwt.Token) (interface{}, error) {

	clientId := "your client id"
	// AWS Cognito public keys are available at address:
	//https://cognito-idp.{region}.amazonaws.com/{userPoolId}/.well-known/jwks.json
	publicKeysURL := "https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_b9CEhR3UR/.well-known/jwks.json"
	iss := "your iss"

	resp, err := http.Get(publicKeysURL)

	if err != nil {
		return token, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return token, err
	}

	// Verify 'iss' claim
	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
	if !checkIss {
		return token, fmt.Errorf("invalid iss")
	}

	// Verify audience and make sure it matches client id
	aud, _ := token.Claims.(jwt.MapClaims)["client_id"].(string)
	if aud != clientId {
		return token, fmt.Errorf("invalid audience")
	}

	// Validates time based claims "exp, iat, nbf"
	err = token.Claims.(jwt.MapClaims).Valid()
	if err != nil {
		return token, errors.New("token expired")
	}

	checkKid := false
	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			checkKid = true
		}
	}

	if !checkKid {
		return token, errors.New("invalid kid")
	}
	pk, err := getPublicKey(token, jwks)
	if err != nil {
		return nil, errors.New("something went wrong")
	}
	return pk, nil
}

// getPublicKey ... function to return the public key
func getPublicKey(token *jwt.Token, jwks Jwks) (*rsa.PublicKey, error) {
	var pk *rsa.PublicKey

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			// decode the base64 bytes for n
			nb, err := base64.RawURLEncoding.DecodeString(jwks.Keys[k].N)
			if err != nil {
				log.Fatal(err)
			}
			e := 0
			// The default exponent is usually 65537, so just compare the
			// base64 for [1,0,1] or [0,1,0,1]
			if jwks.Keys[k].E == "AQAB" || jwks.Keys[k].E == "AAEAAQ" {
				e = 65537
			} else {
				// need to decode "e" as a big-endian int
				log.Fatal("need to deocde e:", jwks.Keys[k].E)
			}
			pk = &rsa.PublicKey{
				N: new(big.Int).SetBytes(nb),
				E: e,
			}
			return pk, nil
		}
	}
	return pk, errors.New("could not find match")
}
