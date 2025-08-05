package repository

import (
	"wallet-service/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) GetBalance(userID uint) (float64, error) {
	args := m.Called(userID)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockWalletRepository) UpdateBalance(userID uint, amount float64) (float64, error) {
	args := m.Called(userID, amount)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockWalletRepository) CreateTopUp(topUp *model.TopUp) (*model.TopUp, error) {
	args := m.Called(topUp)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.TopUp), args.Error(1)
}