package handler

import (
	"gateway-service/internal/dto"
	gifting_pb "gifting-service/proto"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GiftingHandler struct {
	giftingClient gifting_pb.GiftingServiceClient
}

func NewGiftingHandler(client gifting_pb.GiftingServiceClient) *GiftingHandler {
	return &GiftingHandler{giftingClient: client}
}

// SendGift menangani POST /api/gifts
// SendGift godoc
// @Summary Kirim hadiah buku ke user lain
// @Description Meneruskan permintaan pengiriman hadiah ke Gifting Service
// @Tags Gateway - Gifting
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.SendGiftRequest true "Data pengiriman hadiah"
// @Success 201 {object} dto.TemplateSendGiftResponseApi
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /gifts [post]
func (h *GiftingHandler) SendGift(c echo.Context) error {
	// 1. Ambil ID donor (pengirim) dari token JWT
	donorID, ok := c.Get("user_id").(string)
	if !ok || donorID == "" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			StatusCode: http.StatusUnauthorized,
			Message: "Unauthorized",
		})
	}

	// 2. Bind request body ke DTO
	var req dto.SendGiftRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			StatusCode: http.StatusBadRequest,
			Message: "Invalid Request Body",
			Error: err.Error(),
		})
	}

	// 3. Buat request gRPC
	grpcReq := &gifting_pb.SendGiftRequest{
		DonorId:        donorID,
		RecipientEmail: req.RecipientEmail,
		BookId:         req.BookID,
		Message:        req.Message,
	}

	// 4. Panggil gifting-service
	grpcResp, err := h.giftingClient.SendGift(c.Request().Context(), grpcReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message: "Internal Server Error",
			Error: err.Error(),
		})
	}

	// 5. Kirim response ke client
	// return c.JSON(http.StatusCreated, grpcResp)
	return c.JSON(http.StatusCreated, dto.TemplateSendGiftResponse{
		StatusCode: http.StatusCreated,
		Message: "create data successfully",
		Data: grpcResp,
	})
}
