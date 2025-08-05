package server

import (
	"context"
	"fmt"
	pb "gifting-service/proto"
	"gifting-service/internal/service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	pb.UnimplementedGiftingServiceServer
	giftingService service.GiftingService
}

func NewGrpcServer(s service.GiftingService) *GrpcServer {
	return &GrpcServer{giftingService: s}
}

func (s *GrpcServer) SendGift(ctx context.Context, req *pb.SendGiftRequest) (*pb.SendGiftResponse, error) {
	gift, err := s.giftingService.SendGift(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.SendGiftResponse{
		GiftId:         fmt.Sprintf("%d", gift.GiftID),
		DonorId:        fmt.Sprintf("%d", gift.DonorID),
		RecipientEmail: gift.RecipientEmail,
		BookId:         fmt.Sprintf("%d", gift.BookID),
		Status:         gift.Status,
		GiftDate:       timestamppb.New(gift.CreatedAt),
	}, nil
}