package finance

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	financelib "github.com/renato0307/learning-go-lib/finance"
	"github.com/stretchr/testify/assert"
)

func setupGin(mockInterface *financelib.MockInterface) *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/v1")
	SetRouterGroup(mockInterface, v1)

	return r
}

func TestGetCurrConv(t *testing.T) {
	// arrange
	from := "EUR"
	to := "USD"
	amount := 10.0
	amountConverted := 11.0

	mockInterface := financelib.MockInterface{}
	mockCall := mockInterface.On("ConvertCurrency", from, to, amount)
	mockCall.Return(amountConverted, nil)

	r := setupGin(&mockInterface)
	w := httptest.NewRecorder()

	url := fmt.Sprintf("/v1/finance/currconv?from=%s&to=%s&amount=%f", from, to, amount)
	req, _ := http.NewRequest("GET", url, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusOK)

	output := getCurrConvOutput{}
	err := json.Unmarshal(w.Body.Bytes(), &output)

	assert.Nil(t, err)
	assert.Equal(t, from, output.From)
	assert.Equal(t, to, output.To)
	assert.Equal(t, amount, output.Amount)
	assert.Equal(t, amountConverted, output.ConvertedAmount)

	mockInterface.AssertExpectations(t)
}

func TestGetCurrConvWithMissingFrom(t *testing.T) {
	// arrange
	to := "USD"
	amount := 10.0

	mockInterface := financelib.MockInterface{}
	r := setupGin(&mockInterface)
	w := httptest.NewRecorder()

	url := fmt.Sprintf("/v1/finance/currconv?to=%s&amount=%f", to, amount)
	req, _ := http.NewRequest("GET", url, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func TestGetCurrConvWithMissingTo(t *testing.T) {
	// arrange
	from := "EUR"
	amount := 10.0

	mockInterface := financelib.MockInterface{}
	r := setupGin(&mockInterface)
	w := httptest.NewRecorder()

	url := fmt.Sprintf("/v1/finance/currconv?from=%s&amount=%f", from, amount)
	req, _ := http.NewRequest("GET", url, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func TestGetCurrConvWithMissingAmount(t *testing.T) {
	// arrange
	from := "EUR"
	to := "USD"

	mockInterface := financelib.MockInterface{}
	r := setupGin(&mockInterface)
	w := httptest.NewRecorder()

	url := fmt.Sprintf("/v1/finance/currconv?from=%s&to=%s", from, to)
	req, _ := http.NewRequest("GET", url, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func TestGetCurrConvWithInvalidAmount(t *testing.T) {
	// arrange
	from := "EUR"
	to := "USD"
	amount := "invalid"

	mockInterface := financelib.MockInterface{}
	r := setupGin(&mockInterface)
	w := httptest.NewRecorder()

	url := fmt.Sprintf("/v1/finance/currconv?from=%s&to=%s&amount=%s", from, to, amount)
	req, _ := http.NewRequest("GET", url, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusBadRequest)
}

func TestGetCurrConvWithLibraryError(t *testing.T) {
	// arrange
	from := "EUR"
	to := "USD"
	amount := 10.0

	mockInterface := financelib.MockInterface{}
	mockCall := mockInterface.On("ConvertCurrency", from, to, amount)
	mockCall.Return(0.0, errors.New("fake error"))

	r := setupGin(&mockInterface)
	w := httptest.NewRecorder()

	url := fmt.Sprintf("/v1/finance/currconv?from=%s&to=%s&amount=%f", from, to, amount)
	req, _ := http.NewRequest("GET", url, nil)

	// act
	r.ServeHTTP(w, req)

	// assert
	assert.Equal(t, w.Code, http.StatusInternalServerError)
	mockInterface.AssertExpectations(t)
}
