package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/renato0307/learning-go-api/internal/apierror"
	"github.com/renato0307/learning-go-api/internal/apitesting"
	"github.com/stretchr/testify/assert"
)

const userPool string = "https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_xxxxxxxxxx"

func TestAuthenticatorNoAuthHeader(t *testing.T) {

	// arrange - init gin to use the structured logger middleware
	r := gin.New()
	r.Use(Authenticator(nil))
	r.Use(gin.Recovery())

	// arrange - set the routes
	r.GET("/example", func(c *gin.Context) {})

	// act
	w := apitesting.PerformRequest(r, "GET", "/example?a=100")

	// assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	apierror.AssertIsValid(t, w.Body.Bytes())
}

func TestAuthenticatorRootPathSkipsAuth(t *testing.T) {

	// arrange - init gin to use the structured logger middleware
	r := gin.New()
	r.Use(Authenticator(nil))
	r.Use(gin.Recovery())

	// arrange - set the routes
	r.GET("/", func(c *gin.Context) {})

	// act
	w := apitesting.PerformRequest(r, "GET", "/")

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthenticatorWithJWT(t *testing.T) {
	// arrange - generate key, keyset and JWT
	key := generateKey(t)
	jwks := generateKeySetInJSON(&key, t)

	// arrange - define the several test cases
	testCases := []struct {
		JWT          string
		StatusCode   int
		Purpose      string
		BodyContains string
	}{
		{
			JWT: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9." +
				"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ." +
				"SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			StatusCode:   http.StatusUnauthorized,
			Purpose:      "token with signed with another key",
			BodyContains: "Not authorized",
		},
		{
			JWT:          newJWT(key, true, false, false, t),
			StatusCode:   http.StatusUnauthorized,
			Purpose:      "invalid subject claim",
			BodyContains: "Not authorized",
		},
		{
			JWT:          newJWT(key, false, true, false, t),
			StatusCode:   http.StatusUnauthorized,
			Purpose:      "invalid expiration claim",
			BodyContains: "Not authorized",
		},
		{
			JWT:          newJWT(key, false, false, true, t),
			StatusCode:   http.StatusUnauthorized,
			Purpose:      "invalid token_use claim",
			BodyContains: "Not authorized",
		},
		{
			JWT:          newValidJWT(key, t),
			StatusCode:   http.StatusOK,
			Purpose:      "valid JWT",
			BodyContains: "",
		},
	}

	for _, tc := range testCases {

		// arrange - init gin to use the authenticator middleware
		r := gin.New()
		authConfig := AuthenticatorConfig{
			KeySetJSON: jwks,
			Issuer:     userPool,
		}
		r.Use(Authenticator(&authConfig))
		r.Use(gin.Recovery())

		// arrange - set the routes
		r.GET("/example", func(c *gin.Context) {
			cid, _ := c.Get(ClientIdKey)
			c.JSON(http.StatusOK, cid)
		})

		// arrange - headers
		header := http.Header{}
		header.Add("Authentication", tc.JWT)

		// act
		w := apitesting.PerformRequestWithHeader(r, "GET", "/example?a=100", header)

		// assert
		assert.Equal(t,
			tc.StatusCode,
			w.Code,
			fmt.Sprintf("failed %s", tc.Purpose))

		b := w.Body.String()
		assert.Contains(t, b, tc.BodyContains)
	}
}

func newValidJWT(key jwk.Key, t *testing.T) string {
	return newJWT(key, false, false, false, t)
}

func newJWT(key jwk.Key, noSub, noExp, noTokenUse bool, t *testing.T) string {
	token := jwt.New()
	if !noSub {
		token.Set("sub", "client_id_1234567890")
	}
	if !noTokenUse {
		token.Set("token_use", "access")
	}
	token.Set("scope", "https://learning-go-api.com/all")
	token.Set("auth_time", 1641417382)
	token.Set("iss", userPool)
	if !noExp {
		token.Set("exp", time.Now().Unix()+1000)
	} else {
		token.Set("exp", 1)
	}
	token.Set("iat", 1641417382)
	token.Set("version", 2)
	token.Set("jti", "a6dd28cc-500e-4b49-a510-efda5195d2f4")
	token.Set("client_id", "client_id_1234567890")

	signed, err := signJWT(token, key)
	if err != nil {
		t.Fatal(err)
	}

	return string(signed)
}

func signJWT(token jwt.Token, key jwk.Key) ([]byte, error) {
	return jwt.Sign(token, jwa.RS256, key)
}

func generateKey(t *testing.T) jwk.Key {
	raw, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate new RSA private key: %s\n", err)
	}

	key, err := jwk.New(raw)
	if err != nil {
		t.Fatalf("failed to create symmetric key: %s\n", err)
	}
	if _, ok := key.(jwk.RSAPrivateKey); !ok {
		t.Fatalf("expected jwk.SymmetricKey, got %T\n", key)
	}

	key.Set(jwk.KeyIDKey, "mykey")

	return key
}

func generateKeySetInJSON(key *jwk.Key, t *testing.T) []byte {
	set := jwk.NewSet()
	pubKey, _ := (*key).(jwk.RSAPrivateKey).PublicKey()
	pubKey.Set(jwk.AlgorithmKey, "RS256")
	set.Add(pubKey)

	buf, err := json.MarshalIndent(set, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal key into JSON: %s\n", err)
	}

	return buf
}
