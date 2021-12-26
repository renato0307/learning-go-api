package programming

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	programminglib "github.com/renato0307/learning-go-lib/programming"
	"github.com/stretchr/testify/assert"
)

func setupGin(mockInterface *programminglib.MockInterface) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	SetRouterGroup(mockInterface, v1)

	return r
}

func TestPostUuid(t *testing.T) {
	// arrange
	mockInterface := programminglib.MockInterface{}
	mockCall := mockInterface.On("NewUuid", false)
	mockCall.Return("1ce44be5-fe68-46f7-a153-51c1c91a4ae4")

	r := setupGin(&mockInterface)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/programming/uuid", nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusOK)

	output := postUuidOutput{}
	err := json.Unmarshal(w.Body.Bytes(), &output)

	assert.Nil(t, err)
	assert.Len(t, output.UUID, 36)
	assert.Contains(t, output.UUID, "-")

	mockInterface.AssertExpectations(t)
}

func TestPostUuidWithNoHyphen(t *testing.T) {
	// arrange
	mockInterface := programminglib.MockInterface{}
	mockCall := mockInterface.On("NewUuid", true)
	mockCall.Return("1ce44be5fe6846f7a15351c1c91a4ae4")

	r := setupGin(&mockInterface)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/programming/uuid?no-hyphens=true", nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusOK)

	output := postUuidOutput{}
	err := json.Unmarshal(w.Body.Bytes(), &output)

	assert.Nil(t, err)
	assert.Len(t, output.UUID, 32)
	assert.NotContains(t, output.UUID, "-")

	mockInterface.AssertExpectations(t)
}

func TestPostJwtDebug(t *testing.T) {
	// arrange
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"
	expectedHeader := "{\"alg\":\"HS256\",\"typ\":\"JWT\"}"
	expectedPayload := "{\"iat\":1516239022,\"name\":\"John Doe\",\"sub\":\"1234567890\"}"

	mockInterface := programminglib.MockInterface{}
	mockCall := mockInterface.On("DebugJWT", tokenString)
	mockCall.Return(expectedHeader, expectedPayload, nil)

	r := setupGin(&mockInterface)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/programming/jwt", strings.NewReader(tokenString))

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusOK)

	output := postJwtDebuggerOutput{}
	err := json.Unmarshal(w.Body.Bytes(), &output)

	assert.Nil(t, err)
	assert.Equal(t, expectedHeader, output.Header)
	assert.Equal(t, expectedPayload, output.Payload)

	mockInterface.AssertExpectations(t)
}

func TestPostJwtDebugWithInvalidToken(t *testing.T) {
	// arrange
	tokenString := "xxxxx.yyyyy.zzzzz"
	err := errors.New("invalid token error")

	mockInterface := programminglib.MockInterface{}
	mockCall := mockInterface.On("DebugJWT", tokenString)
	mockCall.Return("", "", err)

	r := setupGin(&mockInterface)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/programming/jwt", strings.NewReader(tokenString))

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusBadRequest)
	assert.Contains(t, w.Body.String(), err.Error())

	mockInterface.AssertExpectations(t)
}
