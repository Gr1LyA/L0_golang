package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	s := New()
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	s.midHandle().ServeHTTP(rec, req)
	// log.Println(rec.Body.String())
	assert.Equal(t, rec.Body.String(), "alo")
}