package service

import (
	"context"
	"errors"
	"strconv"
	"wallet-service/internal/model"
	"wallet-service/internal/repository"
	pb "wallet-service/proto"
)

type WalletService interface {
	GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (float64, error)
	TopUp(ctx context.Context, req *pb.TopUpRequest) (*model.TopUp, error)
	Debit(ctx context.Context, req *pb.DebitRequest) (float64, error)
	Credit(ctx context.Context, req *pb.CreditRequest) (float64, error)
}

type walletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) GetBalance(ctx context.Context, req *pb.GetBalanceRequest) (float64, error) {
	userID, err := strconv.ParseUint(req.UserId, 10, 32)
	if err != nil {
		return 0, errors.New("invalid user id format")
	}
	return s.repo.GetBalance(uint(userID))
}

func (s *walletService) TopUp(ctx context.Context, req *pb.TopUpRequest) (*model.TopUp, error) {
	userID, err := strconv.ParseUint(req.UserId, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}
	if req.Amount <= 0 {
		return nil, errors.New("top-up amount must be positive")
	}

	if _, err := s.repo.UpdateBalance(uint(userID), req.Amount); err != nil {
		return nil, err
	}

	topUp := &model.TopUp{
		UserID: uint(userID),
		Amount: req.Amount,
		Method: req.Method,
		Status: "success",
	}

	return s.repo.CreateTopUp(topUp)
}

func (s *walletService) Debit(ctx context.Context, req *pb.DebitRequest) (float64, error) {
	userID, err := strconv.ParseUint(req.UserId, 10, 32)
	if err != nil {
		return 0, errors.New("invalid user id format")
	}
	if req.Amount <= 0 {
		return 0, errors.New("debit amount must be positive")
	}
	return s.repo.UpdateBalance(uint(userID), -req.Amount)
}

func (s *walletService) Credit(ctx context.Context, req *pb.CreditRequest) (float64, error) {
	userID, err := strconv.ParseUint(req.UserId, 10, 32)
	if err != nil {
		return 0, errors.New("invalid user id format")
	}
	if req.Amount <= 0 {
		return 0, errors.New("credit amount must be positive")
	}
	return s.repo.UpdateBalance(uint(userID), req.Amount)
}