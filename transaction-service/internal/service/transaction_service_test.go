package service

import (
	"context"
	"errors"
	"testing"
	"transaction-service/internal/model"
	"transaction-service/internal/repository"
	"transaction-service/pkg/client"
	"transaction-service/pkg/messagebroker"
	pb "transaction-service/proto"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Skenario 1: Tes CreateTransaction dengan alur Kafka yang sukses
func TestCreateTransaction_KafkaSuccess(t *testing.T) {
	// --- Arrange ---
	mockRepo := new(repository.MockTransactionRepository)
	mockBookClient := new(client.MockBookServiceClient)
	mockProducer := new(messagebroker.MockProducer)

	req := &pb.CreateTransactionRequest{
		UserId: "1",
		Items:  []*pb.BookOrderItem{{BookId: "101", Quantity: 2}},
	}
	mockBook := &client.BookDTO{ID: "101", Status: "available", Price: 50000}
	
	// GORM akan mengisi ID setelah Create, jadi kita siapkan modelnya
	mockSavedTx := &model.Transaction{ID: 99, UserID: 1, Status: "pending"}

	// Program semua mock
	mockBookClient.On("GetBookByID", mock.Anything, "101").Return(mockBook, nil)
	mockRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*model.Transaction")).Return(mockSavedTx, nil)
	mockProducer.On("Publish", mock.Anything, "transaction_created", mock.Anything).Return(nil)

	// Buat instance service dengan mock producer
	transactionService := NewTransactionService(mockRepo, mockBookClient, nil, mockProducer)

	// --- Act ---
	result, err := transactionService.CreateTransaction(context.Background(), req)

	// --- Assert ---
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "99", result.TransactionId)
	assert.Equal(t, "pending", result.Status) // Pastikan statusnya pending
	
	// Verifikasi semua mock dipanggil
	mockBookClient.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
}

// Skenario 2: Tes jika gagal mengirim pesan ke Kafka
func TestCreateTransaction_KafkaPublishFailed(t *testing.T) {
	// --- Arrange ---
	mockRepo := new(repository.MockTransactionRepository)
	mockBookClient := new(client.MockBookServiceClient)
	mockProducer := new(messagebroker.MockProducer)

	req := &pb.CreateTransactionRequest{
		UserId: "1",
		Items:  []*pb.BookOrderItem{{BookId: "101", Quantity: 1}},
	}
	mockBook := &client.BookDTO{ID: "101", Status: "available", Price: 50000}
	mockSavedTx := &model.Transaction{ID: 99, UserID: 1}

	// Program mock
	mockBookClient.On("GetBookByID", mock.Anything, "101").Return(mockBook, nil)
	mockRepo.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*model.Transaction")).Return(mockSavedTx, nil)
	// Program Producer untuk GAGAL
	mockProducer.On("Publish", mock.Anything, "transaction_created", mock.Anything).Return(errors.New("kafka is down"))
	
	transactionService := NewTransactionService(mockRepo, mockBookClient, nil, mockProducer)

	// --- Act ---
	result, err := transactionService.CreateTransaction(context.Background(), req)

	// --- Assert ---
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "failed to queue transaction", err.Error())
	mockProducer.AssertExpectations(t)
}