package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	bookServiceURL string
}

func NewBookHandler(bookServiceURL string) *BookHandler {
	return &BookHandler{bookServiceURL: bookServiceURL}
}

// GetAllBooks godoc
// @Summary Get all books
// @Description Retrieve list of all books
// @Tags books
// @Produce json
// @Success 200 {array} dto.BookGetResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books [get]
func (h *BookHandler) GetBooks(c echo.Context) error {
	return h.proxyToBookService(c)
}

// GetBookByID godoc
// @Summary Get a book by ID
// @Description Retrieve a single book by its ID
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} dto.BookCreateResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /books/{id} [get]
func (h *BookHandler) GetBookByID(c echo.Context) error {
	return h.proxyToBookService(c)
}

// CreateBook godoc
// @Summary Create a new book
// @Description Create a new book record
// @Tags books
// @Accept json
// @Produce json
// @Param request body dto.CreateBookRequest true "Book to create"
// @Success 201 {object} dto.BookCreateResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /admin/books [post]
func (h *BookHandler) CreateBook(c echo.Context) error {
	return h.proxyToBookService(c)
}

// UpdateBook godoc
// @Summary Update a book
// @Description Update book information by ID
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param request body dto.UpdateBookRequest true "Updated book info"
// @Success 200 {object} dto.BookCreateResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /admin/books/{id} [put]
func (h *BookHandler) UpdateBook(c echo.Context) error {
	return h.proxyToBookService(c)
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book by ID
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} dto.DeleteResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security BearerAuth
// @Router /admin/books/{id} [delete]
func (h *BookHandler) DeleteBook(c echo.Context) error {
	return h.proxyToBookService(c)
}

// proxyToBookService adalah fungsi private yang berisi logika proxy
func (h *BookHandler) proxyToBookService(c echo.Context) error {
	requestPath := c.Request().URL.Path

	// Cari posisi "/books" untuk mendapatkan path yang relevan bagi backend
	booksIndex := strings.Index(requestPath, "/books")
	if booksIndex == -1 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid proxy path, /books not found"})
	}
	// Ambil path dari "/books" dan seterusnya.
	// Contoh: "/api/admin/books/123" akan menjadi "/books/123"
	backendPath := requestPath[booksIndex:]

	// Gunakan backendPath yang sudah dibersihkan untuk membuat URL tujuan
	targetURL, _ := url.Parse(fmt.Sprintf("%s%s", h.bookServiceURL, backendPath))

	proxyReq, err := http.NewRequest(c.Request().Method, targetURL.String(), c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create proxy request"})
	}

	// Salin semua header dari request asli
	proxyReq.Header = c.Request().Header
	// Salin query params agar tidak hilang
	proxyReq.URL.RawQuery = c.Request().URL.RawQuery

	// Kirim request ke book-service
	client := http.DefaultClient
	resp, err := client.Do(proxyReq)
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "failed to reach book service"})
	}
	defer resp.Body.Close()

	// Salin status code dan body dari response book-service ke response asli
	c.Response().WriteHeader(resp.StatusCode)
	io.Copy(c.Response().Writer, resp.Body)

	return nil
}
