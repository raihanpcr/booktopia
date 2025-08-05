package handler

import (
	"net/http"

	"gateway-service/internal/dto"
	wallet_pb "wallet-service/proto"

	"github.com/labstack/echo/v4"
)

type WalletHandler struct {
	walletClient wallet_pb.WalletServiceClient
}

func NewWalletHandler(client wallet_pb.WalletServiceClient) *WalletHandler {
	return &WalletHandler{walletClient: client}
}

// GetBalance menangani GET /api/wallet/balance
// GetBalance godoc
// @Summary      Ambil saldo wallet user
// @Description  Mengambil saldo wallet berdasarkan user_id dari token JWT
// @Tags         Gateway - Wallet
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.BalanceResponseApi
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /wallet/balance [get]
func (h *WalletHandler) GetBalance(c echo.Context) error {
	// 1. Ambil user_id dari context (yang diisi oleh middleware JWT)
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// 2. Buat request gRPC
	grpcReq := &wallet_pb.GetBalanceRequest{
		UserId: userID,
	}

	// 3. Panggil wallet-service
	grpcResp, err := h.walletClient.GetBalance(c.Request().Context(), grpcReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message: "Internal server error",
			Error: err.Error(),
		})
	}

	// 4. Buat response DTO dan kirim sebagai JSON
	response := dto.BalanceResponse{
		UserID:  grpcResp.UserId,
		Balance: grpcResp.Balance,
	}

	return c.JSON(http.StatusOK, dto.BalanceResponseApi{
		StatusCode: http.StatusOK,
		Message: "Get data success",
		Data: response,
	})
}

// TopUp menangani POST /api/wallet/topup
// TopUp godoc
// @Summary      Top up saldo wallet
// @Description  Melakukan top up saldo wallet berdasarkan request user
// @Tags         Gateway - Wallet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body  dto.TopUpRequest  true  "Detail top up saldo"
// @Success      200  {object}  dto.TopUpResponseDto
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /wallet/topup [post]
func (h *WalletHandler) TopUp(c echo.Context) error {
	// Ambil user_id dari token
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		// return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user ID in token"})
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message: "Invalid user ID in token",
		})
	}

	// Bind request body
	var req dto.TopUpRequest
	if err := c.Bind(&req); err != nil {
		// return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message: "Invalid request body",
			Error: err.Error(),
		})
	}

	// Buat request gRPC
	grpcReq := &wallet_pb.TopUpRequest{
		UserId: userID,
		Amount: req.Amount,
		Method: req.Method,
	}

	// Panggil wallet-service
	grpcResp, err := h.walletClient.TopUp(c.Request().Context(), grpcReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Kirim response
	return c.JSON(http.StatusOK, dto.TopUpResponseApi{
		StatusCode: http.StatusOK,
		Message: "Top up success",
		Data: grpcResp,
	})
}
