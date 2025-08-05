package repository

import (
	"context"
	"transaction-service/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	args := m.Called(ctx, transaction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetTransactionsByUserID(ctx context.Context, userID uint) ([]model.Transaction, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Transaction), args.Error(1)
}