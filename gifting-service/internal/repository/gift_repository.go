package repository

import (
	"context"
	"gifting-service/internal/model"
	"time"

	"gorm.io/gorm"
)

type GiftingRepository interface {
	CreateGift(ctx context.Context, gift *model.EbookGiftLog) (*model.EbookGiftLog, error)
	ExpiredOldGifts(ctx context.Context, days int) (int64, error)
}

type gormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) GiftingRepository {
	return &gormRepository{db: db}
}

func (r *gormRepository) CreateGift(ctx context.Context, gift *model.EbookGiftLog) (*model.EbookGiftLog, error) {
	err := r.db.WithContext(ctx).Create(gift).Error
	return gift, err
}

func (r *gormRepository) ExpiredOldGifts(ctx context.Context, days int) (int64, error) {
	result := r.db.WithContext(ctx).Model(&model.EbookGiftLog{}).Where("status=?", "pending").Where("created_at < ?", time.Now().AddDate(0, 0, -days)).Update("status", "expired")

	return result.RowsAffected, result.Error
}