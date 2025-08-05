package repository

import (
	"context"
	"book-service/internal/model"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockBookRepository adalah implementasi mock dari BookRepository.
// Ia menyematkan mock.Mock dari library testify untuk mendapatkan fungsionalitas mock.
type MockBookRepository struct {
	mock.Mock
}

// Create adalah implementasi mock untuk menyimpan buku.
func (m *MockBookRepository) Create(ctx context.Context, book *model.Book) error {
	// Merekam bahwa method ini dipanggil dengan argumen yang diberikan.
	args := m.Called(ctx, book)
	// Mengembalikan error yang sudah kita program (jika ada).
	return args.Error(0)
}

// FindAll adalah implementasi mock untuk mengambil semua buku.
func (m *MockBookRepository) FindAll(ctx context.Context) ([]model.Book, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// Mengembalikan slice buku dan error yang sudah diprogram.
	return args.Get(0).([]model.Book), args.Error(1)
}

// FindByID adalah implementasi mock untuk mengambil buku berdasarkan ID.
func (m *MockBookRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	// Mengembalikan satu objek buku dan error yang sudah diprogram.
	return args.Get(0).(*model.Book), args.Error(1)
}

// Update adalah implementasi mock untuk memperbarui buku.
func (m *MockBookRepository) Update(ctx context.Context, book *model.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

// Delete adalah implementasi mock untuk menghapus buku.
func (m *MockBookRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}