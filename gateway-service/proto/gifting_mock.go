package proto

import (
	"context"
	pb "gifting-service/proto"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// MockGiftingServiceClient adalah implementasi mock dari GiftingServiceClient.
type MockGiftingServiceClient struct {
	mock.Mock
}

// SendGift adalah implementasi mock untuk mengirim hadiah.
func (m *MockGiftingServiceClient) SendGift(ctx context.Context, in *pb.SendGiftRequest, opts ...grpc.CallOption) (*pb.SendGiftResponse, error) {
	// Merekam bahwa method ini dipanggil.
	args := m.Called(ctx, in)

	// Mengembalikan response dan error yang sudah diprogram dalam tes.
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.SendGiftResponse), args.Error(1)
}