package middleware

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/renato0307/learning-go-api/internal/apierror"
	"github.com/renato0307/learning-go-api/internal/apitesting"
)

func TestAuthenticatorNoAuthHeader(t *testing.T) {

	// arrange - init gin to use the structured logger middleware
	r := gin.New()
	r.Use(Authenticator([]byte{}))
	r.Use(gin.Recovery())

	// arrange - set the routes
	r.GET("/example", func(c *gin.Context) {})

	// act
	w := apitesting.PerformRequest(r, "GET", "/example?a=100")

	// assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	apierror.AssertIsValid(t, w.Body.Bytes())
}

func TestAuthenticatorWithInvalidJwt(t *testing.T) {
	key, err := generateKey()
	if err != nil {
		t.Fatal(err)
	}

	jwks, err := generateKeySetInJSON(&key)
	if err != nil {
		t.Fatal(err)
	}

	// arrange - init gin to use the structured logger middleware
	r := gin.New()
	r.Use(Authenticator(jwks))
	r.Use(gin.Recovery())

	// arrange - set the routes
	r.GET("/example", func(c *gin.Context) {})

	// arrange - headers
	header := http.Header{}
	header.Add("Authentication", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")

	// act
	w := apitesting.PerformRequestWithHeader(r, "GET", "/example?a=100", header)

	// assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	apierror.AssertIsValid(t, w.Body.Bytes())
}

func TestAuthenticatorWithValidJwt(t *testing.T) {

	key, err := generateKey()
	if err != nil {
		t.Fatal(err)
	}

	token := jwt.New()
	token.Set("sub", "client_id_1234567890")
	token.Set("token_use", "access")
	token.Set("scope", "https://learning-go-api.com/all")
	token.Set("auth_time", 1641417382)
	token.Set("iss", "https://cognito-idp.eu-west-1.amazonaws.com/eu-west-1_b9CEhR3UR")
	token.Set("exp", 1641420982)
	token.Set("iat", 1641417382)
	token.Set("version", 2)
	token.Set("jti", "a6dd28cc-500e-4b49-a510-efda5195d2f4")
	token.Set("client_id", "client_id_1234567890")

	signed, err := jwt.Sign(token, jwa.RS256, key)
	if err != nil {
		t.Fatal(err)
	}

	jwks, err := generateKeySetInJSON(&key)
	if err != nil {
		t.Fatal(err)
	}

	// arrange - init gin to use the structured logger middleware
	r := gin.New()
	r.Use(Authenticator(jwks))
	r.Use(gin.Recovery())

	// arrange - set the routes
	r.GET("/example", func(c *gin.Context) {})

	// arrange - headers
	header := http.Header{}
	header.Add("Authentication", string(signed))

	// act
	w := apitesting.PerformRequestWithHeader(r, "GET", "/example?a=100", header)

	// assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func generateKey() (jwk.Key, error) {
	raw, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new RSA private key: %s\n", err)
	}

	key, err := jwk.New(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to create symmetric key: %s\n", err)
	}
	if _, ok := key.(jwk.RSAPrivateKey); !ok {
		return nil, fmt.Errorf("expected jwk.SymmetricKey, got %T\n", key)
	}

	key.Set(jwk.KeyIDKey, "mykey")

	return key, nil
}

func generateKeySetInJSON(key *jwk.Key) ([]byte, error) {
	set := jwk.NewSet()
	pubKey, _ := (*key).(jwk.RSAPrivateKey).PublicKey()
	pubKey.Set(jwk.AlgorithmKey, "RS256")
	set.Add(pubKey)

	buf, err := json.MarshalIndent(set, "", "  ")
	if err != nil {
		return []byte{}, fmt.Errorf("failed to marshal key into JSON: %s\n", err)
	}

	return buf, nil
}
