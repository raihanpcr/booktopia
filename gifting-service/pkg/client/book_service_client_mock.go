package client

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// MockBookServiceClient adalah implementasi mock dari BookServiceClient.
type MockBookServiceClient struct {
	mock.Mock
}

// GetBookByID adalah implementasi mock yang bisa diprogram dalam tes.
func (m *MockBookServiceClient) GetBookByID(ctx context.Context, bookID string) (*BookDTO, error) {
	// Merekam pemanggilan method dan mengembalikan nilai yang sudah diatur.
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BookDTO), args.Error(1)
}