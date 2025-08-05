package client

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockBookServiceClient struct {
	mock.Mock
}

func (m *MockBookServiceClient) GetBookByID(ctx context.Context, bookID string) (*BookDTO, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*BookDTO), args.Error(1)
}