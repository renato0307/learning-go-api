package programming

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupGin() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	SetRouterGroup(v1)

	return r
}

func TestPostUuid(t *testing.T) {
	// arrange
	r := setupGin()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/programming/uuid", nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusOK)

	body := w.Body.String()
	assert.Len(t, body, 38)
	assert.Contains(t, body, "-")
}

func TestPostUuidWithNoHyphen(t *testing.T) {
	// arrange
	r := setupGin()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/programming/uuid?no-hyphens=true", nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusOK)

	body := w.Body.String()
	assert.Len(t, body, 34)
	assert.NotContains(t, body, "-")
}
