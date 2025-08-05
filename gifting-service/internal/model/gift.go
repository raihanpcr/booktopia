package model

import (
	"time"
	"gorm.io/gorm"
)

type EbookGiftLog struct {
	GiftID          uint           `gorm:"primaryKey"`
	DonorID         uint           `gorm:"not null"`
	RecipientEmail  string         `gorm:"type:varchar(100);not null"`
	BookID          uint           `gorm:"not null"`
	Message         string         `gorm:"type:text"`
	Status          string         `gorm:"type:varchar(50);default:'pending'"`
	RecipientUserID *uint          // Pointer ke uint agar bisa NULL
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}