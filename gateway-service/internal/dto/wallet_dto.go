package dto

import wallet_pb "wallet-service/proto"

// TopUpRequest adalah DTO untuk request body saat melakukan top-up.
type TopUpRequest struct {
	Amount float64 `json:"amount" validate:"required,gt=0"`
	Method string  `json:"method" validate:"required"`
}

// BalanceResponse adalah DTO untuk menampilkan saldo ke client.
type BalanceResponse struct {
	UserID  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}

// TopUpResponse dummy struct untuk Swagger docs
type TopUpResponse struct {
	TransactionID string  `json:"transaction_id" example:"trx123"`
	UserID        string  `json:"user_id" example:"user123"`
	Amount        float64 `json:"amount" example:"50000"`
	Method        string  `json:"method" example:"credit_card"`
	Status        string  `json:"status" example:"success"`
	CreatedAt     string  `json:"created_at" example:"2025-07-25T15:00:00Z"`
}

type TopUpResponseApi struct {
	StatusCode 	int              		`json:"status_code" validate:"required" example:"201"`
	Message    	string           		`json:"message" validate:"required" example:"Create data success"`
	Data 		*wallet_pb.TopUpResponse `json:"data"`
}

type TopUpResponseDto struct {
	StatusCode 	int              	`json:"status_code" validate:"required" example:"201"`
	Message    	string           	`json:"message" validate:"required" example:"Create data success"`
	Data 		TopUpResponse	`json:"data"`
}

type BalanceResponseApi struct {
	StatusCode 	int              	`json:"status_code" validate:"required" example:"200"`
	Message    	string           	`json:"message" validate:"required" example:"Get data success"`
	Data 		BalanceResponse	`json:"data"`
}