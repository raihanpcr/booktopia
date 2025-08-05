package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// AdminOnlyMiddleware memeriksa apakah pengguna memiliki peran 'admin'.
// Middleware ini harus dijalankan SETELAH JwtAuthMiddleware.
func AdminOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Ambil peran pengguna dari context, yang sudah di-set oleh JwtAuthMiddleware
		role := c.Get("role")

		// Periksa apakah peran adalah "admin"
		if role == "admin" {
			// Jika ya, lanjutkan ke handler berikutnya
			return next(c)
		}

		// Jika tidak, kembalikan error 403 Forbidden
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied: admin role required"})
	}
}