package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"google.golang.org/protobuf/types/known/timestamppb"

	"transaction-service/internal/model"
	"transaction-service/internal/repository"
	"transaction-service/pkg/client"
	"transaction-service/pkg/messagebroker"
	pb "transaction-service/proto"
	wallet_pb "wallet-service/proto"
)

// TransactionService adalah interface untuk logika bisnis transaksi.
type TransactionService interface {
	CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.TransactionResponse, error)
	GetUserTransactions(ctx context.Context, req *pb.GetUserTransactionsRequest) (*pb.GetUserTransactionsResponse, error)
}

type transactionService struct {
	repo         repository.TransactionRepository
	bookClient   client.BookServiceClient
	walletClient wallet_pb.WalletServiceClient // gRPC Client untuk wallet-service
	producer     messagebroker.Producer
}

// NewTransactionService adalah constructor untuk service.
func NewTransactionService(
	repo repository.TransactionRepository,
	bookClient client.BookServiceClient,
	walletClient wallet_pb.WalletServiceClient,
	producer messagebroker.Producer,
) TransactionService {
	return &transactionService{
		repo:         repo,
		bookClient:   bookClient,
		walletClient: walletClient,
		producer:     producer,
	}
}

// CreateTransaction mengorkestrasi seluruh proses pembuatan transaksi.
func (s *transactionService) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.TransactionResponse, error) {
	var totalAmount float64
	var transactionDetailsModel []model.TransactionDetail

	userID, err := strconv.ParseUint(req.UserId, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	// 1. Validasi buku & hitung total (sinkron)
	for _, item := range req.Items {
		book, err := s.bookClient.GetBookByID(ctx, item.BookId)
		if err != nil {
			return nil, fmt.Errorf("failed to validate book with id %s", item.BookId)
		}
		if book.Status != "available" {
			return nil, fmt.Errorf("book '%s' is not available", book.Title)
		}

		itemPrice := book.Price * float64(item.Quantity)
		totalAmount += itemPrice

		transactionDetailsModel = append(transactionDetailsModel, model.TransactionDetail{
			BookID:       book.ID,
			Quantity:     int(item.Quantity),
			PricePerUnit: book.Price,
		})
	}
	
	// 2. Simpan transaksi dengan status PENDING
	txModel := &model.Transaction{
		UserID:      uint(userID),
		TotalAmount: totalAmount,
		Status:      "pending",
		Details:     transactionDetailsModel,
	}

	savedTransaction, err := s.repo.CreateTransaction(ctx, txModel)
	if err != nil {
		return nil, errors.New("failed to create initial transaction")
	}

	// 3. Buat dan kirim pesan ke Kafka (asinkron)
	eventPayload := map[string]interface{}{
		"transaction_id": fmt.Sprintf("%d", savedTransaction.ID),
		"user_id":        req.UserId,
		"total_amount":   totalAmount,
	}

	if err := s.producer.Publish(ctx, "transaction_created", eventPayload); err != nil {
		log.Printf("CRITICAL: Failed to publish event for tx_id %d.", savedTransaction.ID)
		// Di sini Anda TIDAK memanggil walletClient.Credit karena debit belum terjadi
		return nil, errors.New("failed to queue transaction")
	}

	// 4. Buat response gRPC untuk dikirim kembali ke client
	transactionDetailsProto := make([]*pb.TransactionDetail, len(savedTransaction.Details))
	for i, detail := range savedTransaction.Details {
		transactionDetailsProto[i] = &pb.TransactionDetail{
			BookId:       detail.BookID,
			Quantity:     int32(detail.Quantity),
			PricePerUnit: detail.PricePerUnit,
		}
	}

	// 4. Kembalikan respons cepat ke pengguna
	return &pb.TransactionResponse{
		TransactionId:   fmt.Sprintf("%d", savedTransaction.ID),
		UserId:          fmt.Sprintf("%d", savedTransaction.UserID),
		TransactionDate: timestamppb.New(savedTransaction.CreatedAt),
		TotalAmount:     savedTransaction.TotalAmount,
		Status:          savedTransaction.Status,
		Details:         transactionDetailsProto,
	}, nil
}

func (s *transactionService) GetUserTransactions(ctx context.Context, req *pb.GetUserTransactionsRequest) (*pb.GetUserTransactionsResponse, error) {
	userID, err := strconv.ParseUint(req.UserId, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	transactions, err := s.repo.GetTransactionsByUserID(ctx, uint(userID))
	if err != nil {
		return nil, err
	}

	var protoTransactions []*pb.TransactionResponse
	// Gunakan nama variabel yang jelas, misal: 'txModel'
	for _, txModel := range transactions {
		// Konversi detail transaksi
		detailsProto := make([]*pb.TransactionDetail, len(txModel.Details))
		for i, detail := range txModel.Details {
			detailsProto[i] = &pb.TransactionDetail{
				BookId:       detail.BookID,
				Quantity:     int32(detail.Quantity),
				PricePerUnit: detail.PricePerUnit,
			}
		}

		// Buat variabel mappedTx dan isi nilainya
		mappedTx := &pb.TransactionResponse{
			TransactionId:   fmt.Sprintf("%d", txModel.ID),
			UserId:          fmt.Sprintf("%d", txModel.UserID),
			TransactionDate: timestamppb.New(txModel.CreatedAt),
			TotalAmount:     txModel.TotalAmount,
			Status:          txModel.Status,
			Details:         detailsProto,
		}

		protoTransactions = append(protoTransactions, mappedTx)
	}

	return &pb.GetUserTransactionsResponse{Transactions: protoTransactions}, nil
}
