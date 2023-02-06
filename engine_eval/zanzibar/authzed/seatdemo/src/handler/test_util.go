package handler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
)

func get(uri string) *http.Request {
	return httptest.NewRequest(http.MethodGet, uri, strings.NewReader(""))
}

func post(uri string, body string) *http.Request {
	return reqWithBody(http.MethodPost, uri, body)
}

func put(uri string, body string) *http.Request {
	return reqWithBody(http.MethodPut, uri, body)
}

func delete(uri string) *http.Request {
	return httptest.NewRequest(http.MethodDelete, uri, strings.NewReader(""))
}

func reqWithBody(method string, uri string, body string) *http.Request {
	req := httptest.NewRequest(method, uri, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func runRequest(req *http.Request) *http.Response {
	e := GetEcho()
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec.Result()
}

func assertHttpErrCodeAndMsg(t *testing.T, statusCode int, message string, resp *http.Response) {
	if assert.NotNil(t, resp) {
		assert.Equal(t, statusCode, resp.StatusCode)

		payload := new(strings.Builder)
		_, err := io.Copy(payload, resp.Body)
		assert.NoError(t, err)

		assert.Contains(t, payload.String(), message)
	}
}

func assertJsonResponse(t *testing.T, resp *http.Response, statusCode int, template string, args ...interface{}) {
	if assert.NotNil(t, resp) {
		assert.Equal(t, statusCode, resp.StatusCode)

		payload := new(strings.Builder)
		_, err := io.Copy(payload, resp.Body)
		assert.NoError(t, err)

		ja := jsonassert.New(t)
		ja.Assertf(payload.String(), template, args...)
	}
}
