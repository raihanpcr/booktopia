package routes

import (
	"book-service/internal/handler"

	"github.com/labstack/echo/v4"
)

// SetupRoutes mendaftarkan semua endpoint API untuk book-service.
// Dengan tidak menggunakan group "/api", endpoint akan lebih sederhana dan
// sesuai dengan yang diharapkan oleh gateway.
func SetupRoutes(e *echo.Echo, bookHandler *handler.BookHandler) {
	// Mendaftarkan endpoint langsung ke instance Echo 'e'
	e.POST("/books", bookHandler.CreateBook)
	e.GET("/books", bookHandler.GetAllBooks)
	e.GET("/books/:id", bookHandler.GetBookByID)
	e.PUT("/books/:id", bookHandler.UpdateBook)
	e.DELETE("/books/:id", bookHandler.DeleteBook)
}