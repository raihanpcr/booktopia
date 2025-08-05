package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gateway-service/internal/dto"
	mock_proto "gateway-service/proto"
	pb "gifting-service/proto"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Skenario 1: Tes jika pengiriman hadiah berhasil
func TestSendGift_Success(t *testing.T) {
	// --- Arrange ---
	// 1. Siapkan request body
	requestBody := dto.SendGiftRequest{
		RecipientEmail: "penerima@example.com",
		BookID:         "book-for-donation-123",
		Message:        "Selamat!",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// 2. Siapkan Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/gifts", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-donor-456") // Simulasi ID donor dari middleware

	// 3. Buat dan program mock gRPC client
	mockClient := new(mock_proto.MockGiftingServiceClient)
	mockResponse := &pb.SendGiftResponse{
		GiftId:         "gift-789",
		Status:         "pending",
		RecipientEmail: "penerima@example.com",
	}
	mockClient.On("SendGift", mock.Anything, mock.AnythingOfType("*proto.SendGiftRequest")).Return(mockResponse, nil)

	// 4. Buat handler dengan mock
	h := NewGiftingHandler(mockClient)

	// --- Act ---
	err := h.SendGift(c)

	// --- Assert ---
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	// Verifikasi isi response
	var resp pb.SendGiftResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "gift-789", resp.GiftId)
	assert.Equal(t, "penerima@example.com", resp.RecipientEmail)

	mockClient.AssertExpectations(t)
}

// Skenario 2: Tes jika gifting-service mengembalikan error
func TestSendGift_ServiceError(t *testing.T) {
	// --- Arrange ---
	// 1. Siapkan request body (sama seperti sebelumnya)
	requestBody := dto.SendGiftRequest{
		RecipientEmail: "penerima@example.com",
		BookID:         "book-for-donation-123",
	}
	jsonBody, _ := json.Marshal(requestBody)

	// 2. Siapkan Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/gifts", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-donor-456")

	// 3. Buat dan program mock gRPC client untuk GAGAL
	mockClient := new(mock_proto.MockGiftingServiceClient)
	// Program mock untuk mengembalikan nil dan sebuah error
	mockClient.On("SendGift", mock.Anything, mock.AnythingOfType("*proto.SendGiftRequest")).Return(nil, errors.New("gifting service is down"))

	// 4. Buat handler
	h := NewGiftingHandler(mockClient)

	// --- Act ---
	err := h.SendGift(c)

	// --- Assert ---
	assert.NoError(t, err)                                    // Fungsi handler-nya sendiri tidak error
	assert.Equal(t, http.StatusInternalServerError, rec.Code) // Tapi status code HTTP-nya 500
	assert.Contains(t, rec.Body.String(), "gifting service is down")
	mockClient.AssertExpectations(t)
}
