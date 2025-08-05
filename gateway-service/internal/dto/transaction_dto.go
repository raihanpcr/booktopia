package dto

import (
	"time"
	pb "transaction-service/proto"
)

// DTO untuk request dari client (JSON)
type CreateTransactionRequest struct {
	Items  []BookOrderItem `json:"items"`
}
type BookOrderItem struct {
	BookID   string `json:"book_id"`
	Quantity int    `json:"quantity"`
}

// DTO untuk response ke client (JSON)
type TransactionResponse struct {
	TransactionID   string                      `json:"transaction_id"`
	UserID          string                      `json:"user_id"`
	TransactionDate time.Time                   `json:"transaction_date"`
	TotalAmount     float64                     `json:"total_amount"`
	Status          string                      `json:"status"`
	Details         []TransactionDetailResponse `json:"details"`
}
type TransactionDetailResponse struct {
	BookID       string  `json:"book_id"`
	Quantity     int     `json:"quantity"`
	PricePerUnit float64 `json:"price_per_unit"`
}

type TransactionListResponse struct {
	StatusCode 	int              	`json:"status_code" validate:"required" example:"200"`
	Message    	string           	`json:"message" validate:"required" example:"Create data success"`
	Data 		[]*TransactionResponse `json:"data"`
}

type TransactionResponseApi struct {
	StatusCode 	int              		`json:"status_code" validate:"required" example:"201"`
	Message    	string           		`json:"message" validate:"required" example:"Create data success"`
	Data 		TransactionResponse 	`json:"data"`
}

// Mapper dari gRPC response ke DTO response
func ToTransactionResponse(grpcResp *pb.TransactionResponse) *TransactionResponse {
	details := make([]TransactionDetailResponse, len(grpcResp.Details))
	for i, d := range grpcResp.Details {
		details[i] = TransactionDetailResponse{
			BookID:       d.BookId,
			Quantity:     int(d.Quantity),
			PricePerUnit: d.PricePerUnit,
		}
	}

	return &TransactionResponse{
		TransactionID:   grpcResp.TransactionId,
		UserID:          grpcResp.UserId,
		TransactionDate: grpcResp.TransactionDate.AsTime(),
		TotalAmount:     grpcResp.TotalAmount,
		Status:          grpcResp.Status,
		Details:         details,
	}
}

func ToTransactionListResponse(grpcResp *pb.GetUserTransactionsResponse) []*TransactionResponse {
	var transactions []*TransactionResponse
	for _, grpcTx := range grpcResp.Transactions {
		transactions = append(transactions, ToTransactionResponse(grpcTx))
	}
	return transactions
}