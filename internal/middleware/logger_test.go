package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestStructuredLogger(t *testing.T) {
	// arrange - create a new logger writing to a buffer
	buffer := new(bytes.Buffer)
	var memLogger = zerolog.New(buffer).With().Timestamp().Logger()

	// arrange - init gin to use the structured logger middleware
	r := gin.New()
	r.Use(StructuredLogger(&memLogger))
	r.Use(gin.Recovery())

	// arrange - set the routes
	r.GET("/example", func(c *gin.Context) {})
	r.GET("/force500", func(c *gin.Context) { panic("forced panic") })

	// act & assert
	PerformRequest(r, "GET", "/example?a=100")
	assert.Contains(t, buffer.String(), "200")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/example")
	assert.Contains(t, buffer.String(), "a=100")

	buffer.Reset()
	PerformRequest(r, "GET", "/notfound")
	assert.Contains(t, buffer.String(), "404")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/notfound")

	buffer.Reset()
	PerformRequest(r, "GET", "/force500")
	assert.Contains(t, buffer.String(), "500")
	assert.Contains(t, buffer.String(), "GET")
	assert.Contains(t, buffer.String(), "/force500")
	assert.Contains(t, buffer.String(), "error")

}
