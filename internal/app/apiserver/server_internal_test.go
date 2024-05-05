package apiserver

import (
	"github.com/stretchr/testify/assert"
	"github.com/yeboka/final-project/internal/app/store/teststore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_HandleUsersCreate(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", nil)
	s := newServer(teststore.New())
	s.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}
