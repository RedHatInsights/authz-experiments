package handler

import (
	"net/http"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	resp := runRequest(get("/"))
	assertHttpErrCodeAndMsg(t, http.StatusOK, "Connection to spiceDB successfully established", resp)
}
