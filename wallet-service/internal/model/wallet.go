package model

import (
	"time"
	"gorm.io/gorm"
)

// User merepresentasikan tabel 'users' di database.
type User struct {
	ID    uint           `gorm:"primaryKey"` // GORM umumnya menggunakan uint untuk ID
	Name  string         `gorm:"type:varchar(100);not null"`
	Email string         `gorm:"type:varchar(100);unique;not null"`
	Saldo float64        `gorm:"type:decimal(12,2);default:0.00"`
	// Kolom lain dari tabel users bisa ditambahkan di sini jika perlu
}

// TopUp merepresentasikan tabel 'top_up' di database.
type TopUp struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null"` // Sesuaikan dengan tipe ID di User
	Amount    float64        `gorm:"type:decimal(12,2);not null"`
	Method    string         `gorm:"type:varchar(50)"`
	Status    string         `gorm:"type:varchar(50);default:'success'"`
	CreatedAt time.Time      // GORM otomatis mengelola `created_at`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}