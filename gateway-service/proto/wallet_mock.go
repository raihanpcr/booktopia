package proto

import (
	"context"
	pb "wallet-service/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockWalletServiceClient adalah implementasi mock dari WalletServiceClient.
type MockWalletServiceClient struct {
	mock.Mock
}

// GetBalance adalah implementasi mock untuk mengambil saldo.
func (m *MockWalletServiceClient) GetBalance(ctx context.Context, in *pb.GetBalanceRequest, opts ...grpc.CallOption) (*pb.GetBalanceResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.GetBalanceResponse), args.Error(1)
}

// TopUp adalah implementasi mock untuk top-up saldo.
func (m *MockWalletServiceClient) TopUp(ctx context.Context, in *pb.TopUpRequest, opts ...grpc.CallOption) (*pb.TopUpResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.TopUpResponse), args.Error(1)
}

// Debit adalah implementasi mock untuk mengurangi saldo.
func (m *MockWalletServiceClient) Debit(ctx context.Context, in *pb.DebitRequest, opts ...grpc.CallOption) (*pb.DebitResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.DebitResponse), args.Error(1)
}

// Credit adalah implementasi mock untuk menambah/mengembalikan saldo.
func (m *MockWalletServiceClient) Credit(ctx context.Context, in *pb.CreditRequest, opts ...grpc.CallOption) (*pb.CreditResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.CreditResponse), args.Error(1)
}