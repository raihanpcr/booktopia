package repository

import (
	"context"
	"gifting-service/internal/model"

	"github.com/stretchr/testify/mock"
)

type MockGiftingRepository struct {
	mock.Mock
}

func (m *MockGiftingRepository) CreateGift(ctx context.Context, gift *model.EbookGiftLog) (*model.EbookGiftLog, error) {
	args := m.Called(ctx, gift)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.EbookGiftLog), args.Error(1)
}

func (m *MockGiftingRepository) ExpiredOldGifts(ctx context.Context, days int) (int64, error) {
	args := m.Called(ctx, days)
	return args.Get(0).(int64), args.Error(1)
}