package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	// Import DTO dan gRPC client
	"gateway-service/internal/dto"
	pb "transaction-service/proto"
)

type TransactionHandler struct {
	transactionClient pb.TransactionServiceClient
}

func NewTransactionHandler(client pb.TransactionServiceClient) *TransactionHandler {
	return &TransactionHandler{transactionClient: client}
}

// CreateTransaction menerima HTTP/JSON dan menerjemahkannya ke gRPC
// CreateTransaction godoc
// @Summary      Buat transaksi pembelian buku
// @Description  Menerima list item buku dari user dan membuat transaksi baru (via gRPC ke transaction-service).
// @Tags         Gateway - Transaction
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.CreateTransactionRequest  true  "Data transaksi"
// @Success      201   {object}  dto.TransactionResponseApi
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /transactions [post]
func (h *TransactionHandler) CreateTransaction(c echo.Context) error {
	// 1. Ambil user_id dari context (yang diisi oleh middleware JWT)
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message: "Invalid user ID in token",
		})
	}

	// 2. Bind request JSON dari klien ke DTO lokal
	var req dto.CreateTransactionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message: "Invalid request body",
			Error: err.Error(),
		})
	}

	// 3. Terjemahkan DTO ke format message gRPC
	grpcItems := make([]*pb.BookOrderItem, len(req.Items))
	for i, item := range req.Items {
		grpcItems[i] = &pb.BookOrderItem{
			BookId:   item.BookID,
			Quantity: int32(item.Quantity),
		}
	}
	grpcReq := &pb.CreateTransactionRequest{
		UserId: userID,
		Items:  grpcItems,
	}

	// 3. Panggil service gRPC
	grpcResp, err := h.transactionClient.CreateTransaction(c.Request().Context(), grpcReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message: "Internal server error",
			Error: err.Error(),
		})
	}

	// 4. Terjemahkan response gRPC ke format JSON DTO
	respDto := dto.ToTransactionResponse(grpcResp)

	// 5. Kirim response JSON ke klien
	return c.JSON(http.StatusCreated, dto.TransactionResponseApi{
		StatusCode: http.StatusCreated,
		Message: "Success create transaction",
		Data: *respDto,
	})
}

// GetTransactions godoc
// @Summary      Ambil semua transaksi milik user
// @Description  Mengambil daftar semua transaksi berdasarkan user_id dari token JWT.
// @Tags         Gateway - Transaction
// @Produce      json
// @Security     BearerAuth
// @Success      200   {object}  dto.TransactionListResponse
// @Failure      401   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /transactions [get]
func (h *TransactionHandler) GetTransactions(c echo.Context) error {
	// Ambil user_id dari token JWT di context
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized,dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message: "Invalid user ID in token",
		})
	}

	// Buat request gRPC
	grpcReq := &pb.GetUserTransactionsRequest{UserId: userID}

	// Panggil transaction-service
	grpcResp, err := h.transactionClient.GetUserTransactions(c.Request().Context(), grpcReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message: "Internal server error",
			Error: err.Error(),
		})
	}

	// Konversi response gRPC ke format JSON DTO
	respDto := dto.ToTransactionListResponse(grpcResp)

	// return c.JSON(http.StatusOK, respDto)
	return c.JSON(http.StatusOK, dto.TransactionListResponse{
		StatusCode: http.StatusOK,
		Message: "Get transaction list successfully",
		Data: respDto,
	})
}
