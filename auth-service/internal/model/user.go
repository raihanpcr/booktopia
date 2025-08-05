package model

import (
	"time"

	"gorm.io/gorm"
)

// User merepresentasikan entitas inti pengguna dalam lapisan domain.
// Entitas ini berisi atribut-atribut yang mendefinisikan seorang pengguna
// dan juga dapat mengandung perilaku (logic bisnis) yang terkait langsung
// dengan data pengguna itu sendiri (sesuai prinsip Domain-Driven Design).

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`             // ID unik pengguna
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`         // Nama lengkap pengguna
	Email     string         `gorm:"type:varchar(100);unique;not null" json:"email"` // Email unik
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`            // Kata sandi yang sudah di-hash
	Role      string         `gorm:"type:varchar(255);not null" json:"role"`
	Saldo     float64        `gorm:"type:decimal(12,2);default:0.00" json:"saldo"` // Saldo pengguna
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`             // Timestamp dibuat
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`             // Timestamp terakhir update
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// --- Logic bisnis domain-level langsung di entitas ---
// Ini adalah contoh perilaku (behavior) yang melekat pada entitas User itu sendiri,
// sesuai dengan prinsip Domain-Driven Design (DDD).
// Metode ini melakukan validasi atau operasi yang hanya relevan dengan atribut User.

// IsEmailValid memeriksa apakah format email pengguna dianggap valid.
// Saat ini, validasi ini sangat sederhana (hanya memeriksa panjang).
// Dalam aplikasi nyata, ini bisa ditingkatkan menggunakan ekspresi reguler (regex)
// untuk validasi format email yang lebih ketat dan sesuai standar.
func (u *User) IsEmailValid() bool {
	return len(u.Email) > 5 // Contoh validasi sederhana: email harus lebih dari 5 karakter
}
