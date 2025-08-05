package server

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"wallet-service/internal/service"
	pb "wallet-service/proto"
)

type GrpcServer struct {
	pb.UnimplementedWalletServiceServer
	walletService service.WalletService
}

func NewGrpcServer(ws service.WalletService) *GrpcServer {
	return &GrpcServer{walletService: ws}
}

func (s *GrpcServer) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error) {
	balance, err := s.walletService.GetBalance(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.GetBalanceResponse{UserId: req.UserId, Balance: balance}, nil
}

func (s *GrpcServer) TopUp(ctx context.Context, req *pb.TopUpRequest) (*pb.TopUpResponse, error) {
	topUp, err := s.walletService.TopUp(ctx, req)
	if err != nil {
		return nil, err
	}
	return &pb.TopUpResponse{
		TopUpId:   fmt.Sprintf("%d", topUp.ID),
		UserId:    fmt.Sprintf("%d", topUp.UserID),
		Amount:    topUp.Amount,
		Status:    topUp.Status,
		TopUpDate: timestamppb.New(topUp.CreatedAt),
	}, nil
}

func (s *GrpcServer) Debit(ctx context.Context, req *pb.DebitRequest) (*pb.DebitResponse, error) {
	newBalance, err := s.walletService.Debit(ctx, req)
	if err != nil {
		return &pb.DebitResponse{Success: false}, err
	}
	return &pb.DebitResponse{Success: true, NewBalance: newBalance}, nil
}

func (s *GrpcServer) Credit(ctx context.Context, req *pb.CreditRequest) (*pb.CreditResponse, error) {
	newBalance, err := s.walletService.Credit(ctx, req)
	if err != nil {
		return &pb.CreditResponse{Success: false}, err
	}
	return &pb.CreditResponse{Success: true, NewBalance: newBalance}, nil
}