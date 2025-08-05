package service

import (
	"context"
	"testing"
	"time"

	"book-service/internal/dto"
	"book-service/internal/model"
	"book-service/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Test GetBookByID ---

func TestGetBookByID_Success(t *testing.T) {
	mockRepo := new(repository.MockBookRepository)
	bookID := primitive.NewObjectID()
	mockBook := &model.Book{ID: bookID, Title: "Buku Sukses"}

	// Arrange: Program mock untuk mengembalikan buku
	mockRepo.On("FindByID", mock.Anything, bookID).Return(mockBook, nil)
	bookService := NewBookService(mockRepo)

	// Act: Panggil service
	result, err := bookService.GetBookByID(context.Background(), bookID.Hex())

	// Assert: Pastikan hasilnya benar
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockBook.Title, result.Title)
	mockRepo.AssertExpectations(t)
}

func TestGetBookByID_NotFound(t *testing.T) {
	mockRepo := new(repository.MockBookRepository)
	bookID := primitive.NewObjectID()

	// Arrange: Program mock untuk tidak mengembalikan apa-apa (nil)
	mockRepo.On("FindByID", mock.Anything, bookID).Return(nil, nil)
	bookService := NewBookService(mockRepo)

	// Act
	result, err := bookService.GetBookByID(context.Background(), bookID.Hex())

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "book not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetBookByID_InvalidID(t *testing.T) {
	mockRepo := new(repository.MockBookRepository)
	bookService := NewBookService(mockRepo)

	// Act
	result, err := bookService.GetBookByID(context.Background(), "id-tidak-valid")

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid book ID format", err.Error())
}

// --- Test GetBooks ---

func TestGetBooks_Success(t *testing.T) {
	mockRepo := new(repository.MockBookRepository)
	mockBooks := []model.Book{
		{ID: primitive.NewObjectID(), Title: "Buku Satu"},
		{ID: primitive.NewObjectID(), Title: "Buku Dua"},
	}

	// Arrange
	mockRepo.On("FindAll", mock.Anything).Return(mockBooks, nil)
	bookService := NewBookService(mockRepo)

	// Act
	results, err := bookService.GetBooks(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, results)
	assert.Len(t, results, 2)
	assert.Equal(t, "Buku Satu", results[0].Title)
	mockRepo.AssertExpectations(t)
}

// --- Test CreateBook ---

func TestCreateBook_Success(t *testing.T) {
	mockRepo := new(repository.MockBookRepository)
	req := dto.CreateBookRequest{Title: "Buku Baru", Author: "Penulis Baru"}

	// Arrange: Program mock agar `Create` berhasil (tidak mengembalikan error)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Book")).Return(nil)
	bookService := NewBookService(mockRepo)

	// Act
	result, err := bookService.CreateBook(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)
	assert.Equal(t, "available", result.Status)
	mockRepo.AssertExpectations(t)
}

// --- Test UpdateBook ---

func TestUpdateBook_Success(t *testing.T) {
	mockRepo := new(repository.MockBookRepository)
	bookID := primitive.NewObjectID()
	req := dto.UpdateBookRequest{Title: "Judul Diupdate"}

	mockBook := &model.Book{ID: bookID, Title: "Judul Lama", CreatedAt: time.Now()}

	// Arrange
	mockRepo.On("FindByID", mock.Anything, bookID).Return(mockBook, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Book")).Return(nil)
	bookService := NewBookService(mockRepo)

	// Act
	result, err := bookService.UpdateBook(context.Background(), bookID.Hex(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.Title, result.Title)
	mockRepo.AssertExpectations(t)
}

func TestUpdateBook_NotFound(t *testing.T) {
	mockRepo := new(repository.MockBookRepository)
	bookID := primitive.NewObjectID()
	req := dto.UpdateBookRequest{Title: "Judul Diupdate"}

	// Arrange: Program FindByID agar tidak menemukan buku
	mockRepo.On("FindByID", mock.Anything, bookID).Return(nil, nil)
	bookService := NewBookService(mockRepo)

	// Act
	result, err := bookService.UpdateBook(context.Background(), bookID.Hex(), req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "book not found", err.Error())
	mockRepo.AssertExpectations(t)
}

// --- Test DeleteBook ---

func TestDeleteBook_Success(t *testing.T) {
	mockRepo := new(repository.MockBookRepository)
	bookID := primitive.NewObjectID()

	// Arrange
	mockRepo.On("Delete", mock.Anything, bookID).Return(nil)
	bookService := NewBookService(mockRepo)

	// Act
	err := bookService.DeleteBook(context.Background(), bookID.Hex())

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
