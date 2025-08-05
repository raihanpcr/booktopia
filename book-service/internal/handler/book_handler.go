package handler

import (
	"net/http"

	"book-service/internal/dto"
	"book-service/internal/service"

	"github.com/labstack/echo/v4"
)

// BookHandler bergantung pada service (tidak ada perubahan di sini)
type BookHandler struct {
	service service.BookService
}

func NewBookHandler(service service.BookService) *BookHandler {
	return &BookHandler{service: service}
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
// @Router /books [post]
func (h *BookHandler) CreateBook(c echo.Context) error {
	// 1. Bind ke Request DTO, bukan Model
	var req dto.CreateBookRequest
	if err := c.Bind(&req); err != nil {
		// 2. Gunakan DTO ErrorResponse yang standar
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}

	// 3. Panggil service dengan DTO
	createdBook, err := h.service.CreateBook(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Details: err.Error(),
		})
	}

	// 4. Kembalikan Response DTO dari service
	// return c.JSON(http.StatusCreated, createdBook)
	return c.JSON(http.StatusCreated, dto.BookCreateResponse{
		StatusCode: http.StatusCreated,
		Message: "Create data successfully",
		Data: *createdBook,
	})
}

// GetAllBooks godoc
// @Summary Get all books
// @Description Retrieve list of all books
// @Tags books
// @Produce json
// @Success 200 {array} dto.BookGetResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books [get]
func (h *BookHandler) GetAllBooks(c echo.Context) error {
	books, err := h.service.GetBooks(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.BookGetResponse{
		StatusCode: http.StatusOK,
		Message: "Get books successfully",
		Data: books,
	})
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
	id := c.Param("id")
	book, err := h.service.GetBookByID(c.Request().Context(), id)
	if err != nil {
		// Jika service mengembalikan error spesifik, kita bisa handle di sini
		return c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Data not found",
			Details: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.BookCreateResponse{
		StatusCode: http.StatusOK,
		Message: "Get book by id successfully",
		Data: *book,
	})
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
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c echo.Context) error {
	id := c.Param("id")
	var req dto.UpdateBookRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}

	updatedBook, err := h.service.UpdateBook(c.Request().Context(), id, req)
	if err != nil {
		if err.Error() == "book not found" {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Data not found",
				Details: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal Server Error",
			Details: err.Error(),
		})
	}
	return c.JSON(http.StatusOK, dto.BookCreateResponse{
		StatusCode: http.StatusOK,
		Message: "Update book by id successfully",
		Data: *updatedBook,
	})
}

// DeleteBook godoc
// @Summary Delete a book
// @Description Delete a book by ID
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} dto.DeleteResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c echo.Context) error {
	id := c.Param("id")
	if err := h.service.DeleteBook(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.DeleteResponse{
		Code:    http.StatusNoContent,
		Message: "Book Deleted Succesfully",
	})
}
