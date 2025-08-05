package model

import (
	"time"
	"gorm.io/gorm"
)

// Transaction merepresentasikan tabel 'transactions' dengan GORM tags.
type Transaction struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null"`
	TotalAmount float64      `gorm:"type:decimal(12,2);not null"`
	Status      string         `gorm:"type:varchar(50);default:'completed'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	
	// Mendefinisikan relasi: satu transaksi memiliki banyak detail
	Details     []TransactionDetail `gorm:"foreignKey:TransactionID"`
}

// TransactionDetail merepresentasikan tabel 'transaction_details' dengan GORM tags.
type TransactionDetail struct {
	ID            uint           `gorm:"primaryKey"`
	TransactionID uint           `gorm:"not null"`
	BookID        string         `gorm:"type:varchar(255);not null"`
	Quantity      int            `gorm:"not null"`
	PricePerUnit  float64        `gorm:"type:decimal(10,2);not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}