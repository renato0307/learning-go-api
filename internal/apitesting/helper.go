package apitesting

import (
	"net/http"
	"net/http/httptest"
)

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {

	return PerformRequestWithHeader(r, method, path, nil)
}

func PerformRequestWithHeader(
	r http.Handler,
	method, path string,
	header http.Header) *httptest.ResponseRecorder {

	req := httptest.NewRequest(method, path, nil)
	req.Header = header
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
