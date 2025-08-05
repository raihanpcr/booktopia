package service

import (
	"context"
	"errors"
	"testing"
	"wallet-service/internal/model"
	"wallet-service/internal/repository"
	pb "wallet-service/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Tes untuk GetBalance
func TestGetBalance_Success(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockWalletRepository)
	req := &pb.GetBalanceRequest{UserId: "1"}
	
	// Program mock untuk mengembalikan saldo
	mockRepo.On("GetBalance", uint(1)).Return(float64(100000), nil)
	
	walletService := NewWalletService(mockRepo)

	// Act
	balance, err := walletService.GetBalance(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, float64(100000), balance)
	mockRepo.AssertExpectations(t)
}

// Tes untuk TopUp
func TestTopUp_Success(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockWalletRepository)
	req := &pb.TopUpRequest{UserId: "1", Amount: 50000, Method: "Transfer"}

	// Program mock
	mockRepo.On("UpdateBalance", uint(1), req.Amount).Return(float64(150000), nil)
	mockRepo.On("CreateTopUp", mock.AnythingOfType("*model.TopUp")).Return(&model.TopUp{ID: 1}, nil)

	walletService := NewWalletService(mockRepo)

	// Act
	result, err := walletService.TopUp(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockRepo.AssertExpectations(t)
}

// Tes untuk Debit
func TestDebit_Success(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockWalletRepository)
	req := &pb.DebitRequest{UserId: "1", Amount: 25000}

	// Program mock untuk mengurangi saldo
	mockRepo.On("UpdateBalance", uint(1), -req.Amount).Return(float64(75000), nil)

	walletService := NewWalletService(mockRepo)

	// Act
	newBalance, err := walletService.Debit(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, float64(75000), newBalance)
	mockRepo.AssertExpectations(t)
}

// Tes untuk Debit dengan saldo tidak cukup
func TestDebit_InsufficientFunds(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockWalletRepository)
	req := &pb.DebitRequest{UserId: "1", Amount: 100000}

	// Program mock untuk mengembalikan error "insufficient funds"
	mockRepo.On("UpdateBalance", uint(1), -req.Amount).Return(float64(0), errors.New("insufficient funds"))

	walletService := NewWalletService(mockRepo)

	// Act
	_, err := walletService.Debit(context.Background(), req)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())
	mockRepo.AssertExpectations(t)
}

// Tes untuk Credit
func TestCredit_Success(t *testing.T) {
	// Arrange
	mockRepo := new(repository.MockWalletRepository)
	req := &pb.CreditRequest{UserId: "1", Amount: 10000}

	// Program mock untuk menambah saldo
	mockRepo.On("UpdateBalance", uint(1), req.Amount).Return(float64(110000), nil)

	walletService := NewWalletService(mockRepo)

	// Act
	newBalance, err := walletService.Credit(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, float64(110000), newBalance)
	mockRepo.AssertExpectations(t)
}