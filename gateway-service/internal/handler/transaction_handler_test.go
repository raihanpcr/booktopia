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
	pb "transaction-service/proto"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Skenario 1: Tes CreateTransaction jika service berhasil merespons
func TestCreateTransaction_Success(t *testing.T) {
	// --- Arrange ---
	// 1. Siapkan request body JSON
	requestBody := dto.CreateTransactionRequest{Items: []dto.BookOrderItem{{BookID: "book-123", Quantity: 1}}}
	jsonBody, _ := json.Marshal(requestBody)

	// 2. Siapkan Echo context
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-456") // Simulasi data dari middleware JWT

	// 3. Buat mock gRPC client dan program perilakunya
	mockClient := new(mock_proto.MockTransactionServiceClient)
	// Program mock untuk mengembalikan response sukses
	mockClient.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*proto.CreateTransactionRequest")).Return(&pb.TransactionResponse{
		TransactionId: "tx-789",
		Status:        "completed",
	}, nil)

	// 4. Buat handler dengan mock client
	h := NewTransactionHandler(mockClient)

	// --- Act ---
	// Panggil method handler yang ingin dites
	err := h.CreateTransaction(c)

	// --- Assert ---
	// Pastikan hasilnya sesuai harapan
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockClient.AssertExpectations(t)
}

// Skenario 2: Tes CreateTransaction jika service mengembalikan error
func TestCreateTransaction_ServiceError(t *testing.T) {
	// --- Arrange ---
	requestBody := dto.CreateTransactionRequest{Items: []dto.BookOrderItem{{BookID: "book-123", Quantity: 1}}}
	jsonBody, _ := json.Marshal(requestBody)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/api/transactions", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-456")

	mockClient := new(mock_proto.MockTransactionServiceClient)
	// Program mock untuk mengembalikan error
	mockClient.On("CreateTransaction", mock.Anything, mock.AnythingOfType("*proto.CreateTransactionRequest")).Return(nil, errors.New("insufficient funds"))
	h := NewTransactionHandler(mockClient)

	// --- Act ---
	err := h.CreateTransaction(c)

	// --- Assert ---
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "insufficient funds")
	mockClient.AssertExpectations(t)
}

// Skenario 3: Tes GetTransactions jika service berhasil merespons
func TestGetTransactions_Success(t *testing.T) {
	// --- Arrange ---
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-456")

	mockClient := new(mock_proto.MockTransactionServiceClient)
	// Program mock untuk mengembalikan daftar transaksi
	mockResponse := &pb.GetUserTransactionsResponse{
		Transactions: []*pb.TransactionResponse{
			{TransactionId: "tx-1", Status: "completed"},
			{TransactionId: "tx-2", Status: "completed"},
		},
	}
	mockClient.On("GetUserTransactions", mock.Anything, mock.AnythingOfType("*proto.GetUserTransactionsRequest")).Return(mockResponse, nil)

	h := NewTransactionHandler(mockClient)

	// --- Act ---
	err := h.GetTransactions(c)

	// --- Assert ---
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp []*dto.TransactionResponse
	json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Len(t, resp, 2)
	assert.Equal(t, "tx-1", resp[0].TransactionID)
	mockClient.AssertExpectations(t)
}

// Skenario 4: Tes GetTransactions jika service mengembalikan error
func TestGetTransactions_ServiceError(t *testing.T) {
	// --- Arrange ---
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/transactions", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user_id", "user-456")

	mockClient := new(mock_proto.MockTransactionServiceClient)
	// Program mock untuk mengembalikan error
	mockClient.On("GetUserTransactions", mock.Anything, mock.AnythingOfType("*proto.GetUserTransactionsRequest")).Return(nil, errors.New("service down"))

	h := NewTransactionHandler(mockClient)

	// --- Act ---
	err := h.GetTransactions(c)

	// --- Assert ---
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "service down")
	mockClient.AssertExpectations(t)
}
