package repository

import (
	"auth-service/internal/model"
	"fmt"
	"gorm.io/gorm"
)

// UserRepository adalah interface yang mendefinisikan kontrak untuk operasi database.
type UserRepository interface {
	Create(user model.User) (model.User, error)
	GetByEmail(email string) (model.User, error)
}

// userRepository adalah implementasi dari UserRepository yang menggunakan GORM.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository adalah constructor untuk userRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// Create menyimpan objek model.User baru ke database.
func (r *userRepository) Create(user model.User) (model.User, error) {
	err := r.db.Create(&user).Error
	// GORM akan otomatis mengisi ID ke dalam struct 'user' setelah dibuat.
	return user, err
}

// GetByEmail mengambil data pengguna berdasarkan alamat email.
func (r *userRepository) GetByEmail(email string) (model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return model.User{}, fmt.Errorf("user with email %s not found", email)
	}
	return user, nil
}