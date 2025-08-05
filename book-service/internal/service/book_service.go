package service

import (
	"context"
	"errors"
	"time"

	"book-service/internal/dto"
	"book-service/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BookService sekarang konsisten menggunakan DTO untuk input dan output
type BookService interface {
	CreateBook(ctx context.Context, req dto.CreateBookRequest) (*dto.BookResponse, error)
	GetBooks(ctx context.Context) ([]dto.BookResponse, error)
	GetBookByID(ctx context.Context, id string) (*dto.BookResponse, error)
	UpdateBook(ctx context.Context, id string, req dto.UpdateBookRequest) (*dto.BookResponse, error)
	DeleteBook(ctx context.Context, id string) error
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{repo: repo}
}

// CreateBook: Menerima DTO Request, mengembalikan DTO Response
func (s *bookService) CreateBook(ctx context.Context, req dto.CreateBookRequest) (*dto.BookResponse, error) {
	// Mapping dari DTO ke Model
	book := req.ToBookModel()

	// Logika Bisnis
	if book.Title == "" {
		return nil, errors.New("title cannot be empty")
	}
	book.ID = primitive.NewObjectID()
	book.Status = "available"
	book.CreatedAt = time.Now()

	// Panggil Repository
	if err := s.repo.Create(ctx, book); err != nil {
		return nil, err
	}

	// Mapping dari Model ke DTO Response
	response := dto.ToBookResponse(*book)
	return &response, nil
}

// GetBooks: Mengembalikan slice dari DTO Response
func (s *bookService) GetBooks(ctx context.Context) ([]dto.BookResponse, error) {
	books, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	// Mapping dari list Model ke list DTO Response
	return dto.ToBookResponseList(books), nil
}

// GetBookByID: Mengembalikan satu DTO Response
func (s *bookService) GetBookByID(ctx context.Context, id string) (*dto.BookResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid book ID format")
	}

	book, err := s.repo.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, errors.New("book not found")
	}

	// Mapping dari Model ke DTO Response
	response := dto.ToBookResponse(*book)
	return &response, nil
}

// UpdateBook: Menerima DTO Request, mengembalikan DTO Response
func (s *bookService) UpdateBook(ctx context.Context, id string, req dto.UpdateBookRequest) (*dto.BookResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid book ID format")
	}

	existingBook, err := s.repo.FindByID(ctx, objectID)
	if err != nil {
		return nil, err
	}
	if existingBook == nil {
		return nil, errors.New("book not found")
	}

	// Mapping dari DTO ke Model untuk update
	updatedData := req.ToBookModel()
	updatedData.ID = existingBook.ID
	updatedData.CreatedAt = existingBook.CreatedAt

	if err := s.repo.Update(ctx, updatedData); err != nil {
		return nil, err
	}

	// Mapping dari Model yang sudah diupdate ke DTO Response
	response := dto.ToBookResponse(*updatedData)
	return &response, nil
}

// DeleteBook: Tidak ada perubahan, karena tidak ada data yang dikembalikan
func (s *bookService) DeleteBook(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid book ID format")
	}
	return s.repo.Delete(ctx, objectID)
}
