package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Tes untuk endpoint Register
func TestRegister_ProxySuccess(t *testing.T) {
	// Arrange: Buat server backend palsu
	mockBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"success"}`))
	}))
	defer mockBackend.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewAuthHandler(mockBackend.URL)

	// Act: Panggil fungsi Register
	err := h.Register(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
}

// Tes untuk endpoint Login
func TestLogin_ProxySuccess(t *testing.T) {
	// Arrange: Buat server backend palsu
	mockBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"token":"jwt-token"}`))
	}))
	defer mockBackend.Close()

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := NewAuthHandler(mockBackend.URL)

	// Act: Panggil fungsi Login
	err := h.Login(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "jwt-token")
}
