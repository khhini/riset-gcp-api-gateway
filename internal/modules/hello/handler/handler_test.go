package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	response := `{"message":"Hello, World!"}`

	r := gin.Default()
	NewHandler(r.Group("/api"))

	req, err := http.NewRequest("GET", "/api/hello", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, response, w.Body.String())
}
