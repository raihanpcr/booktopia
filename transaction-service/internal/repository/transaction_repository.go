package repository

import (
	"context"
	"transaction-service/internal/model"

	"gorm.io/gorm"
)

// TransactionRepository adalah interface untuk operasi database.
type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	GetTransactionsByUserID(ctx context.Context, userID uint) ([]model.Transaction, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository adalah constructor untuk GORM repository.
func NewGormRepository(db *gorm.DB) TransactionRepository {
	return &gormRepository{db: db}
}

// CreateTransaction menyimpan transaksi beserta semua detailnya.
// GORM akan otomatis menggunakan transaction database karena kita membuat data yang berelasi.
func (r *gormRepository) CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	err := r.db.WithContext(ctx).Create(transaction).Error
	return transaction, err
}

func (r *gormRepository) GetTransactionsByUserID(ctx context.Context, userID uint) ([]model.Transaction, error) {
    var transactions []model.Transaction
    // Menggunakan Preload untuk mengambil data relasi 'Details' juga
    err := r.db.WithContext(ctx).Preload("Details").Where("user_id = ?", userID).Order("created_at DESC").Find(&transactions).Error
    return transactions, err
}