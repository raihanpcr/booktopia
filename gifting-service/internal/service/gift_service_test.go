package service

import (
	"context"
	"errors"
	"gifting-service/internal/model"
	"gifting-service/internal/repository"
	"gifting-service/pkg/client"
	pb "gifting-service/proto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Skenario 1: Tes jika pengiriman hadiah berhasil
func TestSendGift_Success(t *testing.T) {
	// --- Arrange ---
	// 1. Buat mock untuk repository dan client
	mockRepo := new(repository.MockGiftingRepository)
	mockBookClient := new(client.MockBookServiceClient)

	// 2. Siapkan data request dan data palsu yang akan dikembalikan mock
	req := &pb.SendGiftRequest{
		DonorId:        "1",
		RecipientEmail: "penerima@example.com",
		BookId:         "101",
		Message:        "Selamat!",
	}
	mockBook := &client.BookDTO{ID: "101", IsDonationOnly: false}
	mockGift := &model.EbookGiftLog{GiftID: 99, Status: "pending"}

	// 3. Program mock
	// Jika GetBookByID dipanggil, kembalikan buku donasi
	mockBookClient.On("GetBookByID", mock.Anything, req.BookId).Return(mockBook, nil)
	// Jika CreateGift dipanggil, kembalikan data hadiah yang sudah disimpan
	mockRepo.On("CreateGift", mock.Anything, mock.AnythingOfType("*model.EbookGiftLog")).Return(mockGift, nil)

	// 4. Buat instance service dengan mock
	giftingService := NewGiftingService(mockRepo, mockBookClient)

	// --- Act ---
	result, err := giftingService.SendGift(context.Background(), req)

	// --- Assert ---
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(99), result.GiftID)
	assert.Equal(t, "pending", result.Status)
	mockRepo.AssertExpectations(t)
	mockBookClient.AssertExpectations(t)
}

// Skenario 2: Tes jika buku yang dikirim bukan buku donasi
func TestSendGift_BookNotForDonation(t *testing.T) {
	// --- Arrange ---
	mockRepo := new(repository.MockGiftingRepository)
	mockBookClient := new(client.MockBookServiceClient)

	req := &pb.SendGiftRequest{BookId: "102"}
	// Program mock untuk mengembalikan buku yang BUKAN donasi
	mockBook := &client.BookDTO{ID: "102", IsDonationOnly: true}
	mockBookClient.On("GetBookByID", mock.Anything, req.BookId).Return(mockBook, nil)

	giftingService := NewGiftingService(mockRepo, mockBookClient)

	// --- Act ---
	result, err := giftingService.SendGift(context.Background(), req)

	// --- Assert ---
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "this book cannot be gifted", err.Error())
	mockBookClient.AssertExpectations(t)
	// Pastikan CreateGift tidak pernah dipanggil
	mockRepo.AssertNotCalled(t, "CreateGift", mock.Anything, mock.Anything)
}

// Skenario 3: Tes jika buku tidak ditemukan
func TestSendGift_BookNotFound(t *testing.T) {
	// --- Arrange ---
	mockRepo := new(repository.MockGiftingRepository)
	mockBookClient := new(client.MockBookServiceClient)

	req := &pb.SendGiftRequest{BookId: "999"}
	// Program mock agar GetBookByID mengembalikan error
	mockBookClient.On("GetBookByID", mock.Anything, req.BookId).Return(nil, errors.New("not found"))

	giftingService := NewGiftingService(mockRepo, mockBookClient)

	// --- Act ---
	result, err := giftingService.SendGift(context.Background(), req)

	// --- Assert ---
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "book not found", err.Error())
	mockBookClient.AssertExpectations(t)
}
