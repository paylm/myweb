package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/paylm/myweb/routers"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	r := routers.InitRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestBlogRoute(t *testing.T) {
	r := routers.InitRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/blog/list", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
