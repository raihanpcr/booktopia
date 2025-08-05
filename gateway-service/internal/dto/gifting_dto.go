package dto

import gifting_pb "gifting-service/proto"

// SendGiftRequest adalah DTO untuk request body saat mengirim hadiah.
type SendGiftRequest struct {
	RecipientEmail string `json:"recipient_email" validate:"required,email"`
	BookID         string `json:"book_id" validate:"required"`
	Message        string `json:"message"`
}

// SendGiftResponse dummy untuk Swagger docs
type SendGiftResponse struct {
	GiftId     string `json:"gift_id" example:"gift123"`
	Recipient  string `json:"recipient_email" example:"jane@example.com"`
	BookId     string `json:"book_id" example:"book001"`
	Message    string `json:"message" example:"Selamat membaca!"`
	CreatedAt  string `json:"created_at" example:"2025-07-25T15:00:00Z"`
}

type TemplateSendGiftResponseApi struct {
	StatusCode int              	`json:"status_code" validate:"required" example:"201"`
	Message    string           	`json:"message" validate:"required" example:"Create data success"`
	Data       SendGiftRequest 	`json:"data"`
}
type TemplateSendGiftResponse struct {
	StatusCode int              	`json:"status_code" validate:"required" example:"201"`
	Message    string           	`json:"message" validate:"required" example:"Create data success"`
	Data       *gifting_pb.SendGiftResponse 	`json:"data"`
}