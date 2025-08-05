package repository

import (
	"errors"
	"wallet-service/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// WalletRepository adalah interface untuk operasi database.
type WalletRepository interface {
	GetBalance(userID uint) (float64, error)
	UpdateBalance(userID uint, amount float64) (float64, error)
	CreateTopUp(topUp *model.TopUp) (*model.TopUp, error)
}

type gormRepository struct {
	db *gorm.DB
}

// NewGormRepository adalah constructor untuk GORM repository.
func NewGormRepository(db *gorm.DB) WalletRepository {
	return &gormRepository{db: db}
}

func (r *gormRepository) GetBalance(userID uint) (float64, error) {
	var user model.User
	// Ambil hanya kolom saldo untuk efisiensi
	err := r.db.Select("saldo").First(&user, userID).Error
	return user.Saldo, err
}

func (r *gormRepository) UpdateBalance(userID uint, amount float64) (float64, error) {
	var finalBalance float64

	// GORM's Transaction method menangani commit/rollback secara otomatis.
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		// Kunci baris untuk mencegah race condition (SELECT ... FOR UPDATE)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&user, userID).Error; err != nil {
			return err
		}

		if user.Saldo+amount < 0 {
			return errors.New("insufficient funds")
		}

		newBalance := user.Saldo + amount
		if err := tx.Model(&user).Update("saldo", newBalance).Error; err != nil {
			return err
		}

		finalBalance = newBalance
		return nil // Return nil akan men-commit transaksi
	})

	return finalBalance, err
}

func (r *gormRepository) CreateTopUp(topUp *model.TopUp) (*model.TopUp, error) {
	err := r.db.Create(topUp).Error
	return topUp, err
}
