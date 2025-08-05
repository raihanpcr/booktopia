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
	pb "wallet-service/proto"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Skenario 1: Tes GetBalance jika service berhasil merespons
func TestGetBalance_Success(t *testing.T) {
	// --- Arrange ---
	// 1. Siapkan Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/wallet/balance", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123") // Simulasi data dari middleware JWT

	// 2. Buat mock gRPC client
	mockClient := new(mock_proto.MockWalletServiceClient)
	// 3. Program mock untuk mengembalikan response sukses
	mockResponse := &pb.GetBalanceResponse{UserId: "user-123", Balance: 50000}
	mockClient.On("GetBalance", mock.Anything, &pb.GetBalanceRequest{UserId: "user-123"}).Return(mockResponse, nil)

	// 4. Buat handler dengan mock client
	h := NewWalletHandler(mockClient)

	// --- Act ---
	// Panggil method handler yang ingin dites
	err := h.GetBalance(c)

	// --- Assert ---
	// Pastikan hasilnya sesuai harapan
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp dto.BalanceResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, float64(50000), resp.Balance)
	mockClient.AssertExpectations(t)
}

// Skenario 2: Tes GetBalance jika service mengembalikan error
func TestGetBalance_ServiceError(t *testing.T) {
	// --- Arrange ---
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/wallet/balance", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	mockClient := new(mock_proto.MockWalletServiceClient)
	// Program mock untuk mengembalikan error
	mockClient.On("GetBalance", mock.Anything, mock.AnythingOfType("*proto.GetBalanceRequest")).Return(nil, errors.New("database connection lost"))
	h := NewWalletHandler(mockClient)

	// --- Act ---
	err := h.GetBalance(c)

	// --- Assert ---
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "database connection lost")
	mockClient.AssertExpectations(t)
}

// Skenario 3: Tes TopUp jika service berhasil merespons
func TestTopUp_Success(t *testing.T) {
	// --- Arrange ---
	// 1. Siapkan request body
	requestBody := dto.TopUpRequest{Amount: 100000, Method: "Transfer"}
	jsonBody, _ := json.Marshal(requestBody)

	// 2. Siapkan Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/wallet/topup", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	// 3. Buat dan program mock
	mockClient := new(mock_proto.MockWalletServiceClient)
	mockResponse := &pb.TopUpResponse{TopUpId: "topup-456", Status: "success"}
	mockClient.On("TopUp", mock.Anything, mock.AnythingOfType("*proto.TopUpRequest")).Return(mockResponse, nil)

	h := NewWalletHandler(mockClient)

	// --- Act ---
	err := h.TopUp(c)

	// --- Assert ---
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp pb.TopUpResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Equal(t, "topup-456", resp.TopUpId)
	mockClient.AssertExpectations(t)
}

// Skenario 4: Tes TopUp jika service mengembalikan error
func TestTopUp_ServiceError(t *testing.T) {
	// --- Arrange ---
	requestBody := dto.TopUpRequest{Amount: 100000, Method: "Transfer"}
	jsonBody, _ := json.Marshal(requestBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/wallet/topup", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-123")

	mockClient := new(mock_proto.MockWalletServiceClient)
	// Program mock untuk mengembalikan error
	mockClient.On("TopUp", mock.Anything, mock.AnythingOfType("*proto.TopUpRequest")).Return(nil, errors.New("failed to process top-up"))
	h := NewWalletHandler(mockClient)

	// --- Act ---
	err := h.TopUp(c)

	// --- Assert ---
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to process top-up")
	mockClient.AssertExpectations(t)
}
