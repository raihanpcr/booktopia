package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// Test skenario ketika pengguna memiliki peran "admin" (sukses)
func TestAdminOnlyMiddleware_Success(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulasikan bahwa middleware JWT sebelumnya telah men-set peran "admin"
	c.Set("role", "admin")

	// Buat handler dummy yang akan dipanggil jika middleware lolos
	dummyHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	// Buat middleware handler
	adminMiddleware := AdminOnlyMiddleware(dummyHandler)

	// Act
	err := adminMiddleware(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "success", rec.Body.String())
}

// Test skenario ketika pengguna TIDAK memiliki peran "admin" (gagal/forbidden)
func TestAdminOnlyMiddleware_Forbidden(t *testing.T) {
	// Arrange
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Simulasikan pengguna dengan peran lain, misalnya "customer"
	c.Set("role", "customer")

	dummyHandler := func(c echo.Context) error {
		// Handler ini seharusnya tidak pernah dieksekusi
		return c.String(http.StatusOK, "should not be reached")
	}

	adminMiddleware := AdminOnlyMiddleware(dummyHandler)

	// Act
	err := adminMiddleware(c)

	// Assert
	// 1. Pastikan c.JSON() tidak mengembalikan error internal
	assert.NoError(t, err) 
	
	// 2. Periksa status code yang ditulis ke recorder
	assert.Equal(t, http.StatusForbidden, rec.Code) 
	
	// 3. Periksa body JSON yang ditulis ke recorder
	expectedJSON := `{"error":"access denied: admin role required"}`
	assert.JSONEq(t, expectedJSON, rec.Body.String())
}