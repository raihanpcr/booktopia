package server

import (
	"context"
	"transaction-service/internal/service"
	pb "transaction-service/proto" // Import kode yang di-generate
)

// server mengimplementasikan TransactionServiceServer
type GrpcServer struct {
	pb.UnimplementedTransactionServiceServer
	// Dependensi ke service/repository
	transactionService service.TransactionService
}

func NewGrpcServer(ts service.TransactionService) *GrpcServer {
	return &GrpcServer{transactionService: ts}
}

// CreateTransaction adalah implementasi dari RPC
func (s *GrpcServer) CreateTransaction(ctx context.Context, req *pb.CreateTransactionRequest) (*pb.TransactionResponse, error) {
	// Di sini Anda akan memanggil logika bisnis yang mirip dengan service sebelumnya:
	// 1. Panggil `book-service` (sekarang via gRPC client) untuk validasi & harga
	// 2. Panggil `wallet-service` (via gRPC client) untuk debit saldo
	// 3. Panggil repository lokal untuk menyimpan ke PostgreSQL
	// 4. Konversi hasilnya ke format response protobuf `pb.TransactionResponse`

	// (Logika detail di-skip untuk keringkasan, tapi alurnya sama)

	// Contoh memanggil service
	// Panggil method di lapisan service untuk menjalankan semua logika bisnis
	response, err := s.transactionService.CreateTransaction(ctx, req)
	if err != nil {
		// Jika ada error dari service (misal: buku tidak ada, saldo kurang),
		// gRPC akan meneruskannya ke client.
		return nil, err
	}

	// Kembalikan response yang sudah dalam format protobuf
	return response, nil
}

func (s *GrpcServer) GetUserTransactions(ctx context.Context, req *pb.GetUserTransactionsRequest) (*pb.GetUserTransactionsResponse, error) {
	return s.transactionService.GetUserTransactions(ctx, req)
}