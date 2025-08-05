package proto

import (
	"context"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	pb "transaction-service/proto"
)

type MockTransactionServiceClient struct {
	mock.Mock
}

// CreateTransaction adalah implementasi mock
func (m *MockTransactionServiceClient) CreateTransaction(ctx context.Context, in *pb.CreateTransactionRequest, opts ...grpc.CallOption) (*pb.TransactionResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.TransactionResponse), args.Error(1)
}

// GetUserTransactions adalah implementasi mock
func (m *MockTransactionServiceClient) GetUserTransactions(ctx context.Context, in *pb.GetUserTransactionsRequest, opts ...grpc.CallOption) (*pb.GetUserTransactionsResponse, error) {
	// Merekam pemanggilan method
	args := m.Called(ctx, in)

	// Cek dan kembalikan nilai yang sudah diprogram dalam tes
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GetUserTransactionsResponse), args.Error(1)
}