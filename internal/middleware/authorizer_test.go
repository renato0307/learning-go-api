package middleware

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/renato0307/learning-go-api/internal/apitesting"
	"github.com/stretchr/testify/assert"
)

func TestAuthorizer(t *testing.T) {

	testCases := []struct {
		Scopes     string
		URL        string
		RequestURL string
		Purpose    string
		StatusCode int
	}{
		{
			Scopes:     "https://learninggolang.com/programming-jwtdebugger",
			URL:        "/v1/programming-jwtdebugger",
			RequestURL: "/v1/programming-jwtdebugger?abcd",
			Purpose:    "scopes match URL",
			StatusCode: http.StatusOK,
		},
		{
			Scopes:     "https://learninggolang.com/programming-uuid",
			URL:        "/v1/programming-jwtdebugger",
			RequestURL: "/v1/programming-jwtdebugger",
			Purpose:    "scopes do not match URL",
			StatusCode: http.StatusForbidden,
		},
		{
			Scopes:     "",
			URL:        "/v1/programming-jwtdebugger",
			RequestURL: "/v1/programming-jwtdebugger",
			Purpose:    "no scopes defined",
			StatusCode: http.StatusForbidden,
		},
		{
			Scopes:     "https://learninggolang.com/programming-jwtdebugger https://learninggolang.com/programming-uuid",
			URL:        "/v1/programming-jwtdebugger",
			RequestURL: "/v1/programming-jwtdebugger",
			Purpose:    "list of scopes",
			StatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		// arrange - init gin to use the middlewares
		r := gin.New()
		r.Use(func(c *gin.Context) { // fake Authenticator
			c.Set(ScopeKey, tc.Scopes)
		})
		r.Use(Authorizer())
		r.Use(gin.Recovery())

		// arrange - set the routes
		r.POST(tc.URL, func(c *gin.Context) {})

		// act
		w := apitesting.PerformRequest(r, "POST", tc.RequestURL)

		// assert
		assert.Equal(t, tc.StatusCode, w.Code, tc.Purpose)
	}
}
